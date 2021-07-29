function getTextContent(body, text = []) {
  if (body["content-type"].startsWith("multipart")) {
    body.content.forEach((child) => getTextContent(child, text));
  } else if (body["content-type"] === "text/plain") {
    text.push(body.content);
  }

  return text;
}

export function quotedText(body) {
  return getTextContent(body)
    .map((content) => content.replace(/^/gm, "> "))
    .join("");
}

export function parseMessageBody(
  body,
  content = { text: [], attachments: [], html: [] }
) {
  if (body["content-type"].startsWith("multipart")) {
    body.content.forEach((child) => parseMessageBody(child, content));
  } else if (body["content-type"] === "text/plain") {
    const blocks = parseTextContent(body.content);
    content.text.push(...blocks);
  } else if (body["content-type"] === "text/html") {
    content.html.push(body);
  } else if (body["content-disposition"] === "attachment") {
    content.attachments.push(body);
  }

  return content;
}

function parseTextContent(text) {
  const blocks = [];
  let curBlock = null;

  text.split("\n").forEach((line) => {
    const lType = lineType(line);
    if (curBlock?.type !== lType) {
      if (curBlock) blocks.push(curBlock);
      curBlock = { type: lType, content: "" };
    }

    if (lType === "quote" || lType === "deepquote") {
      curBlock.content += `${line.replace(/^>\s?/, "")}\n`;
    } else {
      curBlock.content += `${line}\n`;
    }
  });
  blocks.push(curBlock);

  return blocks;
}

function lineType(line) {
  if (line.startsWith("> >")) return "deepquote";
  if (line.startsWith(">")) return "quote";

  return "text";
}
