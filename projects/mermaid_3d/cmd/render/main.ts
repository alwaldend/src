import type { MermaidConfig } from "mermaid/dist/config.type.js";
import type {
    ColorScheme,
    Measurements,
    Palette,
    ParsedDiagram,
    GraphData,
    RenderPayload,
    Point3,
} from "../../internal/types.js";
import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";
import {
    adaptDiagram,
    diagramText,
    supportedDiagramTypes,
} from "../../internal/model.js";
import { renderDocument } from "../../internal/viewer/document.js";

// Only this public graph-layout boundary crosses into the browser.
declare const mermaid: {
    initialize(config: MermaidConfig): void;
    detectType(source: string, config?: MermaidConfig): string;
    registerLayoutLoaders(
        loaders: {
            name: string;
            loader(): Promise<{ render(data: GraphData): Promise<void> }>;
        }[],
    ): void;
    render(id: string, source: string): Promise<unknown>;
};

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer") as typeof import("puppeteer").default;
const flags: Record<string, string> = {};
for (let i = 2; i < process.argv.length; i += 2) {
    if (!process.argv[i].startsWith("--") || process.argv[i + 1] === undefined)
        throw new Error("Expected --name value arguments");
    flags[process.argv[i].slice(2)] = process.argv[i + 1];
}
for (const name of [
    "input",
    "output",
    "browser",
    "config",
    "palette",
    "font",
    "font-license",
])
    if (!flags[name]) throw new Error(`Missing --${name}`);
