import { writable } from "svelte/store";

import { traverseContent, quotedText } from "../lib/email";

const { subscribe, set, update } = writable(null);

function create(from) {
  set({ to: [], cc: [], subject: "", body: "" });
}

function reply(reply) {
  const { To, Cc, Subject, ...rest } = reply["reply-headers"];
  const { headers, body } = reply.original;

  const content = body.map((item) =>
    quotedText(traverseContent(item).text.join("\n"))
  );

  set({
    to: To.split(", "),
    cc: Cc?.split(", ") || [],
    subject: Subject,
    headers: { ...rest },
    body: `
On ${headers.Date}, ${headers.From} wrote:
${content.join("\n")}`,
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

export default { subscribe, create, reply, setField, set: setBody, destroy };
