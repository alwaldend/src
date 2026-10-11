## Why

WebP diagrams should use a 5:2 width-to-height canvas by default, as requested
while refining the infrastructure architecture diagrams.

## What changes

- Default rule, CLI, and encoder raster output to exact 5:2 padding.
- Keep explicit positive-integer ratios and a natural-canvas opt-out.
- Preserve SVG layout, content scale, colors, and fonts.
- Update renderer guidance and show the complete 3x-ui render at this ratio.

## Capabilities

### Modified capabilities

- `mermaid-rendering`: Default raster canvas proportions and explicit opt-out.

## Impact

The Mermaid WebP renderer and its consumers change canvas proportions. SVGs
and explicitly selected ratios retain their existing geometry. No dependencies
or live infrastructure changes are needed.
