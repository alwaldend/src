## Why

The existing Mermaid examples are planar images. The user wants solid nodes,
volumetric groups, spatial edge routing, and interactive inspection while
retaining the maintained diagram theme and source definitions.

## What Changes

- Upgrade Mermaid to the latest stable release verified from the primary registry,
  with a compatible CLI and reconciliation of the existing edge-spacing patch.
- Add a hermetic Three.js renderer for the six existing site gallery examples.
- Preserve their authoritative Mermaid source, theme, palette, fonts, labels,
  relationship semantics, and available group structure in an interactive 3D view.
- Add browser verification and reproducible HTML and static PNG/WebP snapshots.

## Capabilities

### New Capabilities

- `spatial-diagrams`: Themed interactive 3D projections of Mermaid diagrams.

### Modified Capabilities

None. Existing SVG and raster output contracts remain applicable.

## Impact

The renderer, tests, and build rules belong to projects/mermaid_3d as explicitly
requested by the user. Dependency pins and locks
belong to tools/, and the existing six sample sources belong to projects/alwaldend.com.
Shared identity extraction is delivered separately from the site-owned theme inputs used here. The user explicitly approved adding
three@0.186.1 on 2026-10-29. No other new external package
is planned. This work
does not deploy the website or change host configuration.
