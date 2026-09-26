## Context

The renderer makes one occupancy draw and one palette draw per grid cell using
a seeded PCG stream. See proposal.md for the requested grouping control.

## Goals / Non-Goals

Preserve default images, seeded replay, and the sampled occupied-cell count.
Favor four-neighbor contact without wrapping across image edges. A single
connected component and a particular cluster size are not guarantees.

## Decisions

Keep the initial independent occupancy sample. For positive clusterization,
store it in a boolean grid and propose swaps between random cells, accepting
only occupied-to-empty moves that strictly increase the total number of
horizontal and vertical occupied neighbor pairs. Discount the source cell when
counting the destination's neighbors. Use a separate PCG stream for swaps.

The factor selects a prefix of at most eight proposals per cell. Increasing
the factor with other options unchanged therefore cannot reduce contacts or
change the occupied-cell count. Default zero bypasses the extra grid. Palette
draws stay at their original cell positions and remain independent of swaps.

Verdict: proceed with bounded swaps. Neighbor-copying probabilities are simpler
but can shift realized density or collapse to a uniform image at full strength.
Sorting a smoothed noise field preserves counts but costs more memory and
does not guarantee monotonic contact counts. Swaps provide directly measurable
behavior with linear bounded work and at most 16 MiB of additional grid memory.

## Risks / Trade-offs

- Clustering can plateau or leave isolated cells; factor 1 is the strongest
  supported effort, not a promise of complete connectivity.
- Very small or nearly empty/full grids may not change. Check those cases and
  density endpoints through the executable.
- Maximum-size grids add work. Measure a 4096 by 4096, size-1, factor-1 run.
- Increasing contacts does not preserve individual shape colors when cells
  move. Colors remain independently selected from the configured palette.

## Migration Plan

No migration is needed: omitted or zero clusterization keeps prior behavior.
Write command-level acceptance checks before implementation, compare decoded
counts and contacts, inspect repeatable PNG examples, then deliver through the
existing pull request.

## Implementation evidence

On 2026-09-26, acceptance checks were written first and failed because the
executable did not recognize `--clusterization`. After implementation, all
28 selected tests passed, including project E2E, OpenSpec, and repository
quality checks. Affected semantic lint and the Hugo site build passed; the
rendered usage page contains the new option. Final publication gates are
recorded by repo-delivery against its exact prepared candidate.

Verified renderer SHA-256:
`2e93f6488640f5f80e6d27b9dd5029aae6d53dbca47e8d7b93a1d64fea4402d6`.
The three prior example PNGs remained byte-identical, and all four current
example manifest commands reproduced their artifacts exactly.

For 512 by 512 white circles on `#212121`, size 16, density 0.2, and seed
5333116260189394413, factors 0, 0.25, 0.75, and 1 produced respectively 78,
177, 225, and 232 touching cell pairs. All four contained 189 occupied cells.
The independent and clustered images were inspected visually.

A 4096 by 4096, size-1, density-0.2, factor-1, seed-42 run completed in 7.897
seconds with measured peak child RSS of 86,244 KiB on this host. This is an
observed performance check, not a portable timing guarantee. The generated
PNG was retained in ignored task outputs and shown to the user.
