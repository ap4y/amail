<script>
  import { mailboxTitles } from "../config";
  import url, { selectedMailbox } from "../stores/url";
  import mailboxes from "../stores/mailboxes";

  import InboxIcon from "./InboxIcon.svelte";
  import MailboxButton from "./MailboxButton.svelte";

  export let collapsed;
</script>

<ul>
  {#await mailboxes.fetch()}
    {#if collapsed}
      <p class="text-gray-300 text-xl text-center">...</p>
    {:else}
      <p class="text-gray-300 text-lg text-center">...loading</p>
    {/if}
  {:then}
    {#each $mailboxes as mailbox, idx}
      <li
        class={`relative p-1 mb-2 flex flex-row ${
          collapsed ? "justify-center" : "justify-between"
        } ${
          $selectedMailbox === mailbox.id
            ? "text-red-500 font-semibold"
            : "text-gray-300"
        }`}
      >
        <a
          href={`/${mailbox.id}`}
          class="flex flex-row items-center hover:text-red-500 visited:text-current"
          on:click|preventDefault={() => url.selectMailbox(mailbox.id)}
        >
          {#if collapsed}
            <MailboxButton id={mailbox.id} />
          {:else}
            <span class="mr-2 w-5">
              <InboxIcon id={mailbox.id} />
            </span>
            {mailboxTitles[mailbox.id] || mailbox.id}
          {/if}
        </a>

        {#if mailbox.unread > 0}
          {#if collapsed}
            <span
              class="bg-red-600 absolute -top-1 -right-1 w-2 h-2 rounded-sm"
            />
          {:else}
            <span
              class={`inline-flex items-center text-sm px-2 rounded ${
                mailbox.id === "inbox"
                  ? "bg-red-600 text-gray-100"
                  : "bg-gray-300 text-gray-600"
              }`}>{mailbox.unread}</span
            >
          {/if}
        {/if}
      </li>
    {/each}
  {:catch}
    {#if collapsed}
      <p class="text-red-500 text-sm text-center">Error</p>
    {:else}
      <p class="text-red-500 text-lg text-center">Error</p>
    {/if}
  {/await}
</ul>
