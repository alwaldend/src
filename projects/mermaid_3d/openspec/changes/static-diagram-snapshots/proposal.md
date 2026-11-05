## Why

The merged renderer produces interactive HTML, but documentation also needs
static 3D images that retain the same theme and work in ordinary image embeds.
The preserved snapshot implementation can now be delivered independently of core.

## What Changes

- Export a self-contained viewer canvas to exact-size PNG or WebP through a Bazel rule and CLI.
- Let the viewer prepare appearance, camera, dimensions, and complete scene framing for export.
- Publish light/dark 3D images in the existing gallery and architecture documentation, with links to interactive views.
- Reuse the saved strict TypeScript exporter and raster E2E against current graph adapters.

## Capabilities

### New Capabilities

- `static-snapshots`: Offline canvas export, configured framing, static documentation embeds, and repeatable raster verification.

### Modified Capabilities

None.

## Impact

Changes remain in `projects/mermaid_3d` and its existing website consumers.
The renderer continues using the merged public Mermaid graph hook and shared
theme bundle. No dependency addition, sequence adapter, or deployment is needed.
