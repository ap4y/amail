import { writable, get } from "svelte/store";
import ApiClient from "../client";
import { threadsPerPage } from "../config";
import { selectedMailbox } from "./url";

export const currentPage = writable(0);
export const hasMore = writable(false);

const { subscribe, set, update } = writable([]);

const fetch = async ({ id, terms }, currentPage) => {
  if (!terms) {
    set([]);
    return [];
  }

  const { has_more, threads } = await ApiClient.default.threads(
    terms,
    currentPage || 0,
    threadsPerPage
  );

  if (get(selectedMailbox) === id) {
    hasMore.set(has_more);
    set(threads);
  }

  return threads;
};

const setTags = (threadId, tags) => {
  update((threads) =>
    threads.map((thread) =>
      thread.thread === threadId ? { ...thread, tags } : thread
    )
  );
};

const remove = (threadId) => {
  update((threads) => threads.filter((thread) => thread != threadId));
};

export default { subscribe, fetch, setTags, remove };
