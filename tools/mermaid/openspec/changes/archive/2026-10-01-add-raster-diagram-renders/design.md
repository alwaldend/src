## Context

See `proposal.md` for the raster render's purpose. The native SVG renderer
already resolves consumer palettes and embeds the selected font before
rendering; the raster path currently always renders the historical theme.

## Goals / Non-Goals

Support native light WebP through the same SVG renderer and browser encoder.
Preserve the historical default, with no new dependencies or color rewriting.

## Decisions

- Opt-in `native` and `palette` attributes select the existing native renderer.
  Require them together to reject silent theme mismatches.
- Render the light SVG into action scratch and encode it with `renderWebp`.
  A dark SVG is useful for theme-aware site embeds but unnecessary for a fixed
  social image; keep the native SVG rule's two-output contract unchanged.
- Verify the raster against a matching native SVG fixture, while retaining
  the historical raster test. Sharing source alone does not prove palette,
  typography, or natural canvas agreement.

## Risks / Trade-offs

- Native selection could be ignored or post-processed → compare actual native
  output pixels and geometry with the SVG fixture.
- Palette or font could drift → use the declared consumer inputs and test
  light colors and label font.
- Existing raster consumers could change → preserve their default mode and
  run the historical contract checks.

## Verification findings

Visual inspection exposed `native-raster-first-label-height`: native SVG roots
omit a `height` attribute, so the encoder's document-wide replacement removed
the first label's `foreignObject` height. Restrict dimension changes to the
root opening tag and verify labelled-edge pixels independently of that encoder.

The independent first-label image comparison failed before the fix with mean
pixel difference 29.2285 and passed afterwards with difference 0. The existing
historical renderer check also passed. Native E2E outputs preserve both WebP
images and a JSON measurement report.
