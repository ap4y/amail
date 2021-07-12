import { writable, derived } from "svelte/store";
import ApiClient from "../client";
import { selectedMailbox, searchTerms } from "./url";

export const address = writable(null);

const { subscribe, set, update } = writable([]);

const fetch = async () => {
  const res = await ApiClient.default.mailboxes();
  set(res.mailboxes);
  address.set(res.address);
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

const mailboxes = { subscribe, fetch, updateUnreadCounters };

export default mailboxes;

export const currentMailbox = derived(
  [mailboxes, selectedMailbox, searchTerms],
  ([$mailboxes, $selectedMailbox, $searchTerms]) =>
    $searchTerms?.length > 0
      ? { id: "search", terms: $searchTerms }
      : $mailboxes.find(({ id }) => id === $selectedMailbox)
);
