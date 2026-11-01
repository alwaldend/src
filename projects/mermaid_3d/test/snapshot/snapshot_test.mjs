import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { createHash } from "node:crypto";
import fs from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { promisify } from "node:util";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const execute = promisify(execFile);
const [browserPath, command, flowchart, sequence, ...images] = process.argv
    .slice(2)
    .map((file) => path.resolve(file));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "The test must preserve its rendered evidence");
await fs.mkdir(output, { recursive: true });
const cases = [
    [flowchart, "light", "isometric", "png", 960, 720],
    [flowchart, "dark", "isometric", "png", 960, 720],
    [flowchart, "light", "front", "png", 960, 720],
    [flowchart, "light", "front", "webp", 960, 720],
    [sequence, "dark", "front", "png", 480, 960],
    [sequence, "light", "isometric", "webp", 1280, 640],
];
assert.equal(images.length, cases.length);
const hash = (bytes) => createHash("sha256").update(bytes).digest("hex");
const report = { images: [], failures: [], status: "running" };
const browser = await puppeteer.launch({
    executablePath: browserPath,
    headless: "shell",
    args: ["--use-angle=swiftshader", "--enable-unsafe-swiftshader"],
});

// The output itself is decoded by the pinned browser. Compare its pixels to
// a fresh live canvas; also test meaningful image properties independently.
async function inspectRaster(page, bytes, extension, view, referenceBytes) {
    return page.evaluate(
        async ({ encoded, extension, view, referenceEncoded }) => {
            const image = new Image();
            image.src = `data:image/${extension};base64,${encoded}`;
            await image.decode();
            const buffer = document.createElement("canvas");
            buffer.width = image.naturalWidth;
            buffer.height = image.naturalHeight;
            const ctx = buffer.getContext("2d", { willReadFrequently: true });
            ctx.drawImage(image, 0, 0);
            const actual = ctx.getImageData(
                0,
                0,
                buffer.width,
                buffer.height,
            ).data;
            ctx.clearRect(0, 0, buffer.width, buffer.height);
            let referenceImage = document.querySelector("#scene canvas");
            if (referenceEncoded) {
                referenceImage = new Image();
                referenceImage.src = `data:image/png;base64,${referenceEncoded}`;
                await referenceImage.decode();
            }
            ctx.drawImage(referenceImage, 0, 0);
            const reference = ctx.getImageData(
                0,
                0,
                buffer.width,
                buffer.height,
            ).data;
            const background = view.background
                .slice(1)
                .match(/../g)
                .map((n) => parseInt(n, 16));
            let foreground = 0;
            let borderForeground = 0;
            let difference = 0;
            let maximumDifference = 0;
            let opaque = true;
            const bounds = {
                left: buffer.width,
                top: buffer.height,
                right: -1,
                bottom: -1,
            };
            for (let y = 0; y < buffer.height; y += 1) {
                for (let x = 0; x < buffer.width; x += 1) {
                    const i = (y * buffer.width + x) * 4;
                    const visible = background.some(
                        (channel, j) => Math.abs(actual[i + j] - channel) > 12,
                    );
                    if (visible) {
                        foreground += 1;
                        bounds.left = Math.min(bounds.left, x);
                        bounds.right = Math.max(bounds.right, x);
                        bounds.top = Math.min(bounds.top, y);
                        bounds.bottom = Math.max(bounds.bottom, y);
                        if (
                            x < 4 ||
                            y < 4 ||
                            x >= buffer.width - 4 ||
                            y >= buffer.height - 4
                        )
                            borderForeground += 1;
                    }
                    opaque &&= actual[i + 3] === 255;
                    for (let j = 0; j < 4; j += 1) {
                        const delta = Math.abs(
                            actual[i + j] - reference[i + j],
                        );
                        difference += delta;
                        maximumDifference = Math.max(maximumDifference, delta);
                    }
                }
            }
            const cornerColors = [
                0,
                buffer.width - 1,
                buffer.width * (buffer.height - 1),
                buffer.width * buffer.height - 1,
            ].map((pixel) =>
                Array.from(actual.slice(pixel * 4, pixel * 4 + 3)),
            );
            const labelInk = view.labels.map(({ text, bounds }) => {
                const colors = new Set();
                for (
                    let y = Math.max(0, Math.ceil(bounds.top));
                    y < Math.min(buffer.height, Math.floor(bounds.bottom));
                    y += 1
                )
                    for (
                        let x = Math.max(0, Math.ceil(bounds.left));
                        x < Math.min(buffer.width, Math.floor(bounds.right));
                        x += 1
                    ) {
                        const i = (y * buffer.width + x) * 4;
                        colors.add(actual[i] >> 4);
                    }
                return { text, shades: colors.size };
            });
            return {
                width: buffer.width,
                height: buffer.height,
                foreground,
                borderForeground,
                meanDifference: difference / actual.length,
                maximumDifference,
                opaque,
                bounds,
                background,
                cornerColors,
                labelInk,
            };
        },
        {
            encoded: bytes.toString("base64"),
            extension,
            view,
            referenceEncoded: referenceBytes?.toString("base64"),
        },
    );
}

