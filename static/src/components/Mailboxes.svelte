<script>
  import { createEventDispatcher } from "svelte";
  import InboxIcon from "./InboxIcon.svelte";
  import { mailboxIds, mailboxTitles } from "../config";
  import { pushState, selectedMailbox } from "../stores/url";
  import mailboxes from "../stores/mailboxes";

  const dispatch = createEventDispatcher();

  export let collapsed;

  function selectMailbox(mailbox) {
    pushState({ mailbox }, `/${mailbox}`);
  }
</script>

<ul>
  {#each mailboxIds as id, idx}
    <li
      class={`relative p-1 mb-1 flex flex-row ${
        collapsed ? "justify-center" : "justify-between"
      } ${
        $selectedMailbox === id ? "text-red-500 font-semibold" : "text-gray-300"
      }`}
    >
      <a
        href={`/${id}`}
        class="flex flex-row items-center hover:text-red-500 visited:text-current"
        on:click|preventDefault={() => selectMailbox(id)}
      >
        {#if collapsed}
          <span class="w-5">
            <InboxIcon {id} />
          </span>
        {:else}
          <span class="mr-2 w-5">
            <InboxIcon {id} />
          </span>
          {mailboxTitles[id]}
        {/if}
      </a>

      {#if $mailboxes[id].unread > 0}
        {#if collapsed}
          <span
            class="bg-red-600 absolute -top-1 -right-1 w-2 h-2 rounded-sm"
          />
        {:else}
          <span
            class={`inline-flex items-center text-sm px-2 rounded ${
              id === "inbox"
                ? "bg-red-600 text-gray-100"
                : "bg-gray-300 text-gray-600"
            }`}>{$mailboxes[id].unread}</span
          >
        {/if}
      {/if}
    </li>
  {/each}
</ul>
