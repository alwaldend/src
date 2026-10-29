// Sizing and drawing consume the same rows, separators, and measured widths.
export function textLayout(
    lines,
    measure,
    { size = 16, heading = false, compartments = [] } = {},
) {
    const lineHeight = size * 1.5;
    const rules = [];
    let y = 8;
    const rows = lines.map((text, index) => {
        if (compartments.includes(index)) {
            rules.push(y + 4);
            y += 8;
        }
        const row = {
            text,
            bold: heading && index === 0,
            y: y + lineHeight / 2,
        };
        y += lineHeight;
        return row;
    });
    return {
        width:
            Math.max(0, ...rows.map((row) => measure(row.text, row.bold))) +
            20,
        height: y + 8,
        size,
        heading,
        rows,
        rules,
    };
}

export function measuredText(lines, measurements, options = {}) {
    return textLayout(
        lines,
        (text, bold) =>
            ((measurements[bold ? "bold" : "normal"]?.[text] ??
                text.length * 9.6) *
                (options.size ?? 16)) /
            16,
        options,
    );
}
