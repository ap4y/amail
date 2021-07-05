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

  import ApiClient from "./client";
  import { mailboxIds, mailboxTitles } from "./config";
  import url, {
    pushState,
    selectedMailbox,
    selectedThread,
  } from "./stores/url";
  import mailboxes, { currentMailbox } from "./stores/mailboxes";
  import threads from "./stores/threads";
  import thread, { selectedMessage } from "./stores/thread";

  const client = ApiClient.default;

  let address = "";
  let refreshing = false;
  let threadList, messageList;
  let sidebarCollapsed = true;

  onMount(() => {
    refreshMailboxes();
  });

  $: console.log("params", $url.mailbox, $url.thread, $url.message);
  $: console.log(
    "selected",
    $selectedMailbox,
    $selectedThread,
    $selectedMessage
  );

  $: if ($currentMailbox.folder) {
    const title =
      $currentMailbox.unread > 0
        ? `(${$currentMailbox.unread}) ${mailboxTitles[$selectedMailbox]}`
        : mailboxTitles[$selectedMailbox];
    document.title = title;
  }

  $: if ($selectedThread) scrollToThread();
  $: if ($selectedMessage) scrollToMessage();

  async function refreshMailboxes() {
    refreshing = true;
    const start = Date.now();

    try {
      const res = await ApiClient.default.mailboxes();
      mailboxes.set(res.mailboxes);
      address = res.address;
    } catch (e) {
      console.log(`Failed to fetch mailboxes: ${e.message}`);
    }

    setTimeout(() => (refreshing = false), 1000 - Date.now() + start);
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

  async function scrollToMessage() {
    await tick();
    const item = messageList.querySelector(
      `div[data-message="${$selectedMessage}"]`
    );
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
        <h2 class="text-gray-100 font-semibold">{address}</h2>
      {/if}

      <RefreshButton {refreshing} on:click={refreshMailboxes} />
    </div>

    <div class="px-4">
      <ComposeButton collapsed={sidebarCollapsed} />
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
      class="h-14 flex-shrink-0 flex flex-row items-center py-3 pr-3 border-b border-gray-500 bg-gray-600"
    >
      <SearchField />
      <ThreadPages />
    </header>

    <section
      class={`${
        $selectedThread ? "h-60 border-b-8" : "flex-1"
      } flex-shrink-0 border-gray-300 overflow-y-auto`}
      bind:this={threadList}
    >
      <Threads />
    </section>

    {#if $thread}
      <section class="flex-1 w-full overflow-y-auto" bind:this={messageList}>
        <Thread thread={$thread} />
      </section>
    {/if}
  </main>
</div>

<style global>
  body {
    padding: 0;
  }
</style>
