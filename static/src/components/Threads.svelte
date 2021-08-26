<script>
  import url, { selectedMailbox, selectedThread } from "../stores/url";
  import threads from "../stores/threads";
  import selectedThreads from "../stores/selected_threads";

  import Checkbox from "./Checkbox.svelte";
  import TagBadge from "./TagBadge.svelte";

  function unreadClasses({ tags }) {
    return tags.includes("unread")
      ? "bg-gray-50 text-red-500 visited:text-red-500"
      : "bg-gray-100 text-gray-600 visited:text-gray-600";
  }

  function commonTag(tag) {
    return [
      "inbox",
      "sent",
      "trash",
      "archive",
      "personal",
      "draft",
      "unread",
    ].includes(tag);
  }
</script>

{#each $threads as thread (thread.thread)}
  <a
    href={`/${$selectedMailbox}/${thread.thread}`}
    data-thread={thread.thread}
    class={`py-2 flex flex-wrap sm:flex-nowrap sm:flex-row items-center border-b hover:bg-gray-200 ${unreadClasses(
      thread
    )} ${$selectedThread === thread.thread ? "bg-red-100 font-semibold" : ""}`}
    on:click|preventDefault={() =>
      url.selectThread($selectedMailbox, thread.thread)}
  >
    <Checkbox
      class="ml-3"
      checked={$selectedThreads.includes(thread.thread)}
      on:click={() => selectedThreads.toggle(thread)}
    />
    <span class="px-3 w-22 sm:w-28 truncate text-xs sm:text-base">{thread.date_relative}</span>
    <span class="px-3 flex-1 sm:flex-none sm:w-40 truncate text-xs sm:text-base">{thread.authors}</span>
    <span
      class={`sm:truncate w-full sm:flex-1 px-3 pt-2 sm:pt-0 ${
        thread.tags.includes("unread") ? "text-red-500" : "text-gray-800"
      }`}
    >
      {thread.subject}
    </span>
    <div class="hidden sm:inline-flex px-3">
      {#each thread.tags as tag}
        {#if !commonTag(tag)}
          <TagBadge class="mr-2" {tag} />
        {/if}
      {/each}
    </div>
  </a>
{/each}
