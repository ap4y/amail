<script>
  import { updateTags, tagChanges } from "../lib/tagging";
  import url, { selectedMailbox, selectedThread } from "../stores/url";
  import mailboxes from "../stores/mailboxes";
  import selectedMessage from "../stores/message";
  import thread, { findOtherMessage } from "../stores/thread";
  import threads from "../stores/threads";

  import BodyItem from "./BodyItem.svelte";
  import MessageInfo from "./MessageInfo.svelte";
  import MessageToolbar from "./MessageToolbar.svelte";

  export let message;

  function selectNextThread() {
    const idx = $threads.findIndex(({ thread }) => thread === $selectedThread);
    const nextThread = $threads[idx + 1];
    if (nextThread) {
      url.selectThread($selectedMailbox, nextThread.thread);
    } else {
      url.deselectThread();
    }
  }

  export function move(folder) {
    const { changes, fromTags } = tagChanges(
      $mailboxes,
      $selectedMailbox,
      folder
    );
    updateTags($selectedThread, message.id, changes);

    const other = findOtherMessage($thread, message.id, fromTags);
    if (other) {
      selectedMessage.selectMessage(other.id);
    } else {
      selectNextThread();
    }
  }
</script>

<div class="relative bg-gray-100">
  <MessageInfo {message} />

  <MessageToolbar {message} on:move={({ detail }) => move(detail)} />
</div>

<div class="p-5 border-b-2 border-gray-200">
  {#each message.body as body (body.id)}
    <BodyItem {body} messageId={window.btoa(message.id)} wrap="100" />
  {/each}
</div>
