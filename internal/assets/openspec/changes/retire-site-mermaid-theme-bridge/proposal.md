## Why

The neutral Mermaid theme is already owned by `//data/alwaldend:mermaid_theme`.
A compatibility struct and Sass build target still expose it through the site,
leaving an unnecessary site dependency for the separately developed 3D renderer.

## What Changes

- Remove the unused site-owned theme struct, palette target, and Sass entry point.
- Keep native SVG/WebP consumers on the existing shared theme bundle.
- Make the 3D renderer consume that bundle directly in PR #211.

## Capabilities

This build-ownership cleanup preserves appearance and uses `skip_specs: true`.

## Impact

The independent PR changes only the site compatibility bridge and shared-theme
documentation. The 3D consumer update remains in its own feature branch. Existing
colors, configuration, fonts, dependencies, and publication boundaries are retained.
