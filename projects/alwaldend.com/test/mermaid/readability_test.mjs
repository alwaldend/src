import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";
import { isDeepStrictEqual } from "node:util";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [executablePath, ...files] = process.argv.slice(2);
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const names = [
    "class",
    "entity-relationship",
    "flowchart",
    "mindmap",
    "sequence",
    "state",
    "edge_semantics",
];
assert.deepEqual(
    files.map((file) => path.basename(file)).sort(),
    names.flatMap((name) => [`${name}.svg`, `${name}.dark.svg`]).sort(),
);
const report = { measurements: [], failures: [] };
const check = (condition, message) => {
    if (!condition) report.failures.push(message);
};

// Inspect delivered SVGs, including their embedded fonts, in fresh documents.
// No theme source or stylesheet strings participate in the assertions.
function measure() {
    const svg = document.querySelector("svg");
    const all = (selector) => [...svg.querySelectorAll(selector)];
    const style = (element) => getComputedStyle(element);
    const canvas = document.createElement("canvas").getContext("2d");
    const rgba = (color) => {
        canvas.clearRect(0, 0, 1, 1);
        canvas.fillStyle = color === "none" ? "transparent" : color;
        canvas.fillRect(0, 0, 1, 1);
        return [...canvas.getImageData(0, 0, 1, 1).data];
    };
    const blend = (color, background, opacity = 1) => {
        const values = rgba(color),
            alpha = (values[3] / 255) * opacity;
        return values
            .slice(0, 3)
            .map((value, i) => value * alpha + background[i] * (1 - alpha));
    };
    const luminance = (rgb) =>
        rgb
            .map((value) => {
                const channel = value / 255;
                return channel <= 0.04045
                    ? channel / 12.92
                    : ((channel + 0.055) / 1.055) ** 2.4;
            })
            .reduce(
                (sum, value, i) => sum + value * [0.2126, 0.7152, 0.0722][i],
                0,
            );
    const contrast = (a, b) =>
        (Math.max(luminance(a), luminance(b)) + 0.05) /
        (Math.min(luminance(a), luminance(b)) + 0.05);
    const background = blend(style(svg).backgroundColor, [255, 255, 255]);
    const opacity = (element) => {
        let value = 1;
        for (
            let parent = element;
            parent instanceof Element;
            parent = parent.parentElement
        ) {
            const css = style(parent);
            value *=
                css.display === "none" || css.visibility === "hidden"
                    ? 0
                    : Number(css.opacity);
        }
        return value;
    };
    const marker = (reference) => {
        if (reference === "none") return null;
        const element = document.getElementById(
            reference.match(/#([^"')]+)/)?.[1],
        );
        if (!element) return { missing: reference };
        return {
            width: element.markerWidth.baseVal.value,
            height: element.markerHeight.baseVal.value,
            units: element.getAttribute("markerUnits"),
            shape: [...element.children].map((child) => [
                child.tagName,
                ...["d", "points", "r", "x1", "x2", "y1", "y2"].map((key) =>
                    child.getAttribute(key),
                ),
            ]),
        };
    };
    const stroke = (element) => {
        const css = style(element),
            matrix = element.getScreenCTM();
        const width = parseFloat(css.strokeWidth);
        const scale =
            css.vectorEffect === "non-scaling-stroke"
                ? 1
                : Math.min(
                      Math.hypot(matrix.a, matrix.b),
                      Math.hypot(matrix.c, matrix.d),
                  );
        return {
            id: element.dataset.id,
            classes: [...element.classList],
            width,
            screenWidth: width * scale,
            color: css.stroke,
            opacity: opacity(element) * Number(css.strokeOpacity),
            contrast: contrast(
                blend(
                    css.stroke,
                    background,
                    opacity(element) * Number(css.strokeOpacity),
                ),
                background,
            ),
            dash: css.strokeDasharray,
            dashed: (css.strokeDasharray.match(/[\d.]+/g) || []).some(
                (value) => Number(value) > 0,
            ),
            start: marker(css.markerStart),
            end: marker(css.markerEnd),
        };
    };
    const fonts = [...document.fonts].map((face) => ({
        family: face.family.replace(/["']/g, ""),
        status: face.status,
    }));
    const labels = all("p, text")
        .filter((element) => element.textContent.trim())
        .map((element) => {
            const css = style(element),
                family = css.fontFamily
                    .split(",")[0]
                    .replace(/["']/g, "")
                    .trim();
            canvas.font = `${css.fontWeight} ${css.fontSize} ${css.fontFamily}`;
            const advances = [
                ...new Set(element.textContent.match(/[A-Za-z0-9]/g) || []),
            ].map((glyph) => canvas.measureText(glyph).width);
            return {
                text: element.textContent.trim(),
                family,
                size: css.fontSize,
                monospace:
                    advances.length === 0 ||
                    Math.max(...advances) - Math.min(...advances) < 0.01,
                available:
                    fonts.some(
                        (face) =>
                            face.family === family && face.status === "loaded",
                    ) &&
                    document.fonts.check(
                        `${css.fontSize} ${css.fontFamily}`,
                        element.textContent,
                    ),
            };
        });
    return {
        background,
        fonts,
        labels,
        viewBox: svg.getAttribute("viewBox"),
        width: svg.getBoundingClientRect().width,
        geometry: all(
            "path, rect, circle, ellipse, polygon, polyline, line, text, foreignObject",
        )
            .filter((element) => !element.closest("marker, defs"))
            .map((element) => {
                const rect = element.getBoundingClientRect();
                // Classic shapes can contain randomized but collinear Bézier
                // control points. Compare rendered bounds, not path serialization.
                return [
                    element.tagName,
                    ...[rect.x, rect.y, rect.width, rect.height].map(
                        (value) => Math.round(value * 1000) / 1000,
                    ),
                ];
            }),
        edges: all('[data-edge="true"]').map(stroke),
        guides: all(".actor-line").map(stroke),
        messages: all(".messageLine0, .messageLine1").map(stroke),
        nodes: all(".mindmap-node").map((node) => {
            const label =
                node.querySelector(".nodeLabel p") ||
                node.querySelector(".nodeLabel, text");
            const surface = node.querySelector(".label-container, .node-bkg");
            const fill = blend(
                style(surface).fill,
                background,
                opacity(surface) * Number(style(surface).fillOpacity),
            );
            const color =
                label.namespaceURI === svg.namespaceURI
                    ? style(label).fill
                    : style(label).color;
            return {
                root: node.classList.contains("section-root"),
                text: label.textContent,
                fill,
                color,
                contrast: contrast(blend(color, fill, opacity(label)), fill),
            };
        }),
    };
}

const browser = await puppeteer.launch({
    executablePath: path.resolve(executablePath),
    headless: "shell",
});
try {
    for (const file of files) {
        const page = await browser.newPage();
        await page.setRequestInterception(true);
        page.on("request", (request) =>
            request.url().startsWith("data:")
                ? request.continue()
                : request.abort(),
        );
        page.on("pageerror", (error) =>
            report.failures.push(`${file}: ${error.message}`),
        );
        await page.setContent(await fs.readFile(path.resolve(file), "utf8"));
        await page.evaluate(async () => {
            document.body.style.margin = "0";
            await document.fonts.ready;
        });
        for (const width of [728, 360]) {
            await page.setViewport({
                width,
                height: 1000,
                deviceScaleFactor: 1,
            });
            const result = {
                file: path.basename(file),
                viewport: width,
                ...(await page.evaluate(measure)),
            };
            report.measurements.push(result);
            await (
                await page.$("svg")
            ).screenshot({
                path: path.join(output, `${result.file}-${width}.png`),
            });
        }
        await page.close();
    }
    // Collect every screenshot before checking, so failed candidates remain reviewable.
    for (const result of report.measurements) {
        const name = result.file.replace(/(?:\.dark)?\.svg$/, "");
        const context = `${result.file} at ${result.viewport}px`;
        const verify = (condition, message) =>
            check(condition, `${context}: ${message}`);
        verify(
            result.width > 0 && result.width <= result.viewport + 1,
            "image fits viewport",
        );
        verify(
            result.fonts.length > 0 &&
                result.labels.length > 0 &&
                result.labels.every((label) => label.available),
            "embedded label fonts load",
        );
        verify(
            result.labels.every((label) => label.monospace),
            "label glyphs have equal advances",
        );
        const light = report.measurements.find(
            (entry) =>
                entry.file === `${name}.svg` &&
                entry.viewport === result.viewport,
        );
        for (const key of ["viewBox", "geometry", "fonts", "labels"])
            verify(
                isDeepStrictEqual(result[key], light[key]),
                `modes preserve ${key}`,
            );
        for (const edge of [...result.edges, ...result.messages]) {
            const dashed = edge.classes.some((value) =>
                [
                    "edge-pattern-dashed",
                    "edge-pattern-dotted",
                    "messageLine1",
                ].includes(value),
            );
            verify(
                edge.dashed === dashed,
                `${edge.id}: dash semantics survive`,
            );
            for (const marker of [edge.start, edge.end].filter(Boolean))
                verify(
                    marker.width > 0 &&
                        marker.height > 0 &&
                        marker.shape.length > 0,
                    `${edge.id}: marker resolves`,
                );
        }
        const semantic = ({ id, dash, start, end }) => ({
            id,
            dash,
            start,
            end,
        });
        verify(
            isDeepStrictEqual(
                result.edges.map(semantic),
                light.edges.map(semantic),
            ),
            "modes preserve dashes and markers",
        );
        if (name === "edge_semantics") {
            const edge = (id) => result.edges.find((item) => item.id === id);
            const normal = edge("normal"),
                thick = edge("thick"),
                dashed = edge("dashed"),
                invisible = edge("invisible"),
                arrows = edge("arrows");
            verify(
                normal && thick && dashed && invisible && arrows,
                "all authored edge semantics are present",
            );
            if (normal && thick && dashed && invisible && arrows) {
                verify(
                    normal.width > 0 && thick.width > normal.width,
                    "thick edges remain stronger",
                );
                verify(
                    invisible.width === 0 ||
                        invisible.opacity === 0 ||
                        invisible.contrast === 1,
                    "invisible edges stay invisible",
                );
                verify(
                    !normal.dashed && !thick.dashed && dashed.dashed,
                    "solid and dashed edges remain distinct",
                );
                verify(
                    !normal.start &&
                        normal.end &&
                        !thick.start &&
                        thick.end &&
                        !dashed.start &&
                        dashed.end &&
                        arrows.start &&
                        arrows.end,
                    "arrow direction survives",
                );
                verify(
                    isDeepStrictEqual(normal.end, thick.end) &&
                        isDeepStrictEqual(normal.end, dashed.end) &&
                        isDeepStrictEqual(normal.end, arrows.end),
                    "arrow geometry is independent of edge weight",
                );
            }
        }
        if (name === "mindmap") {
            verify(
                result.edges.length > 0 &&
                    result.edges.every(
                        (edge) => edge.width > 0 && edge.contrast >= 3,
                    ),
                "branches contrast at least 3:1",
            );
            verify(
                result.nodes.some((node) => node.root) &&
                    result.nodes.every((node) => node.contrast >= 4.5),
                "all labels, including root, contrast at least 4.5:1",
            );
        }
        if (result.viewport !== 728) continue; // Mobile preserves geometry; dense labels need native-size viewing.
        if (!["mindmap", "sequence"].includes(name)) {
            const normal = result.edges.filter(
                (edge) =>
                    edge.classes.includes("edge-thickness-normal") &&
                    !edge.classes.includes("edge-thickness-thick") &&
                    !edge.classes.includes("edge-thickness-invisible"),
            );
            verify(
                normal.length > 0 &&
                    normal.every(
                        (edge) =>
                            edge.screenWidth >= 1 &&
                            edge.opacity > 0 &&
                            edge.contrast > 1,
                    ),
                "normal graph strokes occupy at least one CSS pixel",
            );
        }
        if (name === "sequence") {
            verify(
                result.guides.length > 0 && result.messages.length > 0,
                "guides and messages are present",
            );
            const messageWidth = Math.min(
                ...result.messages.map((edge) => edge.screenWidth),
            );
            const messageContrast = Math.min(
                ...result.messages.map((edge) => edge.contrast),
            );
            verify(
                result.guides.every(
                    (edge) =>
                        edge.screenWidth >= 0.75 &&
                        edge.screenWidth < messageWidth &&
                        edge.contrast >= 3 &&
                        edge.contrast < messageContrast,
                ),
                "guides remain visible and subordinate to messages",
            );
        }
    }
    assert.equal(report.failures.length, 0, report.failures.join("\n"));
} finally {
    await fs.writeFile(
        path.join(output, "readability-report.json"),
        JSON.stringify(report, null, 2),
    );
    await browser.close();
}
