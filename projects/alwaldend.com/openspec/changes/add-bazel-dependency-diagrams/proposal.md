## Why

The site's build relationships need an inspectable diagram generated from the
actual Bazel graph. Visitors also need documentation for the new analyzer.

## What Changes

- Add bounded Mermaid, SVG, and WebP companion targets analyzing `:site`.
- Register the analyzer landing and package its README and example diagram.

## Capabilities

### New Capabilities

- `site-dependency-diagrams`: Build a reproducible view of the site's dependencies.

### Modified Capabilities

None; landing publication follows the existing contract.

## Impact

Site BUILD and documentation, shared project registry, and analyzer-owned
landing content. Conversion is owned by the analyzer's `add-dependency-diagrams`
change. No deployment, DNS change, or new dependency.
