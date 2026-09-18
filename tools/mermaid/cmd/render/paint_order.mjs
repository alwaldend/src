// Mermaid groups a container's title in `<g class="cluster-label">` inside
// `<g class="clusters">`, and emits that group before the edges, the edge
// labels, and the nodes. Document order is the only layer control SVG honours
// -- `z-index` needs a CSS-positioned box and does not reorder SVG elements,
// and `paint-order` only orders fill, stroke, and markers within one element --
// so an edge crossing a container border paints over the title plate and a
// line appears to run through the title text.
//
// Moving each title group to the end of the root group puts the plates above
// the edges. A title's transform is relative to the root group, which has no
// transform of its own, so the lift moves no geometry.

// matchEnd returns the index just past the group that opens at `start`.
//
// Counting nested `<g` and `</g>` is exact for these documents: the renderer
// writes no markup as an attribute value, as `TestGroupScanIsExact` checks, so
// the angle brackets scanned here always begin a tag.
const matchEnd = (svg, start, name) => {
    const open = `<${name}`;
    const close = `</${name}>`;
    let depth = 0;
    for (let index = start; index < svg.length; ) {
        if (svg.startsWith(open, index)) {
            depth += 1;
            index += open.length;
        } else if (svg.startsWith(close, index)) {
            depth -= 1;
            index += close.length;
            if (depth === 0) return index;
        } else {
            index += 1;
        }
    }
    return -1;
};

// findRoot returns the span of the untransformed `<g class="root">` group,
// which holds every paint layer of a rendered flowchart.
const findRoot = (svg) => {
    const marker = '<g class="root">';
    const start = svg.indexOf(marker);
    if (start < 0) return null;
    const end = matchEnd(svg, start, "g");
    return end < 0
        ? null
        : { start, end, contentStart: start + marker.length };
};

// The title group carries `cluster-label`, and Mermaid's own markup writes the
// class with a trailing space that a renderer may or may not trim, so the
// matcher accepts either spelling.
const titlePattern = /<g class="cluster-label(?:\s|")/g;

// liftClusterLabels returns the document with every container title painted
// last. A document without a root group is returned unchanged, as is one with
// no container title.
export const liftClusterLabels = (svg) => {
    const root = findRoot(svg);
    if (root === null) return svg;

    const lifted = [];
    let body = "";
    let cursor = root.contentStart;
    titlePattern.lastIndex = root.contentStart;
    for (let match = titlePattern.exec(svg); ; ) {
        // A title is a direct child of its cluster, so the search never has to
        // leave the range the root group already bounds.
        if (match === null || match.index > root.end) break;
        const end = matchEnd(svg, match.index, "g");
        if (end < 0) break;
        body += svg.slice(cursor, match.index);
        lifted.push(svg.slice(match.index, end));
        cursor = end;
        titlePattern.lastIndex = end;
        match = titlePattern.exec(svg);
    }
    if (lifted.length === 0) return svg;

    body += svg.slice(cursor, root.end - "</g>".length);
    return (
        svg.slice(0, root.contentStart) +
        body +
        lifted.join("") +
        svg.slice(root.end - "</g>".length)
    );
};
