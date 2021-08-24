<script>
  import { mailboxTitles } from "../config";

  import InboxIcon from "./InboxIcon.svelte";

  export let id;

  let hover = false;
</script>

<span
  class={`relative ${mailboxTitles[id] ? "w-5" : "w-7"}`}
  on:mouseover={() => (hover = true)}
  on:mouseout={() => (hover = false)}
>
  {#if mailboxTitles[id]}
    <InboxIcon {id} />
  {:else}
    <InboxIcon {id} />
    <span
      class="absolute text-sm left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-2/3"
      >{id[0]}</span
    >
  {/if}

  <span
    class={`absolute left-8 z-50 transform bottom-0 translate-y-2 p-2 whitespace-nowrap rounded bg-gray-800 text-white ${
      hover ? "" : "hidden"
    }`}>{mailboxTitles[id] || id}</span
  >
</span>
