import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { cli } from "@mermaid-js/mermaid-cli";
import { writeFontConfig } from "../render/fontconfig.mjs";

const fontAnchor = process.env.MERMAID_FONT_ANCHOR;
if (!fontAnchor)
    throw new Error("MERMAID_FONT_ANCHOR must name a pinned font file");

// The CLI parses process.argv itself, so a caller-supplied puppeteer config is
// detected here and left alone: an explicit config owns the whole launch,
// including which fonts Chrome may use.
const hasPuppeteerConfig = process.argv.some(
    (argument) =>
        argument === "-p" ||
        argument === "--puppeteerConfigFile" ||
        argument.startsWith("--puppeteerConfigFile="),
);

// A CLI run is interactive, so its scratch directory belongs to the operating
// system rather than to the workspace.
const scratch = await fs.mkdtemp(path.join(os.tmpdir(), "mermaid-fonts-"));
let exitCode = 0;
try {
    if (!hasPuppeteerConfig) {
        const fontConfig = await writeFontConfig(
            path.resolve(fontAnchor),
            scratch,
        );
        const puppeteerConfig = path.join(scratch, "puppeteer.json");
        await fs.writeFile(
            puppeteerConfig,
            JSON.stringify({
                env: {
                    ...process.env,
                    FONTCONFIG_FILE: fontConfig.file,
                    FONTCONFIG_PATH: fontConfig.dir,
                },
            }),
        );
        process.argv.push("--puppeteerConfigFile", puppeteerConfig);
    }
    await cli();
} catch (exception) {
    console.error(exception instanceof Error ? exception.stack : exception);
    exitCode = 1;
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
process.exitCode = exitCode;
