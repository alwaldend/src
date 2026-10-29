---
title: Mermaid 3D
description: Themed spatial views of Mermaid diagrams
---

This project owns the three-dimensional Mermaid renderer, scene layout,
interaction, packaging, and browser verification. Existing diagram definitions
and their themes stay with their consumers; Mermaid dependency pins remain in
the repository's tools pnpm workspace.

Build all six examples and open any generated HTML file in a WebGL 2 browser:

```sh
bazel_agent bazel build //projects/mermaid_3d/samples:html
```

Outputs are in `bazel-bin/projects/mermaid_3d/samples/`. Each `.3d.html` embeds
Three.js, camera controls, the scene, and Liberation Mono. Viewing requires no
server or network requests. Drag to orbit, right-drag to pan, and scroll to zoom.
The controls select light or dark appearance, reset the camera, show a front
view, and inspect nodes, groups, connections, and notes. Labels face the camera.
Selecting a connection highlights it and its endpoints; node selection includes
direct connections, and group selection includes nested members and notes.
Elements remain independently selectable when different kinds share an ID.
The default readable view keeps text at least 12 CSS pixels at reset scale;
large scenes scroll within their panel. **Fit diagram** shows the complete
overview, and **Reset view** restores readable framing. On touch screens, swipe
to scroll the readable view; switch to the overview for orbit and pinch controls.
The overview panel follows the scene's proportions up to its desktop or phone
height limit, keeping inspection controls near wide diagrams.
Camera clipping follows the scene's depth as the view changes, including large
flowcharts whose initial camera distance exceeds that of the examples.

The [theme gallery](../alwaldend.com/docs/mermaid-theme/README.md) owns the six
Mermaid definitions. The sample targets consume those files directly and use
the site’s configuration, compiled Sass palette, and pinned diagram font. The gallery links static
SVGs, 3D snapshots, and interactive 3D views.

## Rendering

[![Mermaid 3D architecture: parsing and theme measurement feed spatial adapters, which package an offline Three.js viewer](docs/assets/architecture.3d.png)](docs/assets/architecture.3d.png)

[Explore this diagram in 3D](docs/assets/architecture.3d.html).
Both views are generated from one [Mermaid source](docs/architecture.mmd).

`mermaid_3d` in `defs.bzl` takes a Mermaid `src`, JSON `theme`, CSS `palette`,
and pinned `font`, with `out` for HTML and `model` for inspection JSON.
The browser measures text and resolves the light/dark Sass palette.
The default output is HTML; the named JSON output is available for diagnostics.

## Static snapshots

`mermaid_3d_snapshot` renders an existing HTML target as a PNG or WebP image.
It uses the viewer's scene and framing, includes geometry and labels, and
encodes only the canvas. Documentation can embed the image and link to the
interactive HTML. No controls, server, or remote assets are required.

```starlark
load("//projects/mermaid_3d:defs.bzl", "mermaid_3d_snapshot")

mermaid_3d_snapshot(
    name = "diagram_snapshot",
    src = ":diagram_3d",
    out = "diagram.3d.png",
    color_scheme = "light",
    camera = "isometric",
    width = 1600,
    height = 900,
)
```

Choose `camera = "front"` for a front projection. PNG is lossless; WebP uses
the browser's highest quality setting. Defaults are 1280 by 960 pixels. Each
dimension must be 128–8192 pixels, with at most 16,777,216 pixels in total.
Choose an aspect ratio suited to the diagram; fitting the complete scene can
make text small in a narrow image. Build all six light/dark sample pairs with:

```sh
bazel_agent bazel build //projects/mermaid_3d/samples:snapshots
```

The site recognizes matching `.dark.png` and `.dark.webp` siblings using the
same theme menu as SVG diagrams. The architecture image above is a snapshot.

## Scene construction

Mermaid owns parsing. Its version-pinned `getDiagramFromText` databases supply
graph data for flowcharts, state, class, entity-relationship, and mindmap
diagrams. The separate sequence adapter consumes participants and ordered
events, preserving messages, replies, activation, self-calls, notes, and nested
interaction blocks. Participant declaration order is retained for numeric IDs,
and arrowless messages remain arrowless. Blocks include their notes, including
blocks containing only notes. Mermaid's SVG layout hooks remain part of its SVG
rendering lifecycle; this project instead converts semantic data into a Three.js scene.

The parsing browser normalizes display text before measurement, including Mermaid
entity codes and generic class headings. Invisible graph links constrain placement
without becoming visible or selectable connections. Groups use their declared
direction or inherit the enclosing direction. Sequence activations align with
their message events, with nested intervals offset to remain visible.

