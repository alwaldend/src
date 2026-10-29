## Context

See proposal.md. The site owns six example sources and a flat light/dark theme.
The existing renderer uses pinned Mermaid CLI, Chromium, and embedded fonts.
The registry reports Mermaid 12.2.0, CLI 12.0.0, and Three.js 0.186.1 on
2026-10-29; the existing Node 24.13.0 satisfies their runtime requirements.

## Goals / Non-Goals

The new renderer is an additional output from existing definitions. It must
make real scene objects and spatial routes, not place an entire SVG on a plane.
It does not replace existing SVG exports, promise support for every Mermaid
diagram or shape, or deploy the website.

## Decisions

- Reuse Mermaid parsing and available structured graph data. Keep version-sensitive
  diagram adapters isolated. Native rendering can supply text measurement and
  semantic appearance that is not present in a common graph interface.
- Use deterministic placement with the original flow direction,
  group containment, and usable depth. Route edges with explicit clearance in
  three dimensions; retain terminal semantics and readable labels.
- Three.js owns the scene and controls. Render volumetric nodes, group outlines,
  and spatial connections. Keep text readable at the initial camera and selectable
  through ordinary HTML controls as well as pointer picking.
- The source theme and palette remain authoritative. Resolve the light/dark CSS custom properties
  in the pinned browser, embed the existing font, and package pinned JavaScript
  locally in self-contained HTML. No CDN or new layout dependency is needed.
- Reuse the site's six sources without copies; site macros pass their maintained
  theme inputs to the generic renderer. Generated HTML remains a build output.

Alternatives: extruding a finished SVG does not provide independent graph routing;
an unconstrained force simulation disrupts diagram order and reproducibility;
rewriting Mermaid grammars would duplicate an existing authority. These are not
selected. The verdict is to proceed with a separate spatial renderer and a
version-pinned adapter, validating each supported diagram type end to end.

## Risks / Trade-offs

The visual review found tiny default text, mobile resize clipping, labels over
attributes, detached cardinalities, and unnecessary depth detours. Proceed with
a readable view whose canvas can scroll, plus an explicit fitted overview. A
single fit-to-screen policy cannot preserve readable type on narrow screens;
increasing font sizes alone creates overlap. Fit and reset must recompute bounds
for the current camera and container. Keep ordinary text at a readable pixel
scale in the initial view and retain all scene content in its scroll extent.

Prefer direct clear spatial routes between side ports. Use depth detours only
when another node obstructs the direct route, preserving collision checks.
Place relationship cardinalities near their endpoints. Resolve projected label
collisions against node faces and other labels with bounded displacement and
leader lines, updating after camera changes. Class titles and compartments use
the existing font and palette, and group outlines use the theme's line contrast.
These are refinements of the existing renderer, with no additional dependency.

Connections use port-aligned cubic curves for direct graph routes and rounded
quadratic bends for detours and sequence self-calls. Reduce curvature when
needed to retain node clearance. The sampled route remains the common input
for rendering, collision checks, framing, and marker direction. Render joined
tubes and measure dashes along the whole route, so subdivisions do not introduce
visible joints or reset the dash pattern. Straight sequence messages retain
their time alignment. Shared routes and the monochrome theme remain intentional.
Verify smoothness and clearance in the existing browser acceptance, then inspect
markers, dashed curves, self-calls, and both appearances in rendered artifacts.

Selection expands group containment before evaluating connections. A selected
connection highlights only itself and its endpoints; selecting a node includes
its incident connections and their endpoints without following further edges.
Group selection includes descendant notes and nested blocks, with connections
determined from that fixed membership and the block's ordered events. Browser
acceptance inspects actual material opacity rather than only the selected ID.

Viewer identities include element kind as well as the original Mermaid ID.
Keep raw IDs for semantic endpoints, resolving containment separately so actor,
group, note, and connection namespaces cannot overwrite each other. Self-links
use separated ports on the same side and a planar loop outside the node face;
depth detours remain available for unrelated obstacles. Measured group titles
retain line breaks and contribute to group dimensions. Camera framing includes
billboard label rectangles as well as scene geometry.

The sequence adapter consumes Mermaid's explicit actor key order, preserving
numeric participant IDs. Arrowless messages retain no endpoint marker. Notes
record their containing block and participate in its bounds alongside messages
and nested volumes, including blocks containing only notes. A project-owned
browser fixture exercises these combinations through the full build pipeline.

