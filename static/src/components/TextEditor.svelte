<script>
  import { createEventDispatcher } from "svelte";

  export let blocks = [];

  const dispatch = createEventDispatcher();
  let content = blocks;

  function beforeInput(e) {
    const ranges = e.getTargetRanges();
    if (ranges.length === 0) return;

    const { startContainer, startOffset } = ranges[0];
    const node = startContainer.parentNode;

    if (node.dataset.type !== "quote") return;
    e.preventDefault();

    if (e.inputType !== "insertParagraph") return;

    node.innerHTML = startContainer.wholeText.slice(0, startOffset);
    const otherQuote = `
<p class="whitespace-pre-wrap break-words" data-type="text">\n</p>
<p class="whitespace-pre-wrap break-words text-gray-600 border-l-8 border-gray-400 pl-3" data-type="quote">${startContainer.wholeText.slice(
      startOffset + 1
    )}</p>`;
    node.insertAdjacentHTML("afterend", otherQuote);
    window.getSelection().collapse(node.nextSibling.nextSibling);
  }
</script>

<div
  contenteditable="true"
  style="width: 85ch;"
  class="h-80 border rounded p-3 outline-none border-gray-400 hover:border-gray-500 focus:border-red-300 overflow-y-auto"
  on:input
  on:beforeinput={beforeInput}
>
  {#each content as block}
    {#if block.type === "quote"}
      <p
        class="whitespace-pre-wrap break-words text-gray-600 border-l-8 border-gray-400 pl-3"
        data-type={block.type}
      >
        {block.content.replace(/^>\s?/gm, "")}
      </p>
    {:else}
      <p class="whitespace-pre-wrap break-words" data-type={block.type}>
        {block.content}
      </p>
    {/if}
  {/each}
</div>
