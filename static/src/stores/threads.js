import { writable } from "svelte/store";
import ApiClient from "../client";
import { threadsPerPage } from "../config";

export const currentPage = writable(0);
export const totalThreads = writable(0);

const { subscribe, set, update } = writable([]);

const fetch = async (terms, currentPage) => {
  if (!terms) {
    set([]);
    return [];
  }

  const { total, threads } = await ApiClient.default.threads(
    terms,
    currentPage || 0,
    threadsPerPage
  );

  totalThreads.set(total);
  set(threads);

  return threads;
};

const setTags = (threadId, tags) => {
  update((threads) =>
    threads.map((thread) =>
      thread.thread === threadId ? { ...thread, tags } : thread
    )
  );
};

export default { subscribe, fetch, setTags };
