import { writable, derived } from "svelte/store";
import ApiClient from "../client";
import { threadsPerPage } from "../config";
import { currentMailbox } from "./mailboxes";
import { selectedMailbox } from "./url";

export const currentPage = writable(0);

selectedMailbox.subscribe(() => {
  currentPage.set(0);
});

export const totalThreads = writable(0);

const threads = derived(
  [currentMailbox, currentPage],
  ([$currentMailbox, $currentPage], set) => {
    if (!$currentMailbox?.terms) {
      set([]);
    } else {
      ApiClient.default
        .threads($currentMailbox.terms, $currentPage, threadsPerPage)
        .then(({ total, threads }) => {
          totalThreads.set(total);
          set(threads);
        });
    }
  },
  []
);

export default threads;
