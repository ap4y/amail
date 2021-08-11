<script>
  import InboxIcon from "./InboxIcon.svelte";

  export let tag;
  export let viewOnly = true;

  let className = "";
  export { className as class };

  $: classes = viewOnly ? "" : "cursor-pointer hover:text-red-500";
</script>

<div
  class="inline-flex items-center text-gray-600 text-sm max-w-xs truncate {className} {classes}"
  on:click
>
  <span class="w-4 mr-1 flex-shrink-0">
    <InboxIcon id={tag} iconType="tag" />
  </span>

  <span>{tag}</span>
</div>
