## Why

Native site diagrams each embed the full Liberation Sans font, contributing
about 17 MB compressed to the website archive. The user approved FontTools
4.67.0 as a pinned build dependency to retain only needed glyphs.

## What Changes

- Subset the native diagram font from actual rendered SVG text, including entities.
- Pass the subset through Mermaid's CSS input and preserve standalone output.
- Verify retained glyph coverage and metrics, rendering, and archive size.

## Capabilities

### New Capabilities

- `native-font-subsetting`: Small self-contained native SVG fonts.

## Impact

tools/mermaid owns the renderer and subset command; tools/py owns the approved
Python dependency and generated lock. Historical rendering and theme pins stay
unchanged. No host package installation or external font fetch is needed.
