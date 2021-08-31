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
  let threadList, messageList, searchField, activeMessage;
  let sidebarCollapsed = true;
  let searching = false;

  onMount(() => {
    setInterval(() => refreshMailboxes(), refreshInterval);
  });

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

  function nextThread() {
    let node;

    if ($selectedThread) {
      const row = threadList.querySelector(
        `a[data-thread="${$selectedThread}"]`
      );
      node = row?.nextSibling;
    } else if (threadList.children.length > 0) {
      node = threadList.children[0];
    }

    if (node?.dataset?.thread) {
      url.selectThread($selectedMailbox, node.dataset.thread);
    }
  }

  function prevThread() {
    if (!$selectedThread) {
      return;
    }

    const row = threadList.querySelector(`a[data-thread="${$selectedThread}"]`);
    const node = row?.previousSibling;
    if (node?.dataset?.thread) {
      url.selectThread($selectedMailbox, node.dataset.thread);
    }
  }

  function nextMessage() {
    if (!$selectedThread) {
      return;
    }

    const messages = [...messageList.querySelectorAll(`div[data-message]`)];
    const idx = messages.findIndex(
      (node) => node.dataset.message === $selectedMessage
    );

    if (idx < 0 || idx >= messages.length) {
      return;
    }

    const node = messages[idx + 1];
    if (node?.dataset?.message) {
      selectedMessage.selectMessage(node.dataset.message);
    }
  }

  function prevMessage() {
    if (!$selectedThread) {
      return;
    }

    const messages = [...messageList.querySelectorAll(`div[data-message]`)];
    const idx = messages.findIndex(
      (node) => node.dataset.message === $selectedMessage
    );

    if (idx <= 0 || idx > messages.length) {
      return;
    }

    const node = messages[idx - 1];
    if (node?.dataset?.message) {
      selectedMessage.selectMessage(node.dataset.message);
    }
  }

  function openHtml() {
    const htmlItems = messageList?.querySelectorAll(`a[data-html-body]`);
    if (!htmlItems || htmlItems.length === 0) {
      return;
    }

    htmlItems[0].click();
  }

  const keys = {
    n: nextThread,
    p: prevThread,
    Escape: () => url.deselectThread(),
    N: () => messageList?.scrollBy(0, 200),
    P: () => messageList?.scrollBy(0, -200),
    e: nextMessage,
    a: prevMessage,
    C: () => newMessage.create(),
    r: () => newMessage.reply($selectedMessage, "sender"),
    R: () => newMessage.reply($selectedMessage, "all"),
    f: () => newMessage.forward(findMessage($thread, $selectedMessage)),
    s: () => searchField.focus(),
    V: openHtml,
    A: () => activeMessage?.move("archive"),
    D: () => activeMessage?.move("trash"),
    J: () => activeMessage?.move("spam"),
  };

  document.onkeydown = (e) => {
    if ($newMessage || searching) {
      return;
    }

    const { key } = e;
    if (keys[key]) {
      console.debug(`found hotkey for ${key}`);
      keys[key]();
      e.preventDefault();
    }
  };
</script>

<Tailwind />

<div class="h-screen w-screen flex">
  <aside
    class={`absolute sm:fixed h-screen flex-shrink-0 flex flex-col bg-gray-600 z-10 ${
      sidebarCollapsed ? "hidden sm:flex w-16" : "w-56"
    }`}
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
    class={`flex flex-col w-full ${sidebarCollapsed ? "sm:pl-16" : "sm:pl-56"}`}
  >
    <header
      class="flex sm:flex-row flex-wrap items-center sm:justify-center py-3 pr-3 pl-3 sm:pl-0 border-b border-gray-500 bg-gray-600"
    >
      <div class="sm:hidden mr-3 pb-1">
        <MenuButton
          collapsed={sidebarCollapsed}
          on:click={() => (sidebarCollapsed = !sidebarCollapsed)}
        />
      </div>

      <SearchField
        bind:this={searchField}
        on:focus={() => (searching = true)}
        on:blur={() => (searching = false)}
      />
      <ThreadToolbar disabled={$selectedThreads.length === 0} />
      <ThreadPages />
    </header>

    <section
      class={`${
        $selectedThread ? "h-0 sm:h-60 border-b-8" : "flex-1"
      } flex-shrink-0 border-gray-300 overflow-y-auto`}
      bind:this={threadList}
    >
      {#if currentMailbox}
        <Threads mailbox={currentMailbox} />
      {/if}
    </section>

    {#if $selectedThread}
      <section
        class="relative flex-1 w-full overflow-y-auto"
        bind:this={messageList}
      >
        <div class="sticky top-3 w-full right-3 z-10">
          <button
            class="absolute right-3 bg-gray-200 p-1 rounded text-gray-500 order-0 hover:text-gray-800 active:text-gray-500 focus:outline-none"
            on:click={() => url.deselectThread()}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              class="w-5 fill-current"
              ><path d="M0 0h24v24H0z" fill="none" /><path
                d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
              /></svg
            >
          </button>
        </div>

        <Thread thread={$thread} bind:activeMessage />
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
