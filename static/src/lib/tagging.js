import ApiClient from "../client";
import thread from "../stores/thread";
import threads from "../stores/threads";
import mailboxes from "../stores/mailboxes";

export async function markAsRead(thread, messageId) {
  if (!messageId || !thread) return;

  const message = findMessage(messageId, thread);
  if (!message || !message.tags.includes("unread")) return;

  await ApiClient.default.updateTags(messageId, ["-unread"]);
  message.tags = message.tags.filter((tag) => tag !== "unread");
  mailboxes.updateUnreadCounters();
}

// const tags = [...getThreadTags(new Set(), thread)];
// threads.setTags(threadId, tags);

function findMessage(messageId, thread) {
  const [message, subThreads] = thread;
  if (message.id === messageId) return message;

  for (const subThread of subThreads) {
    const match = findMessage(messageId, subThread);
    if (match) return match;
  }

  return null;
}

function getThreadTags(set, thread) {
  const [message, subThreads] = thread;

  message.tags.forEach((tag) => set.add(tag));
  subThreads.forEach((subThread) => getThreadTags(set, subThread));

  return set;
}
