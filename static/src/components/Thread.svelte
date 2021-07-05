<script>
  import Message from "./Message.svelte";
  import CollapsedMessage from "./CollapsedMessage.svelte";
  import { selectedMessage } from "../stores/thread";

  export let thread;
  export let level = 0;

  $: [message, subthreads] = thread;
</script>

<div>
  <div data-message={message.id}>
    {#if message.id !== $selectedMessage}
      <CollapsedMessage {message} {level} />
    {:else}
      <Message {message} />
    {/if}
  </div>

  {#each subthreads as thread}
    <svelte:self {thread} level={level + 1} />
  {/each}
</div>
