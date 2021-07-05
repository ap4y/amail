import { writable, derived } from "svelte/store";
import { mailboxIds } from "../config";
import { selectedMailbox, searchTerms } from "./url";

const mailboxes = writable(
  mailboxIds.reduce((acc, id) => {
    acc[id] = {};
    return acc;
  }, {})
);

export default mailboxes;

export const currentMailbox = derived(
  [mailboxes, selectedMailbox, searchTerms],
  ([$mailboxes, $selectedMailbox, $searchTerms]) =>
    $searchTerms?.length > 0
      ? { id: "search", terms: $searchTerms }
      : $mailboxes[$selectedMailbox]
);
