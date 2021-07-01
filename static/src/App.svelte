<script>
  import { onMount, tick } from "svelte";

  import Tailwind from "./Tailwind.svelte";
  import RefreshButton from "./components/RefreshButton.svelte";
  import ComposeButton from "./components/ComposeButton.svelte";
  import MenuButton from "./components/MenuButton.svelte";
  import Mailboxes from "./components/Mailboxes.svelte";
  import SearchField from "./components/SearchField.svelte";
  import ThreadPages from "./components/ThreadPages.svelte";
  import Threads from "./components/Threads.svelte";
  import Thread from "./components/Thread.svelte";

  import url, { pushState } from "./stores/url";
  import ApiClient from "./client";
  import { mailboxIds, mailboxTitles } from "./config";

  const client = new ApiClient("/api");

  let address = "";
  let mailboxes = mailboxIds.reduce((acc, id) => {
    acc[id] = {};
    return acc;
  }, {});
  let refreshing = false;
  const perPage = 50;
  let currentPage = 0;
  let threads = [];
  let threadList, messageList;
  let thread = null;
  let sidebarCollapsed = true;

  onMount(() => {
    refreshMailboxes();
  });

  $: console.log("params", $url.mailbox, $url.thread, $url.message);

  $: selectedMailbox = $url.mailbox?.length > 0 ? $url.mailbox : mailboxIds[0];
  $: selectedThread = $url.thread?.length > 0 ? $url.thread : null;
  $: selectedMessage =
    $url.message?.length > 0 ? $url.message : thread ? thread[0].id : null;
  $: console.log("selected", selectedMailbox, selectedThread, selectedMessage);

  $: mailbox = mailboxes[selectedMailbox];
  $: {
    if (mailbox.folder) {
      const title =
        mailbox.unread > 0
          ? `(${mailbox.unread}) ${mailboxTitles[selectedMailbox]}`
          : mailboxTitles[selectedMailbox];
      document.title = title;

      client
        .threads(mailbox.folder, currentPage, perPage)
        .then((res) => (threads = res));
    }
  }

  $: {
    if (selectedThread) {
      client.thread(selectedThread).then((res) => (thread = res));
    }
  }

  $: {
    if (selectedMessage) scrollToMessage();
  }

  async function refreshMailboxes() {
    refreshing = true;
    const start = Date.now();

    try {
      const res = await client.mailboxes();
      mailboxes = res.mailboxes;
      address = res.address;
    } catch (e) {
      console.log(`Failed to fetch mailboxes: ${e.message}`);
    }

    setTimeout(() => (refreshing = false), 1000 - Date.now() + start);
  }

  function selectMailbox({ detail }) {
    currentPage = 0;
    thread = null;
    pushState({ mailbox: detail }, `/${detail}`);
  }

  async function selectThread({ detail: { mailbox, thread } }) {
    pushState({ mailbox, thread }, `/${mailbox}/${thread}`);

    const row = threadList.querySelector(`a[data-thread="${thread}"]`);
    await tick();
    row.scrollIntoView({ behavior: "smooth", block: "center" });
  }

  function selectMessage({ detail }) {
    pushState({
      mailbox: selectedMailbox,
      thread: selectedThread,
      message: detail,
    });
  }

  async function scrollToMessage() {
    await tick();
    const item = messageList.querySelector(
      `div[data-message="${selectedMessage}"]`
    );
    item.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function closeThread() {
    thread = null;
    pushState(
      { mailbox: selectedMailbox, thread: null, message: null },
      `/${selectedMailbox}`
    );
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
        <h2 class="text-gray-100 font-semibold">{address}</h2>
      {/if}

      <RefreshButton {refreshing} on:click={refreshMailboxes} />
    </div>

    <div class="px-4">
      <ComposeButton collapsed={sidebarCollapsed} />
    </div>

    <div class="p-4 flex-1">
      <Mailboxes
        {mailboxes}
        {selectedMailbox}
        collapsed={sidebarCollapsed}
        on:click={selectMailbox}
      />
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
      class="h-14 flex-shrink-0 flex flex-row items-center py-3 pr-3 border-b border-gray-500 bg-gray-600"
    >
      <SearchField />
      <ThreadPages
        {currentPage}
        lastPage={Math.ceil(mailbox.total / perPage)}
        on:previous-page={() => (currentPage = currentPage - 1)}
        on:next-page={() => (currentPage = currentPage + 1)}
      />
    </header>

    <section
      class={`${
        selectedThread ? "h-60 border-b-8" : "flex-1"
      } flex-shrink-0 border-gray-300 overflow-y-auto`}
      bind:this={threadList}
    >
      <Threads
        mailbox={selectedMailbox}
        {threads}
        {selectedThread}
        on:click={selectThread}
      />
    </section>

    {#if thread}
      <section class="flex-1 w-full overflow-y-auto" bind:this={messageList}>
        <Thread
          {thread}
          {selectedMessage}
          on:click={selectMessage}
          on:close={closeThread}
        />
      </section>
    {/if}
  </main>
</div>

<style global>
  body {
    padding: 0;
  }
</style>
