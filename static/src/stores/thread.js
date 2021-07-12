import { writable } from "svelte/store";
import ApiClient from "../client";
import url from "./url";

const { subscribe, set, update } = writable(null);

const fetch = async (selectedThread) => {
  if (!selectedThread) {
    set(null);
    return null;
  }

  const thread = await ApiClient.default.thread(selectedThread);
  set(thread);

  return thread;
};

export default { subscribe, fetch };
