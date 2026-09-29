## Why

The tool renders a diagram to SVG, which is the right format for a page and the
wrong one for every consumer that must upload the image to a service. The X
Articles media upload endpoint accepts only a fixed set of image media types and
`image/svg+xml` is not among them, so a diagram cannot reach that consumer at
all without a raster render.

The pipeline already launches a pinned headless Chrome to produce the SVG, and
that browser already encodes WebP. A raster render therefore costs a render
mode rather than a new image-conversion toolchain, and it can keep every
property the SVG render has: the same theme, the same pinned fonts, the same
paint-order pass, and the same hermetic inputs.

## What Changes

- Add a `mermaid_webp` rule beside `mermaid_svg` that renders a diagram source
  to a WebP image.
- Keep the maintained appearance, pinned browser and fonts, and the paint-order
  pass unchanged between the two formats, so one diagram has one appearance.
- Encode WebP through the pinned browser rather than adding an image-conversion
  dependency.
- Extend the tool README and the Mermaid authoring skill to document the raster
  rule and when a consumer needs it.

## Capabilities

### New Capabilities

- `mermaid-rendering`: the diagram render surface this tool publishes to
  consumers, including the raster render a service-upload consumer needs.

### Modified Capabilities

None.

## Impact

- `tools/mermaid/**`: a new rule, its renderer support, and documentation.
- Consumers that upload a diagram now have a hermetic raster render to reference
  instead of committing a raster conversion by hand.
- The browser test surface gains the raster rule; the existing SVG rule and
  every current consumer of it keep their behavior.
