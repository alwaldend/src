## Context

The site publishes hermetically rendered Mermaid SVG images. The user requested
clean typography and rounded shapes that follow the site theme, then explicitly
excluded blog diagrams. Blog figures are historical publication assets and also
illustrate the original hand-drawn style.

## Goals / Non-Goals

**Goals:** rounded nodes and containers, sans-serif labels, adaptive colors,
readable actor icons, and verified selector, system, mobile, and print behavior
for documentation diagrams. Preserve all blog diagrams.

**Non-Goals:** alter diagram topology, blog files or appearance, default Mermaid
rendering, infrastructure configuration, or deploy the site.

## Decisions

Decision-review verdict: proceed with a site render preset and generated SVG
image resources. SVG media queries inherit the embedding element's color
scheme (MDN's `prefers-color-scheme` embedded-elements contract). This keeps
native image loading, accessibility, caching, and isolated IDs. Inline SVG would
expose Mermaid's repeated IDs and foreignObject contents to the document and
print rewriter. JavaScript source swapping would add synchronization state.

`mermaid_site_svg` merges `site-theme.json` over the base configuration and
embeds the pinned Liberation Sans font used during layout. It adds rounded
rectangle corners, straight connectors, and geometric actor icons. Host Bot
exposes both site and original renders from one `.mmd`; its blog consumer keeps
the original target. The default renderer and all blog SVG bytes are unchanged.

The renderer reads the source palette from `tools/mermaid/theme.json` and
converts paint colors to CSS variables with standalone fallbacks. Actor icons
use a CSS mask painted with `currentColor`; no nested SVG palette is needed.
Hugo matches the explicit site-preset CSS marker and appends the adaptive
palette compiled from `_shared_accent.scss`, without parsing colors or decoding
embedded images. Markdown resources and the SVG shortcode use one cached
partial with the page section in its cache key. Blog pages bypass adaptation.

PR review decision: revise the initial Hugo color adapter. Its strongest
advantage was keeping the renderer consumer-neutral, but it coupled templates
to Mermaid's output encoding and required rewriting nested actor SVGs. Keeping
normalization in the renderer and using color-inheriting masks removes those
template responsibilities while retaining native SVG images and shared tokens.

## Risks / Trade-offs

- SVG isolation prevents inheriting page custom properties; the generated
  stylesheet consumes the owning Sass tokens at build time.
- Font changes require re-rendering, so selecting the site preset also selects
  its matching embedded font. Existing publication targets retain their font.
- The browser test needs Docsy's existing SRI-pinned jQuery and Lunr scripts to
  exercise the actual selector. It permits only those two external URLs; the
  diagram build and assets remain self-contained.

## Failure cases to verify

- Explicit light/dark selection against the opposite system preference, saved
  mode after reload, and live system-preference changes.
- Canvas, text, lines, caption plates, or actor icons retaining unreadable colors.
- Theme switching changes image dimensions, clips labels, or breaks resource URLs.
- Rounded nodes or nested containers revert to square or hand-drawn outlines.
- Blog image bytes or URLs change; documentation print images fail to load;
  mobile diagrams overflow the page.

## Validation and session review

The browser check exercises the real theme menu against opposite operating
system preferences, persisted mode after reload, system changes, mobile, print,
rounded corners, label/line/caption colors, and preserved blog image URLs and
bytes. It produces screenshots and a JSON report in Bazel test outputs. The
five original blog freshness checks and existing renderer/CLI checks pass.

Visual inspection caught `MERMAID-RGB-CAPTIONS`: Mermaid normalizes custom CSS
colors to RGB while retaining hex elsewhere. The renderer handles both, and
the browser verifies caption backgrounds. `MERMAID-PRINT-IMAGE-SCHEME` requires
setting color-scheme on the embedding image for paper; image media queries do
not inherit the parent document's print media type.

Ergonomics review: one initial unbounded SVG/worktree search produced excessive
output, and several guessed output paths added avoidable reads. Restrict future
searches by source extension and obtain output paths from build results. The
initial browser launch also needed npm resolution from the existing tools
workspace and explicit access to Docsy's two CDN scripts. These are captured in
the checked-in test. No changes to shared agent policy are proposed.
