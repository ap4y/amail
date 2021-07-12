import { writable } from "svelte/store";

const { subscribe, set, update } = writable(null);

function findMessage(thread, messageId) {
  for (const [message, subThread] of thread) {
    if (message.id === messageId) return message;

    const match = findMessage(subThread, messageId);
    if (match) return match;
  }

  return null;
}

function selectMessage(thread, messageId) {
  if (thread) {
    set(findMessage(thread, messageId));
  }
}

function updateTags(tags) {
  update((message) => ({ ...message, tags }));
}

export default { subscribe, selectMessage, updateTags };
