<script>
  import { createEventDispatcher } from "svelte";

  import ApiClient from "../client";

  export let id;
  export let value;

  const dispatch = createEventDispatcher();
  let input;
  let completions;
  let focused = false;
  let selectedOptionIdx = -1;
  let completeOptions = [];

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

  function selectOption(idx) {
    if (idx !== -1) {
      dispatch("change", [...value, completeOptions[idx]]);
    }
    input.value = null;
    selectedOptionIdx = -1;
    completeOptions = [];
  }

  function completeOption() {
    if (selectedOptionIdx === -1) {
      dispatch("change", [...value, input.value]);
      input.value = null;
      selectedOptionIdx = -1;
      completeOptions = [];
    } else {
      selectOption(selectedOptionIdx);
    }
  }

  function onKeyUp(e) {
    const { key, target } = e;

    if (key === "ArrowDown") {
      selectedOptionIdx += 1;
      if (selectedOptionIdx >= completeOptions.length) {
        selectedOptionIdx = 0;
      }
      completions.children[selectedOptionIdx].scrollIntoView({
        block: "nearest",
      });
    } else if (key === "ArrowUp") {
      selectedOptionIdx -= 1;
      if (selectedOptionIdx < 0) {
        selectedOptionIdx = completeOptions.length - 1;
      }
      completions.children[selectedOptionIdx].scrollIntoView({
        block: "nearest",
      });
    } else if (key === "Enter") {
      completeOption();
    } else if (key === "Backspace") {
      if (input.value.length > 0) return;
      dispatch("change", value.slice(0, value.length - 1));
    }
  }

  function onKeyDown({ key }) {
    if (key === "Tab") {
      completeOption();
    }
  }

  let timer;
  const fetchCompletions = ({ target }) => {
    clearTimeout(timer);
    timer = setTimeout(async () => {
      if (target.value.length <= 2) return;
      const addresses = await ApiClient.default.addresses(target.value);
      completeOptions = addresses.map((addr) => addr["name-addr"]);
    }, 500);
  };
</script>

<div
  class="relative flex-1 flex flex-row flex-wrap px-3 p-1 bg-white text-gray-700 rounded border {classes}"
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
    on:blur={() => (blur = true)}
    on:keyup={onKeyUp}
    on:keydown={onKeyDown}
    on:input={fetchCompletions}
  />

  {#if completeOptions.length > 0}
    <ul
      tabindex="-1"
      class="absolute left-0 right-0 bottom-0 max-h-40 overflow-y-auto block flex-col bg-white border-2 border-gray-600 z-10 transform translate-y-full overflow-x-hidden rounded"
      size="5"
      bind:this={completions}
      on:mouseout={() => (selectedOptionIdx = -1)}
    >
      {#each completeOptions as option, idx}
        <li
          class={`block truncate px-4 py-2 cursor-pointer ${
            idx === selectedOptionIdx ? "bg-red-500 text-white" : ""
          }`}
          on:click={() => selectOption(idx)}
          on:mouseover={() => (selectedOptionIdx = idx)}
        >
          {option}
        </li>
      {/each}
    </ul>
  {/if}
</div>
