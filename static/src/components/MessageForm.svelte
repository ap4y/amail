<script>
  import { onMount, createEventDispatcher } from "svelte";
  import ApiClient from "../client";
  import { address, name } from "../stores/mailboxes";
  import newMessage from "../stores/new_message";
  import { selectedThread } from "../stores/url";
  import thread from "../stores/thread";

  import AddressField from "./AddressField.svelte";
  import TextEditor from "./TextEditor.svelte";

  export let wrap = 80;

  const dispatch = createEventDispatcher();
  let submiting = false;
  let error = null;
  let toInput;
  let fileInput;

  onMount(() => {
    toInput.focus();
  });

  async function saveMessage() {
    error = null;
    submiting = true;
    try {
      await ApiClient.default.saveMessage($newMessage);
      submiting = false;
      newMessage.destroy();
    } catch (e) {
      error = e.message;
      submiting = false;
    }
  }

  async function submitMessage() {
    error = null;
    submiting = true;
    try {
      await ApiClient.default.sendMessage($newMessage);
      submiting = false;
      newMessage.destroy();
      thread.fetch($selectedThread);
    } catch (e) {
      error = e.message;
      submiting = false;
    }
  }

  let blocks = [{ type: "text", content: "" }];
  if ($newMessage?.reply) {
    const { Date, From } = $newMessage?.originalHeaders;
    blocks.push(
      { type: "text", content: `On ${Date}, ${From} wrote:` },
      { type: "quote", content: $newMessage?.body }
    );
  }

  function onInput({ detail }) {
    blocks = detail;
    console.debug(
      "text content",
      "blocks:",
      blocks,
      "content:\n",
      blocks.map(({ content }) => content).join("")
    );
    newMessage.setField({
      body: blocks.map(({ content }) => content).join(""),
    });
  }

  function formatSize(size) {
    if (size > 1000000) {
      return `${(size / 1000000).toFixed(2)}MB`;
    } else if (size > 1000) {
      return `${(size / 1000).toFixed(2)}KB`;
    } else {
      return `${size}B`;
    }
  }

  function attachFile() {
    for (const file of fileInput.files) {
      newMessage.setField({ attachments: [...$newMessage.attachments, file] });
    }
  }

  function removeAttachment(index) {
    newMessage.setField({
      attachments: $newMessage.attachments.filter((_, idx) => idx !== index),
    });
  }

  function fillBlock() {
    const selection = document.getSelection();
    if (selection.rangeCount === 0) return;

    const { startContainer, startOffset } = selection.getRangeAt(0);
    const node = startContainer.parentNode;
    if (node.dataset.type !== "text") return;

    const text = node.innerText.replaceAll("\n", "");
    const words = text.split(/\s/);
    const lines = words.reduce(
      (acc, word) => {
        const line = acc[acc.length - 1];
        if (line.length + word.length > wrap) {
          acc.push(`${word} `);
        } else {
          acc[acc.length - 1] += `${word} `;
        }

        return acc;
      },
      [""]
    );

    node.innerHTML = lines.join("\n").trim();
    selection.collapse(
      node.firstChild,
      Math.min(node.firstChild.nodeValue.length, startOffset)
    );
  }

  function onKeyDown(e) {
    if (e.key === "k" && e.altKey) {
      e.preventDefault();
      dispatch("close");
    } else if (e.key === "c" && e.altKey) {
      e.preventDefault();
      submitMessage();
    } else if (e.key === "q" && e.altKey) {
      fillBlock();
    }
  }
</script>

<div
  class="absolute top-0 sm:top-auto left-0 sm:left-auto bottom-0 sm:bottom-5 right-0 sm:right-5 flex flex-col bg-white shadow rounded border-2 border-gray-500 z-50"
  style="max-width: calc(85ch + 2rem)"
  on:keydown={onKeyDown}