const viewDirection = new Map<string, Point3>([
    ["front-right", [0.28, 0.2, 1.8]],
    ["front-left", [-0.28, 0.1, 1.8]],
]).get(flags.view ?? "front-right");
if (!viewDirection) throw new Error(`Unsupported view: ${flags.view}`);
const source = await fs.readFile(path.resolve(flags.input), "utf8");
const config: MermaidConfig = JSON.parse(
    await fs.readFile(path.resolve(flags.config), "utf8"),
);
const appearances: Record<ColorScheme, Palette> = JSON.parse(
    await fs.readFile(path.resolve(flags.palette), "utf8"),
);
const font = await fs.readFile(path.resolve(flags.font));
const fontLicense = await fs.readFile(
    path.resolve(flags["font-license"]),
    "utf8",
);
const fontFamily = config.fontFamily ?? config.themeVariables?.fontFamily;
if (!fontFamily) throw new Error("The theme must declare its font family");
const fontName = fontFamily
    .split(",")[0]
    .trim()
    .replace(/^['"]|['"]$/g, "");
const fontCSS = `@font-face{font-family:${JSON.stringify(fontName)};src:url(data:font/ttf;base64,${font.toString("base64")}) format('truetype');font-display:block;}`;
const browser = await puppeteer.launch({
    executablePath: path.resolve(flags.browser),
    headless: "shell",
});
let raw: ParsedDiagram;
let measurements: Measurements;
try {
    const page = await browser.newPage();
    await page.addStyleTag({ content: fontCSS });
    await page.evaluate(async (family) => {
        await document.fonts.load(`16px ${family}`);
        await document.fonts.ready;
    }, JSON.stringify(fontName));
    for (const mode of ["light", "dark"] as const)
        if (!appearances[mode]?.background)
            throw new Error(`Palette must define ${mode}.background`);
    await page.addScriptTag({
        path: require.resolve("mermaid/dist/mermaid.js"),
    });
    raw = await page.evaluate(
        async ({ source, config, palette, diagramTypes }) => {
            // Normalize display fields once, before measurement. Mermaid's
            // parser protects entity references with these version-pinned tokens;
            // the browser decodes entities and strips markup in one inert parse.
            const plain = (value: unknown): string => {
                const template = document.createElement("template");
                template.innerHTML = (
                    Array.isArray(value)
                        ? value.join("\n")
                        : String(value ?? "")
                )
                    .replace(/<br\s*\/?\s*>/gi, "\n")
                    .replace(/\\([+\-#~])/g, "$1")
                    .replace(/ﬂ°°/g, "&#")
                    .replace(/ﬂ°/g, "&")
                    .replace(/¶ß/g, ";");
                return template.content.textContent!;
            };
            mermaid.initialize({
                ...config,
                startOnLoad: false,
                securityLevel: "strict",
                layout: "mermaid_3d",
                // Source directives must not replace the 3D layout or these
                // build-time safety settings. Other diagram options still apply.
                secure: [
                    ...(config.secure ?? []),
                    "layout",
                    "secure",
                    "securityLevel",
                    "startOnLoad",
                    "maxTextSize",
                    "maxEdges",
                    "suppressErrorRendering",
                ],
                themeVariables: { ...config.themeVariables, ...palette },
            });
            const type = mermaid.detectType(source, config);
            if (!diagramTypes.includes(type))
                throw new Error(`Unsupported Mermaid diagram type: ${type}`);
            let parsed: ParsedDiagram | undefined;
            mermaid.registerLayoutLoaders([
                {
                    name: "mermaid_3d",
                    loader: async () => ({
                        render: async (data) => {
                            parsed = {
                                type,
                                direction: data.direction,
                                nodes: data.nodes.map((node) => ({
                                    id: node.id,
                                    label: plain(
                                        (node.shape === "classBox"
                                            ? node.text
                                            : null) ??
                                            node.label ??
                                            node.id,
                                    ),
                                    shape: node.shape,
                                    parentId: node.parentId,
                                    isGroup: node.isGroup,
                                    dir: node.dir,
                                    annotations: node.annotations?.map(plain),
                                    members: node.members?.map((member) => ({
                                        text: plain(member.text),
                                    })),
                                    methods: node.methods?.map((method) => ({
                                        text: plain(method.text),
                                    })),
                                    attributes: node.attributes?.map(
                                        (attribute) => ({
                                            type: plain(attribute.type),
                                            name: plain(attribute.name),
                                            keys: attribute.keys,
                                        }),
                                    ),
                                })),
                                edges: data.edges.map((edge) => ({
                                    id: edge.id,
                                    start: edge.start,
                                    end: edge.end,
                                    label: plain(edge.label),
                                    startLabelRight: plain(
                                        edge.startLabelRight,
                                    ),
                                    startLabelLeft: plain(edge.startLabelLeft),
                                    endLabelLeft: plain(edge.endLabelLeft),
                                    endLabelRight: plain(edge.endLabelRight),
                                    arrowTypeStart: edge.arrowTypeStart,
                                    arrowTypeEnd: edge.arrowTypeEnd,
                                    pattern: edge.pattern,
                                    thickness: edge.thickness,
                                    isLayoutOnly: edge.isLayoutOnly,
                                })),
                            };
                        },
                    }),
                },
            ]);
            // The public SVG lifecycle supplies semantic graph data to our
            // layout. Its temporary SVG is discarded; Three.js paints the viewer.
            await mermaid.render("spatial-diagram", source);
            if (!parsed)
                throw new Error(`Mermaid did not supply graph data: ${type}`);
            return parsed;
        },
        {
            source,
            config,
            palette: appearances.light,
            diagramTypes: supportedDiagramTypes,
        },
    );
    measurements = (await page.evaluate(
        ({ strings, fontFamily }) => {
            const context = document.createElement("canvas").getContext("2d")!;
            return Object.fromEntries(
                ["normal", "bold"].map((weight) => {
                    context.font = `${weight} 16px ${fontFamily}`;
                    return [
                        weight,
                        Object.fromEntries(
                            [...new Set(strings)].map((value) => [
                                value,
                                context.measureText(value).width,
                            ]),
                        ),
                    ];
                }),
            );
        },
        { strings: diagramText(raw), fontFamily },
    )) as Measurements;
} finally {
    await browser.close();
}

const model = adaptDiagram(raw, measurements, viewDirection);
const title = path.basename(flags.input, ".mmd").replaceAll("-", " ");
const payload: RenderPayload = {
    model,
    appearances,
    fontFamily,
    fontName,
    fontLicense,
    source,
    title,
    mermaidVersion: require("mermaid/package.json").version,
};
const dataURL = (value: string) =>
    `data:text/javascript;base64,${Buffer.from(value).toString("base64")}`;
const threeRoot = path.dirname(path.dirname(require.resolve("three")));
const core = dataURL(
    await fs.readFile(path.join(threeRoot, "build/three.core.js"), "utf8"),
);
const module = (
    await fs.readFile(path.join(threeRoot, "build/three.module.js"), "utf8")
).replaceAll("./three.core.js", "three/core");
const controls = await fs.readFile(
    path.join(threeRoot, "examples/jsm/controls/OrbitControls.js"),
    "utf8",
);
const importMap: { imports: Record<string, string> } = {
    imports: {
        three: dataURL(module),
        "three/core": core,
        "three/addons/controls/OrbitControls.js": dataURL(controls),
        "three/addons/utils/BufferGeometryUtils.js": dataURL(
            await fs.readFile(
                path.join(
                    threeRoot,
                    "examples/jsm/utils/BufferGeometryUtils.js",
                ),
                "utf8",
            ),
        ),
    },
};
for (const [name, file] of [
    ["diagram/text", "../../internal/text.js"],
    ["diagram/viewer", "../../internal/viewer/main.js"],
    ["diagram/scene", "../../internal/viewer/scene.js"],
    ["diagram/labels", "../../internal/viewer/labels.js"],
])
    importMap.imports[name] = dataURL(
        await fs.readFile(new URL(file, import.meta.url), "utf8"),
    );
const styles = await fs.readFile(
    new URL("../../internal/viewer/viewer.css", import.meta.url),
    "utf8",
);
const html = renderDocument({ title, fontCSS, styles, importMap, payload });
await fs.writeFile(path.resolve(flags.output), html);
if (flags.model)
    await fs.writeFile(
        path.resolve(flags.model),
        JSON.stringify(payload, null, 2) + "\n",
    );
