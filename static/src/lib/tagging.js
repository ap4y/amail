import ApiClient from "../client";
import mailboxes from "../stores/mailboxes";
import threads from "../stores/threads";
import thread from "../stores/thread";

export async function updateTags(threadId, messageId, updates) {
  if (!messageId || !threadId) return [];

  const newTags = await ApiClient.default.updateTags(messageId, updates);

  const threadTags = thread.updateTags(messageId, newTags);
  threads.setTags(threadId, threadTags);
  mailboxes.updateUnreadCounters();

  return newTags;
}

export async function markAsRead(threadId, messageId) {
  return await updateTags(threadId, messageId, ["-unread"]);
}
