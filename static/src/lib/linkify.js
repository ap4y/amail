export function linkify(text, matchFn) {
  const rx = /http[s]?:\/\/\S*\b[=/]?/g;

  const items = [];
  let start = 0;
  text.replaceAll(rx, (match, offset) => {
    items.push(text.substring(start, offset));
    try {
      items.push(new URL(match));
    } catch (e) {
      items.push(match);
      console.error(`failed to linkify: ${e.message}`);
    }
    start = offset + match.length;
  });

  items.push(text.substring(start, text.length));
  return items;
}
