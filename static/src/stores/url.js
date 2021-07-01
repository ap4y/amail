import { derived, writable } from "svelte/store";

const url = new URL(window.location.href);
const href = writable(url);
const [, mailbox, thread] = url.pathname.split("/");
url.mailbox = mailbox;
url.thread = thread;

const update = ({ state }) => {
  const url = new URL(window.location.href);
  url.mailbox = state?.mailbox;
  url.thread = state?.thread;
  url.message = state?.message;
  href.set(url);
};

window.addEventListener("popstate", update);

export function pushState(state, url) {
  history.pushState(state, "", url);
  update({ state });
}

export default derived(href, ($href) => $href);
