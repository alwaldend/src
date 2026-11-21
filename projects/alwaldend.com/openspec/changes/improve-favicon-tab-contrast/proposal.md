## Why

The graphite S8 favicon loses definition on dark browser tabs. The icon needs
contrast against both light and dark tabs, including custom browser themes.

## What Changes

- Render the transparent S8 favicon in white with a thin charcoal outline.
- Verify the icon against light and dark backgrounds at browser tab size.

## Capabilities

### Modified Capabilities

- `project-alwaldend-com` covers favicon contrast across the site.

## Impact

The shared Hugo partial continues to publish a fingerprinted SVG favicon using
the canonical S8 artwork.
