<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let mailbox;
  export let threads;
  export let selectedThread;

  function unreadClasses({ tags }) {
    return tags.includes("unread")
      ? "bg-gray-50 text-red-500 visited:text-red-500"
      : "bg-gray-100 text-gray-600 visited:text-gray-600";
  }
</script>

{#each threads as thread (thread.thread)}
  <a
    href={`/${mailbox}/${thread.thread}`}
    data-thread={thread.thread}
    class={`h-10 flex flex-row items-center border-b hover:bg-gray-200 ${unreadClasses(
      thread
    )} ${selectedThread === thread.thread ? "bg-red-100 font-semibold" : ""}`}
    on:click|preventDefault={() =>
      dispatch("click", { mailbox, thread: thread.thread })}
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
