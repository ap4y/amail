import { writable, derived } from "svelte/store";
import ApiClient from "../client";

export const address = writable("");
export const name = writable(null);
export const predefinedTags = writable([]);

const { subscribe, set, update } = writable([]);

const fetch = async () => {
  const res = await ApiClient.default.mailboxes();
  set(res.mailboxes);
  address.set(res.address);
  name.set(res.name);
  predefinedTags.set(res.tags);
  return res.mailboxes;
};

const updateUnreadCounters = async () => {
  const { mailboxes } = await ApiClient.default.mailboxes();
  update((current) =>
    current.map((mailbox, idx) => ({
      ...mailbox,
      unread: mailboxes[idx].unread,
    }))
  );
};

export default { subscribe, fetch, updateUnreadCounters };
