import { writable, derived } from "svelte/store";
import { mailboxIds } from "../config";
import { selectedMailbox } from "./url";

const mailboxes = writable(
  mailboxIds.reduce((acc, id) => {
    acc[id] = {};
    return acc;
  }, {})
);

export default mailboxes;

export const currentMailbox = derived(
  [mailboxes, selectedMailbox],
  ([$mailboxes, $selectedMailbox]) => $mailboxes[$selectedMailbox]
);