>
  <div
    class="px-3 py-2 flex flex-row items-center justify-between bg-gray-500 text-white"
  >
    <span>{$newMessage?.subject || "New message"}</span>

    <button
      class="w-5 text-white hover:text-gray-200 active:text-gray-400 focus:outline-none"
      on:click={() => dispatch("close")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-full fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
        /></svg
      >
    </button>
  </div>

  <div class="p-3 bg-gray-100 text-gray-500">
    <div class="flex flex-row items-center mb-2">
      <label class="w-20" for="from">From: </label>
      <span class="text-gray-600">
        <strong class="font-semibold ">{$name}</strong>
        &lt;{$address}&gt;
      </span>
    </div>

    <div class="flex flex-row items-center mb-2">
      <label class="w-20 flex-shrink-0" for="to">To: </label>
      <AddressField
        id="to"
        value={$newMessage?.to}
        on:change={({ detail }) => newMessage.setField({ to: detail })}
        bind:this={toInput}
      />
    </div>

    <div class="flex flex-row items-center mb-2">
      <label class="w-20 flex-shrink-0" for="cc">CC: </label>
      <AddressField
        id="cc"
        value={$newMessage?.cc}
        on:change={({ detail }) => newMessage.setField({ cc: detail })}
      />
    </div>

    <div class="flex flex-row items-center">
      <label class="w-20" for="subject">Subject: </label>
      <input
        type="text"
        id="subject"
        class="flex-1 outline-none px-3 py-1 bg-white text-gray-700 rounded border border-gray-400 hover:border-gray-500 focus:border-red-300"
        value={$newMessage?.subject}
        on:change={({ target }) =>
          newMessage.setField({ subject: target.value })}
      />
    </div>
  </div>

  <div class="p-3 flex-1">
    <TextEditor {blocks} on:input={(e) => onInput(e)} />
  </div>

  <div class="px-3 mb-3 flex flex-row justify-between items-center">
    <div class="flex flex-row flex-wrap">
      <input
        multiple
        type="file"
        class="hidden"
        bind:this={fileInput}
        on:change={attachFile}
      />
      <button
        class="h-10 p-2 mr-3 mb-1 rounded hover:border-red-500 bg-white text-gray-700 active:bg-red-300 focus:outline-none border"
        on:click={() => fileInput.click()}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="h-5 fill-current"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M16.5 6v11.5c0 2.21-1.79 4-4 4s-4-1.79-4-4V5c0-1.38 1.12-2.5 2.5-2.5s2.5 1.12 2.5 2.5v10.5c0 .55-.45 1-1 1s-1-.45-1-1V6H10v9.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V5c0-2.21-1.79-4-4-4S7 2.79 7 5v12.5c0 3.04 2.46 5.5 5.5 5.5s5.5-2.46 5.5-5.5V6h-1.5z"
          /></svg
        >
      </button>

      {#if $newMessage?.attachments}
        {#each $newMessage.attachments as attach, idx}
          <button
            style="max-width: 10rem;"
            class="flex h-10 border border-gray-600 p-2 mr-1 mb-1 rounded text-gray-600 text-sm hover:border-red-500 active:bg-red-300 focus:outline-none"
            on:click={() => removeAttachment(idx)}
            ><strong class="mr-1 truncate">{attach.name}</strong> ({formatSize(
              attach.size
            )})</button
          >
        {/each}
      {/if}
    </div>

    <div class="flex flex-row items-center">
      {#if error}
        <p class="mr-3 text-red-500 font-semibold">Failed to send email</p>
      {/if}

      <button
        class="h-10 p-2 mr-3 rounded hover:border-red-500 bg-white text-gray-700 active:bg-red-300 focus:outline-none border"
        disabled={submiting}
        on:click={saveMessage}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="h-5 fill-current"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M17 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"
          /></svg
        >
      </button>

      <button
        class={`${
          submiting
            ? "bg-red-300"
            : "bg-red-500 hover:bg-red-600 active:bg-red-700"
        } focus:outline-none px-3 text-white h-10 rounded border-0 font-semibold`}
        disabled={submiting}
        on:click={submitMessage}
      >
        {submiting ? "Sending" : "Send"}
      </button>
    </div>
  </div>
</div>
