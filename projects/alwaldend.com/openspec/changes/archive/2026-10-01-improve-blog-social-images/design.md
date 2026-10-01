## Context

See `proposal.md`. Existing sources, themes, and generated asset rules remain authoritative.

## Decisions

Select the native color mode through a rule parameter, using the existing media-aware palette resolver. The uploader post explicitly selects dark and uses a small canvas margin rather than a forced social aspect ratio. Its one WebP serves the body and social metadata. The older diagram retains its existing body SVG while a generated WebP supplies its social card. Posts without an opening image receive no invented asset.

## Verification

Compare actual light and dark WebP pixels with independently rendered SVG references. Check dark page/background agreement, useful diagram occupancy with visible margins, generated asset freshness, both posts' social metadata, and unchanged uploader image counts. No X requests are authorized.

## Validation evidence

The native light/dark raster, historical raster, post freshness, and site
browser checks passed (10 targets). Full uploader E2E and converter command
checks passed (2 targets), including the local HTTP banner scenarios. The
regenerated real-post artifact contains 39 blocks, one body image locator,
and a separate banner locator sharing its digest. The rendered dark diagram
was inspected for readable text, visible arrows, and modest canvas margins.
No live X request or site deployment was performed.
