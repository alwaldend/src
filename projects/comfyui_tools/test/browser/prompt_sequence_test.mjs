/** Exercise the field editor in Chromium through ComfyUI's DOM-widget contract. */
import assert from "node:assert/strict";
import { once } from "node:events";
import fs from "node:fs/promises";
import http from "node:http";
import { createRequire } from "node:module";
import path from "node:path";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [script, chrome] = process.argv
    .slice(2)
    .map((file) => path.resolve(file));
const source = await fs.readFile(script, "utf8");
const app = `
export const app = { registerExtension(extension) { window.extension = extension; } };
window.createNode = (raw, separator = "line") => {
  document.querySelector("main").replaceChildren();
  const node = {
    comfyClass: "ImagegenPromptSequence", widgets: [
      {name: "prompts", value: raw, onRemove() {}},
      {name: "context", value: "Shared robot reference."},
      {name: "separator", value: separator},
      {name: "filename_prefix", value: "comic"}
    ],
    addDOMWidget(name, type, element, options) {
      document.querySelector("main").append(element);
      const widget = {name, type, options, element};
      Object.defineProperty(widget, "value", {
        get() { return options.getValue(); },
        set(value) { options.setValue(value); }
      });
      this.widgets.push(widget);
      return widget;
    },
    setSize() {}, computeSize() { return [400, 400]; }, setDirtyCanvas() {}
  };
  window.extension.nodeCreated(node);
  window.node = node;
  return node;
};`;
const server = http.createServer((request, response) => {
    const routes = {
        "/": [
            "text/html",
            '<main></main><script type="module" src="/extensions/prompt_sequence/index.js"></script>',
        ],
        "/scripts/app.js": ["text/javascript", app],
        "/extensions/prompt_sequence/index.js": ["text/javascript", source],
    };
    const resource = routes[request.url];
    if (!resource) {
        response.writeHead(404);
        response.end();
        return;
    }
    response.setHeader("Content-Type", resource[0]);
    response.end(resource[1]);
});
server.listen(0, "127.0.0.1");
await once(server, "listening");
const browser = await puppeteer.launch({
    executablePath: chrome,
    headless: true,
    args: ["--no-sandbox", "--disable-dev-shm-usage", "--disable-gpu"],
});
try {
    const page = await browser.newPage();
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    await page.goto(`http://127.0.0.1:${server.address().port}/`);
    await page.waitForFunction(() => window.extension);
    await page.evaluate(() =>
        window.createNode("A robot rides a bicycle.\nA robot reads a book."),
    );
    assert.equal(await page.$$eval("textarea", (fields) => fields.length), 3);
    await page.focus('textarea[aria-label="Panel 3"]');
    await page.keyboard.type("A robot watches snow.");
    assert.equal(await page.$$eval("textarea", (fields) => fields.length), 4);
    assert.equal(
        await page.evaluate(() =>
            document.activeElement.getAttribute("aria-label"),
        ),
        "Panel 3",
    );
    await page.evaluate(() => {
        const fields = document.querySelectorAll("textarea");
        fields[0].value = "A robot.\n\nSnow falls. <robot>";
        fields[0].dispatchEvent(new Event("input", { bubbles: true }));
        fields[1].value = "";
        fields[1].dispatchEvent(new Event("input", { bubbles: true }));
    });
    const saved = await page.evaluate(() =>
        window.node.widgets.map((widget) => widget.value),
    );
    assert.deepEqual(JSON.parse(saved[0]), [
        "A robot.\n\nSnow falls. <robot>",
        "A robot watches snow.",
    ]);
    assert.equal(saved[1], "Shared robot reference.");
    assert.equal(saved[2], "JSON");
    assert.equal(saved[3], "comic");
    // LiteGraph cloning configures widgets in order without loadedGraphNode.
    await page.evaluate((saved) => {
        const node = window.createNode("default");
        node.widgets.forEach((widget, index) => {
            widget.value = saved[index];
        });
    }, saved);
    assert.deepEqual(
        await page.$$eval("textarea", (fields) =>
            fields.map((field) => field.value),
        ),
        [...JSON.parse(saved[0]), ""],
        "a duplicate must restore separate panels after its separator is restored",
    );
    await page.evaluate(() => {
        const first = document.querySelector("textarea");
        first.value = "A robot sketches a snowy park.";
        first.dispatchEvent(new Event("input", { bubbles: true }));
    });
    assert.deepEqual(
        await page.evaluate(() => JSON.parse(window.node.widgets[0].value)),
        ["A robot sketches a snowy park.", "A robot watches snow."],
        "editing a duplicate must preserve its other panels",
    );
    await page.evaluate((saved) => {
        const node = window.createNode("default");
        node.widgets.forEach((widget, index) => {
            widget.value = saved[index];
        });
        window.extension.loadedGraphNode(node);
    }, saved);
    assert.deepEqual(
        await page.$$eval("textarea", (fields) =>
            fields.map((field) => field.value),
        ),
        [...JSON.parse(saved[0]), ""],
    );
    assert.equal(
        await page.$("robot"),
        null,
        "prompt markup must remain literal text",
    );
    await page.evaluate(() =>
        window.createNode(
            "A robot.\nSnow falls.\n\nA robot reads.",
            "blank line",
        ),
    );
    assert.deepEqual(
        await page.$$eval("textarea", (fields) =>
            fields.map((field) => field.value),
        ),
        ["A robot.\nSnow falls.", "A robot reads.", ""],
    );
    await page.evaluate(() =>
        window.extension.nodeCreated({ comfyClass: "OtherNode" }),
    );
    assert.deepEqual(errors, []);
    const artifacts = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
    assert.ok(artifacts);
    await page.screenshot({
        path: path.join(artifacts, "prompt-sequence.png"),
    });
    await fs.writeFile(
        path.join(artifacts, "prompt-sequence-result.json"),
        JSON.stringify(
            {
                evidence:
                    "browser fixture with ComfyUI DOM-widget lifecycle; no rental or GPU",
                saved,
                checks: [
                    "automatic growth",
                    "focus retention",
                    "empty fields omitted",
                    "paragraph preservation",
                    "ordered save/reload",
                    "node duplication without graph-load hooks",
                    "editing duplicated panels",
                    "existing separators",
                    "literal prompt markup",
                    "unrelated nodes untouched",
                ],
            },
            null,
            2,
        ),
    );
} finally {
    await browser.close();
    server.closeAllConnections();
    await new Promise((resolve) => server.close(resolve));
}
