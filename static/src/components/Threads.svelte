<script>
  import url, {
    pushState,
    selectedMailbox,
    selectedThread,
    searchTerms,
  } from "../stores/url";
  import threads from "../stores/threads";

  function unreadClasses({ tags }) {
    return tags.includes("unread")
      ? "bg-gray-50 text-red-500 visited:text-red-500"
      : "bg-gray-100 text-gray-600 visited:text-gray-600";
  }

  function selectThread({ mailbox, thread }) {
    if ($searchTerms) {
      pushState(
        { mailbox, thread, searchTerms: $searchTerms },
        `/${mailbox}/${thread}?terms=${escape($searchTerms)}`
      );
    } else {
      pushState({ mailbox, thread }, `/${mailbox}/${thread}`);
    }
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
      selectThread({ mailbox: $selectedMailbox, thread: thread.thread })}
  >
    <span class="px-3 w-28">{thread.date_relative}</span>
    <span class="pr-6 w-40 truncate">{thread.authors}</span>
    <span
      class={`truncate flex-1 ${
        thread.tags.includes("unread") ? "text-red-500" : "text-gray-800"
      }`}
    >
      {thread.subject}
    </span>
  </a>
{/each}
