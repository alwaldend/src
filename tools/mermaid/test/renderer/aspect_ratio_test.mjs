import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import puppeteer from "puppeteer";
import { renderWebp } from "../../cmd/render/raster.mjs";

const [browserPath, naturalDarkPath, ...fixtures] = process.argv.slice(2);
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const puppeteerConfig = {
    executablePath: path.resolve(browserPath),
    headless: "shell",
};
const reports = {};
const naturalDark = await fs.readFile(path.resolve(naturalDarkPath));

// Reject malformed inputs before launching Chrome or writing an output. These
// exercise the same encoder used by the complete Bazel rendering actions below.
const invalidSvg = path.join(output, "invalid.svg");
const invalidWebp = path.join(output, "invalid.webp");
const options = { svgPath: invalidSvg, output: invalidWebp, puppeteerConfig };
await fs.writeFile(invalidSvg, '<svg viewBox="0 0 100 50"/>');
for (const aspectRatio of [
    [],
    [5],
    [5, 2, 1],
    [0, 2],
    [-5, 2],
    [5, 2.5],
    [5, Infinity],
    [5, NaN],
    [5, Number.MAX_SAFE_INTEGER + 1],
    ["5", 2],
]) {
    await assert.rejects(
        renderWebp({ ...options, aspectRatio }),
        /aspect ratio.*two positive.*integers/i,
    );
}
for (const deviceScaleFactor of [0, -1, 1.5, NaN, Infinity]) {
    await assert.rejects(
        renderWebp({ ...options, deviceScaleFactor }),
        /scale.*positive.*integer/i,
    );
}
await assert.rejects(
    renderWebp({ ...options, aspectRatio: [16384, 1] }),
    /16383/,
);
await assert.rejects(
    renderWebp({ ...options, deviceScaleFactor: Number.MAX_SAFE_INTEGER }),
    /16383/,
);
for (const viewBox of ["0 0 20000 1", "0 0 1 20000"]) {
    await fs.writeFile(invalidSvg, `<svg viewBox="${viewBox}"/>`);
    await assert.rejects(renderWebp(options), /16383/);
}
await fs.writeFile(invalidSvg, '<svg viewBox="0 0 0 10"/>');
await assert.rejects(renderWebp(options), /usable viewBox/);
await assert.rejects(fs.stat(invalidWebp), { code: "ENOENT" });
await fs.rm(invalidSvg);

