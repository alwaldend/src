## Why

The selected S8 graduated square monogram needs a durable source asset in
the repository's data tree and a generated dark WebP preview.

## What Changes

- Save the approved S8 geometry as an SVG under `data/alwaldend/assets/`.
- Center the complete mark within its square canvas without altering its
  internal geometry or spacing.
- Generate a 512 by 512 lossless WebP preview with off-white ink on charcoal
  under ignored output. The SVG is the sole source artwork.
- Document the source, palette, and rendering procedure.

## Capabilities

No runtime behavior changes. This asset-only record uses `skip_specs: true`.

## Impact

Adds repository-internal identity assets without new external dependencies.
Site integration is tracked by the site-owned `use-square-monogram-favicon`
change. A dedicated WebP build rule is deferred at the user's request.
