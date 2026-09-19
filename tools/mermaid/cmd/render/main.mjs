import fs from "node:fs/promises";
import path from "node:path";
import { run } from "@mermaid-js/mermaid-cli";
import { writeFontConfig } from "./fontconfig.mjs";
import { liftClusterLabels } from "./paint_order.mjs";

// A flag may be repeated; every occurrence accumulates so one flag can carry a
// list of inputs, such as several pinned font files.
const parseFlags = (argv) => {
    const flags = {};
    for (let index = 0; index < argv.length; index += 2) {
        const name = argv[index];
        if (!name.startsWith("--"))
            throw new Error(`Unexpected argument: ${name}`);
        const value = argv[index + 1];
        if (value === undefined || value.startsWith("--"))
            throw new Error(`Flag ${name} needs a value`);
        const key = name.slice(2);
        if (Array.isArray(flags[key])) flags[key].push(value);
        else if (flags[key] !== undefined) flags[key] = [flags[key], value];
        else flags[key] = value;
    }
    return flags;
};

const asList = (value) =>
    value === undefined ? [] : Array.isArray(value) ? value : [value];

const flags = parseFlags(process.argv.slice(2));
for (const required of [
    "input",
    "output",
    "config",
    "browser",
    "font-directory",
])
    if (flags[required] === undefined)
        throw new Error(`Missing required flag --${required}`);

// Resolve against the process working directory. The action runs in the
// execroot and the launcher may supply paths relative to it, while Chrome is a
// separate process that needs absolute paths to open the same files.
const input = path.resolve(flags.input);
const output = path.resolve(flags.output);
const config = path.resolve(flags.config);
const browser = path.resolve(flags.browser);
// Each directory holds only pinned font files, and they become the sole font
// sources, so no family can resolve from the host.
const fontDirectories = asList(flags["font-directory"]).map((directory) =>
    path.resolve(directory),
);
const mermaidConfig = JSON.parse(await fs.readFile(config, "utf8"));

// Render scratch lives beside the action's working directory and is removed
// below, so a rendered diagram never depends on host Fontconfig state.
const scratch = await fs.mkdtemp(path.resolve(".mermaid-render-"));
try {
    const fontConfig = await writeFontConfig(fontDirectories, scratch);
    await run(input, output, {
        puppeteerConfig: {
            executablePath: browser,
            headless: "shell",
            // Puppeteer replaces the browser environment when `env` is set,
            // so inherited values have to be preserved explicitly.
            env: {
                ...process.env,
                FONTCONFIG_FILE: fontConfig.file,
                FONTCONFIG_PATH: fontConfig.dir,
            },
        },
        parseMMDOptions: {
            backgroundColor: flags.background ?? "white",
            mermaidConfig,
            viewport: { width: 800, height: 600, deviceScaleFactor: 1 },
        },
    });
    // Mermaid paints container titles before the edges, so an edge that
    // crosses a container border draws over the title plate. Lifting the
    // plates to the end of the root group fixes the paint order without
    // moving any geometry.
    const rendered = await fs.readFile(output, "utf8");
    await fs.writeFile(output, liftClusterLabels(rendered));
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
