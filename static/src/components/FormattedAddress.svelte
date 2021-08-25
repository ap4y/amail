<script>
  let className = "";
  export { className as class };

  export let address;
  export let removed = false;

  $: [, name, mail] = address.match(/(.*?)\s(\<.*?\>)/) || [null, address];
</script>

<span class={`${className} ${removed ? "line-through" : ""}`}>
  {name}

  {#if mail}
    <span class="hidden sm:inline text-gray-600 pl-1">{mail}</span>
  {/if}
</span>
