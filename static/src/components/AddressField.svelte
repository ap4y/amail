<script>
  import { createEventDispatcher } from "svelte";

  export let id;
  export let value;

  const dispatch = createEventDispatcher();
  let input;
  let focused = false;

  $: classes = focused
    ? "border-red-300"
    : "border-gray-400 hover:border-gray-500";

  export function focus() {
    input.focus();
  }

  function formattedAddress(address) {
    const matches = /(.*)\s\<(.*)\>/.exec(address);
    if (!matches) return address;

    return matches[1];
  }

  function onKeyDown({ keyCode }) {
    if (keyCode !== 8) return;
    if (input.value.length > 0) return;

    dispatch("change", value.slice(0, value.length - 1));
  }

  function onChange({ target }) {
    dispatch("change", [...value, target.value]);
    target.value = null;
  }
</script>

<div
  class="flex-1 flex flex-row flex-wrap px-3 p-1 bg-white text-gray-700 rounded border {classes}"
>
  {#each value as address}
    <span
      class="mr-2 my-1 px-2 py-1 text-xs rounded bg-gray-600 text-white truncate max-w-md"
      >{formattedAddress(address)}</span
    >
  {/each}

  <input
    {id}
    type="text"
    class="flex-1 outline-none"
    bind:this={input}
    on:focus={() => (focused = true)}
    on:blur={() => (focused = false)}
    on:keydown={onKeyDown}
    on:change={onChange}
  />
</div>
