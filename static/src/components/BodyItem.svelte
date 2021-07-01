<script>
  export let body;
  export let messageId;
  export let wrap = 80;

  $: content = traverseContent(body);
  $: console.log(body, content);
  $: blocks = content.text
    ?.join("\n")
    ?.split("\n")
    .reduce((acc, line) => {
      const currentBlock = (acc[acc.length - 1] || { lines: [] }).lines;
      const lastLine = currentBlock[currentBlock.length - 1] || "";

      const currentType = lineType(line);
      if (currentType === lineType(lastLine)) {
        if (acc.length === 0) {
          acc.push({ type: currentType, lines: [line] });
        } else {
          acc[acc.length - 1].lines.push(line);
        }
      } else {
        acc.push({ type: currentType, lines: [line] });
      }

      return acc;
    }, []);

  function traverseContent(item, acc = { text: [], html: [], attach: [] }) {
    if (item["content-type"].startsWith("multipart")) {
      item.content.forEach((child) => traverseContent(child, acc));
    } else if (item["content-type"] === "text/plain") {
      acc.text.push(item.content);
    } else if (item["content-type"] === "text/html") {
      acc.html.push(item.content);
    } else if (item["content-disposition"] === "attachment") {
      acc.attach.push(item);
    }

    return acc;
  }

  function lineType(line) {
    if (line.startsWith(">")) return "quote";

    return "text";
  }

  function openHTML() {
    content.html.forEach((html) => {
      const win = window.open();
      const doc = win.document;
      doc.open();
      doc.write(html);
      doc.close();
    });
  }
</script>

{#if content.text.length > 0}
  {#each blocks as block}
    {#if block.type === "quote"}
      <p
        style={`width: ${wrap}ch;`}
        class="whitespace-pre-line break-words text-gray-600 border-l-8 border-gray-400 pl-3"
      >
        {block.lines.join("\n")}
      </p>
    {:else}
      <p style={`width: ${wrap}ch;`} class="whitespace-pre-line break-words">
        {block.lines.join("\n")}
      </p>
    {/if}
  {/each}
{:else}
  <p class="text-lg">Email has no plaintext content.</p>
{/if}

{#if content.html.length > 0}
  <button
    class="mt-2 px-4 py-1 bg-red-400 text-white rounded border-none"
    on:click={openHTML}>Open HTML</button
  >
{/if}

<div class="flex flex-row flex-wrap mt-4 pt-3">
  {#each content.attach as attach}
    <a
      href={`/api/messages/${messageId}/parts/${attach.id}`}
      download={attach.filename}
      class="flex flex-row items-center mr-2 mb-2 p-3 rounded border-2 border-gray-500 text-gray-600 font-semibold hover:border-gray-700 hover:text-gray-800 visited:text-gray-600 hover:no-underline"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        class="fill-current mr-1 w-5"
        ><path d="M0 0h24v24H0z" fill="none" /><path
          d="M16.5 6v11.5c0 2.21-1.79 4-4 4s-4-1.79-4-4V5c0-1.38 1.12-2.5 2.5-2.5s2.5 1.12 2.5 2.5v10.5c0 .55-.45 1-1 1s-1-.45-1-1V6H10v9.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V5c0-2.21-1.79-4-4-4S7 2.79 7 5v12.5c0 3.04 2.46 5.5 5.5 5.5s5.5-2.46 5.5-5.5V6h-1.5z"
        /></svg
      >
      {attach.filename}</a
    >
  {/each}
</div>
