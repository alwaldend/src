import assert from "node:assert/strict";
import fs from "node:fs/promises";
import http from "node:http";
import path from "node:path";
import { createRequire } from "node:module";

// The repository's shared npm workspace owns the pinned browser dependency.
const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");

const root = path.resolve(process.argv[2]);
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
const server = http.createServer(async (request, response) => {
    const pathname = decodeURIComponent(
        new URL(request.url, "http://localhost").pathname,
    );
    const filename = path.join(
        root,
        pathname.endsWith("/") ? `${pathname}index.html` : pathname,
    );
    try {
        const types = {
            ".html": "text/html",
            ".css": "text/css",
            ".js": "text/javascript",
            ".svg": "image/svg+xml",
        };
        response.setHeader(
            "Content-Type",
            types[path.extname(filename)] || "application/octet-stream",
        );
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
const report = [];
try {
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    // Read exact CDN pins from the built page, whose locked Docsy theme owns
    // them. Only declared scripts carrying SRI may leave the local server.
    const html = await fs.readFile(
        path.join(root, "docs/users/simeonwarren/host_bot/index.html"),
        "utf8",
    );
    const scripts = new Set(
        [...html.matchAll(/<script\b[^>]*>/g)]
            .map(
                ([tag]) =>
                    tag.includes('integrity="') &&
                    tag.match(/src="(https:\/\/[^"]+)"/)?.[1],
            )
            .filter(Boolean),
    );
    page.on("request", (request) => {
        const url = request.url();
        if (
            url.startsWith(`${base}/`) ||
            url.startsWith("data:") ||
            scripts.has(url)
        )
            request.continue();
        else request.abort();
    });
    page.on("pageerror", (error) =>
        report.push({ browserError: error.message }),
    );
    page.on("requestfailed", (request) => {
        if (new URL(request.url()).pathname.endsWith(".svg"))
            report.push({
                failedDiagram: request.url(),
                reason: request.failure().errorText,
            });
    });
    await page.setViewport({ width: 1280, height: 1000 });
    await page.emulateMediaFeatures([
        { name: "prefers-color-scheme", value: "light" },
    ]);
    await page.goto(`${base}/docs/users/simeonwarren/host_bot/`, {
        waitUntil: "networkidle0",
    });
    const visibleImage = async () => {
        const images = await page.$$(
            'img[alt="Host Bot architecture diagram"]',
        );
        const visible = [];
        for (const image of images)
            if (
                await image.evaluate(
                    (element) => getComputedStyle(element).display !== "none",
                )
            )
                visible.push(image);
        assert.equal(
            visible.length,
            1,
            "Exactly one accessible diagram is displayed",
        );
        return visible[0];
    };
    let image = await visibleImage();
    assert.ok(
        image,
        "Documentation diagram is accessible by its alternative text",
    );
    await image.scrollIntoView();
    await image.evaluate((element) => element.decode());
    const initialBox = await image.boundingBox();

    const select = async (mode) => {
        // Exercise the real menu; its options are hidden until the dropdown opens.
        const option = await page.$(`[data-bs-theme-value="${mode}"]`);
        assert.ok(option, `Theme option ${mode} exists`);
        await option.evaluate((element) =>
            element
                .closest(".dropdown")
                .querySelector('[data-bs-toggle="dropdown"]')
                .click(),
        );
        await option.click();
    };
    const check = async (name, mode) => {
        await page.waitForFunction(
            (expected) =>
                document.documentElement.dataset.bsTheme === expected,
            {},
            mode,
        );
        image = await visibleImage();
        await image.scrollIntoView();
        await page.waitForFunction(
            (element) => element.complete && element.naturalWidth > 0,
            {},
            image,
        );
        await image.evaluate((element) => element.decode());
        // Let the newly displayed image paint after the selector changes.
        await page.evaluate(
            () =>
                new Promise((resolve) =>
                    requestAnimationFrame(() =>
                        requestAnimationFrame(resolve),
                    ),
                ),
        );
        const png = await image.screenshot();
        await fs.writeFile(path.join(output, `${name}.png`), png);
        const colors = await page.evaluate(async (encoded) => {
            const sample = new Image();
            sample.src = `data:image/png;base64,${encoded}`;
            await sample.decode();
            const canvas = document.createElement("canvas");
            canvas.width = sample.width;
            canvas.height = sample.height;
            const context = canvas.getContext("2d");
            context.drawImage(sample, 0, 0);
            const pixels = context.getImageData(
                0,
                0,
                canvas.width,
                canvas.height,
            ).data;
            const style = getComputedStyle(document.body);
            const rgb = (value) => value.match(/\d+/g).slice(0, 3).map(Number);
            const background = rgb(style.backgroundColor);
            const foreground = rgb(style.color);
            const count = (color) => {
                let matches = 0;
                for (let i = 0; i < pixels.length; i += 4)
                    if (
                        color.every(
                            (channel, j) =>
                                Math.abs(pixels[i + j] - channel) <= 2,
                        )
                    )
                        matches++;
                return matches;
            };
            let contrastingPixels = 0;
            for (let i = 0; i < pixels.length; i += 4)
                if (
                    background.some(
                        (channel, j) => Math.abs(pixels[i + j] - channel) > 32,
                    )
                )
                    contrastingPixels++;
            return {
                background,
                foreground,
                canvasPixels: count(background),
                inkPixels: count(foreground),
                contrastingPixels,
                area: canvas.width * canvas.height,
            };
        }, Buffer.from(png).toString("base64"));
        assert.ok(
            colors.canvasPixels > colors.area / 2,
            `${name}: diagram canvas matches page ${JSON.stringify(colors)}`,
        );
        assert.ok(
            colors.contrastingPixels > 20,
            `${name}: text and lines contrast with the canvas ${JSON.stringify(colors)}`,
        );
        report.push({ name, mode, ...colors });
    };
    await check("light", "light");
    await select("dark");
    await check("explicit-dark-system-light", "dark");
    const darkBox = await image.boundingBox();
    assert.equal(
        darkBox.width,
        initialBox.width,
        "Theme changes preserve diagram width",
    );
    assert.equal(
        darkBox.height,
        initialBox.height,
        "Theme changes preserve diagram height",
    );
    await page.reload({ waitUntil: "networkidle0" });
    await page.waitForFunction(
        () => document.documentElement.dataset.bsTheme === "dark",
    );
    image = await visibleImage();
    await check("saved-dark", "dark");
    await page.emulateMediaFeatures([
        { name: "prefers-color-scheme", value: "dark" },
    ]);
    await select("light");
    await check("explicit-light-system-dark", "light");
    await select("auto");
    await check("system-dark", "dark");
    await page.emulateMediaFeatures([
        { name: "prefers-color-scheme", value: "light" },
    ]);
    await check("system-light", "light");

    // Each standalone SVG must already contain its final colors and font.
    const diagrams = await page.$$eval(".mermaid-diagram img", (images) =>
        images.map((element) => ({
            src: element.src,
            alt: element.alt,
            mode: element.classList.contains("mermaid-dark")
                ? "dark"
                : "light",
        })),
    );
    assert.ok(
        diagrams.length === 2,
        "Documentation diagram has both native variants",
    );
    const svgPage = await browser.newPage();
    for (const diagram of diagrams) {
        await svgPage.emulateMediaFeatures([
            {
                name: "prefers-color-scheme",
                value: diagram.mode === "dark" ? "light" : "dark",
            },
        ]);
        await svgPage.goto(diagram.src);
        await svgPage.evaluate(() => document.fonts.ready);
        const details = await svgPage.evaluate(() => ({
            geometry: (() => {
                const svg = document.documentElement;
                const parent = document.querySelector('.cluster[id$="-dc1"]');
                const child = document.querySelector(
                    '.cluster[id$="-host_bot"]',
                );
                const outer = parent
                    .querySelector("rect")
                    .getBoundingClientRect();
                const inner = child
                    .querySelector("rect")
                    .getBoundingClientRect();
                const parentTitle = parent
                    .querySelector(".cluster-label")
                    .getBoundingClientRect();
                const childTitle = child
                    .querySelector(".cluster-label")
                    .getBoundingClientRect();
                const scale = svg.getScreenCTM().d;
                return {
                    width: svg.viewBox.baseVal.width,
                    height: svg.viewBox.baseVal.height,
                    nestingInset: (inner.top - outer.top) / scale,
                    titleClearance:
                        (childTitle.top - parentTitle.bottom) / scale,
                    childContained:
                        inner.left >= outer.left &&
                        inner.right <= outer.right &&
                        inner.top >= outer.top &&
                        inner.bottom <= outer.bottom,
                    curvedEdges: [
                        ...document.querySelectorAll(".flowchart-link"),
                    ].filter((edge) => /[CQ]/.test(edge.getAttribute("d")))
                        .length,
                };
            })(),
            font: getComputedStyle(document.documentElement).fontFamily,
            background: getComputedStyle(document.documentElement)
                .backgroundColor,
            labels: [...document.querySelectorAll(".nodeLabel")].map(
                (element) => ({
                    color: getComputedStyle(element).color,
                    text: element.textContent,
                }),
            ),
            lines: [...document.querySelectorAll(".flowchart-link")].map(
                (element) => getComputedStyle(element).stroke,
            ),
            labelBounds: [...document.querySelectorAll(".node")].flatMap(
                (node) => {
                    const label = node.querySelector(".nodeLabel");
                    const shape = node.querySelector("rect");
                    if (!label || !shape || !label.textContent.trim())
                        return [];
                    const text = label.getBoundingClientRect();
                    const box = shape.getBoundingClientRect();
                    return [
                        {
                            text: label.textContent,
                            fits:
                                text.left >= box.left - 1 &&
                                text.right <= box.right + 1 &&
                                text.top >= box.top - 1 &&
                                text.bottom <= box.bottom + 1,
                        },
                    ];
                },
            ),
            captions: [
                ...document.querySelectorAll(".edgeLabel p, .cluster-label p"),
            ]
                .filter((element) => element.textContent.trim())
                .map((element) => {
                    const style = getComputedStyle(element);
                    const bounds = element.getBoundingClientRect();
                    const frame = element
                        .closest("foreignObject")
                        .getBoundingClientRect();
                    return {
                        text: element.textContent,
                        border: parseFloat(style.borderTopWidth),
                        radius: parseFloat(style.borderRadius),
                        background: style.backgroundColor,
                        fits:
                            bounds.left >= frame.left - 1 &&
                            bounds.right <= frame.right + 1 &&
                            bounds.top >= frame.top - 1 &&
                            bounds.bottom <= frame.bottom + 1,
                    };
                }),
            corners: [
                ...document.querySelectorAll(".node > rect, .cluster > rect"),
            ].map((element) => ({
                x: parseFloat(getComputedStyle(element).rx),
                y: parseFloat(getComputedStyle(element).ry),
            })),
        }));
        // The previous 2120 × 790 diagram crowded nested titles. Require a
        // smaller canvas and enough separation for two-line title plates.
        assert.ok(
            details.geometry.width < 1700 && details.geometry.height < 790,
            JSON.stringify(details.geometry),
        );
        assert.ok(
            details.geometry.nestingInset >= 40 &&
                details.geometry.titleClearance >= 8 &&
                details.geometry.childContained,
            JSON.stringify(details.geometry),
        );
        assert.ok(
            details.geometry.curvedEdges > 0,
            "Connectors have smooth bends",
        );
        assert.ok(details.font.includes("Liberation Sans"));
        const palette = report.find((entry) => entry.mode === diagram.mode);
        const ink = `rgb(${palette.foreground.join(", ")})`;
        assert.equal(
            details.background,
            `rgb(${palette.background.join(", ")})`,
        );
        assert.ok(
            details.labelBounds.length > 0 &&
                details.labelBounds.every((label) => label.fits),
            JSON.stringify(details.labelBounds),
        );
        assert.deepEqual(
            details.labels.map((label) => label.color),
            details.labels.map(() => ink),
        );
        assert.ok(details.captions.length > 0);
        assert.ok(
            details.captions.every(
                (caption) =>
                    caption.border > 0 &&
                    caption.radius > 0 &&
                    caption.fits &&
                    caption.background === details.background,
            ),
            JSON.stringify(details.captions),
        );
        assert.ok(details.lines.length > 0);
        assert.deepEqual(
            details.lines,
            details.lines.map(() => ink),
        );
        assert.equal(
            await svgPage.$('[data-look="handDrawn"], .rough-node'),
            null,
            "Published shapes use clean outlines",
        );
        assert.ok(
            details.corners.length > 0 &&
                details.corners.every(
                    (corner) => corner.x > 0 && corner.y > 0,
                ),
            "Nodes and containers have rounded corners",
        );
        report.push({
            diagram: diagram.alt,
            ...details,
            mode: diagram.mode,
        });
        await svgPage.screenshot({
            path: path.join(output, `diagram-${diagram.mode}.png`),
            fullPage: true,
        });
    }
    await svgPage.close();

    await page.setViewport({ width: 390, height: 844 });
    await select("dark");
    await check("mobile-dark", "dark");
    assert.ok(
        await page.evaluate(
            () => document.documentElement.scrollWidth <= innerWidth,
        ),
        "No mobile page overflow",
    );
    await page.goto(`${base}/_print/docs/users/simeonwarren/host_bot/`, {
        waitUntil: "networkidle0",
    });
    image = await visibleImage();
    assert.ok(
        image,
        "Combined print page retains the diagram image and alternative text",
    );
    await page.emulateMediaType("print");
    image = await visibleImage();
    await image.scrollIntoView();
    assert.ok(
        await image.evaluate((element) =>
            element.classList.contains("mermaid-light"),
        ),
        "Print displays the light SVG even when dark is selected",
    );
    report.push(
        await image.evaluate((element) => ({
            printSource: element.src,
            documentURL: location.href,
        })),
    );
    await page.waitForFunction(
        (element) => element.complete && element.naturalWidth > 0,
        {},
        image,
    );
    await image.screenshot({ path: path.join(output, "print.png") });
    report.push({ print: "diagram loaded" });

    await page.emulateMediaType("screen");
    await page.goto(`${base}/docs/misc/old_diagrams/`, {
        waitUntil: "networkidle0",
    });
    const shortcodeImages = await page.$$(".td-content img");
    assert.equal(shortcodeImages.length, 2);
    assert.equal(await page.$(".mermaid-diagram"), null);
    for (const image of shortcodeImages) {
        await image.scrollIntoView();
        await page.waitForFunction(
            (element) => element.complete && element.naturalWidth > 0,
            {},
            image,
        );
        assert.ok(await image.evaluate((element) => element.alt));
    }
    report.push({ shortcodeImagesLoaded: shortcodeImages.length });

    await page.goto(`${base}/blog/diagrams-in-ac-era/`, {
        waitUntil: "networkidle0",
    });
    const blogImages = await page.$$eval(".td-content img", (images) =>
        images.map((element) => ({ src: element.src, alt: element.alt })),
    );
    assert.equal(await page.$(".mermaid-diagram"), null);
    assert.ok(blogImages.length >= 6, "Historical blog images remain present");
    for (const image of blogImages) {
        assert.ok(
            !image.src.includes("/diagrams/"),
            "Blog images keep their original publication URLs",
        );
        const original = await fs.readFile(
            path.join(root, new URL(image.src).pathname),
        );
        assert.deepEqual(
            Buffer.from(await (await fetch(image.src)).arrayBuffer()),
            original,
        );
        assert.ok(image.alt);
    }
    report.push({ blogImagesPreserved: blogImages.length });
} finally {
    await fs.writeFile(
        path.join(output, "report.json"),
        JSON.stringify(report, null, 2),
    );
    await browser.close();
    await new Promise((resolve) => server.close(resolve));
}
