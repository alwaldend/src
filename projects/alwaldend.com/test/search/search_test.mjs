// Exercise the built site, its real search worker, navigation, and failure recovery.
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import http from "node:http";
import path from "node:path";
import { createRequire } from "node:module";
const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const root = path.resolve(process.argv[2]);
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
let failIndex = false;
let indexDelay = 0;
const requests = [];
const pause = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
const server = http.createServer(async (request, response) => {
    const pathname = decodeURIComponent(
        new URL(request.url, "http://localhost").pathname,
    );
    requests.push(pathname);
    const filename = path.resolve(
        root,
        "." + pathname,
        pathname.endsWith("/") ? "index.html" : "",
    );
    if (!filename.startsWith(root + path.sep)) {
        response.writeHead(403).end();
        return;
    }
    if (pathname.includes("prebuilt-search-index.")) {
        await pause(indexDelay);
        if (failIndex) {
            response.writeHead(503).end();
            return;
        }
    }
    try {
        response.setHeader(
            "Content-Type",
            {
                ".html": "text/html",
                ".js": "text/javascript",
                ".json": "application/json",
                ".css": "text/css",
                ".svg": "image/svg+xml",
            }[path.extname(filename)] || "application/octet-stream",
        );
        response.setHeader("Cache-Control", "no-store");
        response.end(await fs.readFile(filename));
    } catch {
        response.writeHead(404).end();
    }
});
await new Promise((resolve) => server.listen(0, "127.0.0.1", resolve));
const base = `http://127.0.0.1:${server.address().port}`;
const browser = await puppeteer.launch({
    executablePath: path.resolve(process.argv[3]),
    headless: "shell",
});
const report = { cases: [], errors: [], externalRequests: [], warnings: [] };
let page;
try {
    page = await browser.newPage();
    await page.setViewport({ width: 1280, height: 1000 });
    await page.setRequestInterception(true);
    page.on("request", (request) => {
        if (
            request.url().startsWith(base + "/") ||
            request.url().startsWith("data:")
        )
            request.continue();
        else {
            report.externalRequests.push(request.url());
            request.abort();
        }
    });
    page.on("console", (message) => {
        if (message.type() === "warn") report.warnings.push(message.text());
    });
    page.on("pageerror", (error) => report.errors.push(error.message));
    const input = ".td-navbar .td-search input";
    const results = ".td-offline-search-results .mt-4 a";
    const waitText = async (text) =>
        page.waitForFunction(
            (value) =>
                document
                    .querySelector(".td-offline-search-results")
                    ?.textContent.includes(value),
            {},
            text,
        );
    await page.goto(base + "/", { waitUntil: "networkidle0" });
    await pause(250);
    assert.equal(
        requests.some((url) =>
            /prebuilt-search-index|search-worker/.test(url),
        ),
        false,
        "Navigation does not load the search engine",
    );
    assert.equal(
        await page.evaluate(() => typeof window.lunr),
        "undefined",
        "Lunr is absent from the UI thread",
    );
    assert.match(
        await fs.readFile(path.join(root, "licenses/lunr.txt"), "utf8"),
        /Copyright.*Oliver Nightingale/,
        "The packaged worker retains its library license",
    );
    report.cases.push("No search work or CDN on initial load");

    indexDelay = 800;
    await page.type(input, "download");
    await page.keyboard.press("Enter");
    await waitText("Loading search");
    // Input remains usable while the worker's index request is held by the server.
    await page.keyboard.down("Control");
    await page.keyboard.press("A");
    await page.keyboard.up("Control");
    await page.keyboard.press("Backspace");
    await pause(1200);
    assert.equal(
        await page.$(".td-offline-search-results"),
        null,
        "A cleared query stays dismissed after loading",
    );
    report.cases.push("Clear pending search without stale results");
    indexDelay = 0;
    await page.type(input, "download");
    await page.keyboard.press("Enter");
    await page.waitForSelector(results);
    const entries = await page.$$eval(results, (links) =>
        links.map((link) => ({
            title: link.textContent,
            path: new URL(link.href).pathname,
        })),
    );
    assert(
        entries.some((entry) => entry.path === "/docs/infra/download/"),
        "Search includes the download documentation",
    );
    assert(
        entries.length > 0 && entries.length <= 10,
        "Search respects the default result limit",
    );
    await page.waitForFunction(
        () =>
            getComputedStyle(
                document.querySelector(".td-offline-search-results"),
            ).opacity === "1",
    );
    await page.screenshot({ path: path.join(output, "search-results.png") });
    report.cases.push({ search: "download", entries });
    await Promise.all([
        page.waitForNavigation({ waitUntil: "networkidle0", timeout: 5000 }),
        page.click(
            '.td-offline-search-results a[href="/docs/infra/download/"]',
        ),
    ]);
    assert.equal(new URL(page.url()).pathname, "/docs/infra/download/");
    report.cases.push("First click on a search result navigates");
    await page.goBack({ waitUntil: "networkidle0" });

    await page.click(input);
    await page.keyboard.down("Control");
    await page.keyboard.press("A");
    await page.keyboard.up("Control");
    await page.keyboard.press("Backspace");
    await page.type(
        input,
        "zzzzuniquesearchpayloadzzzz<img/onerror=window.searchInjection=true>",
    );
    await page.keyboard.press("Enter");
    await waitText("No results found");
    assert.equal(await page.evaluate(() => window.searchInjection), undefined);
    assert.equal(await page.$(".td-offline-search-results img"), null);
    report.cases.push("Query markup is rendered as text");

    failIndex = true;
    await page.reload({ waitUntil: "networkidle0" });
    await page.click(input);
    await page.keyboard.down("Control");
    await page.keyboard.press("A");
    await page.keyboard.up("Control");
    await page.keyboard.press("Backspace");
    await page.type(input, "download");
    await page.keyboard.press("Enter");
    await waitText("Search is unavailable");
    failIndex = false;
    await page.keyboard.press("Enter");
    await page.waitForSelector(results);
    report.cases.push("Failed index request retries without page reload");

    await page.goto(base + "/docs/projects/al/", {
        waitUntil: "networkidle0",
    });
    const tagPaths = await page.$$eval('a[href*="/tags/"]', (links) => [
        ...new Set(
            links.map((link) =>
                decodeURIComponent(new URL(link.href).pathname),
            ),
        ),
    ]);
    const tags = (
        await fs.readdir(path.join(root, "tags"), { withFileTypes: true })
    )
        .filter((entry) => entry.isDirectory() && entry.name !== "page")
        .map((entry) => `/tags/${entry.name}/`);
    for (const tag of tags)
        assert(
            tagPaths.includes(tag),
            `The global tag catalogue retains ${tag}`,
        );
    const navigation = await page.$$eval("#td-section-nav a", (links) =>
        links.map((link) => new URL(link.href).pathname),
    );
    assert(
        navigation.includes("/docs/projects/"),
        "Scoped navigation retains the parent link",
    );
    assert(
        !navigation.includes("/docs/projects/android_launcher/"),
        "Unrelated owner trees remain scoped out",
    );
    report.cases.push({
        navigation: "owner scope and full tag catalogue",
        tags: tags.length,
    });

    await page.setViewport({
        width: 390,
        height: 844,
        isMobile: true,
        hasTouch: true,
    });
    await page.goto(base + "/", { waitUntil: "networkidle0" });
    const links = await page.$$(".home-links a");
    let docs;
    for (const link of links)
        if (
            await link.evaluate(
                (element) => new URL(element.href).pathname === "/docs/",
            )
        )
            docs = link;
    assert(docs, "Homepage has a documentation link");
    await Promise.all([
        page.waitForNavigation({ waitUntil: "networkidle0" }),
        docs.tap(),
    ]);
    assert.equal(new URL(page.url()).pathname, "/docs/");
    await page.screenshot({
        path: path.join(output, "mobile-navigation.png"),
    });
    report.cases.push("First mobile tap opens documentation");
    assert.deepEqual(report.errors, []);
    assert.deepEqual(report.externalRequests, []);
} finally {
    report.requests = requests;
    if (page) {
        report.lastPage = await page.evaluate(() => ({
            url: location.href,
            input: document.querySelector(".td-search input")?.value,
            results: document.querySelector(".td-offline-search-results")
                ?.textContent,
        }));
        await page.screenshot({ path: path.join(output, "last-page.png") });
    }
    await fs.writeFile(
        path.join(output, "report.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
    await browser.close();
    await new Promise((resolve) => server.close(resolve));
}
