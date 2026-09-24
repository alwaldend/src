## Why

The current documentation renderer rewrites SVG colors and label paint order,
and inherits CSS workarounds intended for the historical hand-drawn theme.
The user approved native Mermaid rendering with normal label placement and
plain actor nodes, while preserving blog diagrams.

## What Changes

- Render light and dark SVGs with native Mermaid configuration and CSS input.
- Derive both palettes from the site's existing Sass tokens.
- Select the image using the site's existing theme attribute and print CSS.
- Remove documentation SVG rewriting, custom label geometry, and actor masks.
- Preserve historical blog rendering and source diagrams.

## Capabilities

### New Capabilities

### Modified Capabilities

- `project-alwaldend-com`: documentation diagrams use native rendering without
  custom SVG post-processing.

## Impact

Mermaid build rules and renderer, the site's palette and image templates,
Host Bot documentation packaging, and browser verification. No dependency
upgrades or deployment.
