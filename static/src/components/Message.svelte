<script>
  import { updateTags } from "../lib/tagging";

  import url, { selectedMailbox, selectedThread } from "../stores/url";
  import mailboxes from "../stores/mailboxes";
  import thread from "../stores/thread";
  import selectedMessage from "../stores/message";

  import BodyItem from "./BodyItem.svelte";
  import ToolbarButton from "./ToolbarButton.svelte";
  import MessageInfo from "./MessageInfo.svelte";

  export let message;

  function findOtherMessage(thread, messageId) {
    if (!thread) return null;

    for (const [message, subThread] of thread) {
      if (message.id != messageId && message.tags.includes($selectedMailbox))
        return message;

      const match = findOtherMessage(subThread, messageId);
      if (match) return match;
    }

    return null;
  }

  function markUnread() {
    updateTags($selectedThread, message.id, ["+unread"]);
  }

  function move(folder) {
    const changes = [];
    $mailboxes.forEach(({ id }) => id !== folder && changes.push(`-${id}`));
    updateTags($selectedThread, message.id, [...changes, `+${folder}`]);

    const other = findOtherMessage($thread, message.id);
    if (other) {
      selectedMessage.selectMessage(other.id);
    } else {
      url.deselectThread();
    }
  }
</script>

<div class="relative bg-gray-100" data-message={message.id}>
  <MessageInfo {message} />

  <button
    class="absolute w-5 top-3 right-3 text-gray-500 order-0 hover:text-gray-800 active:text-gray-500 focus:outline-none"
    on:click={() => url.deselectThread()}
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

  <div class="flex flex-row p-3 border-b border-t">
    <ToolbarButton
      tooltip="Mark as unread"
      tooltipPosition="left"
      class="mr-3"
      on:click={() => markUnread()}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-full fill-current"
        ><path d="M0 0h24v24H0V0z" fill="none" /><path
          d="M18.83 7h-2.6L10.5 4 4 7.4V17c-1.1 0-2-.9-2-2V7.17c0-.53.32-1.09.8-1.34L10.5 2l7.54 3.83c.43.23.73.7.79 1.17zM20 8H7c-1.1 0-2 .9-2 2v9c0 1.1.9 2 2 2h13c1.1 0 2-.9 2-2v-9c0-1.1-.9-2-2-2zm0 3.67L13.5 15 7 11.67V10l6.5 3.33L20 10v1.67z"
        /></svg
      >
    </ToolbarButton>

    <ToolbarButton
      tooltip="Move to archive"
      class="mr-1"
      on:click={() => move("archive")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="fill-current w-full"
        ><path d="M0 0h24v24H0V0z" fill="none" /><path
          d="M19 3H4.99c-1.11 0-1.98.89-1.98 2L3 19c0 1.1.88 2 1.99 2H19c1.1 0 2-.9 2-2V5c0-1.11-.9-2-2-2zm0 12h-4c0 1.66-1.35 3-3 3s-3-1.34-3-3H4.99V5H19v10z"
        /></svg
      >
    </ToolbarButton>
    <ToolbarButton
      tooltip="Move to inbox"
      class="mr-1"
      on:click={() => move("inbox")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="fill-current w-full"
        ><path d="M0 0h24v24H0V0z" fill="none" /><path
          d="M19 3H4.99c-1.11 0-1.98.89-1.98 2L3 19c0 1.1.88 2 1.99 2H19c1.1 0 2-.9 2-2V5c0-1.11-.9-2-2-2zm0 12h-4c0 1.66-1.35 3-3 3s-3-1.34-3-3H4.99V5H19v10z"
        /></svg
      >
    </ToolbarButton>
    <ToolbarButton
      tooltip="Move to spam"
      class="mr-1"
      on:click={() => move("spam")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-full fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8 0-1.85.63-3.55 1.69-4.9L16.9 18.31C15.55 19.37 13.85 20 12 20zm6.31-3.1L7.1 5.69C8.45 4.63 10.15 4 12 4c4.42 0 8 3.58 8 8 0 1.85-.63 3.55-1.69 4.9z"
        /></svg
      >
    </ToolbarButton>
    <ToolbarButton
      tooltip="Move to trash"
      class="mr-3"
      on:click={() => move("trash")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-full fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"
        /></svg
      >
    </ToolbarButton>

    <ToolbarButton tooltip="Tag" class="mr-3">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="w-full fill-current"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M21.41 11.58l-9-9C12.05 2.22 11.55 2 11 2H4c-1.1 0-2 .9-2 2v7c0 .55.22 1.05.59 1.42l9 9c.36.36.86.58 1.41.58.55 0 1.05-.22 1.41-.59l7-7c.37-.36.59-.86.59-1.41 0-.55-.23-1.06-.59-1.42zM5.5 7C4.67 7 4 6.33 4 5.5S4.67 4 5.5 4 7 4.67 7 5.5 6.33 7 5.5 7z"
        /></svg
      >
    </ToolbarButton>

    <div class="flex flex-row ml-auto">
      <ToolbarButton tooltip="Reply" class="mr-1">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="w-full fill-current"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M10 9V5l-7 7 7 7v-4.1c5 0 8.5 1.6 11 5.1-1-5-4-10-11-11z"
          /></svg
        >
      </ToolbarButton>
      <ToolbarButton tooltip="Reply all" class="mr-1">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="w-full fill-current"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M7 8V5l-7 7 7 7v-3l-4-4 4-4zm6 1V5l-7 7 7 7v-4.1c5 0 8.5 1.6 11 5.1-1-5-4-10-11-11z"
          /></svg
        >
      </ToolbarButton>
      <ToolbarButton tooltip="Forward" tooltipPosition="right">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          class="w-full fill-current"
          ><path d="M0 0h24v24H0z" fill="none" /><path
            d="M12 8V4l8 8-8 8v-4H4V8z"
          /></svg
        >
      </ToolbarButton>
    </div>
  </div>
</div>

<div class="p-5 border-b-2 border-gray-200">
  {#each message.body as body (body.id)}
    <BodyItem {body} messageId={btoa(message.id)} wrap="100" />
  {/each}
</div>
