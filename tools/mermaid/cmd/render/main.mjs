import fs from "node:fs/promises";
import path from "node:path";
import { run } from "@mermaid-js/mermaid-cli";
import { writeFontConfig } from "./fontconfig.mjs";

const parseFlags = (argv) => {
    const flags = {};
    for (let index = 0; index < argv.length; index += 2) {
        const name = argv[index];
        if (!name.startsWith("--"))
            throw new Error(`Unexpected argument: ${name}`);
        flags[name.slice(2)] = argv[index + 1];
    }
    return flags;
};

const flags = parseFlags(process.argv.slice(2));
for (const required of ["input", "output", "config", "browser", "font-anchor"])
    if (!flags[required] || flags[required].startsWith("--"))
        throw new Error(`Missing required flag --${required}`);

// Resolve against the process working directory. The action runs in the
// execroot and the launcher may supply paths relative to it, while Chrome is a
// separate process that needs absolute paths to open the same files.
const input = path.resolve(flags.input);
const output = path.resolve(flags.output);
const config = path.resolve(flags.config);
const browser = path.resolve(flags.browser);
const fontAnchor = path.resolve(flags["font-anchor"]);
const mermaidConfig = JSON.parse(await fs.readFile(config, "utf8"));

console.error("[diag] cwd=" + process.cwd());
console.error("[diag] anchor=" + flags["font-anchor"]);
console.error("[diag] resolved=" + fontAnchor);
console.error(
    "[diag] exists=" +
        (await fs
            .stat(fontAnchor)
            .then(() => "yes")
            .catch(() => "NO")),
);
console.error("[diag] RUNFILES=" + (process.env.RUNFILES ?? "(unset)"));
console.error("[diag] dirname=" + path.dirname(fontAnchor));

// Render scratch lives beside the action's working directory and is removed
// below, so a rendered diagram never depends on host Fontconfig state.
const scratch = await fs.mkdtemp(path.resolve(".mermaid-render-"));
try {
    const fontConfig = await writeFontConfig(fontAnchor, scratch);
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
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
