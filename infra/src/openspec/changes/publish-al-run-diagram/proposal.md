## Why

The standalone AL module needs the repository site adapter to render and publish
its lifecycle diagram using the existing Mermaid toolchain.

## What Changes

- Render the authoritative AL diagram through tools/al.
- Package light/dark SVG assets for the AL lifecycle documentation page.
- Provide a manual raster preview for visual verification.

## Capabilities

None. Documentation packaging only; skip_specs is enabled.

## Impact

The tools/al documentation targets. Diagram content is owned by
projects/al/openspec/changes/document-al-run-phases.
