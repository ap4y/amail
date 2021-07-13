import { writable } from "svelte/store";
import ApiClient from "../client";
import selectedMessage from "./message";

const { subscribe, set, update } = writable(null);

export function findMessage(thread, messageId) {
  if (!thread) return null;

  for (const [message, subThread] of thread) {
    if (message.id === messageId) return message;

    const match = findMessage(subThread, messageId);
    if (match) return match;
  }

  return null;
}

function getThreadTags(set, thread) {
  for (const [message, subThread] of thread) {
    message.tags.forEach((tag) => set.add(tag));
    getThreadTags(set, subThread);
  }

  return set;
}

const fetch = async (selectedThread) => {
  if (!selectedThread) {
    set(null);
    return null;
  }

  const thread = await ApiClient.default.thread(selectedThread);
  set(thread);
  selectFirstMessage(thread);

  return thread;
};

function updateTags(messageId, tags) {
  const threadTags = new Set();

  update((thread) => {
    const message = findMessage(thread, messageId);
    if (message) message.tags = tags;
    getThreadTags(threadTags, thread);
    return thread;
  });

  return [...threadTags];
}

function selectFirstMessage(thread) {
  const [firstThread] = thread;
  if (!firstThread) return null;

  const [message] = firstThread;
  if (!message) return null;

  selectedMessage.selectMessage(message.id);
  return message;
}

export default { subscribe, fetch, updateTags, selectFirstMessage };
