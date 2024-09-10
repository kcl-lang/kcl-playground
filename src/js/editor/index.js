import escape from 'escape-html';
import LZString from 'lz-string';
import CodeMirror from 'codemirror';
import { SHARE_QUERY_KEY, SNIPPETS } from '../constants';
import { date2time } from './utils';
import './format';
import { snippet } from './snippet';
import './share';
import './codemirror-kcl';
import { load, invokeKCLRun, invokeKCLFmt } from "../module";

const response = await fetch("kcl.wasm");
const wasm = await response.arrayBuffer();
const inst = await load({ data: wasm });
const source = document.getElementById('source');
const run = document.getElementById('run');
const outputContainer = document.getElementById('output-container');
const output = document.getElementById('output');
const lastUpdated = document.getElementById('last-updated');

const noop = () => {};

const editor = CodeMirror.fromTextArea(source, {
  mode: 'kcl',
  tabSize: 4,
  lineNumbers: true,
  lineWrapping: true,
  smartIndent: true,
});

class Mutex {
  constructor() {
    this.queue = [];
    this.locked = false;
  }

  async lock() {
    return new Promise((resolve) => {
      if (!this.locked) {
        this.locked = true;
        resolve();
      } else {
        this.queue.push(resolve);
      }
    });
  }

  unlock() {
    if (this.queue.length > 0) {
      const next = this.queue.shift();
      next();
    } else {
      this.locked = false;
    }
  }
}
const mutex = new Mutex();
export const Command = {
  async run() {
    await mutex.lock();
    try {
      const code = editor.getValue();
      const result = invokeKCLRun(inst, {
        filename: "test.k",
        source: code,
      });
      if (result.startsWith("ERROR:")) {
        Command.clear();
        Command.print(result.replace(/^ERROR:\s*/, ''));
      } else {
        Command.clear();
        Command.print(result);
      }
    } finally {
      mutex.unlock();
    }
  },

  async format() {
    await mutex.lock();
    try {
      const code = editor.getValue();
      const result = invokeKCLFmt(inst, {
        source: code,
      });
      if (result.startsWith("ERROR:")) {
        Command.clear();
        Command.print(result.replace(/^ERROR:\s*/, ''));
      } else {
        editor.setValue(result)
      }
    } finally {
      mutex.unlock();
    }
  },
  
  async set(value) {
    await mutex.lock();
    try {
      editor.setValue(value);
    } finally {
      mutex.unlock();
    }
  },

  getValue: () => editor.getValue(),
  setValue: (value) => editor.setValue(value),

  print: (str) => {
    const now = new Date();
    const time = date2time(now);
    lastUpdated.textContent = `LAST UPDATE: ${time}`;
    lastUpdated.dateTime = now.toISOString();
    lastUpdated.style.display = 'block';

    output.innerHTML += escape(`${str}\n`);

    outputContainer.scrollTop = outputContainer.scrollHeight - outputContainer.clientHeight;
  },

  clear: () => {
    lastUpdated.dateTime = '';
    lastUpdated.style.display = 'none';
    output.innerHTML = '';
    outputContainer.scrollTop = 0;
  },
};

editor.on("change", Command.run)
const query = new window.URLSearchParams(window.location.search);

if (query.has(SHARE_QUERY_KEY)) {
  Command.setValue(LZString.decompressFromEncodedURIComponent(query.get(SHARE_QUERY_KEY)));
  snippet.selectedIndex = 0;
} else {
  Command.setValue(SNIPPETS[0].value);
}

editor.addKeyMap({
  'Ctrl-Enter': noop,
  'Shift-Enter': noop,
  'Ctrl-L': noop,
});

document.addEventListener(
  'keydown',
  (e) => {
    const { ctrlKey, metaKey, shiftKey, key: raw } = e;
    const key = raw.toLowerCase();

    // Ctrl + Enter
    if (ctrlKey && key === 'enter') {
      Command.run();
    }

    // Shift + Enter
    if (shiftKey && key === 'enter') {
      Command.format();
    }

    // Ctrl + L
    if (ctrlKey && key === 'l') {
      Command.clear();
    }
  },
  false,
);

run.addEventListener(
  'click',
  async (e) => {
    e.preventDefault();
    await Command.run();
  },
  false,
);
