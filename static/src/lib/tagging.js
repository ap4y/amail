import ApiClient from "../client";
import mailboxes from "../stores/mailboxes";
import selectedMessage from "../stores/message";

export async function markAsRead(message) {
  if (!message || !message.tags.includes("unread")) return;

  await ApiClient.default.updateTags(message.id, ["-unread"]);
  selectedMessage.updateTags(message.tags.filter((tag) => tag !== "unread"));
  mailboxes.updateUnreadCounters();
}

function getThreadTags(set, thread) {
  for (const [message, subThread] of thread) {
    message.tags.forEach((tag) => set.add(tag));
    getThreadTags(set, subThread);
  }

  return set;
}
