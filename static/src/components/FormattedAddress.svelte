<script>
  export let address;

  $: [, name, mail] = address.match(/(.*?)\s(\<.*?\>)/) || [null, address];
</script>

<span>{name}</span>

{#if mail}
  <span class="text-gray-600 pl-1">{mail}</span>
{/if}
