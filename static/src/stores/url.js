import { derived, writable } from "svelte/store";
import { mailboxIds } from "../config";

const initialURL = new URL(window.location.href);
const href = writable(initialURL);
const [, mailbox, thread] = initialURL.pathname.split("/");
initialURL.mailbox = mailbox;
if (mailbox === "search") {
  initialURL.searchTerms = initialURL.searchParams.get("terms");
} else {
  initialURL.thread = thread;
}

const update = ({ state }) => {
  const url = new URL(window.location.href);
  url.mailbox = state?.mailbox;
  url.thread = state?.thread;
  url.message = state?.message;
  url.searchTerms = state?.searchTerms;
  href.set(url);
};

window.addEventListener("popstate", update);

export function pushState(state, url) {
  history.pushState(state, "", url);
  update({ state });
}

const url = derived(href, ($href) => $href);

export default url;

export const selectedMailbox = derived(url, ($url) =>
  $url.mailbox?.length > 0 ? $url.mailbox : mailboxIds[0]
);

export const selectedThread = derived(url, ($url) => $url.thread);

export const searchTerms = derived(url, ($url) => $url.searchTerms);
