<script>
  import url from "../stores/url";

  import BodyItem from "./BodyItem.svelte";
  import MessageInfo from "./MessageInfo.svelte";
  import MessageToolbar from "./MessageToolbar.svelte";

  export let message;
</script>

<div class="relative bg-gray-100" data-message={message.id}>
  <MessageInfo {message} />

  <button
    class="absolute w-5 top-3 right-3 text-gray-500 order-0 hover:text-gray-800 active:text-gray-500 focus:outline-none"
    on:click={() => url.deselectThread()}
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

  <div class="flex flex-row p-3 border-b border-t">
    <MessageToolbar {message} />
  </div>
</div>

<div class="p-5 border-b-2 border-gray-200">
  {#each message.body as body (body.id)}
    <BodyItem {body} messageId={btoa(message.id)} wrap="100" />
  {/each}
</div>
