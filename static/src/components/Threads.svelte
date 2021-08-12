<script>
  import url, { selectedMailbox, selectedThread } from "../stores/url";
  import threads, { currentPage } from "../stores/threads";
  import selectedThreads from "../stores/selected_threads";

  import Checkbox from "./Checkbox.svelte";
  import TagBadge from "./TagBadge.svelte";

  export let mailbox;

  $: if (mailbox) threads.fetch(mailbox.terms, $currentPage);

  function unreadClasses({ tags }) {
    return tags.includes("unread")
      ? "bg-gray-50 text-red-500 visited:text-red-500"
      : "bg-gray-100 text-gray-600 visited:text-gray-600";
  }

  function commonTag(tag) {
    return ["inbox", "sent", "trash", "archive", "personal", "draft", "unread"].includes(
      tag
    );
  }
</script>

{#each $threads as thread (thread.thread)}
  <a
    href={`/${$selectedMailbox}/${thread.thread}`}
    data-thread={thread.thread}
    class={`h-10 flex flex-row items-center border-b hover:bg-gray-200 ${unreadClasses(
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
    <span class="px-3 w-28">{thread.date_relative}</span>
    <span class="pr-6 w-40 truncate">{thread.authors}</span>
    <span
      class={`truncate flex-1 ${
        thread.tags.includes("unread") ? "text-red-500" : "text-gray-800"
      }`}
    >
      {thread.subject}
    </span>
    <div class="inline-flex px-3">
      {#each thread.tags as tag}
        {#if !commonTag(tag)}
          <TagBadge class="mr-2" {tag} />
        {/if}
      {/each}
    </div>
  </a>
{/each}
