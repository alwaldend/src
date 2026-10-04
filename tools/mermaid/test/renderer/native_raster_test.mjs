import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import puppeteer from "puppeteer";

// Exercise the two Bazel render actions independently, then compare the native
// light and dark rasters with their SVGs as displayed by the browser. The
// reference does not pass through the raster encoder, so losing SVG content
// during encoding fails.
const [
    lightSvg,
    lightRaster,
    darkSvg,
    darkRaster,
    palettePath,
    executablePath,
] = process.argv.slice(2).map((filename) => path.resolve(filename));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const puppeteerConfig = { executablePath, headless: "shell" };
const reports = {};

const browser = await puppeteer.launch(puppeteerConfig);
try {
    for (const [colorScheme, svgPath, rasterPath] of [
        ["light", lightSvg, lightRaster],
        ["dark", darkSvg, darkRaster],
    ]) {
        const referencePath = path.join(
            output,
            `native-${colorScheme}-reference.webp`,
        );
        const raster = await fs.readFile(rasterPath);
        await fs.copyFile(
            rasterPath,
            path.join(output, `native-${colorScheme}-actual.webp`),
        );
        const page = await browser.newPage();
        await page.emulateMediaFeatures([
            { name: "prefers-color-scheme", value: colorScheme },
        ]);
        await page.addStyleTag({
            content: await fs.readFile(palettePath, "utf8"),
        });
        const palette = await page.evaluate(() => {
            const style = getComputedStyle(document.documentElement);
            const color = (name) => {
                const canvas = document.createElement("canvas");
                const context = canvas.getContext("2d");
                context.fillStyle = style.getPropertyValue(
                    `--mermaid-${name}`,
                );
                context.fillRect(0, 0, 1, 1);
                return [...context.getImageData(0, 0, 1, 1).data].slice(0, 3);
            };
            return {
                background: color("background"),
                ink: color("textColor"),
            };
        });
        const svg = await fs.readFile(svgPath, "utf8");
        assert.ok(svg.includes("font-family:'Liberation Sans'"));
        assert.ok(svg.includes("src:url(data:font/ttf;base64,"));
        assert.ok(!svg.includes("Architects Daughter"));
        assert.ok(!svg.includes('data-look="handDrawn"'));
        await page.setContent(svg);
        const reference = await page.evaluate(() => {
            const svg = document.querySelector("svg");
            const label = document.querySelector(".nodeLabel");
            const width = Math.ceil(svg.viewBox.baseVal.width);
            const height = Math.ceil(svg.viewBox.baseVal.height);
            document.body.style.margin = "0";
            svg.setAttribute("width", width);
            svg.setAttribute("height", height);
            svg.style.width = `${width}px`;
            svg.style.height = `${height}px`;
            svg.style.maxWidth = "none";
            return {
                width: svg.viewBox.baseVal.width,
                height: svg.viewBox.baseVal.height,
                font: getComputedStyle(label).fontFamily,
                color: getComputedStyle(label).color,
            };
        });
        assert.ok(reference.font.includes("Liberation Sans"));
        assert.equal(reference.color, `rgb(${palette.ink.join(", ")})`);
        await page.setViewport({
            width: Math.ceil(reference.width),
            height: Math.ceil(reference.height),
            deviceScaleFactor: 2,
        });
        await page.evaluate(() => document.fonts.ready);
        const containerTitles = await page.evaluate(() =>
            [...document.querySelectorAll(".cluster")].map((container) => {
                const border = container
                    .querySelector("rect")
                    .getBoundingClientRect();
                const plate = container
                    .querySelector(".cluster-label p")
                    .getBoundingClientRect();
                return {
                    title: container.querySelector(".cluster-label p")
                        .textContent,
                    centerOffset: plate.y + plate.height / 2 - border.y,
                    containedHorizontally:
                        plate.x >= border.x && plate.right <= border.right,
                };
            }),
        );
        assert.ok(containerTitles.length >= 2);
        for (const title of containerTitles) {
            assert.ok(
                Math.abs(title.centerOffset) <= 1,
                `${title.title} caption must straddle its top border; offset ${title.centerOffset}`,
            );
            assert.ok(title.containedHorizontally);
        }
        const firstLabel = await page.$eval(".edgeLabel p", (element) => {
            const bounds = element.getBoundingClientRect();
            return {
                text: element.textContent,
                x: Math.floor(bounds.x * 2),
                y: Math.floor(bounds.y * 2),
                width: Math.ceil(bounds.width * 2),
                height: Math.ceil(bounds.height * 2),
            };
        });
        assert.equal(firstLabel.text, "HTTPS");
        assert.ok(firstLabel.width > 0 && firstLabel.height > 0);
        const referenceRaster = await page.screenshot({
            path: referencePath,
            type: "webp",
            quality: 90,
        });
        const actual = await page.evaluate(
            async (encoded, expected, label, ink) => {
                const decode = async (data) => {
                    const image = new Image();
                    image.src = `data:image/webp;base64,${data}`;
                    await image.decode();
                    const canvas = document.createElement("canvas");
                    canvas.width = image.naturalWidth;
                    canvas.height = image.naturalHeight;
                    const context = canvas.getContext("2d");
                    context.drawImage(image, 0, 0);
                    return { image, context };
                };
                const rendered = await decode(encoded);
                const reference = await decode(expected);
                const pixels = ({ context }) =>
                    context.getImageData(
                        label.x,
                        label.y,
                        label.width,
                        label.height,
                    ).data;
                const actualPixels = pixels(rendered);
                const referencePixels = pixels(reference);
                const labelDifference =
                    actualPixels.reduce(
                        (difference, channel, index) =>
                            difference +
                            Math.abs(channel - referencePixels[index]),
                        0,
                    ) / actualPixels.length;
                let inkPixels = 0;
                for (let index = 0; index < actualPixels.length; index += 4) {
                    if (
                        ink.every(
                            (channel, offset) =>
                                Math.abs(
                                    actualPixels[index + offset] - channel,
                                ) <= 25,
                        )
                    )
                        inkPixels += 1;
                }
                return {
                    inkPixels,
                    width: rendered.image.naturalWidth,
                    height: rendered.image.naturalHeight,
                    background: [
                        ...rendered.context.getImageData(0, 0, 1, 1).data,
                    ].slice(0, 3),
                    firstLabelMeanDifference: labelDifference,
                };
            },
            raster.toString("base64"),
            Buffer.from(referenceRaster).toString("base64"),
            firstLabel,
            palette.ink,
        );
        reports[colorScheme] = {
            reference,
            actual,
            palette,
            firstLabel,
            containerTitles,
            scale: 2,
        };
        await fs.writeFile(
            path.join(output, "native-raster-report.json"),
            JSON.stringify(reports, null, 2),
        );
        assert.ok(
            actual.inkPixels > 20,
            "Edge label must contain visible palette ink",
        );
        assert.ok(
            actual.firstLabelMeanDifference <= 1,
            `First edge label must remain visible; pixel difference ${actual.firstLabelMeanDifference}`,
        );
        assert.equal(actual.width, Math.ceil(reference.width) * 2);
        assert.equal(actual.height, Math.ceil(reference.height) * 2);
        actual.background.forEach((channel, index) => {
            assert.ok(
                Math.abs(channel - palette.background[index]) <= 3,
                `Raster background must use the ${colorScheme} palette, allowing WebP compression`,
            );
        });
        await page.close();
    }
    assert.deepEqual(reports.dark.firstLabel, reports.light.firstLabel);
    assert.equal(reports.dark.actual.width, reports.light.actual.width);
    assert.equal(reports.dark.actual.height, reports.light.actual.height);
    assert.notDeepEqual(reports.dark.palette, reports.light.palette);
    const intensity = (color) =>
        color.reduce((sum, channel) => sum + channel, 0);
    assert.ok(
        intensity(reports.dark.palette.background) <
            intensity(reports.dark.palette.ink),
    );
    assert.ok(
        intensity(reports.light.palette.background) >
            intensity(reports.light.palette.ink),
    );
    console.log(
        "Native light/dark raster palettes, embedded font, labels, and scaled canvas verified",
    );
} finally {
    await browser.close();
}
