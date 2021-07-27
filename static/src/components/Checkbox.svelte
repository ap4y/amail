<script>
  let className = "";
  export { className as class };
  export let checked = false;

  $: borderClass = checked ? "border-gray-500" : "border-gray-300";
</script>

<div
  class="w-5 h-5 rounded border-2 text-gray-500 {borderClass} {className}"
  on:click|stopPropagation|preventDefault
>
  {#if checked}
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current w-full"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4L9 16.2z"
      /></svg
    >
  {/if}
</div>
