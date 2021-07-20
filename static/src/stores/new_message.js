import { writable } from "svelte/store";

const { subscribe, set, update } = writable(null);

function create(from) {
  set({ to: [], cc: [], subject: "", body: "" });
}

function setField(updates) {
  update((message) => ({ ...message, ...updates }));
}

function destroy() {
  set(null);
}

export default { subscribe, create, setField, destroy };
