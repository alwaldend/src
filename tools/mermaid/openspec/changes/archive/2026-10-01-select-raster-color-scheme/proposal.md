## Why

Native WebP output always selects the light palette even when its consumer
needs the site's dark appearance for a standalone image.

## What Changes

- Add explicit light/dark selection for native WebP using the existing palette.
- Preserve light as the default and the SVG rule's paired output contract.
- Verify both color modes against independent browser renderings of the SVGs.

## Capabilities

### Modified Capabilities

- `mermaid-rendering`: native raster color scheme selection.

## Impact

The raster rule, native renderer, documentation, skill, and renderer E2E
checks change. No external dependency or historical-theme change is needed.
