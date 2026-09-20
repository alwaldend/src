// Mermaid groups a container's title in `<g class="cluster-label">` inside
// `<g class="clusters">`, and emits that group before the edges, the edge
// labels, and the nodes. Document order is the only layer control SVG honours
// -- `z-index` needs a CSS-positioned box and does not reorder SVG elements,
// and `paint-order` only orders fill, stroke, and markers within one element --
// so an edge crossing a container border paints over the title plate and a
// line appears to run through the title text.
//
// Moving each title group to the end of the root group puts the plates above
// the edges. Dagre nests translated groups, so each lifted title retains its
// ancestors' transforms in addition to its own.

const groupTags = () => /<\/?g(?=[\s>])[^>]*>/g;
const transformOf = (tag) => tag.match(/\btransform="([^"]*)"/)?.[1] ?? "";

// matchEnd returns the index just past the group that opens at `start`.
//
// Mermaid escapes label text and attribute values. Empty paint layers may be
// self-closing, and must not increase the nesting depth.
const matchEnd = (svg, start) => {
    const tags = groupTags();
    tags.lastIndex = start;
    let depth = 0;
    for (let match = tags.exec(svg); match; match = tags.exec(svg)) {
        const tag = match[0];
        if (tag.startsWith("</")) depth -= 1;
        else if (!tag.endsWith("/>")) depth += 1;
        if (depth === 0) return tags.lastIndex;
    }
    return -1;
};

// findRoot returns the span of the outermost `<g class="root">` group,
// which holds every paint layer of a rendered flowchart.
const findRoot = (svg) => {
    const marker = '<g class="root">';
    const start = svg.indexOf(marker);
    if (start < 0) return null;
    const end = matchEnd(svg, start);
    return end < 0
        ? null
        : { start, end, contentStart: start + marker.length };
};

// The title group carries `cluster-label`, and Mermaid's own markup writes the
// class with a trailing space that a renderer may or may not trim, so the
// matcher accepts either spelling.
const titlePattern = /\bclass="cluster-label(?:\s|")/;

const retainTransforms = (title, ancestors) => {
    const inherited = ancestors.filter(Boolean).join(" ");
    if (!inherited) return title;
    const openingEnd = title.indexOf(">");
    const opening = title.slice(0, openingEnd);
    const own = transformOf(opening);
    const transform = [inherited, own].filter(Boolean).join(" ");
    const updated = own
        ? opening.replace(/\btransform="[^"]*"/, `transform="${transform}"`)
        : opening.replace(/\btransform=""/, "") + ` transform="${transform}"`;
    return updated + title.slice(openingEnd);
};

// liftClusterLabels returns the document with every container title painted
// last. A document without a root group is returned unchanged, as is one with
// no container title.
export const liftClusterLabels = (svg) => {
    const root = findRoot(svg);
    if (root === null) return svg;

    const lifted = [];
    let body = "";
    let cursor = root.contentStart;
    const ancestors = [];
    const tags = groupTags();
    tags.lastIndex = root.contentStart;
    for (let match = tags.exec(svg); match; match = tags.exec(svg)) {
        if (match.index >= root.end - "</g>".length) break;
        const tag = match[0];
        if (tag.startsWith("</")) {
            ancestors.pop();
        } else if (titlePattern.test(tag)) {
            const end = matchEnd(svg, match.index);
            if (end < 0) return svg;
            body += svg.slice(cursor, match.index);
            lifted.push(
                retainTransforms(svg.slice(match.index, end), ancestors),
            );
            cursor = end;
            tags.lastIndex = end;
        } else if (!tag.endsWith("/>")) {
            ancestors.push(transformOf(tag));
        }
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
