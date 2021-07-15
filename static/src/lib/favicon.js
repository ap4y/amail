export default function favicon(hasUnread) {
  if (!hasUnread) {
    return "/favicon.png";
  }

  const canvas = document.createElement("canvas");
  canvas.height = 64;
  canvas.width = 64;

  const ctx = canvas.getContext("2d");
  ctx.font = "64px serif";
  ctx.fillText("✉️", 0, 64);

  return canvas.toDataURL();
}
