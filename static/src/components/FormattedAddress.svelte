<script>
  export let address;
  export let removed = false;

  $: [, name, mail] = address.match(/(.*?)\s(\<.*?\>)/) || [null, address];
</script>

<span class={`${removed ? "line-through" : ""}`}>{name}</span>

{#if mail}
  <span class={`text-gray-600 pl-1 ${removed ? "line-through" : ""}`}
    >{mail}</span
  >
{/if}