const browser = await puppeteer.launch(puppeteerConfig);
try {
    for (let index = 0; index < fixtures.length; index += 6) {
        const [name, svgFile, rasterFile, width, height, scaleValue] =
            fixtures.slice(index, index + 6);
        const ratio = [Number(width), Number(height)];
        const scale = Number(scaleValue);
        const svg = await fs.readFile(path.resolve(svgFile), "utf8");
        const raster = await fs.readFile(path.resolve(rasterFile));
        await fs.copyFile(
            path.resolve(rasterFile),
            path.join(output, `${name}-actual.webp`),
        );
        const page = await browser.newPage();
        await page.setContent(svg);
        const dimensions = await page.evaluate(() => {
            const svg = document.querySelector("svg");
            return {
                width: Math.ceil(svg.viewBox.baseVal.width),
                height: Math.ceil(svg.viewBox.baseVal.height),
                background: getComputedStyle(svg).backgroundColor,
            };
        });
        const contentWidth = dimensions.width * scale;
        const contentHeight = dimensions.height * scale;
        // Independent integer search proves both exact ratio and minimal size;
        // the reference does not reuse the encoder's reduction or sizing code.
        let canvasWidth = contentWidth;
        while (
            (canvasWidth * ratio[1]) % ratio[0] !== 0 ||
            (canvasWidth * ratio[1]) / ratio[0] < contentHeight
        )
            canvasWidth += 1;
        const canvasHeight = (canvasWidth * ratio[1]) / ratio[0];
        await page.setViewport({
            width: canvasWidth,
            height: canvasHeight,
            deviceScaleFactor: 1,
        });
        const reference = await page.evaluate(
            ({
                contentWidth,
                contentHeight,
                canvasWidth,
                canvasHeight,
                background,
            }) => {
                document.body.style.cssText = `margin:0;display:grid;place-items:center;width:${canvasWidth}px;height:${canvasHeight}px;background:${background}`;
                const svg = document.querySelector("svg");
                svg.style.width = `${contentWidth}px`;
                svg.style.height = `${contentHeight}px`;
                svg.style.maxWidth = "none";
                const bounds = svg.getBoundingClientRect();
                return {
                    x: bounds.x,
                    y: bounds.y,
                    width: bounds.width,
                    height: bounds.height,
                };
            },
            {
                contentWidth,
                contentHeight,
                canvasWidth,
                canvasHeight,
                background: dimensions.background,
            },
        );
        await page.evaluate(() => document.fonts.ready);
        const referenceRaster = await page.screenshot({
            path: path.join(output, `${name}-reference.webp`),
            type: "webp",
            quality: 90,
        });
        const actual = await page.evaluate(
            async (encoded, expected, background, bounds, naturalEncoded) => {
                const decode = async (data) => {
                    const image = new Image();
                    image.src = `data:image/webp;base64,${data}`;
                    await image.decode();
                    const canvas = document.createElement("canvas");
                    canvas.width = image.naturalWidth;
                    canvas.height = image.naturalHeight;
                    const context = canvas.getContext("2d");
                    context.drawImage(image, 0, 0);
                    return {
                        width: canvas.width,
                        height: canvas.height,
                        context,
                    };
                };
                const actual = await decode(encoded);
                const expectedImage = await decode(expected);
                const samples = actual.context.getImageData(
                    0,
                    0,
                    actual.width,
                    actual.height,
                ).data;
                const referenceSamples = expectedImage.context.getImageData(
                    0,
                    0,
                    actual.width,
                    actual.height,
                ).data;
                const canvas = document.createElement("canvas");
                const context = canvas.getContext("2d");
                context.fillStyle = background;
                context.fillRect(0, 0, 1, 1);
                const expectedBackground = [
                    ...context.getImageData(0, 0, 1, 1).data,
                ];
                let difference = 0;
                let contentPixels = 0;
                for (let pixel = 0; pixel < samples.length; pixel += 4) {
                    let contrast = 0;
                    for (let channel = 0; channel < 3; channel += 1) {
                        difference += Math.abs(
                            samples[pixel + channel] -
                                referenceSamples[pixel + channel],
                        );
                        contrast += Math.abs(
                            samples[pixel + channel] -
                                expectedBackground[channel],
                        );
                    }
                    if (contrast > 100) contentPixels += 1;
                }
                const contentDifference = (first, second) => {
                    let difference = 0;
                    for (let index = 0; index < first.length; index += 1)
                        difference += Math.abs(first[index] - second[index]);
                    return difference / first.length;
                };
                const content = actual.context.getImageData(
                    Math.floor(bounds.x),
                    Math.floor(bounds.y),
                    bounds.width,
                    bounds.height,
                ).data;
                const expectedContent = expectedImage.context.getImageData(
                    Math.floor(bounds.x),
                    Math.floor(bounds.y),
                    bounds.width,
                    bounds.height,
                ).data;
                let naturalMeanDifference;
                if (naturalEncoded !== null) {
                    // Existing default action uses the original natural CSS
                    // dimensions at DPR 2. Compare just the content so padding
                    // cannot dilute a change in scale, strokes, or label ink.
                    const natural = await decode(naturalEncoded);
                    naturalMeanDifference = contentDifference(
                        content,
                        natural.context.getImageData(
                            0,
                            0,
                            natural.width,
                            natural.height,
                        ).data,
                    );
                }
                return {
                    width: actual.width,
                    height: actual.height,
                    meanDifference: difference / (samples.length * 0.75),
                    contentMeanDifference: contentDifference(
                        content,
                        expectedContent,
                    ),
                    naturalMeanDifference,
                    contentPixels,
                    background: [
                        ...actual.context.getImageData(0, 0, 1, 1).data,
                    ],
                    expectedBackground,
                };
            },
            raster.toString("base64"),
            Buffer.from(referenceRaster).toString("base64"),
            dimensions.background,
            reference,
            name === "dark" ? naturalDark.toString("base64") : null,
        );
        reports[name] = { ratio, scale, dimensions, reference, actual };
        await fs.writeFile(
            path.join(output, "aspect-ratio-report.json"),
            JSON.stringify(reports, null, 2),
        );
        assert.equal(actual.width, canvasWidth);
        assert.equal(actual.height, canvasHeight);
        assert.equal(actual.width * ratio[1], actual.height * ratio[0]);
        assert.equal(reference.width, contentWidth);
        assert.equal(reference.height, contentHeight);
        assert.equal(reference.x, (canvasWidth - contentWidth) / 2);
        assert.equal(reference.y, (canvasHeight - contentHeight) / 2);
        assert.ok(
            actual.meanDifference <= 1,
            `${name}: content changed or moved (difference ${actual.meanDifference})`,
        );
        assert.ok(
            actual.contentPixels > 200,
            `${name}: visible content missing`,
        );
        assert.ok(
            actual.contentMeanDifference <= 1,
            `${name}: content region changed`,
        );
        if (name === "dark")
            assert.ok(
                actual.naturalMeanDifference <= 2,
                `Original DPR2 content changed: ${actual.naturalMeanDifference}`,
            );
        actual.background.forEach((channel, offset) => {
            assert.ok(
                Math.abs(channel - actual.expectedBackground[offset]) <= 3,
                `${name}: padding changed theme background`,
            );
        });
        await page.close();
    }
    assert.ok(
        reports.classic.reference.y > 0,
        "Wide diagram needs vertical padding",
    );
    assert.ok(
        reports.dark.reference.x > 0,
        "Tall diagram needs horizontal padding",
    );
    console.log(
        "Aspect ratio, minimal canvas, centered original content, and theme padding verified",
    );
} finally {
    await browser.close();
}
