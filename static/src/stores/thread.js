import { writable } from "svelte/store";
import ApiClient from "../client";

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

export function findOtherMessage(
  thread,
  messageId,
  withTags = [],
  withoutTags = []
) {
  if (!thread) return null;

  for (const [message, subThread] of thread) {
    if (
      message.id != messageId &&
      withTags.every((tag) => message.tags.includes(tag)) &&
      withoutTags.every((tag) => !message.tags.includes(tag))
    ) {
      return message;
    }

    const match = findOtherMessage(subThread, messageId, withTags, withoutTags);
    if (match) return match;
  }

  return null;
}

export function findLastMessage(
  thread,
  messageId,
  withTags = [],
  withoutTags = []
) {
  const messages = matchedMessages(
    thread,
    messageId,
    withTags,
    withoutTags
  ).sort((a, b) => b.timestamp - a.timestamp);

  return messages.length > 0 ? messages[0] : null;
}

function matchedMessages(thread, messageId, withTags, withoutTags, acc = []) {
  if (!thread) return acc;

  for (const [message, subThread] of thread) {
    if (
      message.id != messageId &&
      withTags.every((tag) => message.tags.includes(tag)) &&
      withoutTags.every((tag) => !message.tags.includes(tag))
    ) {
      acc.push(message);
    }

    matchedMessages(subThread, messageId, withTags, withoutTags, acc);
  }

  return acc;
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

export function getFirstMessage(thread) {
  const [firstThread] = thread;
  if (!firstThread) return null;

  const [message] = firstThread;
  if (!message) return null;

  return message;
}

export default { subscribe, fetch, updateTags };
