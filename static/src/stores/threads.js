import { writable, derived } from "svelte/store";
import ApiClient from "../client";
import { threadsPerPage } from "../config";
import { currentMailbox } from "./mailboxes";
import { selectedMailbox } from "./url";

export const currentPage = writable(0);

selectedMailbox.subscribe(() => {
  currentPage.set(0);
});

const threads = derived(
  [currentMailbox, currentPage],
  ([$currentMailbox, $currentPage], set) => {
    ApiClient.default
      .threads($currentMailbox.folder, $currentPage, threadsPerPage)
      .then((res) => set(res));
  },
  []
);

export default threads;
