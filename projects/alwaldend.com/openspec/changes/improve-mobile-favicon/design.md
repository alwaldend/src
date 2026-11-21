## Context

The Hugo partial selects a favicon presentation of the canonical S8 SVG and
fingerprints the result. The same SVG supplies the inline header logo.

## Decisions

Render the original monogram in solid charcoal. Use a rounded white tile
that blends into a white-backed mobile tab and provides contrast on a dark
tab. Keep the same paths for the header and favicon, and retain the existing
fingerprinted resource pipeline.

## Risks and Trade-offs

Browser tab strips may add their own background around the icon. Inspect the
16-pixel rendering against both light and dark surfaces and within a white
circular backing like the one in the mobile screenshot.
