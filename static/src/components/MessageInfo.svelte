<script>
  import FormattedAddress from "./FormattedAddress.svelte";
  import TagBadge from "./TagBadge.svelte";

  export let message;
</script>

<div class="p-3">
  <p class="mb-2">
    <span class="inline-flex w-10 text-gray-500">From: </span>
    <FormattedAddress address={message.headers.From} />
  </p>
  <p class="mb-2 text-gray-500">
    <span class="inline-flex w-10">To: </span>
    {message.headers.To}
  </p>
  {#if message.headers.Cc}
    <p class="mb-2 text-gray-500">
      <span class="inline-flex w-10">Cc: </span>
      {message.headers.Cc}
    </p>
  {/if}

  <h3 class="mb-3">
    {message.headers.Subject}
  </h3>

  <div class="flex">
    <div class="inline-flex items-center mr-4 text-gray-600 text-sm">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-4 fill-current mr-1"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z"
        /><path d="M12.5 7H11v6l5.25 3.15.75-1.23-4.5-2.67z" /></svg
      >
      {message.headers.Date}
    </div>

    {#each message.tags as tag}
      <TagBadge class="mr-2" {tag} />
    {/each}
  </div>
</div>
