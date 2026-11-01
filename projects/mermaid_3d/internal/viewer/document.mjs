export function renderDocument({
    title,
    fontCSS,
    styles,
    importMap,
    payload,
}) {
    const escapedTitle = title
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;");
    const safeJSON = (value) =>
        JSON.stringify(value).replaceAll("<", "\\u003c");
    return `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>${escapedTitle} — 3D diagram</title><style>${fontCSS}\n${styles}</style><script type="importmap">${safeJSON(importMap)}</script></head><body>
    <header><div><p class="eyebrow">Interactive diagram</p><h1>${escapedTitle}<span>3D</span></h1></div><div class="controls"><label>Appearance<select id="appearance"><option value="light">Light</option><option value="dark">Dark</option></select></label><button id="reset">Reset view</button><button id="front">Front view</button><button id="fit" aria-pressed="false">Fit diagram</button></div></header>
    <main><div class="diagram-panel"><p id="view-status" aria-live="polite"></p><div id="viewport" tabindex="0" aria-label="Scrollable interactive three-dimensional diagram"><div id="scene"></div></div></div><aside><label for="selection">Inspect an element</label><select id="selection"><option value="">Entire diagram</option></select><p id="details" aria-live="polite">Select a node, group, or connection to trace its relationships.</p><div id="legend"><p>Drag to orbit · Wheel to zoom</p><p>Right-drag to pan · Select to trace</p><p>On touch screens: scroll the readable view; use Fit diagram to orbit and pinch.</p></div><details><summary>Mermaid source</summary><pre id="source"></pre></details></aside></main><footer><span id="summary"></span></footer><script type="application/json" id="diagram-data">${safeJSON(payload)}</script><script type="module">import "diagram/viewer";</script></body></html>`;
}
