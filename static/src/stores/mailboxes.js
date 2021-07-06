import { writable, derived } from "svelte/store";
import { mailboxIds } from "../config";
import { selectedMailbox, searchTerms } from "./url";

const mailboxes = writable([]);

export default mailboxes;

export const currentMailbox = derived(
  [mailboxes, selectedMailbox, searchTerms],
  ([$mailboxes, $selectedMailbox, $searchTerms]) =>
    $searchTerms?.length > 0
      ? { id: "search", terms: $searchTerms }
      : $mailboxes.find(({ id }) => id === $selectedMailbox)
);