The next visual review found fixed-height sequence cards overflowing on multiline
labels, internal diamond facets crossing text, and duplicate nested encoding of
the Three.js core. Share measured text layout between card sizing and drawing;
use those dimensions for participant separation and sequence event spacing.
Use a flat-front diamond extrusion and one import-map entry per embedded module.
Keep the separate graph and sequence semantic adapters.

Use one path renderer for connections and guides, merging the tube sections for
each path so dashes do not require individual draw calls. Three.js already owns
the geometry utility. Derive quieter rear group strokes from the existing
palette and current camera. Fit mode should use only the panel height needed
for its projected content, bounded by the existing desktop/phone maximums.
Browser acceptance covers offline file size, draw-call cost, label containment,
and compact overview framing, alongside the existing interaction checks.

Normalize Mermaid display fields in the parsing browser before text measurement.
Use the browser's HTML entity decoder after restoring Mermaid's encoded entity
tokens, preserving literal generic type brackets. Keep identifiers separate from
display names. Retain invisible graph edges as placement constraints, excluding
them from routes and interactive relationships. Recursive group placement uses
each group's declared direction or inherits its enclosing direction. Diamond
dimensions must contain the measured text rectangle within the tapered face.

Sequence activation intervals use the actual message event coordinates rather
than offsets from the next message slot. Track nesting depth explicitly and
offset nested bars and their message attachments so each interval stays visible.
Browser fixtures exercise these boundary cases through parsing, measurement,
packaging, rendering, and interaction, with repeatable screenshot artifacts.

- Version-specific databases and DOM details: inspect the new release and verify
  all six examples, nested groups, cycles, self-links, and parallel connections.
- 3D occlusion: provide selection, camera reset, readable anchored labels, and
  group outlines instead of opaque enclosure walls.
- Route collisions: independently check swept segments against unrelated node
  volumes in browser acceptance; fail clearly instead of silently omitting data.
- Theme drift: derive appearance from the site configuration and computed Sass palette;
  verify both modes and the loaded font in browser acceptance.
- Upgrade changes to SVG: run existing renderer and site diagram checks and
  reconcile obsolete dependency patches or generated assets through their owners.

## Migration Plan

Keep the build command under `cmd/render` and move implementation under
`internal`. Model adaptation and sequence semantics remain separate, with shared
curve and graph-routing code in `routing.mjs`. The browser implementation owns
its document, styles, scene construction, label placement, and controller under
`internal/viewer`. Scene and label state have explicit owners; the controller
renders after an update instead of relying on nested rendering side effects.
The existing import-map packaging keeps all browser modules offline. Browser
fixture targets and runfiles derive from common lists. Existing end-to-end
checks and output comparisons establish that the refactor preserves behavior.

Upgrade and validate dependencies, add the optional spatial output, build all
samples, inspect representative rendered artifacts, and deliver a reviewed
feature branch. Existing static diagrams remain the default documentation view.

The user included sequence diagrams after reviewing the adapter approach. A
separate event adapter preserves participants, message order, activation, notes,
and nested interaction blocks. All spatial code is owned by projects/mermaid_3d.

The architecture documentation uses one project-owned Mermaid source for the
README, landing page, static light/dark renders, and an interactive 3D view.
The docs package's tool-consumer exception permits the existing Mermaid render
actions and their font inputs; published outputs contain diagrams and the viewer,
not the build tools. Both pages consume the same packaged assets.

## Static snapshots

The viewer owns snapshot preparation, including exact canvas size, appearance,
selection reset, camera orientation, clipping planes, and bounded fitting of
actual geometry and placed labels. A separate hermetic command loads the same
HTML offline and encodes its canvas as PNG or WebP. The build rule consumes
that HTML target; it does not parse Mermaid or reconstruct the scene.

Both camera views and appearances, encoded dimensions, actual raster content,
fonts, controls exclusion, offline failures, and distant camera depth are
verified end to end. Samples expose light/dark PNG pairs. The project README
and landing demonstrate a snapshot with its interactive HTML link retained.

Fit camera clipping planes from geometry and label depth during every render,
sharing the calculation between interactive navigation and snapshots. Retreat
along the view direction when the scene crosses the near plane. This preserves
visibility when scene extent moves the initial camera beyond a fixed far plane,
and after orbit, reset, appearance changes, or resize. A fifty-node flowchart
exercises the production pipeline and verifies actual canvas pixels through
those interactions, with images and measurements retained as test artifacts.
