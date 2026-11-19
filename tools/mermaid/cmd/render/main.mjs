import fs from "node:fs/promises";
import path from "node:path";
import { run } from "@mermaid-js/mermaid-cli";
import { embedFont } from "./embed_font.mjs";
import { writeFontConfig } from "./fontconfig.mjs";
import { liftClusterLabels } from "./paint_order.mjs";

// Switches carry no value; every other flag may be repeated, and every
// occurrence accumulates so one flag can carry a list of inputs, such as
// several pinned font files.
const switches = new Set(["plain"]);

const parseFlags = (argv) => {
    const flags = {};
    for (let index = 0; index < argv.length; ) {
        const name = argv[index];
        if (!name.startsWith("--"))
            throw new Error(`Unexpected argument: ${name}`);
        const key = name.slice(2);
        if (switches.has(key)) {
            flags[key] = true;
            index += 1;
            continue;
        }
        const value = argv[index + 1];
        if (value === undefined || value.startsWith("--"))
            throw new Error(`Flag ${name} needs a value`);
        if (Array.isArray(flags[key])) flags[key].push(value);
        else if (flags[key] !== undefined) flags[key] = [flags[key], value];
        else flags[key] = value;
        index += 2;
    }
    return flags;
};

const asList = (value) =>
    value === undefined ? [] : Array.isArray(value) ? value : [value];

const flags = parseFlags(process.argv.slice(2));
// A plain render asks Mermaid for its own defaults: no theme, no post-processing,
// and no font pinning, so it shows the upstream appearance rather than this
// repository's. Input, output, and browser are still required.
for (const required of ["input", "output", "browser"])
    if (flags[required] === undefined)
        throw new Error(`Missing required flag --${required}`);
const plain = flags.plain === true;
if (!plain && flags.config === undefined)
    throw new Error("Missing required flag --config");

// Resolve against the process working directory. The action runs in the
// execroot and the launcher may supply paths relative to it, while Chrome is a
// separate process that needs absolute paths to open the same files.
const input = path.resolve(flags.input);
const output = path.resolve(flags.output);
const browser = path.resolve(flags.browser);
const config =
    flags.config === undefined ? undefined : path.resolve(flags.config);
// Each directory holds only pinned font files, and they become the sole font
// sources, so no family can resolve from the host.
const fontDirectories = asList(flags["font-directory"]).map((directory) =>
    path.resolve(directory),
);
// The handwriting face is embedded in the output so the document carries the
// glyphs its geometry was measured with; see embed_font.mjs.
const embedFontFile =
    flags.font === undefined ? undefined : path.resolve(flags.font);
const mermaidConfig = plain
    ? undefined
    : JSON.parse(await fs.readFile(config, "utf8"));

// Render scratch lives beside the action's working directory and is removed
// below, so a rendered diagram never depends on host Fontconfig state.
const scratch = await fs.mkdtemp(path.resolve(".mermaid-render-"));
try {
    // A plain render leaves Fontconfig alone, so Chrome resolves whatever the
    // host provides and the output is not hermetic.
    let fontConfig;
    if (!plain && fontDirectories.length > 0)
        fontConfig = await writeFontConfig(fontDirectories, scratch);
    await run(input, output, {
        puppeteerConfig: {
            executablePath: browser,
            headless: "shell",
            // Puppeteer replaces the browser environment when `env` is set,
            // so inherited values have to be preserved explicitly.
            ...(fontConfig === undefined
                ? {}
                : {
                      env: {
                          ...process.env,
                          FONTCONFIG_FILE: fontConfig.file,
                          FONTCONFIG_PATH: fontConfig.dir,
                      },
                  }),
        },
        parseMMDOptions: {
            backgroundColor: flags.background ?? "white",
            ...(mermaidConfig === undefined ? {} : { mermaidConfig }),
            viewport: { width: 800, height: 600, deviceScaleFactor: 1 },
        },
    });
    // Mermaid paints container titles before the edges, so an edge that
    // crosses a container border draws over the title plate. Lifting the
    // plates to the end of the root group fixes the paint order without
    // moving any geometry.
    if (!plain) {
        const rendered = await fs.readFile(output, "utf8");
        const lifted = liftClusterLabels(rendered);
        await fs.writeFile(
            output,
            embedFontFile === undefined
                ? lifted
                : embedFont(lifted, await fs.readFile(embedFontFile)),
        );
    }
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
