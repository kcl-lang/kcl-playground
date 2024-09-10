import { Command } from '.';

const format = document.getElementById('format');

format.addEventListener(
  'click',
  async (e) => {
    e.preventDefault();
    await Command.format();
  },
  false,
);
