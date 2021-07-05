import { derived } from "svelte/store";
import ApiClient from "../client";
import url, { selectedThread } from "./url";

const thread = derived(selectedThread, ($selectedThread, set) => {
  if (!$selectedThread) {
    set(null);
  } else {
    ApiClient.default.thread($selectedThread).then((res) => set(res));
  }
});

export default thread;

export const selectedMessage = derived([url, thread], ([$url, $thread]) =>
  $url.message?.length > 0 ? $url.message : $thread ? $thread[0].id : null
);
