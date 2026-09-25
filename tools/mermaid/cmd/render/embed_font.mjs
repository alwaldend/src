// A documentation page embeds a rendered diagram as an <img>, and an SVG
// loaded that way is an isolated document: it cannot see the page's webfonts,
// and it loads no network font. Text then falls back to a host font whose
// metrics differ from the pinned face the renderer measured, so a label can
// outgrow the box Mermaid sized for it and appear clipped.
//
// Embedding the pinned face as a data URI makes the document self-contained:
// the font travels with the SVG, so the glyphs and the measured geometry match
// on every consumer, including one that can load no webfont at all.

// The `name` table is the authoritative source for the family the renderer
// must request, so a font swap does not silently stop matching the theme.
const readFamilyName = (data) => {
    const view = new DataView(data.buffer, data.byteOffset, data.byteLength);
    const tableCount = view.getUint16(4);
    for (let index = 0; index < tableCount; index += 1) {
        const record = 12 + index * 16;
        if (
            String.fromCharCode(...data.subarray(record, record + 4)) !==
            "name"
        )
            continue;
        const table = view.getUint32(record + 8);
        const recordCount = view.getUint16(table + 2);
        const stringOffset = table + view.getUint16(table + 4);
        // Windows platform, Unicode BMP, English: the record a renderer
        // expects for a family name, and the one these pinned fonts carry.
        let fallback;
        for (let entry = 0; entry < recordCount; entry += 1) {
            const start = table + 6 + entry * 12;
            const platform = view.getUint16(start);
            const language = view.getUint16(start + 4);
            const nameId = view.getUint16(start + 6);
            if (nameId !== 1) continue;
            const length = view.getUint16(start + 8);
            const offset = stringOffset + view.getUint16(start + 10);
            const bytes = data.subarray(offset, offset + length);
            const text =
                platform === 3 || platform === 0
                    ? new TextDecoder("utf-16be").decode(bytes)
                    : new TextDecoder("latin1").decode(bytes);
            if (platform === 3 && language === 0x409) return text;
            fallback ??= text;
        }
        if (fallback !== undefined) return fallback;
        throw new Error("font carries no family name record");
    }
    throw new Error("font carries no name table");
};

// Native renders pass this stylesheet through the CLI's supported CSS input.
// The historical renderer embeds the same bytes below to preserve its output.
export const fontFace = (font) => {
    const family = readFamilyName(font);
    const encoded = Buffer.from(font).toString("base64");
    return (
        `@font-face{font-family:'${family}';` +
        `src:url(data:font/ttf;base64,${encoded}) format('truetype');}`
    );
};

export const embedFont = (svg, font) => {
    const rule = fontFace(font);
    const style = "<style>";
    const start = svg.indexOf(style);
    if (start < 0) return svg;
    return (
        svg.slice(0, start + style.length) +
        rule +
        svg.slice(start + style.length)
    );
};
