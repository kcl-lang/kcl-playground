import escape from 'escape-html';
import LZString from 'lz-string';
import CodeMirror from 'codemirror';
import { SHARE_QUERY_KEY, SNIPPETS } from '../constants';
import { date2time } from './utils';
import './format';
import { snippet } from './snippet';
import './share';
import './codemirror-kcl';

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

const playgroundOptions = {
  'compileURL': '/-/play/compile',
  'fmtURL': '/-/play/fmt',
};

function getBackendUrl() {
  var backendUrl = window.plutoEnv.BACKEND_URL;
  if (/http:\/\/localhost/g.test(backendUrl)) {
    var parts = backendUrl.split(":");
    var port = parts[parts.length - 1];

    var protocol = window.location.protocol;
    var host = window.location.hostname;
    backendUrl = `${protocol}//${host}:${port}`;
  }
  return backendUrl;
}

export const Command = {
  set: (value) => {
    editor.setValue(value);
  },

  run: () => {
    var data = { "body": Command.getValue() };
    fetch(getBackendUrl() + playgroundOptions.compileURL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        Command.clear();
        Command.print(data.error);
      } else {
        Command.clear();
        Command.print(data.body);
      }
    })
    .catch(error => {
      console.error('Error:', error);
    });
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

  format: () => {
    var data = { "body": Command.getValue() };
    fetch(getBackendUrl() + playgroundOptions.fmtURL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        Command.clear();
        Command.print(data.error);
      } else {
        Command.setValue(data.body);
      }
    })
    .catch(error => {
      console.error('Error:', error);
    });
  },
};

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
  (e) => {
    e.preventDefault();
    Command.run();
  },
  false,
);
