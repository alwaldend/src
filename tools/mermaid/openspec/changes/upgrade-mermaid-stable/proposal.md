## Why

The user requested the latest stable Mermaid before adding spatial diagrams.

## What Changes

- Pin Mermaid 12.2.0 and CLI 12.0.0 and regenerate the pnpm lock.
- Port the maintained Dagre spacing patch and validate existing consumers.
- Migrate native font input to CLI 12's licensed `customFontCSS` API and patch
  its `.tff` extension typo. Preserve the previous flowchart width defaults
  explicitly in the maintained themes.

## Capabilities

No behavior contract changes are intended; this dependency update skips specs.
The spatial renderer is tracked by projects/mermaid_3d's owner-local change.

## Impact

The tools pnpm workspace, Mermaid patch, and existing SVG/raster consumers.
