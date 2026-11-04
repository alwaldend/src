import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";
import { parseArgs } from "node:util";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const origin = "http://127.0.0.1:1313";
const contentTypes = {
    ".html": "text/html",
    ".css": "text/css",
    ".js": "application/javascript",
    ".json": "application/json",
    ".svg": "image/svg+xml",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".webp": "image/webp",
    ".gif": "image/gif",
    ".ico": "image/x-icon",
    ".woff": "font/woff",
    ".woff2": "font/woff2",
    ".ttf": "font/ttf",
};

async function capture() {
    const [siteArgument, browserArgument, ...args] = process.argv.slice(2);
    const { values } = parseArgs({
        args,
        options: {
            "output-dir": { type: "string" },
            path: { type: "string", multiple: true },
            width: { type: "string", default: "1440" },
            height: { type: "string", default: "1000" },
            "color-scheme": { type: "string", default: "light" },
            help: { type: "boolean", short: "h" },
        },
    });
    if (values.help) {
        console.log(`Capture full-page PNGs of the built site without a server.

bazel run //projects/alwaldend.com:site_screenshots -- [options]

  --output-dir DIR       Output directory (default: out/site-screenshots)
  --path PATH            Site-relative page path; repeat to select pages
  --width PIXELS         Viewport width (default: 1440)
  --height PIXELS        Viewport height (default: 1000)
  --color-scheme MODE    light or dark (default: light)

Default pages: /, /projects/, /docs/projects/
Output: one full-page PNG per path and screenshots.json`);
        return;
    }
    const width = Number(values.width);
    const height = Number(values.height);
    if (
        ![width, height].every(
            (value) => Number.isSafeInteger(value) && value > 0,
        )
    ) {
        throw new Error("Viewport dimensions must be positive integers");
    }
    const colorScheme = values["color-scheme"];
    if (!["light", "dark"].includes(colorScheme))
        throw new Error("Color scheme must be light or dark");
    const pages = (values.path ?? ["/", "/projects/", "/docs/projects/"]).map(
        (input) => {
            const url = new URL(input, origin);
            if (
                !input.startsWith("/") ||
                url.origin !== origin ||
                url.search ||
                url.hash
            ) {
                throw new Error(
                    "Page paths must be site-relative pathnames without queries or fragments",
                );
            }
            if (!path.extname(url.pathname) && !url.pathname.endsWith("/"))
                url.pathname += "/";
            const name =
                url.pathname.split("/").filter(Boolean).join("-") || "root";
            return {
                path: url.pathname,
                file: `${name.replace(/[^a-zA-Z0-9_.%-]/g, "_")}.png`,
            };
        },
    );
    if (new Set(pages.map((page) => page.file)).size !== pages.length) {
        throw new Error(
            "Page paths must produce distinct screenshot filenames",
        );
    }
    const site = await fs.realpath(siteArgument);
    const output = path.resolve(
        process.env.BUILD_WORKSPACE_DIRECTORY || process.cwd(),
        values["output-dir"] || "out/site-screenshots",
    );
    await fs.mkdir(output, { recursive: true });
    const profile = await fs.mkdtemp(path.join(output, ".browser-"));
    let browser;
    try {
        browser = await puppeteer.launch({
            executablePath: path.resolve(browserArgument),
            headless: "shell",
            userDataDir: profile,
        });
        const page = await browser.newPage();
        await page.setViewport({ width, height, deviceScaleFactor: 1 });
        await page.emulateMediaFeatures([
            { name: "prefers-color-scheme", value: colorScheme },
        ]);
        let blockedExternalRequests = 0;
        await page.setRequestInterception(true);
        page.on("request", async (request) => {
            if (request.url().startsWith("data:")) return request.continue();
            const url = new URL(request.url());
            if (url.origin !== origin) {
                blockedExternalRequests++;
                return request.abort();
            }
            try {
                const pathname = decodeURIComponent(url.pathname);
                const filename = await fs.realpath(
                    path.join(
                        site,
                        pathname.endsWith("/")
                            ? `${pathname}index.html`
                            : pathname,
                    ),
                );
                const relative = path.relative(site, filename);
                if (
                    relative === ".." ||
                    relative.startsWith(`..${path.sep}`) ||
                    path.isAbsolute(relative)
                ) {
                    return request.respond({
                        status: 403,
                        body: "Outside site directory",
                    });
                }
                await request.respond({
                    status: 200,
                    contentType:
                        contentTypes[path.extname(filename)] ||
                        "application/octet-stream",
                    body: await fs.readFile(filename),
                });
            } catch {
                await request.respond({ status: 404, body: "Not found" });
            }
        });
        const manifest = {
            viewport: { width, height },
            colorScheme,
            fullPage: true,
            pages: [],
        };
        for (const selected of pages) {
            blockedExternalRequests = 0;
            const response = await page.goto(`${origin}${selected.path}`, {
                waitUntil: "networkidle0",
            });
            if (!response?.ok())
                throw new Error(
                    `Requested page ${selected.path} returned ${response?.status() ?? "no response"}`,
                );
            await page.evaluate(() => document.fonts.ready);
            const png = Buffer.from(
                await page.screenshot({
                    path: path.join(output, selected.file),
                    fullPage: true,
                }),
            );
            manifest.pages.push({
                ...selected,
                title: await page.title(),
                image: {
                    width: png.readUInt32BE(16),
                    height: png.readUInt32BE(20),
                },
                blockedExternalRequests,
            });
        }
        await fs.writeFile(
            path.join(output, "screenshots.json"),
            JSON.stringify(manifest, null, 2) + "\n",
        );
        console.log(JSON.stringify({ outputDirectory: output, ...manifest }));
    } finally {
        try {
            if (browser) await browser.close();
        } finally {
            await fs.rm(profile, { recursive: true, force: true });
        }
    }
}

capture().catch((error) => {
    console.error(`Site screenshots: ${error.message}`);
    process.exitCode = 1;
});
