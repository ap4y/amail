<script>
  import Message from "./Message.svelte";
  import CollapsedMessage from "./CollapsedMessage.svelte";
  import selectedMessage from "../stores/message";

  export let thread;
  export let level = 0;

  function hiddenMessage({ tags }) {
    return tags.includes("draft") || tags.includes("trash");
  }
</script>

<div>
  {#if thread}
    {#each thread as [message, subthread]}
      {#if !hiddenMessage(message)}
        <div data-message={message.id}>
          {#if message.id !== $selectedMessage}
            <CollapsedMessage {message} {level} />
          {:else}
            <Message {message} />
          {/if}
        </div>
      {/if}

      <svelte:self thread={subthread} level={level + 1} />
    {/each}
  {/if}
</div>
