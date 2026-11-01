## Context

The neutral site preset uses an embedded Liberation face. The previous
appearance study records four named light/dark palettes separately from its
font and geometry proposals. See proposal.md for the requested outcome.

## Decisions

- Match the top-level/theme font families and embedded file to Liberation
  Mono. Preserve existing font sizes and sequence-specific defaults.
- Retain the active Sass palette. Move the four proposal palettes into the
  reusable Mermaid package as `color-schemes.json`, with named semantic roles.
  Keep typography/geometry proposals separate and reference the catalog by ID.
  A data catalog preserves the requested choices without adding a selector or
  a second palette rendering implementation.
- Reuse existing browser, font-subsetting, and raster checks. Verify actual
  equal glyph advances and embedded-face availability before implementation.
- Preserve the existing displayed-edge-width thresholds. If the wider font
  reduces graph strokes below one CSS pixel, increase only normal edge weight
  as required; preserve thick/invisible distinctions and all palette colors.
- Correct current site typography references in the renderer's consumer
  documentation and fixtures. Historical assessments retain their original
  measurements and fonts.

## Failure cases

- Requested Mono family differs from embedded face, or some grammars retain
  proportional labels.
- Wider labels clip captions, overlap container titles, or reduce connectors
  below the established article-width threshold.
- Saving alternatives accidentally changes the active neutral palette.
- Native SVG/WebP fixtures retain a Sans face while consuming the Mono preset.
- Palette values drift during promotion from the proposal record.

## Verification

Host Bot's current T3 Code and ComfyUI diagram targets use this native preset.
The earlier diagram-migration blog keeps its existing image and presentation
through Host Bot's explicitly named `rendered_historical` target. When the
normal target switches to a light/dark pair, update that one consumer label
and verify its existing generated-image freshness test. No article content,
publication asset, shared preset, or palette changes are needed.

Render all six gallery families in both modes and inspect representative
output. Retain browser reports/screenshots; compare saved palette values to
the previous committed proposal. Run site, native font/raster, formatting,
OpenSpec, repository-quality, and affected-package checks before publication.
