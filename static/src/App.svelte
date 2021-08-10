<script>
  import { onMount, tick } from "svelte";

  import Tailwind from "./Tailwind.svelte";
  import RefreshButton from "./components/RefreshButton.svelte";
  import ComposeButton from "./components/ComposeButton.svelte";
  import MenuButton from "./components/MenuButton.svelte";
  import Mailboxes from "./components/Mailboxes.svelte";
  import SearchField from "./components/SearchField.svelte";
  import ThreadToolbar from "./components/ThreadToolbar.svelte";
  import ThreadPages from "./components/ThreadPages.svelte";
  import Threads from "./components/Threads.svelte";
  import Thread from "./components/Thread.svelte";
  import MessageForm from "./components/MessageForm.svelte";

  import { mailboxTitles, refreshInterval } from "./config";
  import { markAsRead } from "./lib/tagging";
  import favicon from "./lib/favicon";

  import url, {
    selectedMailbox,
    selectedThread,
    searchTerms,
  } from "./stores/url";
  import mailboxes, { address } from "./stores/mailboxes";
  import thread, {
    findMessage,
    findOtherMessage,
    findLastMessage,
  } from "./stores/thread";
  import selectedMessage from "./stores/message";
  import newMessage from "./stores/new_message";
  import selectedThreads from "./stores/selected_threads";

  let refreshing = false;
  let threadList, messageList;
  let sidebarCollapsed = true;

  onMount(() => {
    setInterval(() => refreshMailboxes(), refreshInterval);
    refreshMailboxes();
  });

  $: console.log(
    "selected",
    $selectedMailbox,
    $selectedThread,
    $selectedMessage,
    $searchTerms
  );

  $: currentMailbox =
    $searchTerms?.length > 0
      ? { id: "search", terms: $searchTerms, tags: [] }
      : $mailboxes.find(({ id }) => id === $selectedMailbox);
  $: document.title =
    currentMailbox?.unread > 0
      ? `(${currentMailbox.unread}) ${mailboxTitles[$selectedMailbox]}`
      : mailboxTitles[$selectedMailbox] || $searchTerms;
  $: document.head.querySelector('link[rel="icon"]').href = favicon(
    currentMailbox?.unread > 0
  );

  $: if ($selectedThread) {
    loadThread();
    scrollToThread();
  }
  $: if ($selectedMessage)
    scrollToMessage(findMessage($thread, $selectedMessage));

  function refreshMailboxes() {
    refreshing = true;
    const start = Date.now();
    mailboxes.updateUnreadCounters();
    setTimeout(() => (refreshing = false), 1000 - Date.now() + start);
  }

  async function loadThread() {
    const res = await thread.fetch($selectedThread);
    await tick();

    const messages = messageList.querySelectorAll(`div[data-message]`);
    if (messages.length === 0) return;

    const messageId = messages[messages.length - 1].dataset.message;
    const message =
      findOtherMessage($thread, null, ["unread"]) ||
      findLastMessage($thread, null, [], ["trash"]) ||
      findMessage(res, messageId);

    selectedMessage.selectMessage(message.id);
    markAsRead($selectedThread, message);
  }

  async function scrollToThread() {
    await tick();
    let interval;

    const check = () => {
      const row = threadList.querySelector(
        `a[data-thread="${$selectedThread}"]`
      );
      if (!row) return;

      row.scrollIntoView({ behavior: "smooth", block: "center" });
      clearInterval(interval);
    };

    interval = setInterval(check, 100);
    check();
  }

  async function scrollToMessage(message) {
    if (!message) return;

    await tick();
    const item = messageList.querySelector(`div[data-message="${message.id}"]`);
    item.scrollIntoView({ behavior: "smooth", block: "start" });
  }
</script>

<Tailwind />

<div class="h-screen w-screen flex">
  <aside
    style={`width: ${sidebarCollapsed ? 4 : 14}rem`}
    class="flex-shrink-0 flex flex-col bg-gray-600"
  >
    <div
      class={`h-14 flex flex-row ${
        sidebarCollapsed ? "justify-center" : "justify-between"
      } items-center p-4`}
    >
      {#if !sidebarCollapsed}
        <h2 class="text-gray-100 font-semibold">{$address}</h2>
      {/if}

      <RefreshButton {refreshing} on:click={refreshMailboxes} />
    </div>

    <div class="px-4">
      <ComposeButton
        collapsed={sidebarCollapsed}
        on:click={() => newMessage.create()}
      />
    </div>

    <div class="p-4 flex-1">
      <Mailboxes collapsed={sidebarCollapsed} />
    </div>

    <div class="p-4">
      <MenuButton
        collapsed={sidebarCollapsed}
        on:click={() => (sidebarCollapsed = !sidebarCollapsed)}
      />
    </div>
  </aside>

  <main
    class="flex flex-col"
    style={`width: calc(100% - ${sidebarCollapsed ? 4 : 14}rem)`}
  >
    <header
      class="h-14 flex-shrink-0 flex flex-row items-center justify-center py-3 pr-3 border-b border-gray-500 bg-gray-600"
    >
      <SearchField />
      <ThreadToolbar disabled={$selectedThreads.length === 0} />
      <ThreadPages />
    </header>

    <section
      class={`${
        $selectedThread ? "h-60 border-b-8" : "flex-1"
      } flex-shrink-0 border-gray-300 overflow-y-auto`}
      bind:this={threadList}
    >
      <Threads mailbox={currentMailbox} />
    </section>

    {#if $selectedThread}
      <section class="flex-1 w-full overflow-y-auto" bind:this={messageList}>
        <Thread thread={$thread} />
      </section>
    {/if}
  </main>

  {#if $newMessage}
    <MessageForm on:close={() => newMessage.destroy()} />
  {/if}
</div>

<style global>
  body {
    padding: 0;
  }
</style>
