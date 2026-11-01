import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { createHash } from "node:crypto";
import { createRequire } from "node:module";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");

// These expectations describe the maintained examples, not the renderer's
// implementation. Geometry checks below independently inspect exported routes.
const expected = {
    "graph-boundaries": {
        nodes: 2,
        edges: 2,
        labels: ["Source", "Target", "Retry"],
    },
    "group-title": {
        nodes: 2,
        edges: 1,
        labels: [
            "An extensive processing group title",
            "Including review and publication",
        ],
    },
    "sequence-identifiers": {
        nodes: 2,
        edges: 3,
        labels: ["Client", "Service", "Request", "Response", "Inspect"],
    },
    "class-text": {
        nodes: 2,
        edges: 1,
        labels: ["Box<T>", "+List<T> items", "Repository", "stores"],
    },
    "graph-semantics": {
        nodes: 4,
        edges: 2,
        labels: ["This is line three", "quoted", "Layout target", "Finished"],
    },
    "sequence-activations": {
        nodes: 2,
        edges: 7,
        labels: ["Nested request", "Nested response", "Explicit response"],
    },
    "sequence-text": {
        nodes: 2,
        edges: 3,
        labels: [
            "First",
            "Document processing service",
            "One",
            "Six",
            "Ready",
            "Complete",
        ],
    },
    flowchart: {
        nodes: 8,
        edges: 9,
        labels: ["Document processing", "Search index", "Revise", "Audit log"],
    },
    state: {
        minimumNodes: 8,
        edges: 11,
        labels: ["Processing", "Validate", "Cancelled", "retry"],
    },
    class: {
        nodes: 5,
        edges: 4,
        labels: ["Repository", "interface", "find(id)", "1..*", "implements"],
    },
    "entity-relationship": {
        nodes: 4,
        edges: 3,
        labels: ["AUTHOR", "author_id", "FK", "produces"],
    },
    mindmap: {
        nodes: 10,
        edges: 9,
        labels: ["Document service", "Operations", "Metrics"],
    },
    sequence: {
        nodes: 4,
        edges: 8,
        labels: [
            "Client",
            "Submit document",
            "Until acknowledged",
            "Valid document",
            "Invalid document",
            "Delivery may be repeated",
        ],
    },
    "sequence-semantics": {
        nodes: 3,
        edges: 6,
        labels: [
            "First",
            "Second",
            "Observer",
            "In outer loop",
            "In nested block",
            "Note without a message",
        ],
    },
};

function intersects(start, end, node, margin = 0.8) {
    let enter = 0;
    let leave = 1;
    for (let axis = 0; axis < 3; axis += 1) {
        const lower = node.position[axis] - node.size[axis] / 2 - margin;
        const upper = node.position[axis] + node.size[axis] / 2 + margin;
        const delta = end[axis] - start[axis];
        if (Math.abs(delta) < 1e-8) {
            if (start[axis] < lower || start[axis] > upper) return false;
        } else {
            const times = [
                (lower - start[axis]) / delta,
                (upper - start[axis]) / delta,
            ].sort((a, b) => a - b);
            enter = Math.max(enter, times[0]);
            leave = Math.min(leave, times[1]);
            if (enter > leave) return false;
        }
    }
    return enter <= leave;
}

// Inspect the actual rendered objects, not just layout-model positions. The
// earlier DOM-overflow check missed clipping inside the WebGL canvas.
function checkView(view, name, readable) {
    assert.ok(view.labels.length > 0, `${name}: rendered labels`);
    for (const label of view.labels) {
        if (readable)
            assert.ok(
                label.fontPixels >= 12,
                `${name}: readable ${label.text}`,
            );
        const b = label.bounds;
        assert.ok(
            b.left >= -1 &&
                b.top >= -1 &&
                b.right <= view.canvas.width + 1 &&
                b.bottom <= view.canvas.height + 1,
            `${name}: label stays in canvas: ${label.text}`,
        );
        for (const other of view.labels) {
            if (label === other || label.kind === "node") continue;
            const a = other.bounds;
            const overlapX =
                Math.min(b.right, a.right) - Math.max(b.left, a.left);
            const overlapY =
                Math.min(b.bottom, a.bottom) - Math.max(b.top, a.top);
            assert.ok(
                overlapX <= 1 || overlapY <= 1,
                `${name}: ${label.text} overlaps ${other.text}`,
            );
        }
    }
    for (const point of view.sceneCorners)
        assert.ok(
            point[0] >= -1 &&
                point[1] >= -1 &&
                point[0] <= view.canvas.width + 1 &&
                point[1] <= view.canvas.height + 1,
            `${name}: fitted geometry stays in canvas`,
        );
}

