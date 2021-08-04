import { writable } from "svelte/store";

import { quotedText } from "../lib/email";

const { subscribe, set, update } = writable(null);

function create() {
  set({ to: [], cc: [], subject: "", body: "", headers: {}, attachments: [] });
}

function reply(reply) {
  const { To, Cc, Subject, ...rest } = reply["reply-headers"];
  const { headers, body } = reply.original;

  const content = body.map((item) => quotedText(item)).join("\n");
  set({
    to: To.split(", "),
    cc: Cc?.split(", ") || [],
    subject: Subject,
    headers: { ...rest },
    body: content,
    originalHeaders: headers,
    reply: true,
    attachments: [],
  });
}

function forward({ body, headers }) {
  const content = body.map((item) => quotedText(item)).join("\n");
  set({
    to: [],
    cc: [],
    subject: `Fwd: ${headers.Subject}`,
    headers: {},
    body: content,
    originalHeaders: headers,
    reply: true,
    attachments: [],
  });
}

function setField(updates) {
  update((message) => ({ ...message, ...updates }));
}

function setBody({ body }) {
  update((message) => ({ ...message, body }));
}

function destroy() {
  set(null);
}

export default {
  subscribe,
  create,
  reply,
  forward,
  setField,
  set: setBody,
  destroy,
};
