<script>
  import { createEventDispatcher } from "svelte";

  export let tags = [];

  let className = "";
  export { className as class };

  const dispatch = createEventDispatcher();
  let newTag = "";

  function onKeyPress({ charCode }) {
    if (charCode === 13) {
      dispatch("add", newTag);
      newTag = "";
    }
  }
</script>

<div
  class="{className} w-72 p-4 bg-white border-2 border-gray-600 rounded shadow"
>
  <button
    class="absolute w-5 top-2 right-2 text-gray-500 order-0 hover:text-gray-800 active:text-gray-500 focus:outline-none"
    on:click={() => dispatch("close")}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="w-full fill-current"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
      /></svg
    >
  </button>

  <h3 class="mb-1 font-semibold text-gray-700">Labels</h3>
  <input
    type="text"
    placeholder="Type in a new tag"
    class="w-full outline-none mb-4 px-3 py-1 bg-white text-sm text-gray-700 rounded border border-gray-400 hover:border-gray-500 focus:border-red-300"
    on:keypress={onKeyPress}
    bind:value={newTag}
  />

  <div class="flex flex-row flex-wrap">
    {#each tags as tag}
      <div
        class="inline-flex text-gray-600 text-sm mr-2 max-w-xs truncate cursor-pointer hover:text-red-500"
        on:click={() => dispatch("remove", tag)}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="w-4 fill-current mr-1 flex-shrink-0"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M21.41 11.58l-9-9C12.05 2.22 11.55 2 11 2H4c-1.1 0-2 .9-2 2v7c0 .55.22 1.05.59 1.42l9 9c.36.36.86.58 1.41.58.55 0 1.05-.22 1.41-.59l7-7c.37-.36.59-.86.59-1.41 0-.55-.23-1.06-.59-1.42zM5.5 7C4.67 7 4 6.33 4 5.5S4.67 4 5.5 4 7 4.67 7 5.5 6.33 7 5.5 7z"
          /></svg
        >
        <span>{tag}</span>
      </div>
    {/each}
  </div>
</div>
