---
title: Mermaid CLI
description: Compatibility patch for the pinned Mermaid CLI
---

The pnpm workspace pins `@mermaid-js/mermaid-cli` 12.0.0.
`ttf-extension.patch` corrects `.tff` to `.ttf` in its local resource loader,
allowing the supported `customFontCSS` API to load the repository's pinned
TrueType fonts and embed their subsets and license in native SVG output.

Source: [Mermaid CLI 12.0.0](https://github.com/mermaid-js/mermaid-cli/tree/12.0.0).
Remove the patch when an upstream release includes the correction. Existing
native font, raster, and site readability tests exercise the patched path.
