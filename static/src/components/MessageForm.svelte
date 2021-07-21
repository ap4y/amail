<script>
  import { createEventDispatcher } from "svelte";
  import ApiClient from "../client";
  import { address } from "../stores/mailboxes";
  import newMessage from "../stores/new_message";

  const dispatch = createEventDispatcher();
  let submiting = false;
  let error = null;

  async function submitMessage() {
    error = null;
    submiting = true;
    try {
      await ApiClient.default.sendMessage($newMessage);
      submiting = false;
      newMessage.destroy();
    } catch (e) {
      error = e.message;
      submiting = false;
    }
  }
</script>

<div
  class="absolute bottom-5 right-5 bg-white shadow rounded border-2 border-gray-500"
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
      <span class="font-semibold text-gray-600">{$address}</span>
    </div>

    <div class="flex flex-row items-center mb-2">
      <label class="w-20" for="to">To: </label>
      <input
        type="text"
        id="to"
        class="w-full outline-none px-3 py-1 bg-white text-gray-700 rounded border border-gray-400 hover:border-gray-500 focus:border-red-300"
        value={$newMessage?.to?.join(" ")}
        on:change={({ target }) =>
          newMessage.setField({ to: target.value.split(" ") })}
      />
    </div>

    <div class="flex flex-row items-center mb-2">
      <label class="w-20" for="cc">CC: </label>
      <input
        type="text"
        id="cc"
        class="w-full outline-none px-3 py-1 bg-white text-gray-700 rounded border border-gray-400 hover:border-gray-500 focus:border-red-300"
        value={$newMessage?.cc?.join(" ")}
        on:change={({ target }) =>
          newMessage.setField({ cc: target.value.split(" ") })}
      />
    </div>

    <div class="flex flex-row items-center">
      <label class="w-20" for="subject">Subject: </label>
      <input
        type="text"
        id="subject"
        class="w-full outline-none px-3 py-1 bg-white text-gray-700 rounded border border-gray-400 hover:border-gray-500 focus:border-red-300"
        value={$newMessage?.subject}
        on:change={({ target }) =>
          newMessage.setField({ subject: target.value })}
      />
    </div>
  </div>

  <div class="p-3">
    <textarea
      style="width: 85ch;"
      class="border rounded p-3 outline-none border-gray-400 hover:border-gray-500 focus:border-red-300"
      placeholder="Type in message body"
      rows="10"
      value={$newMessage?.body}
      on:change={({ target }) => newMessage.setField({ body: target.value })}
    />
  </div>

  <div class="px-3 mb-3 flex flex-row justify-end items-center">
    <button
      class="h-10 p-2 mr-3 rounded hover:border-gray-500 bg-white text-gray-700 active:bg-gray-300 focus:outline-none border"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="h-full fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"
        /></svg
      >
    </button>

    <button
      class="bg-red-500 hover:bg-red-600 active:bg-red-700 focus:outline-none px-3 text-white h-10 rounded border-0 font-semibold"
      disabled={submiting}
      on:click={() => submitMessage()}
    >
      Send
    </button>
  </div>
</div>
