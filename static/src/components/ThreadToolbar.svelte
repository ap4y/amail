<script>
  import { updateThreadTags, tagChanges } from "../lib/tagging";

  import mailboxes from "../stores/mailboxes";
  import { selectedMailbox } from "../stores/url";
  import selectedThreads from "../stores/selected_threads";

  import ToolbarButton from "./ToolbarButton.svelte";

  export let disabled = false;

  async function move(folder) {
    for (const thread of $selectedThreads) {
      const { changes } = tagChanges($mailboxes, $selectedMailbox, folder);
      await updateThreadTags(thread, [...changes, "-unread"]);
      selectedThreads.toggle({ thread });
    }
  }
</script>

<div class="h-full mx-4 flex flex-row">
  <ToolbarButton
    variant="toolbar"
    tooltip="Move to archive"
    class="mr-1"
    {disabled}
    on:click={() => move("archive")}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      enable-background="new 0 0 24 24"
      viewBox="0 0 24 24"
      class="fill-current h-full"
      ><g><rect fill="none" height="24" width="24" /></g><g
        ><path
          d="M20,2H4C3,2,2,2.9,2,4v3.01C2,7.73,2.43,8.35,3,8.7V20c0,1.1,1.1,2,2,2h14c0.9,0,2-0.9,2-2V8.7c0.57-0.35,1-0.97,1-1.69V4 C22,2.9,21,2,20,2z M15,14H9v-2h6V14z M20,7H4V4h16V7z"
        /></g
      ></svg
    >
  </ToolbarButton>

  <ToolbarButton
    variant="toolbar"
    tooltip="Move to spam"
    class="mr-1"
    {disabled}
    on:click={() => move("spam")}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="h-full fill-current"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8 0-1.85.63-3.55 1.69-4.9L16.9 18.31C15.55 19.37 13.85 20 12 20zm6.31-3.1L7.1 5.69C8.45 4.63 10.15 4 12 4c4.42 0 8 3.58 8 8 0 1.85-.63 3.55-1.69 4.9z"
      /></svg
    >
  </ToolbarButton>
  <ToolbarButton
    variant="toolbar"
    tooltip="Move to trash"
    class="mr-1"
    {disabled}
    on:click={() => move("trash")}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="h-full fill-current"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"
      /></svg
    >
  </ToolbarButton>

  <ToolbarButton variant="toolbar" tooltip="Tag" {disabled}>
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      class="h-full fill-current"
      ><path d="M0 0h24v24H0z" fill="none" /><path
        d="M21.41 11.58l-9-9C12.05 2.22 11.55 2 11 2H4c-1.1 0-2 .9-2 2v7c0 .55.22 1.05.59 1.42l9 9c.36.36.86.58 1.41.58.55 0 1.05-.22 1.41-.59l7-7c.37-.36.59-.86.59-1.41 0-.55-.23-1.06-.59-1.42zM5.5 7C4.67 7 4 6.33 4 5.5S4.67 4 5.5 4 7 4.67 7 5.5 6.33 7 5.5 7z"
      /></svg
    >
  </ToolbarButton>
</div>
