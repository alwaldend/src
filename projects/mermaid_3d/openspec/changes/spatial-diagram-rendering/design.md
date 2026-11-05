## Context

The site owns six Mermaid examples. Shared identity data owns the light/dark
palette and diagram font.
Dependency pins live in the tools workspace; the compiler and Three.js setup
are delivered separately. This project owns the optional spatial renderer.

## Goals / Non-Goals

Render independent solid nodes, group volumes, and spatial connections from
existing Mermaid definitions. Retain source meaning, theme, offline operation,
and usable desktop/phone interaction. Full Mermaid syntax/shape parity,
replacement of existing SVG exports, and website deployment are outside scope.

## Decisions

### Parsing and layout

Use the pinned Mermaid parser and normalize display text before measurement.
Register a public Mermaid layout loader and call its public render API. The
loader receives graph data and projects only the fields needed by the spatial
model. The temporary SVG lifecycle completes without painting graph geometry;
the packaged viewer owns Three.js rendering. Sequence diagrams are excluded
because their rendering path bypasses this graph extension point.

Graph placement respects flow direction and nested group directions, condenses
cycles before ranking, and retains invisible links as layout constraints.
Route from side ports, preferring direct curves and using depth corridors when
nodes obstruct them. Sample curves once for drawing, markers, bounds, and
collision verification. Self-loops stay outside their node face.
Distribute decision ports along each side, then intersect their connection axis
with the diamond face before adding the small endpoint clearance.

Extruding SVG would retain planar routing. A force simulation would disrupt
source order and reproducibility. Deterministic model adaptation supplies the
required spatial behavior without another grammar or layout dependency.

### Rendering and interaction

Share measured text rows between card sizing and label drawing. Use flat-front
diamond prisms, volume outlines for groups, and one tube renderer for paths
and guides. Merge dash sections to keep draw-call cost bounded. Derive quieter
rear strokes from the maintained palette. Edge captions and group titles share
bordered background plates using the theme's border color.

Use readable scrolling plus a fitted overview; one fit-to-screen policy cannot
keep dense diagrams legible on phones. Place annotations with bounded screen
space displacement and leader lines. Frame geometry and anchored label extents,
and fit camera clipping to scene depth whenever the view changes. Size the
drawing buffer in one operation with a pixel ratio bounded by WebGL texture,
renderbuffer, and viewport limits; CSS size and readable scrolling are retained.
The renderer command owns initial front-right/front-left view presets. The
selected direction informs both projected route clearance and the viewer's
saved reset perspective; the architecture target selects front-left.

Element identities include kind and Mermaid ID. Selection expands containment
before adding direct connections, without following further graph relationships.
One registration path owns object identities and picking metadata.

### Build and ownership

Compile ordinary `.ts` sources with strict checking and separate Node/browser
contexts. Keep one scene model type, initializing graph routes before filling
them. Consume compiler declarations during the build and package emitted `.js`
ES modules, project metadata, styles, palette, and embedded font into offline
HTML. One `MermaidThemeInfo` target supplies configuration, JSON appearances,
font, and license; the renderer consumes those files directly. Each import-map entry embeds its module once, including the Three.js core.

Keep command orchestration in `cmd/render`, model/routing/text code in
`internal`, and scene/labels/controller/document code in `internal/viewer`.
Reuse the existing gallery sources. One project-owned architecture definition
feeds the themed SVG embed and interactive HTML on the README and landing.

### Independent snapshots

Core is delivered in PR 211. The complete typed snapshot implementation and
its passing E2E artifacts are preserved separately. Static export, snapshot
preparation, raster website support, and image acceptance belong to the
follow-up after core merges. Core retains interactive geometry/depth fitting.

## Risks / Trade-offs

- Mermaid graph-data changes: exercise every maintained adapter on upgrades.
- Dense scenes and unusual camera angles: retain selection, reset, scrolling,
  and overview; bounded label placement cannot guarantee an optimal layout.
- Route collisions: independently inspect swept segments against unrelated
  nodes and fail if no clear route exists.
- Theme or font drift: load the bundled JSON palettes and verify both
  appearances and actual canvas fonts without remote requests.

## Verification

Existing browser E2Es cover all five supported gallery types, semantic and layout boundary
fixtures, selection opacity, phone interaction, and a fifty-node flowchart.
Preserve screenshots and artifact hashes, compare maintained scene models and
pixels for refactors, and run native Mermaid/site consumers, repository quality,
and affected semantic lint. Delivery receipts own exact publication state.
