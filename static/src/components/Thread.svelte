<script>
  import Message from "./Message.svelte";
  import CollapsedMessage from "./CollapsedMessage.svelte";

  export let thread;
  export let selectedMessage;
  export let level = 0;

  $: [message, subthreads] = thread;
</script>

<div>
  <div data-message={message.id}>
    {#if message.id !== selectedMessage}
      <CollapsedMessage {message} {level} on:click />
    {:else}
      <Message {message} on:click on:close />
    {/if}
  </div>

  {#each subthreads as thread}
    <svelte:self {thread} {selectedMessage} level={level + 1} on:click />
  {/each}
</div>
