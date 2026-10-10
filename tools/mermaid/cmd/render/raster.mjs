import fs from "node:fs/promises";
import puppeteer from "puppeteer";

// Rasterize a rendered SVG through the pinned browser.
//
// The SVG is the maintained artifact: it already carries its selected theme and
// pinned fonts. Encoding it to WebP is therefore a projection of that document
// rather than a second appearance, and Chrome
// performs the encode, so the rule needs no image-conversion toolchain.
//
// The raster is authored at the diagram's natural size, then scaled by
// `deviceScaleFactor` when a consumer requests higher pixel density.

// readDimensions returns the drawing's own size. Mermaid writes the size into
// `viewBox`; taking it there ignores any `width`/`height` the document also
// carries, so the raster never inherits a consumer's layout width.
const readViewBox = (svg) => {
    const match = svg.match(/\bviewBox="([^"]+)"/);
    if (match === null) return undefined;
    const [minX, minY, width, height] = match[1]
        .trim()
        .split(/\s+/)
        .map(Number);
    if (
        ![minX, minY, width, height].every(Number.isFinite) ||
        width <= 0 ||
        height <= 0
    )
        return undefined;
    return { minX, minY, width, height };
};

// readBackground returns the canvas color Mermaid painted into the document.
const readBackground = (svg) => {
    const match = svg.match(/background-color:\s*([^;"]+)/);
    return match === null ? "white" : match[1].trim();
};

/**
 * Render an SVG file to a WebP image through the pinned browser.
 *
 * @param {object} options
 * @param {string} options.svgPath Rendered SVG on disk, with its theme and font.
 * @param {string} options.output Destination `.webp` path.
 * @param {object} options.puppeteerConfig Launch options, including pinned fonts.
 * @param {number} options.deviceScaleFactor Raster scale; 1 uses natural size.
 * @param {number[]} [options.aspectRatio] Optional positive integer width/height.
 */
export async function renderWebp({
    svgPath,
    output,
    puppeteerConfig,
    deviceScaleFactor = 1,
    aspectRatio,
}) {
    if (!Number.isSafeInteger(deviceScaleFactor) || deviceScaleFactor <= 0)
        throw new Error("Raster scale must be a positive safe integer");
    if (
        aspectRatio !== undefined &&
        (!Array.isArray(aspectRatio) ||
            aspectRatio.length !== 2 ||
            !aspectRatio.every(
                (term) => Number.isSafeInteger(term) && term > 0,
            ))
    )
        throw new Error(
            "Raster aspect ratio must contain two positive safe integers",
        );
    const svg = await fs.readFile(svgPath, "utf8");
    const viewBox = readViewBox(svg);
    if (viewBox === undefined)
        throw new Error(`rendered SVG ${svgPath} carries no usable viewBox`);
    const background = readBackground(svg);

    // Author the document at the diagram's natural size with no margin, so the
    // screenshot's clip box is exactly the drawing.
    const width = Math.ceil(viewBox.width);
    const height = Math.ceil(viewBox.height);
    let canvasWidth = width * deviceScaleFactor;
    let canvasHeight = height * deviceScaleFactor;
    if (aspectRatio !== undefined) {
        let [numerator, denominator] = aspectRatio;
        while (denominator !== 0)
            [numerator, denominator] = [denominator, numerator % denominator];
        const ratioWidth = aspectRatio[0] / numerator;
        const ratioHeight = aspectRatio[1] / numerator;
        const multiplier = Math.ceil(
            Math.max(canvasWidth / ratioWidth, canvasHeight / ratioHeight),
        );
        canvasWidth = ratioWidth * multiplier;
        canvasHeight = ratioHeight * multiplier;
    }
    // WebP limits each dimension to 16383 pixels. Fail before starting Chrome
    // rather than returning an empty or differently sized screenshot.
    if (
        ![canvasWidth, canvasHeight].every(
            (dimension) =>
                Number.isSafeInteger(dimension) && dimension <= 16383,
        )
    )
        throw new Error(
            "Raster dimensions must not exceed WebP's 16383-pixel limit",
        );
    // Native Mermaid SVGs omit the root height. Restrict dimension changes to
    // the root tag so an absent root attribute never removes a label's
    // foreignObject height and makes that label disappear.
    const document = svg.replace(/<svg\b[^>]*>/, (root) =>
        root
            .replace(/\swidth="[^"]*"/, "")
            .replace(/\sheight="[^"]*"/, "")
            .replace(/<svg\b/, `<svg width="${width}" height="${height}"`),
    );
    const html =
        `<html><head><style>html,body{margin:0;padding:0;` +
        `background:${background}}</style></head><body>${document}</body></html>`;

    const browser = await puppeteer.launch(puppeteerConfig);
    try {
        const page = await browser.newPage();
        // A ratio applies to final pixels, not CSS pixels. Use a physical-pixel
        // viewport for this opt-in path so odd dimensions remain exact even at
        // scale 2 or 3. The SVG keeps the same scaled geometry. The default
        // path retains its existing device scale, layout, and capture bounds.
        const padded = aspectRatio !== undefined;
        const renderScale = padded ? deviceScaleFactor : 1;
        const viewportScale = padded ? 1 : deviceScaleFactor;
        await page.setViewport({
            width: padded ? canvasWidth : width,
            height: padded ? canvasHeight : height,
            deviceScaleFactor: viewportScale,
        });
        await page.setContent(html);
        await page.evaluate(
            (naturalWidth, naturalHeight, padding) => {
                const svg = document.querySelector("svg");
                // Keep the explicit `viewBox`-derived size through measurement
                // and capture. Removing it would let the SVG fall back to the
                // browser's default replaced-element layout, so a diagram whose
                // natural canvas differs from that default would be captured at
                // the wrong size.
                svg.style.width = `${naturalWidth}px`;
                svg.style.height = `${naturalHeight}px`;
                svg.style.maxWidth = "none";
                svg.style.background = "transparent";
                if (padding !== null) {
                    svg.style.position = "absolute";
                    svg.style.left = `${padding.x}px`;
                    svg.style.top = `${padding.y}px`;
                }
            },
            width * renderScale,
            height * renderScale,
            padded
                ? {
                      x: (canvasWidth - width * renderScale) / 2,
                      y: (canvasHeight - height * renderScale) / 2,
                  }
                : null,
        );
        const clip = padded
            ? { x: 0, y: 0, width: canvasWidth, height: canvasHeight }
            : await page.$eval("svg", (element) => {
                  const bounds = element.getBoundingClientRect();
                  return {
                      x: Math.floor(bounds.left),
                      y: Math.floor(bounds.top),
                      width: Math.ceil(bounds.width),
                      height: Math.ceil(bounds.height),
                  };
              });
        await page.setViewport({
            width: clip.x + clip.width,
            height: clip.y + clip.height,
            deviceScaleFactor: viewportScale,
        });
        const data = await page.screenshot({
            type: "webp",
            quality: 90,
            clip,
            omitBackground: background === "transparent",
        });
        await fs.writeFile(output, data);
    } finally {
        await browser.close();
    }
}
