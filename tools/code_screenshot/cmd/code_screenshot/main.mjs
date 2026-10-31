import fs from "node:fs/promises";
import path from "node:path";
import { parseArgs } from "node:util";
import puppeteer from "puppeteer";
import hljs from "highlight.js";

async function main() {
    const { values } = parseArgs({
        options: {
            input: { type: "string" },
            output: { type: "string" },
            language: { type: "string" },
            theme: { type: "string", default: "light" },
            width: { type: "string", default: "960" },
            scale: { type: "string", default: "2" },
            help: { type: "boolean" },
        },
    });
    if (values.help) {
        console.log(
            "code_screenshot --input FILE --output IMAGE.png [--theme light|dark] [--width 960] [--scale 2] [--language NAME]",
        );
        return;
    }
    if (!values.input || !values.output)
        throw new Error("--input and --output are required");
    if (!["light", "dark"].includes(values.theme))
        throw new Error("--theme must be light or dark");
    const width = Number(values.width);
    const scale = Number(values.scale);
    if (!Number.isInteger(width) || width < 100 || width > 10000)
        throw new Error("--width must be an integer between 100 and 10000");
    if (!Number.isFinite(scale) || scale <= 0 || scale > 4)
        throw new Error("--scale must be greater than zero and at most 4");
    const base = process.env.BUILD_WORKSPACE_DIRECTORY || process.cwd();
    const input = path.resolve(base, values.input);
    const output = path.resolve(base, values.output);
    if (input === output)
        throw new Error("Input and output must be different files");
    if (path.extname(output).toLowerCase() !== ".png")
        throw new Error("--output must have a .png extension");
    const source = await fs.readFile(input, "utf8");
    const normalized = source.replace(/\r\n?|\n/g, "\n");
    const extension = path.extname(input).slice(1).toLowerCase();
    const language =
        values.language ||
        ({
            mjs: "javascript",
            cjs: "javascript",
            bazel: "python",
            bzl: "python",
            h: "cpp",
        }[extension] ??
            extension);
    if (values.language && !hljs.getLanguage(language))
        throw new Error(`Unsupported language: ${language}`);
    const highlighted = hljs.getLanguage(language)
        ? hljs.highlight(normalized, { language, ignoreIllegals: true }).value
        : normalized
              .replace(/&/g, "&amp;")
              .replace(/</g, "&lt;")
              .replace(/>/g, "&gt;");
    const css = await fs.readFile(process.env.CODE_SCREENSHOT_THEME, "utf8");
    const font = await fs.readFile(process.env.CODE_SCREENSHOT_FONT);
    const browser = await puppeteer.launch({
        executablePath: path.resolve(process.env.PUPPETEER_EXECUTABLE_PATH),
        headless: true,
        args: ["--no-sandbox", "--disable-dev-shm-usage"],
    });
    try {
        const page = await browser.newPage();
        await page.setJavaScriptEnabled(false);
        await page.setRequestInterception(true);
        page.on("request", (request) => request.abort());
        await page.setViewport({
            width,
            height: 600,
            deviceScaleFactor: scale,
        });
        await page.setContent(
            '<!doctype html><html><head><meta charset="utf-8"></head><body><main><header></header><pre></pre></main></body></html>',
        );
        await page.addStyleTag({ content: css });
        await page.addStyleTag({
            content: `
            @font-face { font-family: CaptureMono; src: url(data:font/ttf;base64,${font.toString("base64")}); }
            body { margin: 0; padding: 32px; width: ${width}px; }
            main { border: 1px solid var(--bs-border-color); border-radius: .625rem; overflow: hidden; }
            header { font: 14px/1.65 CaptureMono, monospace; padding: 12px 20px; border-bottom: 1px solid var(--bs-border-color); color: var(--bs-secondary-color); }
            pre { margin: 0; padding: 20px; overflow: visible; background: var(--bs-tertiary-bg); color: var(--bs-body-color); font: 14px/1.65 CaptureMono, monospace; tab-size: 4; white-space: pre-wrap; overflow-wrap: anywhere; }
            .source-line { display: grid; grid-template-columns: 5ch minmax(0, 1fr); }
            .line-number { color: var(--bs-secondary-color); text-align: right; padding-right: 2ch; user-select: none; }
            .line-text { min-width: 0; }
            .hljs-comment, .hljs-quote { color: var(--bs-secondary-color); }
            .hljs-keyword, .hljs-selector-tag { color: var(--bs-emphasis-color); font-weight: 600; }
        `,
        });
        // DevTools evaluation populates inert text nodes; source is never evaluated.
        await page.evaluate(
            (text, title, theme) => {
                document.documentElement.dataset.bsTheme = theme;
                document.querySelector("header").textContent = title;
                const pre = document.querySelector("pre");
                // Parse only highlighter-produced span markup in a detached template.
                // Then rebuild every source fragment using textContent.
                const template = document.createElement("template");
                template.innerHTML = text;
                const lines = [[]];
                function visit(node, classes = []) {
                    if (node.nodeType === Node.TEXT_NODE) {
                        const pieces = node.textContent.split("\n");
                        for (const [index, piece] of pieces.entries()) {
                            if (index) lines.push([]);
                            if (piece)
                                lines.at(-1).push({ text: piece, classes });
                        }
                    } else {
                        const next =
                            node.nodeType === Node.ELEMENT_NODE
                                ? [...classes, ...node.classList]
                                : classes;
                        for (const child of node.childNodes)
                            visit(child, next);
                    }
                }
                visit(template.content);
                if (lines.length > 1 && lines.at(-1).length === 0) lines.pop();
                for (const [index, fragments] of lines.entries()) {
                    const row = document.createElement("span");
                    row.className = "source-line";
                    const number = document.createElement("span");
                    number.className = "line-number";
                    number.textContent = String(index + 1);
                    const content = document.createElement("span");
                    content.className = "line-text";
                    for (const fragment of fragments) {
                        const span = document.createElement("span");
                        span.className = fragment.classes.join(" ");
                        span.textContent = fragment.text;
                        content.append(span);
                    }
                    if (!fragments.length) content.textContent = "\u200b";
                    row.append(number, content);
                    pre.append(row);
                }
            },
            highlighted,
            path.basename(input),
            values.theme,
        );
        await page.evaluate(() => document.fonts.ready);
        const bounds = await page.evaluate(() => ({
            width: document.body.scrollWidth,
            height: document.body.scrollHeight,
        }));
        if (
            bounds.width * scale > 30000 ||
            bounds.height * scale > 30000 ||
            bounds.width * bounds.height * scale * scale > 100000000
        )
            throw new Error(
                "Source exceeds the supported image dimensions; capture a smaller file or reduce --scale",
            );
        await fs.mkdir(path.dirname(output), { recursive: true });
        await page.screenshot({
            path: output,
            type: "png",
            clip: { x: 0, y: 0, ...bounds },
            captureBeyondViewport: true,
        });
        console.log(output);
    } finally {
        await browser.close();
    }
}
main().catch((error) => {
    console.error(`code_screenshot: ${error.message}`);
    process.exitCode = 1;
});
