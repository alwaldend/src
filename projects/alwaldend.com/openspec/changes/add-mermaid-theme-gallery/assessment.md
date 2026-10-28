# Flat Mermaid theme readability

Assessed on 2026-10-28 from base revision
`a1c368565d7e2ea6affa4a734a51cd3f53e6cbe5`, using the six sources in
[`docs/mermaid-theme`](../../../docs/mermaid-theme/README.md).
The shared preset, palette, and renderer were unchanged at that baseline. The
alternatives below were rendered separately under ignored
`out/mermaid-readability/`. The approved combination is implemented in the
follow-on [readability change](../improve-mermaid-theme-readability/proposal.md).

## Findings

The current flowchart, state, class, and ER connectors already use strong
contrast: `#212121` on white and `#ededed` on `#212121`, nominally 16.10:1 and
13.75:1. Their normal strokes are only 1px at native size. Downscaling reduces
their apparent weight; changing the ink color alone would miss that problem.

| Sample               | Native width | Scale at 728px | What the example exposes                                                    |
| -------------------- | -----------: | -------------: | --------------------------------------------------------------------------- |
| Flowchart            |       1192px |            61% | Thin retry/return curves compete with boxed captions                        |
| Sequence             |        871px |            84% | Bright lifelines compete with message arrows                                |
| State                |        699px |           100% | Transition captions and composite boundaries have similar emphasis          |
| Class                |       1204px |            60% | Fine relationships, hollow markers, and 11px multiplicities shrink together |
| Entity relationships |        956px |            76% | Crow's feet and dashed relationships need more weight than table rules      |
| Mindmap              |        781px |            93% | Dark branches are black on charcoal; the root label is also dark            |

Scale assumes a maximum 728px image width without upscaling. Image padding and
the actual article column can reduce it further. At 360px the flowchart's
native 1px stroke is about 0.30px and its 16px text about 4.8px: no restrained
stroke adjustment makes that complete diagram comfortably readable on a phone.
The gallery links every image to its full-size SVG.

The mindmap is a separate palette defect, confirmed in the rendered output.
Its branch colors come from Mermaid's derived `cScale` values, rather than
the site's `lineColor`. The first dark branches are black on `#212121`
(approximately 1.30:1 nominal contrast). This was visible in both the SVG
stylesheet and the preview.

## Proposed improvements

1. **Give normal graph connectors 1.75px strokes.** The flowchart's `Revise`
   and `Retry` routes become easier to follow at article width, and class/ER
   relationships separate more clearly from their table rules. Keep node
   outlines thin. Scope the change to normal flowchart, state, class, and ER
   edges so thick and invisible edges retain their meanings. Keep the current
   dash patterns and marker geometry. At 728px, the flowchart's proposed
   stroke is approximately 1.07px instead of 0.61px. A 1.5px setting would be
   a more delicate alternative, but has not been rendered in this comparison.
2. **Use quieter caption borders.** Retain the opaque canvas fill, text,
   padding, and 3px margins; change only the edge-caption border from ink to
   the container-border token. The flowchart's `Submit`, `Validate`, and
   `Write` captions then read as annotations instead of extra nodes. The
   frame is intentionally subtle; the fill still separates text from lines.
3. **Subordinate sequence lifelines.** Keep messages at their existing 1.5px
   and full ink. Set lifelines to 1px with the existing site's gray-600 in
   light mode and gray-500 in dark mode. In the queue-delivery sample this
   brings the calls and replies forward while preserving participant guides.
   Use the native `actorLineColor` variable: sequence SVGs do not expose the
   `.root` palette scope used by the flowchart's CSS variables.
4. **Make mindmap colors explicit.** Use site secondary surfaces and ink for
   nodes, and 1.75px ink branches in both modes. The proposed dark mindmap
   retains a flat monochrome appearance while making every branch and the
   root label visible. This uniform weight gives up Mermaid's pronounced
   depth-weight hierarchy; indentation and topology still show grouping.
