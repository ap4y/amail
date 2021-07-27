import { writable } from "svelte/store";

const { subscribe, update } = writable([]);

function toggle({ thread }) {
  update((threads) =>
    threads.includes(thread)
      ? threads.filter((id) => id !== thread)
      : [...threads, thread]
  );
}

export default { subscribe, toggle };