async function exportImage(input, destination, extra = []) {
    return execute(
        command,
        [
            "--browser",
            browserPath,
            "--input",
            input,
            "--output",
            destination,
            "--width",
            "960",
            "--height",
            "720",
            "--color-scheme",
            "light",
            "--camera",
            "isometric",
            ...extra,
        ],
        { timeout: 60000, maxBuffer: 1024 * 1024 },
    );
}

async function reject(name, input, extra, expected, extension = "png") {
    const destination = path.join(output, `${name}.${extension}`);
    await assert.rejects(exportImage(input, destination, extra), (error) => {
        assert.match(error.stderr, expected, name);
        report.failures.push({
            name,
            exitCode: error.code,
            diagnostic: error.stderr.trim(),
        });
        return true;
    });
    await assert.rejects(fs.stat(destination), { code: "ENOENT" });
}

try {
    for (const [
        index,
        [input, colorScheme, camera, extension, width, height],
    ] of cases.entries()) {
        const bytes = await fs.readFile(images[index]);
        const name = path.basename(images[index]);
        if (extension === "png") {
            assert.deepEqual(
                bytes.subarray(0, 8),
                Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]),
            );
            assert.equal(
                bytes.readUInt32BE(16),
                width,
                `${name}: encoded width`,
            );
            assert.equal(
                bytes.readUInt32BE(20),
                height,
                `${name}: encoded height`,
            );
        } else {
            assert.equal(bytes.toString("ascii", 0, 4), "RIFF");
            assert.equal(bytes.toString("ascii", 8, 12), "WEBP");
            assert.equal(bytes.readUInt32LE(4) + 8, bytes.length);
        }
        const page = await browser.newPage();
        const network = [];
        const errors = [];
        await page.setViewport({ width, height, deviceScaleFactor: 1 });
        // Deliberately start in the opposite appearance: the requested option
        // must win over both the browser preference and the initial viewer.
        await page.emulateMediaFeatures([
            {
                name: "prefers-color-scheme",
                value: colorScheme === "light" ? "dark" : "light",
            },
        ]);
        await page.evaluateOnNewDocument(() => {
            window.snapshotFonts = [];
            const fillText = CanvasRenderingContext2D.prototype.fillText;
            CanvasRenderingContext2D.prototype.fillText = function (...args) {
                window.snapshotFonts.push(this.font);
                return fillText.apply(this, args);
            };
        });
        page.on("pageerror", (error) => errors.push(error.message));
        await page.setRequestInterception(true);
        const url = pathToFileURL(input).href;
        page.on("request", (request) => {
            if (request.url() === url || /^(data|blob):/.test(request.url()))
                void request.continue();
            else {
                network.push(request.url());
                void request.abort();
            }
        });
        await page.goto(url);
        await page.waitForFunction(() => window.diagramView?.ready, {
            timeout: 30000,
        });
        const view = await page.evaluate(
            async (options) => {
                await document.fonts.ready;
                return window.diagramView.prepareSnapshot(options);
            },
            { width, height, colorScheme, camera },
        );
        assert.equal(
            view.appearance,
            colorScheme,
            `${name}: explicit appearance`,
        );
        assert.deepEqual(view.canvas, { width, height });
        assert.equal(view.selected, "");
        assert.ok(view.fontLoaded);
        const fonts = await page.evaluate(() => ({
            faces: [...document.fonts].map(({ family, status }) => ({
                family,
                status,
            })),
            drawing: window.snapshotFonts,
        }));
        assert.ok(
            fonts.faces.length > 0 &&
                fonts.faces.every((face) => face.status === "loaded"),
        );
        assert.ok(
            fonts.drawing.length > 0 &&
                fonts.drawing.every((font) =>
                    font.includes("Liberation Mono"),
                ),
        );
        assert.ok(view.labels.length > 0 && view.sceneCorners.length > 0);
        for (const label of view.labels) {
            assert.ok(
                label.bounds.left >= 2 &&
                    label.bounds.top >= 2 &&
                    label.bounds.right <= width - 2 &&
                    label.bounds.bottom <= height - 2,
                `${name}: clipped label ${label.text}`,
            );
        }
        for (const [x, y] of view.sceneCorners)
            assert.ok(
                x >= 2 && x <= width - 2 && y >= 2 && y <= height - 2,
                `${name}: clipped scene geometry`,
            );
        const direction = view.camera.map(
            (value, i) => value - view.target[i],
        );
        const desired =
            camera === "front" ? [0, 0, 1] : view.model.viewDirection;
        const cosine =
            direction.reduce((sum, value, i) => sum + value * desired[i], 0) /
            Math.hypot(...direction) /
            Math.hypot(...desired);
        assert.ok(cosine > 1 - 1e-8, `${name}: requested camera direction`);
        const raster = await inspectRaster(page, bytes, extension, view);
        assert.equal(raster.width, width);
        assert.equal(raster.height, height);
        assert.equal(raster.opaque, true);
        assert.ok(
            raster.foreground > width * height * 0.005,
            `${name}: nonblank raster`,
        );
        assert.equal(
            raster.borderForeground,
            0,
            `${name}: raster clipping or UI at image boundary`,
        );
        for (const corner of raster.cornerColors)
            for (let i = 0; i < 3; i += 1)
                assert.ok(
                    Math.abs(corner[i] - raster.background[i]) <= 3,
                    `${name}: appearance in encoded pixels`,
                );
        for (const label of raster.labelInk)
            assert.ok(
                label.shades > 1,
                `${name}: visible label pixels for ${label.text}`,
            );
        if (extension === "png")
            assert.equal(
                raster.maximumDifference,
                0,
                `${name}: lossless canvas pixels without UI`,
            );
        else
            assert.ok(
                raster.meanDifference < 2,
                `${name}: WebP scene fidelity`,
            );
        assert.deepEqual(network, []);
        assert.deepEqual(errors, []);
        await fs.copyFile(images[index], path.join(output, name));
        await page.screenshot({
            path: path.join(output, name + ".viewer.png"),
        });
        report.images.push({
            name,
            sha256: hash(bytes),
            inputSha256: hash(await fs.readFile(input)),
            colorScheme,
            camera,
            view,
            raster,
        });
        await page.close();
    }
    assert.notEqual(
        report.images[0].sha256,
        report.images[1].sha256,
        "Appearance must change raster bytes",
    );
    assert.notEqual(
        report.images[0].sha256,
        report.images[2].sha256,
        "Front camera must change raster bytes",
    );

    const source = await fs.readFile(flowchart, "utf8");
    const altered = async (name, html) => {
        const input = path.join(output, name + ".html");
        await fs.writeFile(input, html);
        return input;
    };
    // A full-page or element screenshot would capture this opaque overlay.
    // Encoding the actual scene canvas preserves precisely the original image.
    const controls = await altered(
        "controls-overlay",
        source.replace(
            "</body>",
            '<button style="position:fixed;inset:0;background:#ff00ff;z-index:99999">Viewer controls</button></body>',
        ),
    );
    const controlsImage = path.join(output, "controls-overlay.png");
    await exportImage(controls, controlsImage);
    assert.equal(
        hash(await fs.readFile(controlsImage)),
        report.images[0].sha256,
        "DOM controls must not appear in canvas output",
    );

    // Orthographic appearance must not change when only camera distance grows
    // beyond the old fixed 20,000-unit far plane. Preserve orientation, scene
    // geometry, labels, and all embedded inputs in this full-pipeline fixture.
    const payloadPattern =
        /(<script type="application\/json" id="diagram-data">)([\s\S]*?)(<\/script>)/;
    const payloadMatch = source.match(payloadPattern);
    assert.ok(
        payloadMatch,
        "The generated view must expose its model payload",
    );
    const distantPayload = JSON.parse(payloadMatch[2]);
    distantPayload.model.viewDirection =
        distantPayload.model.viewDirection.map((value) => value * 100);
    const distantHTML = source.replace(
        payloadPattern,
        () =>
            payloadMatch[1] +
            JSON.stringify(distantPayload).replaceAll("<", "\\u003c") +
            payloadMatch[3],
    );
    const distantInput = await altered("distant-camera", distantHTML);
    const distantImage = path.join(output, "distant-camera.png");
    await exportImage(distantInput, distantImage);
    const distantBytes = await fs.readFile(distantImage);
    const comparisonPage = await browser.newPage();
    const distantRaster = await inspectRaster(
        comparisonPage,
        distantBytes,
        "png",
        report.images[0].view,
        await fs.readFile(images[0]),
    );
    await comparisonPage.close();
    assert.equal(distantRaster.width, 960);
    assert.equal(distantRaster.height, 720);
    assert.equal(distantRaster.opaque, true);
    assert.equal(
        distantRaster.borderForeground,
        0,
        "Distant scene keeps full raster margins",
    );
    assert.ok(
        distantRaster.foreground >= report.images[0].raster.foreground * 0.98,
        "Distant camera retains visible scene content",
    );
    assert.ok(
        distantRaster.meanDifference < 1,
        "Camera distance preserves orthographic raster appearance",
    );
    for (const name of ["left", "right", "top", "bottom"])
        assert.ok(
            Math.abs(
                distantRaster.bounds[name] -
                    report.images[0].raster.bounds[name],
            ) <= 2,
            `Distant camera retains ${name} content extent`,
        );
    report.images.push({
        name: "distant-camera.png",
        sha256: hash(distantBytes),
        inputSha256: hash(await fs.readFile(distantInput)),
        cameraDistanceMultiplier: 100,
        raster: distantRaster,
    });

    await reject(
        "unsupported-format",
        flowchart,
        [],
        /output.*(?:png|webp)/i,
        "jpg",
    );
    await reject("zero-height", flowchart, ["--height", "0"], /height/i);
    await reject("oversize-width", flowchart, ["--width", "8193"], /width/i);
    await reject(
        "fractional-width",
        flowchart,
        ["--width", "640.5"],
        /width/i,
    );
    await reject("unknown-camera", flowchart, ["--camera", "rear"], /camera/i);
    await reject(
        "unknown-appearance",
        flowchart,
        ["--color-scheme", "auto"],
        /color.scheme/i,
    );
    const noFont = source.replace(/@font-face\s*\{[^}]*\}/, "");
    assert.notEqual(
        noFont,
        source,
        "Failure fixture must remove the embedded font",
    );
    await reject(
        "missing-font",
        await altered("missing-font", noFont),
        [],
        /font/i,
    );
    const remoteFont = source.replace(
        /data:font\/ttf;base64,[A-Za-z0-9+/=]+/,
        "https://snapshot.invalid/font.ttf",
    );
    assert.notEqual(remoteFont, source);
    await reject(
        "remote-font",
        await altered("remote-font", remoteFont),
        [],
        /(?:offline|network|request)/i,
    );
    await reject(
        "remote-image",
        await altered(
            "remote-image",
            source.replace(
                "</body>",
                '<img src="https://snapshot.invalid/image.png"></body>',
            ),
        ),
        [],
        /(?:offline|network|request)/i,
    );
    await reject(
        "local-dependency",
        await altered(
            "local-dependency",
            source.replace(
                "</body>",
                '<img src="file:///snapshot-external-image.png"></body>',
            ),
        ),
        [],
        /(?:offline|self.contained|request)/i,
    );
    const neverReady = source.replace('import "diagram/viewer";', "");
    assert.notEqual(neverReady, source);
    await reject(
        "readiness-timeout",
        await altered("readiness-timeout", neverReady),
        ["--timeout-ms", "1000"],
        /(?:ready|readiness).*tim(?:e|ed).*out/i,
    );
    report.status = "passed";
} finally {
    await fs.writeFile(
        path.join(output, "verification.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
    await browser.close();
}
