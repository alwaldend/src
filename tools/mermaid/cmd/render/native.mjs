import fs from "node:fs/promises";
import { execFileSync } from "node:child_process";
import puppeteer from "puppeteer";
import { renderMermaid } from "@mermaid-js/mermaid-cli";
import { fontFace } from "./embed_font.mjs";

// Resolve the consumer's stylesheet in Chrome, then pass concrete theme
// variables to Mermaid. Generated SVG bytes are written without modification.
export async function renderNative({
    input,
    output,
    darkOutput,
    colorScheme = "light",
    palette,
    config,
    font,
    fontSubsetter,
    puppeteerConfig,
}) {
    const browser = await puppeteer.launch(puppeteerConfig);
    try {
        const page = await browser.newPage();
        await page.addStyleTag({
            content: await fs.readFile(palette, "utf8"),
        });
        const definition = await fs.readFile(input, "utf8");
        let myCSS = fontFace(await fs.readFile(font));
        let subsetted = false;
        const variants =
            darkOutput === undefined
                ? [[colorScheme, output]]
                : [
                      ["light", output],
                      ["dark", darkOutput],
                  ];
        for (const [mode, destination] of variants) {
            await page.emulateMediaFeatures([
                { name: "prefers-color-scheme", value: mode },
            ]);
            const themeVariables = await page.evaluate(() => {
                const style = getComputedStyle(document.documentElement);
                const prefix = "--mermaid-";
                return Object.fromEntries(
                    [...style]
                        .filter((name) => name.startsWith(prefix))
                        .map((name) => [
                            name.slice(prefix.length),
                            style.getPropertyValue(name).trim(),
                        ]),
                );
            });
            if (!themeVariables.background)
                throw new Error("Palette must define --mermaid-background");
            // Make the same resolved palette available to Mermaid's input
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
                myCSS,
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
                myCSS = fontFace(
                    execFileSync(fontSubsetter, [font], {
                        input: data,
                        maxBuffer: 8 * 1024 * 1024,
                        timeout: 60000,
                    }),
                );
                subsetted = true;
                ({ data } = await renderMermaid(browser, definition, "svg", {
                    ...options,
                    myCSS,
                }));
            }
            await fs.writeFile(destination, data);
        }
    } finally {
        await browser.close();
    }
}
