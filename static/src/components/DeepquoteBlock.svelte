<script>
  export let content;
  export let wrap = 80;

  let folded = true;
</script>

<a
  href="#fold"
  style={`width: ${wrap}ch;`}
  class="block whitespace-pre-line break-words text-gray-500 border-l-8 border-gray-400 pl-3"
  on:click|preventDefault={() => (folded = !folded)}
>
  {#if folded}
    <span class="text-red-500 font-semibold">...</span>
  {:else}
    {content}
  {/if}
</a>
