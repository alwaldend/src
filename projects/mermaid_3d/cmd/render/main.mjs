import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";
import { adaptDiagram, diagramText } from "../../internal/model.mjs";
import { renderDocument } from "../../internal/viewer/document.mjs";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const flags = {};
for (let i = 2; i < process.argv.length; i += 2) {
    if (!process.argv[i].startsWith("--") || process.argv[i + 1] === undefined)
        throw new Error("Expected --name value arguments");
    flags[process.argv[i].slice(2)] = process.argv[i + 1];
}
for (const name of ["input", "output", "browser", "config", "palette", "font"])
    if (!flags[name]) throw new Error(`Missing --${name}`);
const source = await fs.readFile(path.resolve(flags.input), "utf8");
const config = JSON.parse(
    await fs.readFile(path.resolve(flags.config), "utf8"),
);
const appearances = JSON.parse(
    await fs.readFile(path.resolve(flags.palette), "utf8"),
);
for (const mode of ["light", "dark"])
    if (
        typeof appearances?.[mode]?.background !== "string" ||
        !appearances[mode].background.trim()
    )
        throw new Error(`Palette must define a background color for ${mode}`);
const font = await fs.readFile(path.resolve(flags.font));
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
let raw;
let measurements;
try {
    const page = await browser.newPage();
    await page.addStyleTag({ content: fontCSS });
    await page.evaluate(async (family) => {
        await document.fonts.load(`16px ${family}`);
        await document.fonts.ready;
    }, JSON.stringify(fontName));
    await page.addScriptTag({
        path: require.resolve("mermaid/dist/mermaid.js"),
    });
    raw = await page.evaluate(
        async ({ source, config, palette }) => {
            mermaid.initialize({
                ...config,
                startOnLoad: false,
                securityLevel: "strict",
                themeVariables: { ...config.themeVariables, ...palette },
            });
            // parse() only validates. Pin this small database adapter and exercise
            // every supported diagram on each Mermaid upgrade.
            const diagram =
                await mermaid.mermaidAPI.getDiagramFromText(source);
            const db = diagram.db;
            const values = {};
            for (const name of [
                "getData",
                "getDirection",
                "getActors",
                "getActorKeys",
                "getMessages",
                "getBoxes",
                "getCreatedActors",
                "getDestroyedActors",
            ])
                if (typeof db[name] === "function") values[name] = db[name]();
            if (db.LINETYPE) values.lineTypes = db.LINETYPE;
            if (values.getData) delete values.getData.config;
            const raw = JSON.parse(
                JSON.stringify(
                    { type: diagram.type, values },
                    (_key, value) =>
                        value instanceof Map
                            ? Object.fromEntries(value)
                            : value,
                ),
            );
            // Normalize display fields once, before measurement. Mermaid's
            // parser protects entity references with these version-pinned tokens;
            // the browser decodes entities and strips markup in one inert parse.
            const plain = (value) => {
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
                return template.content.textContent;
            };
            for (const node of raw.values.getData?.nodes ?? []) {
                node.label = plain(
                    (node.shape === "classBox" ? node.text : null) ??
                        node.label ??
                        node.id,
                );
                node.annotations = (node.annotations ?? []).map(plain);
                for (const member of [
                    ...(node.members ?? []),
                    ...(node.methods ?? []),
                ])
                    member.text = plain(member.text);
                for (const attribute of node.attributes ?? []) {
                    attribute.type = plain(attribute.type);
                    attribute.name = plain(attribute.name);
                }
            }
            for (const edge of raw.values.getData?.edges ?? [])
                for (const key of [
                    "label",
                    "startLabelRight",
                    "startLabelLeft",
                    "endLabelLeft",
                    "endLabelRight",
                ])
                    if (edge[key] != null) edge[key] = plain(edge[key]);
            for (const actor of Object.values(raw.values.getActors ?? {}))
                actor.description = plain(actor.description ?? actor.name);
            for (const event of raw.values.getMessages ?? [])
                event.message = plain(event.message);
            return raw;
        },
        { source, config, palette: appearances.light },
    );
    measurements = await page.evaluate(
        ({ strings, fontFamily }) => {
            const context = document.createElement("canvas").getContext("2d");
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
    );
} finally {
    await browser.close();
}

const model = adaptDiagram(raw, measurements);
const title = path.basename(flags.input, ".mmd").replaceAll("-", " ");
const payload = {
    model,
    appearances,
    fontFamily,
    fontName,
    source,
    title,
    mermaidVersion: require("mermaid/package.json").version,
};
const dataURL = (value) =>
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
const importMap = {
    imports: {
        three: dataURL(module),
        "three/core": core,
        "three/controls": dataURL(controls),
        "three/geometry": dataURL(
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
    ["diagram/text", "../../internal/text.mjs"],
    ["diagram/viewer", "../../internal/viewer/main.mjs"],
    ["diagram/scene", "../../internal/viewer/scene.mjs"],
    ["diagram/labels", "../../internal/viewer/labels.mjs"],
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