5. **Offer angular routing selectively.** The separate `linear` flowchart
   comparison makes the outgoing and retry routes more geometric. It also
   introduces sharp corners and changes the established smooth style. Keep
   `basis` as the default; use this alternative only if clearer routing is
   preferred for a particular diagram family. It does not remove crossings
   or solve dense topology.

The recommended prototype combines 1–4 and keeps curved routing. It changes
no fonts, node geometry, layout spacing, shadows, gradients, or accent colors.
Inspect markers at natural size before promoting it: their dimensions are
fixed independently of edge width. Native-size viewing and smaller diagrams
remain necessary for dense content, especially class multiplicities.

## Reproduce the prototypes

Build the gallery's `:samples` target first. Derive each temporary configuration
from `projects/alwaldend.com/pkg/mermaid/theme.json`, then append only the
corresponding CSS below to `themeCSS`. Use the same native renderer, compiled
site palette, pinned browser, and embedded Liberation Sans as the baseline;
do not rewrite SVG output. The combined prototype appends all three blocks.

Connector weight:

```css
.flowchart-link.edge-thickness-normal,
.transition.edge-thickness-normal,
.relation.edge-thickness-normal,
.relationshipLine.edge-thickness-normal {
  stroke-width: 1.75px;
}
```

Caption and lifeline hierarchy:

```css
.edgeLabel p {
  border-color: var(--mermaid-clusterBorder);
}
.actor-line {
  stroke-width: 1px;
}
```

For that variant, append `--mermaid-actorLineColor: #666666` in the temporary
palette's light `:root` and `#999999` in its dark media rule. These are
projections of the existing gray-600/gray-500 Sass tokens; implementation
should reference those tokens in the owning palette mixin.

Mindmap appearance:

```css
.mindmap-node rect,
.mindmap-node circle,
.mindmap-node polygon,
.mindmap-node path {
  fill: var(--mermaid-secondaryColor);
  stroke: var(--mermaid-nodeBorder);
  stroke-width: 1px;
}
.mindmap-node .nodeLabel,
.mindmap-node text,
.mindmap-node span {
  color: var(--mermaid-textColor);
  fill: var(--mermaid-textColor);
}
.edge[class*="section-edge-"] {
  stroke: var(--mermaid-lineColor);
  stroke-width: 1.75px;
}
```

The routing-only variant changes `flowchart.curve` from `basis` to `linear`.
All other inputs remain the same. Compare both color modes at identical
widths, then open the SVGs at natural size. The task-local `compare.mjs` and
`variants.json` retain the executed rendering recipe and exact overrides.

## Evidence and limits

All six baselines and the combined prototype were rendered in both modes;
isolated comparisons cover edge weight (flowchart, class, ER), hierarchy
(flowchart, sequence), routing (flowchart), and palette (mindmap). Preview
images show a visual improvement, not a user study or universal readability
guarantee. A theme implementation would need to verify other grammars,
intentional thick/invisible edges, and dense real diagrams before promotion.

Pinned Mermaid 11.17.2 source explains the independent
[normal edge weights](https://github.com/mermaid-js/mermaid/blob/mermaid%4011.17.2/packages/mermaid/src/styles.ts#L79),
[sequence styling](https://github.com/mermaid-js/mermaid/blob/mermaid%4011.17.2/packages/mermaid/src/diagrams/sequence/styles.js#L28),
and [mindmap colors](https://github.com/mermaid-js/mermaid/blob/mermaid%4011.17.2/packages/mermaid/src/diagrams/mindmap/styles.ts#L19).
The current [theming guide](https://mermaid.js.org/config/theming.html) and
[flowchart guide](https://mermaid.js.org/syntax/flowchart.html) describe the
customization surfaces, but their newer defaults do not override the pinned
renderer or the explicit site preset used here.
