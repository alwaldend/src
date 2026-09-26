## Context

See proposal.md for the visual problem. Existing algorithms share a sampled
occupancy grid and coordinate-stable palette stream. All rendering uses Go's
standard library. The user approved one-child extensions in a separate algorithm.

## Goals / Non-Goals

**Goals:** persistent curves, visibly unequal branches, open gaps, deterministic
bounded growth, and unchanged existing algorithms.

**Non-Goals:** biological simulation, guaranteed connectivity, new dependencies,
or another collection of tuning flags before visual comparison.

## Decisions

- Add `internal/generator/tips.go`. Reuse independent initial sampling and PNG
  rendering. Density controls roots, not final coverage; recommend much smaller
  densities than replication. Compact stays default.
- Use a continuous position, heading, and slowly varying curvature, rasterized
  onto eight-connected cells without antialiasing. Grid-direction resampling was
  rejected because it produces visibly jagged independent turns.
- Clusterization scales the root growth budget. Zero preserves independent
  placement. A fork splits the remaining budget unequally; the main tip keeps
  the larger share. Side branches begin at an angle and inherit curvature.
- Retain all sampled roots, shuffle their processing order, and grow each root's
  main tip before its deferred branches. A compact root index list avoids an
  entire canvas of floating-point tip state. Each successful step occupies a
  new cell, and budgets and shrinking branch allocations bound work.
- Stop on occupied destinations, boundaries, or contact with existing strands.
  Exempt the recent local trail and a small fork junction so rasterized curves
  and branches can leave their own stem. There is no wraparound or restart.
- Keep replication options exclusive to replication. Compare densities and
  growth lengths before considering separate tip controls.

## Risks / Trade-offs

- Organic appearance is subjective → inspect full-size renders; report visible
  structure rather than treating command success as visual acceptance.
- Dense starting populations block tips immediately → document sparse examples
  and retain initial cells rather than silently changing density semantics.
- Main tips claim space before deferred branches → seeded root shuffling avoids
  systematic scan-direction priority; collisions can shorten any branch.
- Continuous paths may repeat a raster cell → bound substeps per advance and
  never spend an unbounded loop waiting for a vacancy.
- More length does not guarantee more coverage → collisions and branching change
  placements; report actual occupied coverage with each render.
