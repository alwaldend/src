## Context

See proposal.md. The site exports an offline application stylesheet compiler and the repository already pins Puppeteer and Chrome.

## Goals / Non-Goals

Use site-owned tokens without palette duplication. Capture literal source; highlighting uses the existing Highlight.js pin; line selection remains outside this command.

## Decisions

Compile the existing application stylesheet and add tool-owned layout only. Use textContent rather than HTML interpolation for source and filename. Embed a pinned monospace font for stable geometry. Reject requests in Chrome and disable JavaScript. Normalize caller paths against BUILD_WORKSPACE_DIRECTORY while resolving bundled resources in runfiles.

## Risks / Trade-offs

Very large files can exceed Chrome image limits; fail explicitly above a bounded canvas size. Comments use the site secondary text color; keywords use emphasis and medium weight. Browser sandbox support varies by host; use the existing repository headless launch flags.

Each original source line is a grid row with a number gutter and an inert text cell. CSS pre-wrap preserves indentation and wraps long lines at the selected canvas width.

Highlight the full normalized file on the Node side. Walk the resulting span/text DOM to project token classes into numbered lines without losing multiline token scope. Source is escaped by the highlighter and inserted only as text nodes during projection. Infer language by extension with an explicit override and plain-text fallback.
