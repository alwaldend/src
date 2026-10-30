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

`scss/application.scss` exports the site's Bootstrap typography, neutral palette,
light/dark variables, and shared controls for embedded application pages. It
imports the same `_shared_accent.scss`, `_shared_palette.scss`, and
`_shared_controls.scss` used by the Hugo theme; applications supply only their
own layout rules. The `site_stylesheet` rule in `defs.bzl` compiles these sources
with the pinned Sass tool and existing Bootstrap archive. The profile cleaner's
local browser interface embeds that output, without a runtime site dependency.
The Mermaid stylesheet compiler remains available through its existing alias.

Shared controls own filled and outline primary buttons and checkbox selection
and keyboard-focus styles in both palettes. Application consumers use Bootstrap's
component classes rather than defining a second primary-button variant.
