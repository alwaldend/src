# Flat appearance options

The implemented theme retains the neutral site palette. These four alternatives
are design proposals, rendered from the approved preset across all six gallery
types in light and dark modes. They use existing pinned Liberation fonts and
Mermaid's classic look, native shapes, and embedded-font rendering.

**Recommend Teal surfaces.** Its pale node fills separate content from the
canvas, while green connectors and quieter outlines establish a useful visual
hierarchy. Cobalt drafting is the strongest alternative for an architectural
appearance. Warm editorial and Violet technical deliberately explore typography
beyond the study's original sans-serif baseline. The subsequent
[monospace change](../use-mermaid-monospace/proposal.md) preserves these
proposals while changing the active neutral theme's typography.

| Direction        | Treatment                                                                                                                   | Evidence in the samples                                                                                                                                   | Trade-off                                                                                                                                      |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Teal surfaces    | Sans 16px; 14px captions; green connectors; pale teal nodes; 6px ordinary rectangle corners; 8px containers; basis curves   | Flowchart nodes separate from their canvas. Mindmap branches and labels remain clear in dark mode. Class and ER tables share the same restrained palette. | Filled tables have more visual weight. Keep large containers canvas-filled.                                                                    |
| Cobalt drafting  | Sans 16px; 14px captions; blue connectors; cool neutral surfaces; 2px ordinary rectangles; square containers; linear routes | Blue class relations stand apart from quieter compartment outlines. The flowchart looks precise and architectural.                                        | Angular retry paths feel busier. Linear follows layout waypoints; it does not guarantee orthogonal routes.                                     |
| Warm editorial   | Serif 17px; 15px captions; rust connectors; ivory/umber canvases; square ordinary rectangles; basis curves                  | Sequence and state diagrams resemble book illustrations. Warm mindmap surfaces give a coherent grouping hierarchy.                                        | Serif method signatures and ER attributes feel less like code. Best suited to explanatory articles.                                            |
| Violet technical | Mono 15px; 14px captions; violet surfaces/connectors; 2px ordinary rectangles; square containers; linear routes             | Class methods gain a consistent monospaced rhythm, and the dark palette has a distinctive technical character.                                            | The class sample expands to 1326px. Its normal paths occupy only 0.96 CSS px at 728px; test 1.9px paths or provide more width before adoption. |

Typography and geometry inputs are in [appearance-presets.json](appearance-presets.json).
Their `colorScheme` IDs reference the maintained
[color catalog](../../../pkg/mermaid/color-schemes.json), whose light/dark
colors have named canvas, surface, ink, outline, accent, and muted-border roles.
These are original design choices, not claims of measured user preference.

## Render interpretation

- Apply each proposal over the shared approved preset. Preserve the 1.75px
  normal-edge treatment, exclusions for thick/invisible edges, 1px gray sequence
  guides, caption padding/margins, and explicit mindmap styling.
- Use ink for text, accent for graph/message/branch connectors and container
  titles, outline for node borders, muted border for containers and captions,
  and canvas for caption fills. Use surface for nodes, notes, and mindmaps;
  Warm editorial uses canvas for ordinary nodes.
- Resolve native palette variables for sequence actors, messages, activation
  bars, notes, and composite-state surfaces. Sequence styles need concrete
  palette colors because their SVG lacks the graph's `.root` variable scope.
- Match `fontFamily` and the embedded `Liberation{Sans,Serif,Mono}-Regular.ttf`
  face. Set sequence actor/message/note fonts explicitly. Apply the listed
  caption/title metrics and rectangle radii through `themeCSS`.
- Preserve path-defined semantic shapes. State/class shapes and terminals do
  not all follow ordinary rectangle radius rules. Mindmap layout retains its
  native branch curves. Avoid rounding ER attribute rows.

All 48 proposal SVGs loaded their embedded fonts. Measured connector-to-canvas
contrast was at least 5.55:1; mindmap label-to-surface contrast was at least
10.46:1. These are nominal computed-color ratios, not antialiasing measurements.
Thin strokes can still look weak after scaling, as illustrated by Violet's
class sample. [W3C's guidance](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html)
also distinguishes nominal contrast from the appearance of thin graphics.

ER's native optional-cardinality circles retain a white fill in these renders.
That fill is conspicuous on tinted dark canvases and needs diagram-specific
palette treatment before adopting any colored direction. Semantic markers must
remain hollow or filled according to their meaning.

## Evidence and reproduction

The task-local `out/mermaid-appearance/` contains each proposal's resolved
`theme.json`, `palette.css`, twelve native SVGs, and rendering logs. The
`render-options.mjs` driver uses the checked-in gallery sources and the pinned
renderer action's browser, font directory, and subsetter. It does not rewrite
generated SVGs. `appearance-measurements.json` records all 48 computed checks.

`comparison-standalone.html` embeds all 48 proposals and the twelve implemented
baseline SVGs. Its controls select diagram, mode, and fit/article/native size;
the baseline can be expanded for comparison. Desktop and mobile previews were
inspected, along with every diagram family in both modes. These ignored
artifacts are review evidence, not site assets or a second maintained theme.
