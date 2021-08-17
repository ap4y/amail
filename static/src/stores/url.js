import { derived, writable } from "svelte/store";
import { mailboxIds } from "../config";
import { currentPage } from "./threads";
import selectedMessage from "./message";

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
  url.searchTerms = state?.searchTerms;
  href.set(url);
};

window.addEventListener("popstate", update);

function pushState(state, url) {
  history.pushState(state, "", url);
  update({ state });
  console.debug(`navigating to ${url}, state: `, state);
}

const url = derived(href, ($href) => $href);

function selectMailbox(mailbox) {
  currentPage.set(0);
  pushState({ mailbox }, `/${mailbox}`);
}

function selectThread(mailbox, thread) {
  const state = history.state;
  const newState = { mailbox, thread, searchTerms: state.searchTerms };
  let path = `/${mailbox}/${thread}`;
  if (state.searchTerms) {
    path += `?terms=${escape(state.searchTerms)}`;
  }

  pushState(newState, path);
}

function deselectThread() {
  const state = history.state;

  selectedMessage.selectMessage(null);
  pushState({ mailbox: state.mailbox, thread: null }, `/${state.mailbox}`);
}

function search(terms) {
  pushState(
    { mailbox: "search", searchTerms: terms },
    `/search/?terms=${escape(terms)}`
  );
}

export default {
  subscribe: url.subscribe,
  selectMailbox,
  selectThread,
  deselectThread,
  search,
};

export const selectedMailbox = derived(url, ($url) =>
  $url.mailbox?.length > 0 ? $url.mailbox : mailboxIds[0]
);

export const selectedThread = derived(url, ($url) => $url.thread);

export const searchTerms = derived(url, ($url) => $url.searchTerms);
