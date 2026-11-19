import fs from "node:fs/promises";
import { createHash } from "node:crypto";
import path from "node:path";
import { pathToFileURL } from "node:url";
import puppeteer from "puppeteer";

const [webapp, browserPath, input, manifestPath, fontAnchor] =
    process.argv.slice(2);
if (!fontAnchor)
    throw new Error(
        "Expected webapp, browser, input, output manifest and font anchor",
    );
const xml = await fs.readFile(input, "utf8");
const outputs = JSON.parse(await fs.readFile(manifestPath, "utf8"));
const scratch = await fs.mkdtemp(path.resolve(".drawio-render-"));
const fontConfig = path.join(scratch, "fonts.conf");
const escapeXml = (s) =>
    s
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll('"', "&quot;");
const aliases = {
    Helvetica: "Liberation Sans",
    Arial: "Liberation Sans",
    "sans-serif": "Liberation Sans",
    Times: "Liberation Serif",
    "Times New Roman": "Liberation Serif",
    serif: "Liberation Serif",
    Courier: "Liberation Mono",
    monospace: "Liberation Mono",
};
await fs.writeFile(
    fontConfig,
    `<?xml version="1.0"?><fontconfig><dir>${escapeXml(path.dirname(path.resolve(fontAnchor)))}</dir><cachedir>${escapeXml(path.join(scratch, "font-cache"))}</cachedir>${Object.entries(
        aliases,
    )
        .map(
            ([family, target]) =>
                `<alias><family>${family}</family><prefer><family>${target}</family></prefer></alias>`,
        )
        .join("")}</fontconfig>`,
);
let browser;
try {
    browser = await puppeteer.launch({
        executablePath: path.resolve(browserPath),
        headless: true,
        args: ["--no-sandbox", "--disable-dev-shm-usage"],
        env: {
            ...process.env,
            FONTCONFIG_FILE: fontConfig,
            FONTCONFIG_PATH: scratch,
        },
    });
    for (const [pageName, output] of Object.entries(outputs)) {
        const page = await browser.newPage();
        const errors = [];
        page.on("pageerror", (error) => errors.push(error.message));
        await page.setRequestInterception(true);
        page.on("request", (request) => {
            const url = request.url();
            if (url.startsWith("file:") || url.startsWith("data:"))
                request.continue();
            else {
                errors.push(`Network resource is unsupported: ${url}`);
                request.abort();
            }
        });
        await page.setUserAgent(
            "Mozilla/5.0 Chrome/149.0.0.0 Electron/42.0.0 draw.io/30.2.6",
        );
        await page.evaluateOnNewDocument(() => {
            window.drawioListeners = {};
            window.drawioMessages = {};
            window.electron = {
                registerMsgListener: (name, callback) => {
                    window.drawioListeners[name] = callback;
                },
                sendMessage: (name, value) => {
                    window.drawioMessages[name] = value;
                },
            };
        });
        await page.goto(
            pathToFileURL(path.resolve(webapp, "export3.html")).href,
            { waitUntil: "load" },
        );
        await page.waitForFunction(
            () => typeof window.drawioListeners.render === "function",
        );
        await page.evaluate(
            ({ xml, pageName }) => {
                const doc = new DOMParser().parseFromString(
                    xml,
                    "application/xml",
                );
                if (doc.querySelector("parsererror"))
                    throw new Error("Invalid Drawio XML");
                const pages = [...doc.querySelectorAll("mxfile > diagram")];
                const matching = pages.filter(
                    (p) => p.getAttribute("name") === pageName,
                );
                if (matching.length !== 1)
                    throw new Error(
                        `Expected exactly one page named ${pageName}`,
                    );
                const index = pages.indexOf(matching[0]);
                window.drawioListeners.render({
                    xml,
                    format: "svg",
                    from: index,
                    to: index,
                    scale: 1,
                    border: 8,
                    theme: "light",
                    embedXml: "0",
                    embedImages: "0",
                    embedFonts: "0",
                });
            },
            { xml, pageName },
        );
        await page.waitForFunction(
            () =>
                "render-finished" in window.drawioMessages ||
                "export-error" in window.drawioMessages,
        );
        const info = await page.evaluate(() => window.drawioMessages);
        if (info["export-error"] || !info["render-finished"])
            throw new Error(`Drawio render failed: ${pageName}`);
        await page.evaluate(() => window.drawioListeners["get-svg-data"]());
        await page.waitForFunction(() => "svg-data" in window.drawioMessages);
        let svg = await page.evaluate(() => window.drawioMessages["svg-data"]);
        // Drawio's SVG exporter ignores the PDF-only bg option. Paint an
        // opaque canvas so light-theme labels remain readable on dark pages
        // and when the image is opened or embedded outside these docs.
        svg = svg.replace(
            /(<svg\b[^>]*>)/,
            '$1<rect width="100%" height="100%" fill="#ffffff"/>',
        );
        // Drawio gives the SVG root a random ID used by its CSS selector.
        // Stabilize that ID and its references without changing diagram cells.
        const rootId = svg.match(/<svg\b[^>]*\bid="(ge-svg-[A-Za-z0-9_-]+)"/);
        if (rootId) {
            const stableId =
                "ge-svg-" +
                createHash("sha256")
                    .update(xml)
                    .update(pageName)
                    .digest("hex")
                    .slice(0, 20);
            svg = svg.replaceAll(rootId[1], stableId);
        }
        if (errors.length) throw new Error(errors.join("\n"));
        if (
            !svg.includes("<svg") ||
            !/<(?:path|rect|ellipse|polygon)\b/.test(svg)
        )
            throw new Error(`Empty SVG: ${pageName}`);
        await fs.mkdir(path.dirname(output), { recursive: true });
        await fs.writeFile(
            output,
            svg.replace(
                "<svg",
                `<!-- Source SHA-256: ${createHash("sha256").update(xml).digest("hex")} -->\n<svg`,
            ) + "\n",
        );
        await page.close();
    }
} finally {
    await browser?.close();
    await fs.rm(scratch, { recursive: true, force: true });
}
