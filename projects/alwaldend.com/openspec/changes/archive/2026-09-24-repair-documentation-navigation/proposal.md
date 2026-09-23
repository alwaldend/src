## Why

OpenSpec Markdown has a heading but no front matter, leaving blank page titles and navigation rows. Missing section pages also flatten its directory hierarchy. The shared UI needs review in collapsed and expanded states and across all rendered page types.

## What Changes

- Derive missing documentation titles and section pages from packaged source during the site build.
- Preserve canonical Markdown, existing metadata, source paths, and anchors.
- Fix confirmed navigation, spacing, and responsive defects across page types.
- Verify real browser behavior with reproducible scripts, screenshots, and measurements.

## Capabilities

### Modified Capabilities

- `project-alwaldend-com`: Nonempty documentation titles, navigable hierarchy, and responsive navigation states.

## Impact

The main site's packaging, shared styles, and heading rendering. No canonical infrastructure or OpenSpec content changes are required.
