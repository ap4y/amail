<script>
  import Message from "./Message.svelte";
  import CollapsedMessage from "./CollapsedMessage.svelte";
  import selectedMessage from "../stores/message";

  export let thread;
  export let level = 0;
</script>

<div>
  {#if thread}
    {#each thread as [message, subthread]}
      <div data-message={message.id}>
        {#if message.id !== $selectedMessage}
          <CollapsedMessage {message} {level} />
        {:else}
          <Message {message} />
        {/if}
      </div>

      <svelte:self thread={subthread} level={level + 1} />
    {/each}
  {/if}
</div>
