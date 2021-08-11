<script>
  import { createEventDispatcher } from "svelte";

  import TagBadge from "./TagBadge.svelte";

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
      <TagBadge
        class="mb-2 mr-2"
        viewOnly={false}
        {tag}
        on:click={() => dispatch("remove", tag)}
      />
    {/each}
  </div>
</div>
