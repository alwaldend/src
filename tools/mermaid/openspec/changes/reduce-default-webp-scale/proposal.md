## Why

The default twofold raster scale makes Mermaid WebP images larger than needed. Use natural-size rasters by default while retaining explicit higher-resolution renders.

## What Changes

- Change the WebP rule, render CLI, and encoder defaults from scale 2 to 1.
- Retain positive integer scale overrides and exact aspect-ratio padding.
- Update documentation and regenerate the two checked-in blog WebPs through their owning update targets.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `mermaid-rendering`: Specify natural-size default raster geometry and explicit scale overrides.

## Impact

Mermaid rendering defaults, existing renderer fixtures, and canonical skill guidance. The blog consumers retain their source, appearance, and 5:2 padding; their generated WebPs shrink. No new dependencies or deployment.
