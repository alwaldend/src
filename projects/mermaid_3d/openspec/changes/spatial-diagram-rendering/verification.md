## Core acceptance

The core renderer builds flowchart, state, class, entity-relationship, and
mindmap definitions as self-contained Three.js HTML. Mermaid's public layout
loader and render APIs supply typed graph data. Sequence diagrams are rejected
explicitly; no internal Mermaid database API or sequence adapter is used.
The browser normalizes and measures text with the shared font before spatial
placement and curved routing.

All implementation and browser E2E sources use ordinary strict TypeScript with
`noEmitOnError`, separate Node/browser contexts, and emitted `.js` modules.
There are no unchecked `any`, checking suppressions, or library-check exemptions.
One `MermaidThemeInfo` target supplies configuration, JSON light/dark palettes,
font, and license; each viewer embeds its assets and works offline.

## Verification scope

The browser E2E covers the five graph examples, five graph boundary fixtures,
and a fifty-node flowchart. It checks source meaning, geometry, projected
labels, relationship selection, both appearances, embedded fonts, desktop and
phone interaction, payload size, and draw-call budgets. At device pixel ratio
2 it verifies actual WebGL backing dimensions and GPU limits, raster content
for every node, and restored native density in the fitted overview. Artifacts
include screenshots, HTML hashes, visibility measurements, and drawing-buffer
metrics. The high-density failure was reproduced before the viewer correction.
Diamond-port acceptance measures endpoint distance to the polygon's sides for
horizontal and vertical fan-in/fan-out and self-loops. Its focused fixture
reproduced detached ports before the routing correction; front-view images and
per-port distances identify the failed and corrected geometry.

Compare every retained graph field and all 18 light/dark scene PNGs with the
prior reviewed baseline. Only empty sequence-specific model fields are removed.
Run existing native Mermaid and site consumers, TypeScript setup acceptance,
root repository quality, affected DNS registry checks, and semantic lint. Build
the five sample HTML files, shared SVG/HTML architecture, and complete site;
verify generated asset links and unsupported-sequence CLI failure. Strict
OpenSpec validation covers the maintained change artifacts. Ignored delivery
evidence owns exact candidates, test results, reviews, and publication state.

## Independent snapshots

The earlier complete TypeScript snapshot implementation is preserved on the
task-owned snapshot branch. Derive only its exporter, rule, preparation API,
raster consumers, and image E2E changes after core merges, using the current
graph-only model and public Mermaid integration. Its PR targets fresh master.
Preservation refs and raw validation artifacts remain in ignored task evidence.

## Pipeline documentation

The shared architecture flowchart follows the renderer CLI from inputs to
self-contained HTML and optional JSON containing the scene and metadata.
Bazel is a caller of the CLI. Headless Chrome supplies parsing, normalization,
and text measurements; CLI code computes layout/routes and writes the outputs.
These boundaries are checked against the command and rule sources.

The `mermaid_3d` CLI and headless Chrome are separate processes, connected by
Puppeteer calls and returned graph data/font metrics. Mermaid runs inside Chrome.
The architecture target selects the saved front-left perspective. Browser
acceptance checks the CLI to the left and nearer the camera than Chrome and
exercises the normal orbit/reset controls on this generated document.
Font loading and initialization precede the `mermaid.render()` call.
The registered layout hook projects typed graph data and normalizes text;
canvas measurement supplies normal/bold font metrics. CLI code then sizes and
places nodes/groups, routes connections, and packages the outputs.
Browser inspection verifies the processing stages, nested group membership,
connections, loaded fonts, and
unclipped labels in both appearances. Inspect the actual 3D image at chat width
as well as native size. Images provide the inline preview; complete
self-contained HTML retains camera and selection controls. Exact geometry,
label-size metrics, artifacts, and delivery state remain in ignored evidence.
Inspect caption plates in both appearances, including edge captions and endpoint
cardinalities, for the shared themed border. Existing browser E2Es exercise all
supported examples, selection, framing, and interaction after the drawing change.
