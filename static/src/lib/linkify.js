export function linkify(text, matchFn) {
  const rx = /http[s]?:\/\/\S*\b[=/]?/g;

  const items = [];
  let start = 0;
  text.replaceAll(rx, (match, offset) => {
    items.push(text.substring(start, offset));
    items.push(new URL(match));
    start = offset + match.length;
  });

  items.push(text.substring(start, text.length));
  return items;
}
