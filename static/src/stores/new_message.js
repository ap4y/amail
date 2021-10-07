import { writable } from "svelte/store";
import ApiClient from "../client";

import { quotedText, getAttachments } from "../lib/email";

const { subscribe, set, update } = writable(null);

function create() {
  set({ to: [], cc: [], subject: "", body: "", headers: {}, attachments: [] });
}

async function reply(messageId, replyTo) {
  const reply = await ApiClient.default.replyToMessage(messageId, replyTo);

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

function forward({ id, body, headers }) {
  const content = body.map((item) => quotedText(item)).join("\n");
  const attachments = [];
  body.forEach((item) => {
    getAttachments(item).forEach((attach) =>
      attachments.push({
        id: `${id}:${attach.id}`,
        name: attach.filename,
        size: attach["content-length"],
        type: attach["content-type"],
      })
    );
  });
  set({
    to: [],
    cc: [],
    subject: `Fwd: ${headers.Subject}`,
    headers: {},
    body: content,
    originalHeaders: headers,
    reply: true,
    attachments,
  });
}

function edit({ id, body, headers }) {
  const { To, Cc, Subject, ...rest } = headers;

  const content = body.map((item) => quotedText(item)).join("\n");
  const attachments = [];
  body.forEach((item) => {
    getAttachments(item).forEach((attach) =>
      attachments.push({
        id: `${id}:${attach.id}`,
        name: attach.filename,
        size: attach["content-length"],
        type: attach["content-type"],
      })
    );
  });
  set({
    to: To.split(", "),
    cc: Cc?.split(", ") || [],
    subject: Subject,
    headers: { ...rest },
    body: content,
    originalHeaders: headers,
    reply: false,
    attachments,
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
  edit,
  setField,
  set: setBody,
  destroy,
};
