## Context

Variant 10 combines inheritance 0.65, forward bias 0.65, crowding 0.75, and
branching 0.2. The user wants less uniformity across the canvas, so the selected
experiment adds regional reproduction and seeding rather than family traits.

## Goals / Non-Goals

Provide seeded smooth spatial variation while preserving average initial
density, paired births, parent colors, deterministic replay, and off-by-default
compatibility. Avoid new dependencies, full-resolution field buffers, changes
to compact/strands, or guarantees of thin branches at high local fertility.

## Decisions

Use three options: field scale 1..4096 logical cells (default 128), growth
variation 0..1 (default 0), and seeding variation 0..1 (default 0).

One deterministic field combines smoothstep-interpolated hashed lattice values
at two spatial scales, weighted 3:1. Evaluate it on demand. A bounded grid pass
measures its mean and extrema without allocating a full-resolution buffer.
Coordinates are logical grid cells, so pixel-size scales the pattern visually.

For reproduction, center the field and divide by the larger distance to its
observed extrema. Add growth-variation times that value to the base probability
and clamp to [0,1]. Base clusterization zero always disables births. Clamping
near endpoints can change the mean reproduction probability.

For seeding, map the shared raw field to a nonnegative fourth-power weight.
Center these weights around their measured grid mean, then scale the deviation
by the largest gain that keeps probabilities in [0,1]. Interpolate this contrast
with seeding-variation. The mean remains the requested density; zero/full
density and constant fields remain uniform. Individual occupied counts remain
random, not exact. Count clipped edge cells once, like other logical cells.

Regional seeding remains active at zero reproduction probability, allowing the
initial population to be inspected and compared with later growth. It can
therefore change zero-clusterization output when explicitly enabled. Growth-only
variation has no effect at zero base probability. Return a sampled grid even if
regional sampling produces zero/full occupancy, avoiding renderer fallback to
the uniform sample. Both variation strengths zero use the original path.

## Risks / Trade-offs

- Favorable areas can become dense. Inspect restrained and stronger settings
  alongside the original variant and report final coverage.
- One normalization pass adds work only when regional controls are active;
  measure maximum-grid time and retain constant field memory.
- The spatial scale is approximate, not a region diameter or periodic tile size.
- Average density means average probability, not an exact sample count. Validate
  mean behavior statistically over large fixed-seed grids and inspect regional
  contrast separately. Write E2E coverage before implementation.
