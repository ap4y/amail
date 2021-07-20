<script>
  export let collapsed = false;
</script>

{#if collapsed}
  <button
    class="bg-red-600 hover:bg-red-700 active:bg-red-800 focus:outline-none w-8 text-white h-8 p-1 rounded border-0"
    on:click
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current h-full"
    >
      <path d="M19,13h-6v6h-2v-6H5v-2h6V5h2v6h6V13z" />
      ></svg
    >
  </button>
{:else}
  <button
    class="bg-red-600 hover:bg-red-700 active:bg-red-800 focus:outline-none w-full text-white h-10 rounded border-0 font-semibold"
    on:click
  >
    New message
  </button>
{/if}
