<script>
  import { threadsPerPage } from "../config";
  import { currentPage, hasMore } from "../stores/threads";
</script>

<div class="flex flex-row items-center ml-auto text-gray-300">
  <button
    disabled={$currentPage == 0}
    class="border-0 rounded hover:bg-gray-700 active:bg-gray-800 focus:outline-none disabled:opacity-50 disabled:bg-transparent"
    on:click={() => currentPage.update((page) => page - 1)}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current w-7"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"
      /></svg
    >
  </button>
  <span class="font-semibold px-3">{$currentPage + 1}</span>
  <button
    disabled={!$hasMore}
    class="border-0 rounded hover:bg-gray-700 active:bg-gray-800 focus:outline-none disabled:opacity-50 disabled:bg-transparent"
    on:click={() => currentPage.update((page) => page + 1)}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="fill-current w-7"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"
      /></svg
    >
  </button>
</div>
