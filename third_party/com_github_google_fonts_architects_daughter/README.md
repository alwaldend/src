---
title: Architects Daughter
description: Pinned handwriting font used by the Mermaid diagram theme
---

[Architects Daughter](https://fonts.google.com/specimen/Architects+Daughter) by
Kimberly Geswein, licensed under the SIL Open Font License 1.1.

`tools/mermaid/theme.json` requests this family so hand-drawn diagrams keep a
handwritten label voice. The renderer exposes it to Chrome through the pinned
Fontconfig alongside the repository's Liberation fonts, so a render never
resolves the family from the host.

Upstream source:
[`google/fonts`](https://github.com/google/fonts/tree/main/ofl/architectsdaughter).
The pin is the immutable commit that last changed the font file, and the SRI
digest is recorded in `include.MODULE.bazel`.
