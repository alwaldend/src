## Context

The X uploader demonstration post is both a draft blog page and a conversion
fixture. Its single diagram must work in the site, social metadata, and the
uploader. See `proposal.md` for the raster policy.

## Goals / Non-Goals

Reuse the native site appearance and one bundled image across all consumers.
Keep draft state unchanged and make no X API call or site deployment.

## Decisions

- `mermaid_site_webp` uses the native WebP mode owned by
  `tools/mermaid/openspec/changes/archive/2026-10-01-add-raster-diagram-renders`.
  It shares the
  site SVG macro's theme, palette, and label font rather than copying colors
  into a post-specific theme.
- The source diagram uses Mermaid's canvas padding and ends at the uploader's
  actual result, an article draft. Publishing is outside the tool.
- Hugo's existing `images` front matter selects `pipeline.webp` for both
  social templates. A featured-image copy or head override would add no
  capability and would create another resource to maintain.

## Risks / Trade-offs

- Social consumers choose their own crop → retain generous canvas padding
  and inspect the rendered image at typical card dimensions.
- A changed image invalidates the old conversion digest → refresh the offline
  artifact after regeneration; never retry draft creation without authorization.
- Front matter might affect the conversion fixture → run its existing E2E
  check and verify it still has one image locator.

## Verification

The generated image is 2148 × 1072 with all nodes and both edge labels visible.
The site browser E2E check resolves identical Open Graph and Twitter image URLs
and `summary_large_image`; its output includes `pipeline-social.png` and a
JSON report. The uploader E2E check passes with one image. Offline conversion
produces 39 blocks with 13 continuing style-loss diagnostics and no fatal
diagnostic. Draft creation remains paused; no X request or site deployment
was performed.
