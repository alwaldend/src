import { readFileSync, writeFileSync } from "node:fs";

const [monogramFile, outputFile] = process.argv.slice(2);
const monogram = readFileSync(monogramFile, "utf8");
const favicon = monogram.replace("<svg ", '<svg data-variant="favicon" ');
const icon = Buffer.from(favicon).toString("base64");

const sizes = [16, 20, 24, 32, 48, 64];
const samples = sizes
    .map(
        (size) => `<li>
  <div class="sample" style="--icon-size: ${size}px">
    <div class="icon"><img src="data:image/svg+xml;base64,${icon}" width="${size}" height="${size}" alt=""></div>
  </div>
  <span>${size} px</span>
</li>`,
    )
    .join("\n");

const strip = (label, background) => `<section class="${background}">
  <h2>${label}</h2>
  <ul>${samples}</ul>
</section>`;

writeFileSync(
    outputFile,
    `<!doctype html>
<html lang="en">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Alwaldend favicon preview</title>
<style>
  body { font: 14px/1.4 system-ui, sans-serif; }
  main { display: grid; gap: 8px; }
  section { padding: 12px; border-radius: 10px; }
  h2 { margin: 0; font-size: 14px; font-weight: 600; }
  ul {
    display: grid;
    grid-template-columns: repeat(6, minmax(0, 1fr));
    margin: 0;
    padding: 0;
    list-style: none;
  }
  li { display: grid; justify-items: center; gap: 4px; }
  .sample { height: 88px; display: grid; place-items: center; }
  .icon, img { display: block; }
  .light { background: #f0f1f3; color: #212121; }
  .dark { background: #292a2d; color: #fff; }
  .mobile { background: #171717; color: #fff; }
  .mobile .icon { padding: calc(var(--icon-size) * 0.16); background: #fff; border-radius: 50%; }
  @media (max-width: 520px) {
    ul { grid-template-columns: repeat(3, minmax(0, 1fr)); }
  }
</style>
<main aria-label="Alwaldend favicon preview">
  ${strip("Light tab", "light")}
  ${strip("Dark tab", "dark")}
  ${strip("Mobile tab with white backing", "mobile")}
</main>
</html>\n`,
);
