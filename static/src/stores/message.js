import { writable } from "svelte/store";

const { subscribe, set, update } = writable(null);

function selectMessage(messageId) {
  set(messageId);
}

export default { subscribe, selectMessage };