function checkSelection(view, highlighted, name) {
    const expected = new Set(
        highlighted.map(([kind, id]) => JSON.stringify([kind, id])),
    );
    const { model } = view;
    assert.equal(
        view.elements.length,
        model.nodes.length +
            model.groups.length +
            model.edges.length +
            model.notes.length,
    );
    for (const element of view.elements)
        assert.deepEqual(
            element.opacities,
            [
                expected.has(JSON.stringify([element.kind, element.id]))
                    ? 1
                    : 0.14,
            ],
            `${name}: rendered selection opacity for ${element.kind} ${element.id}`,
        );
}

async function selectElement(page, kind, id) {
    const key = await page.evaluate(
        (kind, id) =>
            window.diagramView
                .inspect()
                .elements.find((item) => item.kind === kind && item.id === id)
                ?.key,
        kind,
        id,
    );
    assert.ok(key, `Selectable ${kind} ${id}`);
    await page.select("#selection", key);
    return key;
}

function activationInterval(model, actor, startLabel, endLabel) {
    const start = model.edges.find((edge) => edge.label === startLabel);
    const end = model.edges.find((edge) => edge.label === endLabel);
    const matches = model.guides.filter(
        (guide) =>
            guide.activation &&
            guide.actor === actor &&
            guide.points[0][1] === start.points.at(-1)[1] &&
            guide.points.at(-1)[1] === end.points.at(-1)[1],
    );
    assert.equal(
        matches.length,
        1,
        `${actor}: activation aligns with ${startLabel} and ${endLabel}`,
    );
    return matches[0];
}

async function checkLargeFlowchart(page, output) {
    const frames = [];
    async function capture(name) {
        const { png, ...frame } = await page.$eval(
            "#viewport canvas",
            (canvas) => {
                const sample = document.createElement("canvas");
                sample.width = canvas.width;
                sample.height = canvas.height;
                const context = sample.getContext("2d");
                context.drawImage(canvas, 0, 0);
                const pixels = context.getImageData(
                    0,
                    0,
                    sample.width,
                    sample.height,
                ).data;
                const background = window.diagramView
                    .inspect()
                    .background.match(/[a-f0-9]{2}/gi)
                    .map((value) => parseInt(value, 16));
                let painted = 0;
                for (let offset = 0; offset < pixels.length; offset += 4)
                    if (
                        background.some(
                            (channel, index) =>
                                Math.abs(pixels[offset + index] - channel) > 8,
                        )
                    )
                        painted += 1;
                return {
                    painted,
                    width: canvas.width,
                    height: canvas.height,
                    png: canvas.toDataURL("image/png").split(",")[1],
                };
            },
        );
        await fs.writeFile(
            path.join(output, `large-flowchart.${name}.png`),
            Buffer.from(png, "base64"),
        );
        frames.push({ name, ...frame });
        await fs.writeFile(
            path.join(output, "large-flowchart.visibility.json"),
            JSON.stringify(frames, null, 2) + "\n",
        );
        assert.ok(
            frame.painted > 100,
            `large-flowchart ${name}: visible geometry and labels (${frame.painted} pixels)`,
        );
    }
    await capture("readable");
    await page.click("#fit");
    await capture("overview");
    const panel = await page.$eval("#viewport", (element) => {
        const box = element.getBoundingClientRect();
        return { x: box.x + box.width / 2, y: box.y + box.height / 2 };
    });
    const before = await page.evaluate(
        () => window.diagramView.inspect().camera,
    );
    await page.mouse.move(panel.x, panel.y);
    await page.mouse.down();
    await page.mouse.move(panel.x + 60, panel.y + 20, { steps: 8 });
    await page.mouse.up();
    assert.notDeepEqual(
        await page.evaluate(() => window.diagramView.inspect().camera),
        before,
        "large-flowchart orbit",
    );
    await capture("orbit");
    await page.click("#front");
    await capture("front");
    await page.select("#appearance", "dark");
    await capture("dark");
    await page.setViewport({ width: 390, height: 844 });
    await page.waitForFunction(
        () => window.diagramView.inspect().container.width <= 390,
    );
    await page.click("#reset");
    await capture("mobile-readable");
    await page.click("#fit");
    await capture("mobile-overview");
    await page.setViewport({ width: 1440, height: 1050 });
    await page.waitForFunction(
        () => window.diagramView.inspect().container.width > 390,
    );
    await page.click("#reset");
    await capture("reset");
    return frames;
}

