<script>
  import { traverseContent } from "../lib/email";
  import { linkify } from "../lib/linkify";

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
        {block.lines.map((line) => line.substring(1)).join("\n")}
      </p>
    {:else}
      <p style={`width: ${wrap}ch;`} class="whitespace-pre-line break-words">
        {#each linkify(block.lines.join("\n")) as text}
          {#if text instanceof URL}
            <a
              href={text.toString()}
              class="inline text-gray-600 hover:text-red-600"
              target="_blank"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                class="inline fill-current h-5"
                ><path d="M0 0h24v24H0z" fill="none" /><path
                  d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"
                /></svg
              >
              <span class="text-sm">{text.hostname}</span>
            </a>
          {:else}
            {text}
          {/if}

          {#if block.lines[block.lines.length - 1] === ""}
            {"\n"}
          {/if}
        {/each}
      </p>
    {/if}
  {/each}
{:else}
  <p class="text-lg">Email has no plaintext content.</p>
{/if}

{#each content.html as block}
  <a
    href={`/api/messages/${messageId}/parts/${block.id}`}
    class="block mt-1 text-red-400 hover:text-red-500 font-semibold underline"
    target="_blank">Open HTML</a
  >
{/each}

<div class="flex flex-row flex-wrap mt-3">
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
