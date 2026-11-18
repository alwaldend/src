## Context

The shared theme target already bundles configuration, JSON light/dark palettes,
the pinned font, and its license through `MermaidThemeInfo`. Native SVG/WebP
rules consume it. The old site struct and compiled Sass palette have no trunk
consumers; only the unmerged 3D feature uses that bridge.

## Decisions

Use the existing bundle as the consumer contract. Remove the compatibility
target and its unused Sass source instead of moving a second label list or
duplicating the palette mappings. Site SVG/WebP convenience macros retain
their consumer-specific defaults while selecting the shared data target.

Keep the cleanup PR based on master and the 3D provider update in PR #211.
Neither needs a new project, dependency, palette, or deployment.

## Verification

Search trunk consumers before removing the bridge. Build and test the native
Mermaid and website consumers, compare all twelve gallery SVG hashes, and run
repository quality and affected semantic lint. Validate the 3D consumer's
models and light/dark images on its own branch.
