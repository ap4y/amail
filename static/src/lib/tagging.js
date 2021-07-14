import ApiClient from "../client";
import mailboxes from "../stores/mailboxes";
import threads from "../stores/threads";
import thread from "../stores/thread";

export async function updateTags(threadId, messageId, updates) {
  if (!messageId || !threadId) return [];

  const newTags = await ApiClient.default.updateTags(
    `id:${messageId}`,
    updates
  );

  const threadTags = thread.updateTags(messageId, newTags);
  threads.setTags(threadId, threadTags);
  mailboxes.updateUnreadCounters();

  return newTags;
}

export async function updateThreadTags(threadId, updates) {
  if (!threadId) return [];

  const newTags = await ApiClient.default.updateTags(
    `thread:${threadId}`,
    updates
  );

  threads.setTags(threadId, newTags);
  await mailboxes.updateUnreadCounters();

  return newTags;
}

export async function markAsRead(threadId, message) {
  if (!message.tags.includes("unread")) return [];

  return await updateTags(threadId, message.id, ["-unread"]);
}
