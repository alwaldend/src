## Context

The user selected winding, branching strands with open gaps after reviewing
the compact eight-neighbor clustering. The prior maximum-contact objective
conflicts with that outcome; see proposal.md.

## Goals / Non-Goals

Keep the sampled occupied-cell count, zero-factor image compatibility, and
deterministic replay. At sparse densities favor thin strands and preserve
space between them. Full density necessarily fills every cell.

## Decisions

Use the original random stream to determine the cell budget, then grow strands
on an empty grid using a separate seeded stream. Factor controls maximum
strand length from 1 to 65 cells with a quadratic response. Choose random
starting positions, favoring uncrowded cells, and prefer forward and gradual
turns when choosing among eight neighboring positions.

Growth prefers candidates with at most two occupied neighbors. Occasionally
resume at an earlier point in the current strand with a turned heading to
form a branch. Bound each strand's history and all placement attempts. When
growth cannot continue, start another strand. Bounded random seed probes with
a monotonic free-cell scan preserve progress even near full density.

Verdict: proceed with directional growth. Softening the old contact reward
still optimizes compactness; a smoothed noise threshold still tends toward
islands. Direction and crowding preferences directly control the requested
strands without a new dependency. The prior monotonic contact guarantee is
deliberately superseded by the user's visual goal.

## Risks / Trade-offs

- Tiny or dense canvases cannot maintain open strands; retain exact count and
  density endpoints, and document that gaps shrink as density rises.
- Quantitative morphology checks alone cannot establish visual quality. Inspect
  the user's 512-by-512 circle pattern before delivery.
- Positive factors may change all positions; only zero preserves old pixels.
- Keep extra memory bounded by one boolean grid plus fixed strand history;
  measure maximum-grid runtime and verify skinny/full grids terminate.

## Implementation evidence

The morphology checks were written first and failed the compact renderer: for
seed 0 at density 0.2 and factor 0.75 on a 64-by-48 cell grid, 299 of 625
occupied cells had more than four neighbors, with 205 filled 2-by-2 blocks.
The strand implementation passes those checks for all five selected seeds,
plus the existing density, geometry, palette, replay, and output-failure cases.
Affected semantic lint and strict OpenSpec validation passed.

Source review added a four-consecutive-failure bound so blocked growth cannot
spend an entire strand budget repeatedly trying a crowded origin. The E2E suite
passed again after that bounded-work improvement. Pre-format renderer SHA-256:
`ec0599d4a7154b221725e8018ea60596cbb61044e5ea6c49d006091dce01846f`.

The user's 512-by-512, size-16, density-0.2, white-circle example with seed
5333116260189394413 retains 189 cells at all factors. At factor 0.75 it has
214 eight-neighbor pairs, no isolated cells, no cells with more than four
neighbors, and no filled 2-by-2 blocks. The preview was inspected visually and
shown to the user. Zero factor matches the prior independent image exactly.
All four E2E example commands replay exactly; the three original zero-factor
examples remain byte-identical.

On this host, 4096-by-4096 size-1, factor-1, seed-42 runs took 3.268 seconds at
density 0.2 and 11.020 seconds at density 0.99, with measured peak child RSS
86,244 KiB. These observations establish bounded practical behavior on this
host, not a portable timing guarantee. Final aggregate validation and publication
are owned by repo-delivery receipts for the prepared candidate.