One measured text layout supplies card dimensions and label drawing, including
bold headings and compartments. Sequence participant spacing and note intervals
accommodate multiline text. Decision nodes use diamond prisms with clear front
faces sized to contain multiline labels; rear group strokes recede using the
existing palette as the camera moves.
Group titles retain line breaks and contribute to their frame dimensions;
camera framing includes the labels' visible extents.

Graph layout condenses cycles before ranking, places group members recursively,
and reserves depth corridors for spatial edge routes. Connections leave side
ports and use smooth curves when clear, adding rounded detours around obstructing
nodes. Curvature adapts to available clearance; arrowheads follow the path tangent,
and dash spacing remains continuous through bends.
Self-loops use separate ports on the same side and remain outside their node face.
Connections and guides share one tube renderer, with dash sections merged into
one mesh per path. Each embedded JavaScript module has one import-map entry,
avoiding duplicate copies of the Three.js core.
Their default projection also clears unrelated node faces. Nodes are solid meshes;
groups are volume outlines; paths, markers, and sequence lifelines have thickness.
Sequence time runs downward, with participant depth and separate message lanes.
Relationship cardinalities remain explicit labels anchored near their endpoints.
Class and entity cards distinguish headings and compartments. Annotation labels
avoid projected node faces and other labels using bounded displacement with
leader lines, recomputed when the camera changes.

This is an initial renderer for the maintained examples, not full Mermaid syntax
or shape parity. It preserves the shared font and palette, not arbitrary SVG CSS
or per-node styles. Unsupported sequence events, actor boxes, and dynamic actor
creation/destruction fail explicitly. Arbitrary camera angles or unusually dense
diagrams can exceed the label-placement search area; selecting an element and
changing the camera remains useful. Optimal routing is not guaranteed. Mermaid
database changes require adapter verification on upgrades.

## Source layout

The command builds the model and packages a standalone viewer. Implementation
modules stay under `internal/`; the browser controller owns camera and UI state,
with scene construction and label placement handled separately.

| Path                                                | Responsibility                                                               |
| --------------------------------------------------- | ---------------------------------------------------------------------------- |
| `cmd/render/main.mjs`                               | Parse inputs, measure themed text, and embed the viewer's modules and assets |
| `internal/model.mjs`, `sequence.mjs`, `routing.mjs` | Adapt Mermaid data, place elements, and route connections                    |
| `internal/text.mjs`                                 | Shared text layout for build-time measurement and browser drawing            |
| `internal/viewer/main.mjs`                          | Camera, controls, selection, and rendering updates                           |
| `internal/viewer/scene.mjs`, `labels.mjs`           | Three.js objects, label drawing, placement, and inspection                   |
| `internal/viewer/document.mjs`, `viewer.css`        | Offline HTML document and viewer styles                                      |
| `test/browser/`                                     | Browser verification and behavior fixtures                                   |

`cmd/snapshot` owns offline image export; `internal/snapshot.bzl` exposes its
build rule, and `test/snapshot` checks the actual encoded images end to end.

## Verification

```sh
bazel_agent bazel test //projects/mermaid_3d/test/browser:browser_test
bazel_agent bazel test //projects/mermaid_3d/test/snapshot:snapshot_test
```

The browser test opens every built view with remote requests blocked. It checks
source labels and counts, node volumes, group containment, routes against unrelated
nodes, sequence ordering, theme/font loading, selection, camera orbit/reset, and
mobile controls, readable label sizes, projected label separation, scene framing,
and native touch scrolling. It emits full-scene, light/dark, front, phone, and
overview screenshots and `verification.json` containing
the checked HTML hashes in the test's `test.outputs` directory.
Two sequence fixtures check numeric participant order, arrowless messages,
nested and note-only blocks, multiline cards, and actual selection opacity.
Three adapter fixtures additionally cover activation timing and nesting, entity
decoding, generic class text, diamond containment, invisible links, and local
subgraph direction.
Three boundary fixtures cover shared element identifiers, self-loop projection,
and long multiline group titles on phones.
The fourteen maintained views also enforce budgets of 4 MiB per offline HTML file
and 150 draw calls per scene, recording measured values in the JSON report.
A separate fifty-node flowchart checks visible canvas pixels after opening,
orbiting, resetting, changing view and appearance, and resizing. Its images and
visibility measurements are retained alongside the collection's artifacts.
