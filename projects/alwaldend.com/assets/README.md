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
