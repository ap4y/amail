export function traverseContent(
  item,
  acc = { text: [], html: [], attach: [] }
) {
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

export function quotedText(text) {
  return text
    .split()
    .map((line) => `> ${line}`)
    .join();
}
