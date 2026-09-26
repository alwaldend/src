## Context

Current placement grows strands. The prior eight-neighbor compact algorithm is
available in task history at commit 9ff408267d2c5f4de6bb36fdcee491c377696b62.
See proposal.md for the user's three-way selection and reproduction request.

## Goals / Non-Goals

Restore compact output as the default, retain explicit strand output, and implement
literal two-child reproduction without deleting visible parent cells. Keep
all eight neighboring directions and bounded, deterministic execution.

## Decisions

Use `--clusterization-algorithm compact|strands|replication`, default compact
as explicitly requested by the user.
One initial occupancy sampler feeds all algorithms and preserves the original
random stream's palette draws. Compact applies the prior contact-improving
swaps. Strands clears the sampled grid and grows the sampled cell budget using
the current unchanged random sequence.

For replication, density defines the initial independent population and the
factor is the reproduction probability. A queued cell gets one decision. A
successful reproduction chooses two distinct, uniformly sampled empty neighbors
of that parent, marks them occupied, and queues them. With fewer than two empty
neighbors it adds none. A death adds none and retains the visible parent.
Children receive the same one-time decision. Roots are queued in row order;
processing is FIFO, and occupancy is updated immediately to prevent duplicates.

Verdict: proceed with user-confirmed initial population semantics. A fixed final count would
require extra roots after extinction or partial births near an odd budget,
which changes the requested reproduction/death model.

## Risks / Trade-offs

- High reproduction probabilities can fill most of the canvas. Document the
  distinction between initial density and final coverage and inspect previews.
- Each cell is queued at most once, so the queue is bounded by the grid size.
  A preallocated uint32 queue adds at most 64 MiB to the existing raster and
  boolean-grid allocations. Measure maximum-grid replication cost.
- Preserve strand output and restore compact output by comparing existing PNGs.
- When two adjacent offspring cannot fit, the parent stops without a partial
  birth; tiny-grid E2E cases verify that rule and retain all parents.

## Verification

Write E2E selection, replay, compact adjacency, replication parity, parent
retention, exact tiny-grid births, descendants, and invalid-selector cases
before implementation. Keep the forced disk-write failure case unchanged.
Retain one example per algorithm with replay arguments and pixel hashes.
