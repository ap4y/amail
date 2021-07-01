<script>
  export let collapsed;
</script>

<button
  class="text-gray-300 h-8 p-1 rounded hover:bg-gray-700 focus:outline-none"
  on:click
>
  {#if collapsed}
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current h-full"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"
      /></svg
    >
  {:else}
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current h-full"
      ><path d="M0 0h24v24H0V0z" fill="none" /><path
        d="M3 18h13v-2H3v2zm0-5h10v-2H3v2zm0-7v2h13V6H3zm18 9.59L17.42 12 21 8.41 19.59 7l-5 5 5 5L21 15.59z"
      /></svg
    >
  {/if}
</button>
