import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { cli } from "@mermaid-js/mermaid-cli";
import { embedFont } from "../render/embed_font.mjs";
import { writeFontConfig } from "../render/fontconfig.mjs";
import { liftClusterLabels } from "../render/paint_order.mjs";

// A plain run asks Mermaid for its own defaults: no repository theme, no
// paint-order post-processing, and no font pinning. It is the escape hatch for
// comparing against the upstream appearance.
const plain = process.env.MERMAID_PLAIN === "1";

// The launcher passes the pinned font directories through the environment
// because mermaid-cli owns process.argv and would otherwise need a new flag for
// each.
// Each pinned font is named by its own file so the wrapper can derive the
// directory Fontconfig needs. The two families live in separate directories,
// so exposing one never exposes the other.
const fontFiles = plain
    ? []
    : [process.env.MERMAID_BODY_FONT, process.env.MERMAID_LABEL_FONT]
          .filter(Boolean)
          .map((file) => path.resolve(file));
// The label font is embedded in the rendered document so the output keeps its
// measured text metrics without a webfont, matching a build action.
const labelFont = plain ? undefined : process.env.MERMAID_LABEL_FONT;
const fontDirectories = fontFiles.map((file) => path.dirname(file));
if (!plain && fontDirectories.length === 0)
    throw new Error(
        "MERMAID_BODY_FONT and MERMAID_LABEL_FONT must name fonts",
    );
if (new Set(fontDirectories).size !== fontDirectories.length)
    throw new Error("each pinned font must live in its own directory");

// The CLI parses process.argv itself, so a caller-supplied puppeteer config is
// detected here and left alone: an explicit config owns the whole launch,
// including which fonts Chrome may use.
const hasPuppeteerConfig = process.argv.some(
    (argument) =>
        argument === "-p" ||
        argument === "--puppeteerConfigFile" ||
        argument.startsWith("--puppeteerConfigFile="),
);

// The CLI parses its own output flag, so the path is recovered here to apply
// the same paint-order fix a build action applies. A caller can name the
// output as `-o FILE`, `--output FILE`, or `--output=FILE`; `-` writes to
// stdout and has no file to rewrite.
const outputArgument = () => {
    for (let index = 0; index < process.argv.length; index += 1) {
        const argument = process.argv[index];
        if (argument === "-o" || argument === "--output")
            return process.argv[index + 1];
        if (argument.startsWith("--output="))
            return argument.slice("--output=".length);
    }
    return undefined;
};

// The launcher always supplies the repository theme, so a plain run removes it
// from the argument list before the CLI parses one.
if (plain) {
    const stripped = [];
    for (let index = 0; index < process.argv.length; index += 1) {
        const argument = process.argv[index];
        if (argument === "--configFile" || argument === "-c") {
            index += 1;
            continue;
        }
        if (argument.startsWith("--configFile=")) continue;
        stripped.push(argument);
    }
    process.argv = stripped;
}

// A CLI run is interactive, so its scratch directory belongs to the operating
// system rather than to the workspace.
const scratch = await fs.mkdtemp(path.join(os.tmpdir(), "mermaid-fonts-"));
let exitCode = 0;
try {
    if (!plain && !hasPuppeteerConfig) {
        const fontConfig = await writeFontConfig(fontDirectories, scratch);
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
    // Mermaid paints container titles before the edges, so an edge that
    // crosses a container border draws over the title plate. Applying the same
    // lift as the build action keeps an interactive render identical to the
    // committed one.
    const output = outputArgument();
    if (!plain && output !== undefined && output.endsWith(".svg")) {
        const rendered = await fs.readFile(output, "utf8");
        const lifted = liftClusterLabels(rendered);
        await fs.writeFile(
            output,
            labelFont === undefined
                ? lifted
                : embedFont(
                      lifted,
                      await fs.readFile(path.resolve(labelFont)),
                  ),
        );
    }
} catch (exception) {
    console.error(exception instanceof Error ? exception.stack : exception);
    exitCode = 1;
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
process.exitCode = exitCode;
