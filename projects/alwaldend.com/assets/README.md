---
title: Assets
description: Hugo assets
---

The assets package maps `//internal/assets:monogram` to `icons/favicon.svg`
and Docsy's `icons/logo.svg` header resource in the Hugo source archive.
The [identity asset](../../../internal/assets/README.md)
owns the S8 square monogram, its centered geometry, and adaptive light/dark
colors. The site consumes that SVG directly without a separate artwork copy.

The shared favicon partial selects the SVG's favicon variant and
fingerprints the result to refresh cached tab icons when its contents change.

The Hugo archive mounts the [common web theme](../../../internal/assets/web)
under `assets/scss`. The application stylesheet in that package exports the
same Bootstrap typography, neutral palette, and shared controls. Its
`site_stylesheet` rule uses the pinned Sass tool and Bootstrap archive. The
profile cleaner embeds the compiled stylesheet in its local web interface.

Shared controls own filled and outline primary buttons and checkbox selection
and keyboard-focus styles in both palettes. Application consumers use Bootstrap's
component classes rather than defining a second primary-button variant.

`//internal/assets:theme.json` defines the color values. The common web theme
maps them to Bootstrap and Docsy variables in the generated
`_shared_accent.scss`. Site-specific component rules live in this assets tree.
