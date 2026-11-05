## Why

The existing Mermaid examples are planar images. The user wants solid nodes,
volumetric groups, spatial edge routing, and interactive inspection while
retaining the maintained diagram theme and source definitions.

## What Changes

- Use the Mermaid parser, CLI, and Three.js dependencies delivered separately.
- Add a hermetic Three.js renderer for the five graph examples in the site gallery.
- Preserve their authoritative Mermaid source, theme, palette, fonts, labels,
  relationship semantics, and available group structure in an interactive 3D view.
- Use public Mermaid graph layout/render APIs; exclude sequence diagrams.
- Add browser verification and reproducible offline HTML views.
- Keep static PNG/WebP export and raster website integration in the separate snapshot delivery.
- Use strict TypeScript for the core renderer, consuming the merged setup from
  PR 259; preserve the separately delivered snapshot implementation in TypeScript.

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
