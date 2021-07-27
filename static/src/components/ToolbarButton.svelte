<script>
  export let tooltip;
  export let tooltipPosition = "center";
  export let variant = "message";
  export let disabled = false;

  let className = "";
  export { className as class };

  let tooltipClass = "-translate-x-1/2";
  $: if (tooltipPosition === "left") {
    tooltipClass = "left-0";
  } else if (tooltipPosition === "right") {
    tooltipClass = "right-0";
  } else {
    tooltipClass = "left-4 -translate-x-1/2";
  }

  let hover = false;
</script>

<div
  class="h-full relative {className}"
  on:mouseover={() => (hover = true)}
  on:mouseout={() => (hover = false)}
>
  {#if variant === "toolbar"}
    <button
      {disabled}
      class="h-full p-2 rounded text-gray-300 focus:outline-none hover:text-gray-100 active:text-gray-100 disabled:opacity-25"
      on:click
    >
      <slot />
    </button>
    <span
      class={`absolute -bottom-full z-50 transform ${tooltipClass} translate-y-2 p-2 whitespace-nowrap rounded bg-gray-800 text-white ${
        hover ? "" : "hidden"
      }`}>{tooltip}</span
    >
  {:else}
    <button
      {disabled}
      class="w-9 p-2 rounded hover:border-gray-500 bg-white text-gray-700 active:bg-gray-300 focus:outline-none border"
      on:click
    >
      <slot />
    </button>
    <span
      class={`absolute -top-full z-50 transform ${tooltipClass} -translate-y-2 p-2 whitespace-nowrap rounded bg-gray-800 text-white ${
        hover ? "" : "hidden"
      }`}>{tooltip}</span
    >
  {/if}
</div>
