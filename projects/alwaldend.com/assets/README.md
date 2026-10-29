---
title: Assets
description: Hugo assets
---

The assets package maps `//data/alwaldend:monogram` to `icons/favicon.svg`
in the Hugo source archive. The [identity asset](../../../data/alwaldend/README.md)
owns the S8 square monogram, its centered geometry, and adaptive light/dark
colors. The site consumes that SVG directly without a separate artwork copy.

The shared favicon partial fingerprints the resource to refresh cached tab
icons when the canonical SVG changes.

`scss/_shared_accent.scss.tmpl` maps the authoritative identity data from
`//data/alwaldend:assets/theme.json` to Bootstrap/Docsy variables. The build
packages its generated `_shared_accent.scss` at the usual theme import path.
Component rules stay in this assets tree; shared color values are edited in
the identity data rather than in generated Sass.
