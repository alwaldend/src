## Context

The rule, CLI fallback, and encoder fallback currently use scale 2. Blog WebPs are generated projections with freshness checks.

## Goals / Non-Goals

Use scale 1 consistently across entry points. Preserve SVG rendering, quality 90 encoding, selected themes, and padding semantics.

## Decisions

Change all three defaults together. Keep existing native comparison fixtures at explicit scale 2 so they continue checking higher-resolution overrides; let the existing classic geometry check exercise the new default. Regenerate blog assets through their update rules.

## Risks / Trade-offs

Lower pixel density reduces sharpness on high-density screens; consumers can explicitly request scale 2. Failure cases to verify before implementation: divergent entry-point defaults, incorrect natural geometry, broken explicit scale overrides, altered padding or palettes, stale checked-in WebPs, and clipped labels.
