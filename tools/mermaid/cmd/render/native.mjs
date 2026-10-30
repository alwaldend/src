import fs from "node:fs/promises";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { execFileSync } from "node:child_process";
import puppeteer from "puppeteer";
import { renderMermaid } from "@mermaid-js/mermaid-cli";
import { fontFace } from "./embed_font.mjs";

// Pass the consumer's palette JSON directly to Mermaid. Generated SVG bytes
// are written without modification.
export async function renderNative({
    input,
    output,
    darkOutput,
    colorScheme = "light",
    palette,
    config,
    font,
    fontLicense,
    fontSubsetter,
    puppeteerConfig,
}) {
    const appearances = JSON.parse(await fs.readFile(palette, "utf8"));
    const scratch = await fs.mkdtemp(path.resolve(".mermaid-native-"));
    let browser;
    try {
        browser = await puppeteer.launch(puppeteerConfig);
        const definition = await fs.readFile(input, "utf8");
        const fontBytes = await fs.readFile(font);
        const cssPath = path.join(scratch, "font.css");
        const fontPath = path.join(scratch, "font.ttf");
        await fs.writeFile(fontPath, fontBytes);
        await fs.writeFile(cssPath, fontFace(fontBytes, "font.ttf"));
        await fs.copyFile(fontLicense, path.join(scratch, "LICENSE"));
        let subsetted = false;
        const variants =
            darkOutput === undefined
                ? [[colorScheme, output]]
                : [
                      ["light", output],
                      ["dark", darkOutput],
                  ];
        for (const [mode, destination] of variants) {
            const themeVariables = appearances?.[mode];
            if (
                typeof themeVariables?.background !== "string" ||
                !themeVariables.background.trim()
            )
                throw new Error(
                    `Palette must define a background color for ${mode}`,
                );
            // Make the same palette available to Mermaid's input
            // theme CSS, so caption borders and fills share its native colors.
            const properties = Object.entries(themeVariables)
                .map(([name, value]) => `--mermaid-${name}:${value}`)
                .join(";");
            const options = {
                mermaidConfig: {
                    ...config,
                    themeCSS: `.root{${properties}}${config.themeCSS ?? ""}`,
                    themeVariables: {
                        ...config.themeVariables,
                        ...themeVariables,
                        darkMode: mode === "dark",
                    },
                },
                backgroundColor: themeVariables.background,
                // CLI 12 embeds custom faces together with their license.
                customFontCSS: [{ cssUrl: pathToFileURL(cssPath) }],
            };
            let { data } = await renderMermaid(
                browser,
                definition,
                "svg",
                options,
            );
            if (!subsetted) {
                // Use actual rendered text, including decoded/generated labels.
                // Re-render with the subset through Mermaid's supported CSS input;
                // delivered SVG geometry is never rewritten.
                await fs.writeFile(
                    fontPath,
                    execFileSync(fontSubsetter, [font], {
                        input: data,
                        maxBuffer: 8 * 1024 * 1024,
                        timeout: 60000,
                    }),
                );
                subsetted = true;
                ({ data } = await renderMermaid(
                    browser,
                    definition,
                    "svg",
                    options,
                ));
            }
            await fs.writeFile(destination, data);
        }
    } finally {
        await browser?.close();
        await fs.rm(scratch, { recursive: true, force: true });
    }
}
