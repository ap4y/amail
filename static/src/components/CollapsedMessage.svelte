<script>
  import { pushState, selectedMailbox, selectedThread } from "../stores/url";
  import FormattedAddress from "./FormattedAddress.svelte";

  export let message;
  export let level;

  let drag = false;

  function selectMessage() {
    pushState({
      mailbox: $selectedMailbox,
      thread: $selectedThread,
      message: message.id,
    });
  }
</script>

<div
  class="p-3 flex flex-row justify-between items-center bg-gray-100 border-b-2 border-gray-200 cursor-pointer"
  on:mousedown={() => (drag = false)}
  on:mousemove={() => (drag = true)}
  on:mouseup={() => !drag && selectMessage()}
>
  <div
    class="flex flex-row items-center"
    style={`padding-left: ${10 * Math.max(0, level - 1)}px;`}
  >
    {#if level > 0}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-5 fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"
        /></svg
      >
    {/if}
    <FormattedAddress address={message.headers.From} />
  </div>
  <span class="text-gray-500 text-sm">{message.headers.Date}</span>
</div>
