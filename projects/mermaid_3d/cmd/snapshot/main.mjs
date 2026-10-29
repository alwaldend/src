import fs from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { parseArgs } from "node:util";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");

async function main() {
    const { values } = parseArgs({
        options: {
            input: { type: "string" },
            output: { type: "string" },
            browser: { type: "string" },
            width: { type: "string", default: "1280" },
            height: { type: "string", default: "960" },
            "color-scheme": { type: "string", default: "light" },
            camera: { type: "string", default: "isometric" },
            "timeout-ms": { type: "string", default: "30000" },
        },
    });
    for (const name of ["input", "output", "browser"])
        if (!values[name]) throw new Error(`Missing --${name}`);
    const extension = path.extname(values.output).slice(1);
    if (!["png", "webp"].includes(extension))
        throw new Error("Output must have a .png or .webp extension");
    const integer = (name, minimum, maximum) => {
        const value = Number(values[name]);
        if (
            !/^\d+$/.test(values[name]) ||
            !Number.isSafeInteger(value) ||
            value < minimum ||
            value > maximum
        )
            throw new Error(
                `--${name} must be an integer from ${minimum} to ${maximum}`,
            );
        return value;
    };
    const width = integer("width", 128, 8192);
    const height = integer("height", 128, 8192);
    if (width * height > 16777216)
        throw new Error(
            "Snapshot dimensions must not exceed 16777216 pixels total",
        );
    const timeout = integer("timeout-ms", 1000, 120000);
    const colorScheme = values["color-scheme"];
    if (!["light", "dark"].includes(colorScheme))
        throw new Error("--color-scheme must be light or dark");
    if (!["isometric", "front"].includes(values.camera))
        throw new Error("--camera must be isometric or front");

    const input = path.resolve(values.input);
    const output = path.resolve(values.output);
    await fs.access(input);
    await fs.mkdir(path.dirname(output), { recursive: true });
    // Keep Chrome's disposable profile within the action's output directory.
    // No host profile, installed browser, font path, or network asset is used.
    const profile = await fs.mkdtemp(
        path.join(path.dirname(output), ".snapshot-"),
    );
    let browser;
    let timer;
    let stage = "browser startup";
    try {
        browser = await puppeteer.launch({
            executablePath: path.resolve(values.browser),
            headless: "shell",
            userDataDir: profile,
            timeout,
            env: { ...process.env, TZ: "UTC", LANG: "C.UTF-8" },
            args: [
                "--use-angle=swiftshader",
                "--enable-unsafe-swiftshader",
                "--disable-background-networking",
                "--host-resolver-rules=MAP * ~NOTFOUND",
            ],
        });
        const page = await browser.newPage();
        page.setDefaultTimeout(timeout);
        page.setDefaultNavigationTimeout(timeout);
        const url = pathToFileURL(input).href;
        let fail;
        const failure = new Promise((_, reject) => {
            fail = reject;
        });
        // The listener can fire during setup before Promise.race is attached.
        void failure.catch(() => {});
        page.on("pageerror", (error) =>
            fail(new Error(`Viewer error: ${error.message}`)),
        );
        page.on("error", (error) =>
            fail(new Error(`Browser page failed: ${error.message}`)),
        );
        await page.setBypassServiceWorker(true);
        await page.setCacheEnabled(false);
        await page.setOfflineMode(true);
        await page.setRequestInterception(true);
        let documentRequested = false;
        page.on("request", (request) => {
            const target = request.url();
            if (
                !documentRequested &&
                target === url &&
                request.isNavigationRequest() &&
                request.frame() === page.mainFrame()
            ) {
                documentRequested = true;
                void request.continue().catch(fail);
            } else if (/^(data|blob):/.test(target)) {
                void request.continue().catch(fail);
            } else {
                // Reject local file dependencies as well as remote resources.
                // Report only the scheme: URLs can carry unrelated query data.
                fail(
                    new Error(
                        `Offline snapshot requires self-contained HTML; blocked ${target.split(":", 1)[0]}: request`,
                    ),
                );
                void request.abort("blockedbyclient").catch(() => {});
            }
        });
        await page.setViewport({ width, height, deviceScaleFactor: 1 });
        await page.emulateMediaFeatures([
            { name: "prefers-color-scheme", value: colorScheme },
        ]);
        timer = setTimeout(
            () =>
                fail(
                    new Error(
                        `Snapshot ${stage} timed out after ${timeout} ms`,
                    ),
                ),
            timeout,
        );
        const encoded = await Promise.race([
            failure,
            (async () => {
                stage = "viewer readiness";
                await page.goto(url, { waitUntil: "load" });
                await page.waitForFunction(
                    () => window.diagramView?.ready === true,
                );
                stage = "font loading and rendering";
                return page.evaluate(
                    async (options) => {
                        const payload = JSON.parse(
                            document.getElementById("diagram-data")
                                .textContent,
                        );
                        const faces = await document.fonts.load(
                            `16px ${JSON.stringify(payload.fontName)}`,
                        );
                        await document.fonts.ready;
                        // fonts.check() alone also accepts absent families with a
                        // system fallback. Require a real, successfully loaded face.
                        if (
                            !faces.length ||
                            faces.some((face) => face.status !== "loaded")
                        )
                            throw new Error(
                                "The embedded diagram font is missing or failed to load",
                            );
                        if (
                            typeof window.diagramView.prepareSnapshot !==
                            "function"
                        )
                            throw new Error(
                                "Viewer does not support prepareSnapshot",
                            );
                        const view =
                            window.diagramView.prepareSnapshot(options);
                        const canvas = document.querySelector("#scene canvas");
                        if (
                            !canvas ||
                            canvas.width !== options.width ||
                            canvas.height !== options.height ||
                            view.canvas.width !== options.width ||
                            view.canvas.height !== options.height
                        )
                            throw new Error(
                                "Viewer produced incorrect snapshot dimensions",
                            );
                        if (
                            view.appearance !== options.colorScheme ||
                            view.selected ||
                            !view.fontLoaded
                        )
                            throw new Error(
                                "Viewer did not apply the snapshot appearance or embedded font",
                            );
                        if (!view.labels.length || !view.sceneCorners.length)
                            throw new Error(
                                "Viewer produced an empty snapshot scene",
                            );
                        const inside = (x, y) =>
                            Number.isFinite(x) &&
                            Number.isFinite(y) &&
                            x >= 0 &&
                            x <= options.width &&
                            y >= 0 &&
                            y <= options.height;
                        if (
                            view.sceneCorners.some(
                                ([x, y]) => !inside(x, y),
                            ) ||
                            view.labels.some(
                                ({ bounds }) =>
                                    !inside(bounds.left, bounds.top) ||
                                    !inside(bounds.right, bounds.bottom),
                            )
                        )
                            throw new Error(
                                "Snapshot would clip scene geometry or labels",
                            );
                        const direction = view.camera.map(
                            (value, index) => value - view.target[index],
                        );
                        const expected =
                            options.camera === "front"
                                ? [0, 0, 1]
                                : view.model.viewDirection;
                        const cosine =
                            direction.reduce(
                                (sum, value, index) =>
                                    sum + value * expected[index],
                                0,
                            ) /
                            Math.hypot(...direction) /
                            Math.hypot(...expected);
                        if (!(cosine > 1 - 1e-8))
                            throw new Error(
                                "Viewer did not apply the requested snapshot camera",
                            );
                        // The viewer owns scene creation, framing and drawing.
                        // Direct canvas encoding excludes all DOM controls without
                        // a second renderer, screenshot crop, or image resampling.
                        const data = canvas.toDataURL(
                            `image/${options.extension}`,
                            1,
                        );
                        if (
                            !data.startsWith(
                                `data:image/${options.extension};base64,`,
                            )
                        )
                            throw new Error(
                                `Browser cannot encode ${options.extension} snapshots`,
                            );
                        return data.slice(data.indexOf(",") + 1);
                    },
                    {
                        width,
                        height,
                        colorScheme,
                        camera: values.camera,
                        extension,
                    },
                );
            })(),
        ]);
        clearTimeout(timer);
        const bytes = Buffer.from(encoded, "base64");
        if (!bytes.length)
            throw new Error("Browser returned an empty snapshot image");
        await fs.writeFile(output, bytes);
    } catch (error) {
        throw new Error(`Snapshot ${stage} failed: ${error.message}`, {
            cause: error,
        });
    } finally {
        clearTimeout(timer);
        if (browser) await browser.close();
        await fs.rm(profile, { recursive: true, force: true });
    }
}

try {
    await main();
} catch (error) {
    console.error(error.message);
    process.exitCode = 1;
}
