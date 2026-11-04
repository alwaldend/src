## Why

Capturing the same site indexes required manual browser setup and scratch
automation. A repository command should produce repeatable screenshots directly
from the built site with its existing pinned browser tools.

## What Changes

- Add `//projects/alwaldend.com:site_screenshots` for full-page PNG capture.
- Default to the homepage, project landing index, and project documentation index.
- Support selected paths, viewport dimensions, color scheme, and output directory.
- Render local build assets without starting a server or accessing external URLs.
- Write a manifest and exercise the command with an offline browser fixture.

## Capabilities

### New Capabilities

- `site-screenshots`: Reproducible screenshot artifacts from built Hugo pages.

### Modified Capabilities

None.

## Impact

The screenshot command lives under `projects/alwaldend.com/cmd/screenshots` and
uses the existing Puppeteer and Chrome pins. The owning README documents it.
No new external dependency, deployment, or site presentation change is required.
