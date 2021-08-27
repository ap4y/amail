<script>
  import url from "../stores/url";

  import BodyItem from "./BodyItem.svelte";
  import MessageInfo from "./MessageInfo.svelte";
  import MessageToolbar from "./MessageToolbar.svelte";

  export let message;
</script>

<div class="relative bg-gray-100">
  <MessageInfo {message} />

  <MessageToolbar {message} />
</div>

<div class="p-5 border-b-2 border-gray-200">
  {#each message.body as body (body.id)}
    <BodyItem {body} messageId={btoa(message.id)} wrap="100" />
  {/each}
</div>
