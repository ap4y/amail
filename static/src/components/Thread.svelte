<script>
  import Message from "./Message.svelte";
  import CollapsedMessage from "./CollapsedMessage.svelte";
  import selectedMessage from "../stores/message";

  export let thread;
  export let level = 0;
</script>

<div>
  {#each thread as [message, subthread]}
    <div data-message={message.id}>
      {#if message.id !== $selectedMessage?.id}
        <CollapsedMessage {message} {level} />
      {:else}
        <Message message={$selectedMessage} />
      {/if}
    </div>

    <svelte:self thread={subthread} level={level + 1} />
  {/each}
</div>
