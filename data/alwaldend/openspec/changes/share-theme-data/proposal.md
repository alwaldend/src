## Why

Website Sass currently owns the shared diagram colors, and native diagrams
recover those values through browser-computed CSS. Shared identity data should
have one owner independently of the website implementation.

## What Changes

- Store neutral colors, explicit light/dark roles, and diagram typography in
  `data/alwaldend/assets/theme.json`.
- Generate website Sass and Mermaid palette/configuration with the existing
  JSON template rule. Keep consumer-specific mappings with their consumers.
- Pass diagram configuration, palettes, font, and license through one Bazel
  theme target; preserve the separate historical hand-drawn default.
- Verify current website/application Sass and native SVG/raster appearances without changing colors.

## Capabilities

### New Capabilities

- `shared-theme`: One authoritative Alwaldend identity theme for site and diagrams.

## Impact

Identity data, website assets, Mermaid theme presets and native rendering,
and application stylesheet inputs. PR #211 separately owns the spatial renderer,
Mermaid upgrade, and static snapshot capability. No new dependency,
deployment, or independent publication of the data tree is required.
