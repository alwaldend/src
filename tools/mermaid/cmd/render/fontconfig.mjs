import fs from "node:fs/promises";
import path from "node:path";

// Families Mermaid diagrams request through the theme, mapped to the pinned
// font files. Resolving them here keeps label metrics identical on every host
// instead of following whatever Fontconfig the host provides.
const aliases = {
    Helvetica: "Liberation Sans",
    Arial: "Liberation Sans",
    "sans-serif": "Liberation Sans",
    Times: "Liberation Serif",
    "Times New Roman": "Liberation Serif",
    serif: "Liberation Serif",
    Courier: "Liberation Mono",
    monospace: "Liberation Mono",
};

const escapeXml = (value) =>
    value
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;");

/**
 * Writes a Fontconfig file that exposes only the repository's pinned fonts.
 *
 * @param {string[]} fontDirectories Absolute directories holding pinned fonts.
 *     They are the sole font sources, so no family resolves from the host.
 * @param {string} scratch Directory that receives the generated file.
 * @returns {Promise<{file: string, dir: string}>} Absolute paths to pass to
 *     Chrome through FONTCONFIG_FILE and FONTCONFIG_PATH.
 */
export async function writeFontConfig(fontDirectories, scratch) {
    if (fontDirectories.length === 0)
        throw new Error("at least one pinned font directory is required");
    const file = path.join(scratch, "fonts.conf");
    await fs.writeFile(
        file,
        `<?xml version="1.0"?><fontconfig>${fontDirectories
            .map((directory) => `<dir>${escapeXml(directory)}</dir>`)
            .join("")}<cachedir>${escapeXml(
            path.join(scratch, "font-cache"),
        )}</cachedir>${Object.entries(aliases)
            .map(
                ([family, target]) =>
                    `<alias><family>${family}</family><prefer><family>${target}</family></prefer></alias>`,
            )
            .join("")}</fontconfig>`,
    );
    return { file, dir: scratch };
}
