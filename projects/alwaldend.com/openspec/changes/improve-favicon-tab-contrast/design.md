## Context

Browser color preference can differ from a tab's background. The S8 source
already defines a presentation selected by the site's shared favicon partial.

## Decisions

Use one transparent SVG with white fill and charcoal stroke on the canonical S8
paths. Paint the outline beneath the fill to retain the width of the white
artwork. Limit this styling to the favicon variant and keep the existing Hugo
content fingerprint.

## Risks and Trade-offs

An overly thick stroke can close narrow gaps when rasterized. Inspect the
16-pixel result over both light and dark backgrounds.