const [browserPath, ...inputs] = process.argv.slice(2);
assert.equal(
    inputs.length,
    Object.keys(expected).length + 1,
    "Every maintained sample and adapter fixture must be verified",
);
const output =
    process.env.TEST_UNDECLARED_OUTPUTS_DIR ?? process.env.SPATIAL_TEST_OUTPUT;
assert.ok(output, "An artifact directory is required");
await fs.mkdir(output, { recursive: true });
const browser = await puppeteer.launch({
    executablePath: path.resolve(browserPath),
    headless: "shell",
    args: ["--use-angle=swiftshader", "--enable-unsafe-swiftshader"],
});
const report = [];
try {
    for (const input of inputs) {
        const name = path.basename(input).replace(/\.3d\.html$/, "");
        const contract = expected[name];
        assert.ok(
            contract || name === "large-flowchart",
            `Unknown sample ${name}`,
        );
        const html = await fs.readFile(input);
        assert.ok(
            html.length < 4 * 1024 * 1024,
            `${name}: offline view stays below 4 MiB`,
        );
        const page = await browser.newPage();
        await page.evaluateOnNewDocument(() => {
            window.renderedFonts = [];
            const fillText = CanvasRenderingContext2D.prototype.fillText;
            CanvasRenderingContext2D.prototype.fillText = function (...args) {
                window.renderedFonts.push(this.font);
                return fillText.apply(this, args);
            };
        });
        await page.setViewport({
            width: 1440,
            height: 1050,
            deviceScaleFactor: 1,
        });
        const errors = [];
        const network = [];
        page.on("pageerror", (error) => errors.push(error.message));
        await page.setRequestInterception(true);
        page.on("request", (request) => {
            if (/^https?:/.test(request.url())) {
                network.push(request.url());
                request.abort();
            } else request.continue();
        });
        await page.goto(pathToFileURL(path.resolve(input)).href);
        await page.waitForFunction(() => window.diagramView?.ready === true, {
            timeout: 30000,
        });
        const initial = await page.evaluate(() =>
            window.diagramView.inspect(),
        );
        const { model } = initial;
        if (name === "large-flowchart") {
            assert.equal(model.nodes.length, 50);
            assert.equal(model.edges.length, 49);
            const frames = await checkLargeFlowchart(page, output);
            assert.deepEqual(errors, [], `${name} browser errors`);
            assert.deepEqual(network, [], `${name} remote requests`);
            report.push({
                name,
                frames,
                sha256: createHash("sha256").update(html).digest("hex"),
                passed: true,
            });
            await page.close();
            continue;
        }
        assert.ok(
            initial.drawCalls <= 150,
            `${name}: maintained scenes stay below 150 draw calls`,
        );
        const renderedFonts = await page.evaluate(() => [
            ...new Set(window.renderedFonts),
        ]);
        assert.ok(renderedFonts.length > 0);
        assert.ok(
            renderedFonts.every((font) => font.includes("Liberation Mono")),
            `${name}: canvas labels use the theme font: ${renderedFonts.join(", ")}`,
        );
        if (contract.nodes)
            assert.equal(model.nodes.length, contract.nodes, `${name} nodes`);
        else
            assert.ok(
                model.nodes.length >= contract.minimumNodes,
                `${name} nodes`,
            );
        assert.equal(model.edges.length, contract.edges, `${name} edges`);
        assert.equal(
            initial.elements.length,
            model.nodes.length +
                model.groups.length +
                model.edges.length +
                model.notes.length,
            `${name}: every element has a distinct selectable object`,
        );
        const text = JSON.stringify(model);
        for (const label of contract.labels)
            assert.ok(text.includes(label), `${name} missing ${label}`);
        assert.ok(
            initial.meshCount >= model.nodes.length + model.edges.length,
            `${name} must contain 3D objects`,
        );
        assert.ok(
            new Set(model.nodes.map((node) => node.position[2])).size > 1,
            `${name} depth placement`,
        );
        for (const node of model.nodes) {
            assert.ok(node.position.every(Number.isFinite));
            assert.ok(
                node.size.every(
                    (value) => Number.isFinite(value) && value > 0,
                ),
            );
        }
        let curvedEdges = 0;
        for (const edge of model.edges) {
            assert.ok(edge.points.length >= 2, `${name} ${edge.id} route`);
            assert.ok(edge.points.flat().every(Number.isFinite));
            let bends = 0;
            for (let i = 1; i + 1 < edge.points.length; i += 1) {
                const incoming = edge.points[i].map(
                    (value, axis) => value - edge.points[i - 1][axis],
                );
                const outgoing = edge.points[i + 1].map(
                    (value, axis) => value - edge.points[i][axis],
                );
                const lengths =
                    Math.hypot(...incoming) * Math.hypot(...outgoing);
                assert.ok(
                    lengths > 1e-8,
                    `${name}: ${edge.id} distinct route points`,
                );
                const cosine =
                    incoming.reduce(
                        (sum, value, axis) => sum + value * outgoing[axis],
                        0,
                    ) / lengths;
                const angle = Math.acos(Math.max(-1, Math.min(1, cosine)));
                assert.ok(
                    angle < Math.PI / 9,
                    `${name}: ${edge.id} has a sharp bend of ${(angle * 180) / Math.PI} degrees`,
                );
                if (angle > 0.001) bends += 1;
            }
            if (bends >= 3) curvedEdges += 1;
            for (const node of model.nodes) {
                if (node.id === edge.start || node.id === edge.end) continue;
                for (let i = 1; i < edge.points.length; i += 1)
                    assert.ok(
                        !intersects(edge.points[i - 1], edge.points[i], node),
                        `${name}: ${edge.id} intersects ${node.id}`,
                    );
            }
        }
        assert.ok(curvedEdges > 0, `${name}: visibly curved connections`);
        if (name === "graph-boundaries") {
            const loop = model.edges.find((edge) => edge.start === edge.end);
            const node = model.nodes.find((item) => item.id === loop.start);
            const face = {
                position: [...node.position.slice(0, 2), 0],
                size: [...node.size.slice(0, 2), 1],
            };
            for (let i = 1; i < loop.points.length; i += 1)
                assert.ok(
                    !intersects(
                        [...loop.points[i - 1].slice(0, 2), 0],
                        [...loop.points[i].slice(0, 2), 0],
                        face,
                    ),
                    "Self-loop clears its own face in front projection",
                );
        }
        if (name === "group-title") {
            const group = model.groups[0];
            const title = initial.labels.find(
                (label) => label.ownerId === group.id,
            );
            assert.ok(
                title.text.includes("\n"),
                "Group title retains explicit line breaks",
            );
        }
        for (const group of model.groups) {
            assert.ok(group.size.every((value) => value > 0));
            for (const id of group.members) {
                const child = [
                    ...model.nodes,
                    ...model.groups,
                    ...model.notes,
                ].find((item) => item.id === id);
                assert.ok(child, `${name} group member ${id}`);
                for (let axis = 0; axis < 3; axis += 1)
                    assert.ok(
                        Math.abs(child.position[axis] - group.position[axis]) +
                            child.size[axis] / 2 <=
                            group.size[axis] / 2 + 0.01,
                        `${name}: ${group.id} contains ${id}`,
                    );
            }
        }
        if (name === "graph-semantics") {
            const node = (id) => model.nodes.find((item) => item.id === id);
            assert.equal(node("B").label, '\u2665 & "quoted" \u00a0');
            assert.equal(model.groups[0].label, 'Vertical "group"');
            assert.equal(model.edges[0].label, '"Proceed"');
            assert.ok(node("A").position[1] > node("B").position[1]);
            assert.ok(
                node("B").position[1] > node("C").position[1],
                "Invisible constraint retains vertical ordering",
            );
            assert.equal(
                node("A").position[0],
                node("B").position[0],
                "Local TB direction overrides root LR",
            );
            assert.equal(node("B").position[0], node("C").position[0]);
            assert.ok(
                !model.edges.some(
                    (edge) => edge.start === "B" && edge.end === "C",
                ),
                "Invisible links are not relationships",
            );
            await selectElement(page, "node", "C");
            checkSelection(
                await page.evaluate(() => window.diagramView.inspect()),
                [["node", "C"]],
                "Invisible links do not expand selection",
            );
            await page.select("#selection", "");
        }
        if (name === "sequence-activations") {
            assert.equal(
                model.guides.filter((guide) => guide.activation).length,
                4,
            );
            activationInterval(model, "A", "Request", "Response");
            const outer = activationInterval(
                model,
                "B",
                "Request",
                "Response",
            );
            const inner = activationInterval(
                model,
                "B",
                "Nested request",
                "Nested response",
            );
            activationInterval(
                model,
                "B",
                "Explicit request",
                "Explicit response",
            );
            assert.ok(
                Math.abs(inner.points[0][0] - outer.points[0][0]) >= 8,
                "Nested bars do not cover each other",
            );
            for (const label of ["Nested request", "Nested response"]) {
                const edge = model.edges.find((item) => item.label === label);
                assert.equal(
                    edge.points.at(-1)[0],
                    inner.points[0][0],
                    `${label}: message reaches nested activation`,
                );
            }
        }
        if (name === "sequence-semantics") {
            assert.deepEqual(
                model.nodes.map((node) => node.id),
                ["10", "2", "T"],
                "Participant declaration order, including numeric IDs",
            );
            const markers = Object.fromEntries(
                model.edges.map((edge) => [edge.label, edge.endMarker]),
            );
            assert.equal(
                markers["Plain line"],
                "none",
                "Solid arrowless message",
            );
            assert.equal(
                markers["Dashed line"],
                "none",
                "Dashed arrowless message",
            );
            assert.equal(
                markers.Request,
                "point",
                "Arrowheaded message stays directed",
            );
            assert.equal(
                model.groups.length,
                3,
                "Nested and note-only blocks",
            );
            assert.equal(model.notes.length, 3);
            for (const note of model.notes) {
                const parent = model.groups.find(
                    (group) => group.id === note.parentId,
                );
                assert.ok(
                    parent?.members.includes(note.id),
                    `${note.label}: block membership`,
                );
            }
        }
        if (name === "sequence") {
            const messages = model.edges;
            assert.deepEqual(
                messages.map((edge) => edge.label),
                [
                    "Submit document",
                    "Enqueue job",
                    "Job identifier",
                    "Accepted",
                    "Deliver job",
                    "Build search entry",
                    "Acknowledge",
                    "Reject with reason",
                ],
            );
            for (let i = 1; i < messages.length; i += 1)
                assert.ok(
                    messages[i].points[0][1] < messages[i - 1].points[0][1],
                    "Sequence time proceeds downwards",
                );
            assert.equal(
                model.guides.filter((guide) => guide.activation).length,
                1,
                "API activation interval",
            );
            activationInterval(model, "A", "Submit document", "Accepted");
            assert.equal(
                model.groups.length,
                2,
                "Nested loop and alternative blocks",
            );
            assert.equal(model.notes.length, 1, "Delivery note");
            assert.equal(
                messages.filter((edge) => edge.dashed).length,
                4,
                "Reply styles",
            );
            const selfCall = messages.find((edge) => edge.start === edge.end);
            assert.ok(
                selfCall.points.length >= 4 &&
                    selfCall.points[0][1] > selfCall.points.at(-1)[1],
                "Self-call keeps a visible return path",
            );
        }
        for (const mode of ["light", "dark"]) {
            await page.select("#appearance", mode);
            await page.waitForFunction(
                (value) => window.diagramView.inspect().appearance === value,
                {},
                mode,
            );
            const themed = await page.evaluate(() =>
                window.diagramView.inspect(),
            );
            const colorsMatch = await page.evaluate(
                ({ background, palette }) => {
                    const context = document
                        .createElement("canvas")
                        .getContext("2d");
                    context.fillStyle = background;
                    const actual = context.fillStyle;
                    context.fillStyle = palette.background;
                    return actual === context.fillStyle;
                },
                themed,
            );
            assert.ok(colorsMatch, `${name} ${mode} theme background`);
            assert.equal(themed.fontLoaded, true, `${name} embedded font`);
            assert.equal(themed.fontFamily, "Liberation Mono, monospace");
            checkView(themed, `${name} ${mode}`, true);
            await page.screenshot({
                path: path.join(output, `${name}.${mode}.png`),
            });
            const pixels = await page.$eval(
                "#viewport canvas",
                (canvas) => canvas.toDataURL("image/png").split(",")[1],
            );
            await fs.writeFile(
                path.join(output, `${name}.scene.${mode}.png`),
                Buffer.from(pixels, "base64"),
            );
        }
        await page.click("#front");
        const frontView = await page.evaluate(() =>
            window.diagramView.inspect(),
        );
        for (const node of [...model.nodes, ...model.notes].filter((node) =>
            ["rect", "squareRect", "classBox", "diamond"].includes(node.shape),
        )) {
            const card = frontView.nodeBounds.find(
                (item) => item.id === node.id,
            ).bounds;
            const text = frontView.labels.find(
                (item) => item.ownerId === node.id && item.kind === "node",
            ).bounds;
            assert.ok(
                text.left >= card.left &&
                    text.right <= card.right &&
                    text.top >= card.top &&
                    text.bottom <= card.bottom,
                `${name}: ${node.label} text stays inside its card`,
            );
            if (node.shape === "diamond") {
                const cx = (card.left + card.right) / 2;
                const cy = (card.top + card.bottom) / 2;
                for (const x of [text.left, text.right])
                    for (const y of [text.top, text.bottom])
                        assert.ok(
                            Math.abs(x - cx) / ((card.right - card.left) / 2) +
                                Math.abs(y - cy) /
                                    ((card.bottom - card.top) / 2) <=
                                1,
                            `${name}: text rectangle stays inside the diamond face`,
                        );
            }
        }
        if (name === "sequence-text") {
            const cards = frontView.nodeBounds;
            const a = cards.find((item) => item.id === "A").bounds;
            const b = cards.find((item) => item.id === "B").bounds;
            assert.ok(
                a.right + 30 < b.left,
                "Measured participant widths retain separation",
            );
            const note = frontView.labels.find((item) =>
                item.text.startsWith("One\n"),
            );
            const message = frontView.labels.find(
                (item) => item.text === "Continue",
            );
            assert.ok(
                note.bounds.bottom + 20 < message.bounds.top,
                "Multiline notes reserve space before the next message",
            );
        }
        checkView(
            await page.evaluate(() => window.diagramView.inspect()),
            `${name} front view`,
            true,
        );
        await page.screenshot({
            path: path.join(output, `${name}.front.png`),
        });
        await page.click("#reset");
        const selectedEdge = model.edges[0];
        await selectElement(page, "edge", selectedEdge.id);
        checkSelection(
            await page.evaluate(() => window.diagramView.inspect()),
            [
                ["edge", selectedEdge.id],
                ["node", selectedEdge.start],
                ["node", selectedEdge.end],
            ],
            `${name} connection`,
        );
        assert.ok(
            (
                await page.$eval("#details", (element) => element.textContent)
            ).includes(
                selectedEdge.label ||
                    `${selectedEdge.start} → ${selectedEdge.end}`,
            ),
            `${name}: connection details belong to the selected connection`,
        );
        const selectedNodeKey = await selectElement(
            page,
            "node",
            model.nodes[0].id,
        );
        assert.equal(
            await page.evaluate(() => window.diagramView.inspect().selected),
            selectedNodeKey,
        );
        const all = [
            ["node", model.nodes],
            ["group", model.groups],
            ["edge", model.edges],
            ["note", model.notes],
        ].flatMap(([kind, items]) => items.map((item) => ({ ...item, kind })));
        const ids = (labels) =>
            labels.map((label) => {
                const item = all.find((item) => item.label === label);
                assert.ok(item, `${name}: selection fixture label ${label}`);
                return [item.kind, item.id];
            });
        if (name === "class")
            checkSelection(
                await page.evaluate(() => window.diagramView.inspect()),
                ids([
                    "Repository",
                    "SearchRepository",
                    "Indexer",
                    "Document",
                    "implements",
                    "uses",
                    "stores",
                ]),
                "Class node direct relationships",
            );
        if (name === "sequence-semantics") {
            await selectElement(page, ...ids(["loop: Work"])[0]);
            checkSelection(
                await page.evaluate(() => window.diagramView.inspect()),
                ids([
                    "loop: Work",
                    "opt: Nested work",
                    "First",
                    "Second",
                    "Request",
                    "Inspect",
                    "Response",
                    "In outer loop",
                    "In nested block",
                ]),
                "Nested sequence block",
            );
            const selectedPixels = await page.$eval(
                "#viewport canvas",
                (canvas) => canvas.toDataURL("image/png").split(",")[1],
            );
            await fs.writeFile(
                path.join(output, `${name}.selected-block.png`),
                Buffer.from(selectedPixels, "base64"),
            );
            await selectElement(page, ...ids(["loop: Annotation only"])[0]);
            checkSelection(
                await page.evaluate(() => window.diagramView.inspect()),
                ids(["loop: Annotation only", "Note without a message"]),
                "Note-only sequence block",
            );
            await selectElement(page, ...ids(["In nested block"])[0]);
            checkSelection(
                await page.evaluate(() => window.diagramView.inspect()),
                ids(["In nested block"]),
                "Sequence note",
            );
        }
        const before = await page.evaluate(
            () => window.diagramView.inspect().camera,
        );
        await page.mouse.move(850, 450);
        await page.mouse.down();
        await page.mouse.move(1060, 560, { steps: 12 });
        await page.mouse.up();
        const after = await page.evaluate(
            () => window.diagramView.inspect().camera,
        );
        assert.notDeepEqual(after, before, `${name} orbit`);
        await page.mouse.wheel({ deltaY: -240 });
        await page.waitForFunction(
            () => window.diagramView.inspect().zoom > 1,
        );
        const beforePan = await page.evaluate(
            () => window.diagramView.inspect().target,
        );
        await page.mouse.down({ button: "right" });
        await page.mouse.move(980, 600, { steps: 6 });
        await page.mouse.up({ button: "right" });
        assert.notDeepEqual(
            await page.evaluate(() => window.diagramView.inspect().target),
            beforePan,
            `${name} pan`,
        );
        await page.click("#reset");
        const reset = await page.evaluate(
            () => window.diagramView.inspect().camera,
        );
        assert.ok(
            reset.every(
                (value, i) => Math.abs(value - initial.camera[i]) < 0.01,
            ),
            `${name} reset`,
        );
        assert.equal(
            await page.evaluate(() => window.diagramView.inspect().zoom),
            1,
            `${name} reset zoom`,
        );
        await page.setViewport({ width: 390, height: 844 });
        await page.select("#selection", "");
        await page.waitForFunction(
            () => window.diagramView.inspect().container.width <= 390,
        );
        await page.click("#reset");
        checkView(
            await page.evaluate(() => window.diagramView.inspect()),
            `${name} mobile readable`,
            true,
        );
        await page.screenshot({
            path: path.join(output, `${name}.mobile.png`),
            fullPage: true,
        });
        // A fresh touch device must retain readable text and be able to reach
        // content beyond the visible panel, not merely avoid DOM overflow.
        await page.setViewport({
            width: 390,
            height: 844,
            isMobile: true,
            hasTouch: true,
            deviceScaleFactor: 1,
        });
        await page.goto(pathToFileURL(path.resolve(input)).href);
        await page.waitForFunction(() => window.diagramView?.ready === true);
        checkView(
            await page.evaluate(() => window.diagramView.inspect()),
            `${name} fresh touch view`,
            true,
        );
        const touch = await page.createCDPSession();
        const panel = await page.$eval("#viewport", (element) => {
            const rect = element.getBoundingClientRect();
            return {
                x: rect.left,
                y: rect.top,
                width: element.clientWidth,
                scrollAxis:
                    element.scrollWidth > element.clientWidth + 40
                        ? "x"
                        : element.scrollHeight > element.clientHeight + 40
                          ? "y"
                          : null,
            };
        });
        const x = panel.x + panel.width - 30;
        const y = panel.y + 140;
        if (panel.scrollAxis) {
            await touch.send("Input.dispatchTouchEvent", {
                type: "touchStart",
                touchPoints: [{ x, y }],
            });
            for (let step = 1; step <= 6; step += 1)
                await touch.send("Input.dispatchTouchEvent", {
                    type: "touchMove",
                    touchPoints: [
                        {
                            x: x - (panel.scrollAxis === "x" ? step * 35 : 0),
                            y: y - (panel.scrollAxis === "y" ? step * 35 : 0),
                        },
                    ],
                });
            await touch.send("Input.dispatchTouchEvent", {
                type: "touchEnd",
                touchPoints: [],
            });
            await page.waitForFunction(
                (axis) =>
                    document.getElementById("viewport")[
                        axis === "x" ? "scrollLeft" : "scrollTop"
                    ] > 40,
                {},
                panel.scrollAxis,
            );
        }
        await page.screenshot({
            path: path.join(output, `${name}.mobile-scrolled.png`),
            fullPage: true,
        });
        await touch.detach();
        await page.click("#fit");
        const fitted = await page.evaluate(() => window.diagramView.inspect());
        checkView(fitted, `${name} mobile overview`, false);
        assert.equal(
            fitted.canvas.width,
            fitted.container.width,
            `${name}: overview fits phone width`,
        );
        if (name === "mindmap") {
            const height =
                Math.max(...fitted.sceneCorners.map((point) => point[1])) -
                Math.min(...fitted.sceneCorners.map((point) => point[1]));
            assert.ok(
                fitted.container.height <= height + 120,
                "Phone overview avoids an oversized empty panel",
            );
        }
        await page.screenshot({
            path: path.join(output, `${name}.mobile-fit.png`),
            fullPage: true,
        });
        await page.setViewport({ width: 1440, height: 1050 });
        await page.waitForFunction(
            () => window.diagramView.inspect().container.width > 390,
        );
        checkView(
            await page.evaluate(() => window.diagramView.inspect()),
            `${name} resized overview`,
            false,
        );
        assert.equal(
            await page.evaluate(
                () =>
                    document.documentElement.scrollWidth <= window.innerWidth,
            ),
            true,
            `${name} responsive controls`,
        );
        assert.deepEqual(errors, [], `${name} browser errors`);
        assert.deepEqual(network, [], `${name} remote requests`);
        report.push({
            name,
            nodes: model.nodes.length,
            edges: model.edges.length,
            curvedEdges,
            groups: model.groups.length,
            htmlBytes: html.length,
            drawCalls: initial.drawCalls,
            renderedFonts,
            sha256: createHash("sha256").update(html).digest("hex"),
            passed: true,
        });
        await page.close();
    }
} finally {
    await browser.close();
    await fs.writeFile(
        path.join(output, "verification.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
}
