import { writable, derived } from "svelte/store";
import ApiClient from "../client";
import url from "./url";

const { subscribe, set } = writable(null);

const fetch = async (selectedThread) => {
  if (!selectedThread) {
    set(null);
    return null;
  }

  const thread = await ApiClient.default.thread(selectedThread);
  set(thread);

  return thread;
};

const thread = { subscribe, fetch };

export default thread;

export const selectedMessage = derived([url, thread], ([$url, $thread]) =>
  $url.message?.length > 0 ? $url.message : $thread ? $thread[0].id : null
);
