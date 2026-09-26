## Context

The preceding change was locally prepared but not published when the user
clarified that diagonal neighbors must count. Its four-neighbor evidence
remains historical; this change supersedes that adjacency definition.

## Goals / Non-Goals

Count the eight neighboring grid positions equally while retaining the same
bounded swap process, sampled cell count, and default output.

## Decisions

Extend the neighbor count to both diagonals above and below each cell, clipping
at the grid bounds. Change the executable-level measurements before changing
the renderer. Keep the factor range and semantics unchanged.

## Risks / Trade-offs

Diagonal circle cells are neighbors but their circular shapes do not physically
touch at the corners. Documentation must distinguish adjacency from physical
contact. Stronger factors guarantee nondecreasing total eight-neighbor pairs;
horizontal/vertical contacts alone can decrease.

## Implementation evidence

Eight-direction E2E measurements and additional small grids were written
before the renderer change. Those cases also passed the prior four-neighbor
implementation, so they establish invariants rather than uniquely identifying
the neighbor policy. Source review additionally verified the bounded 3-by-3
traversal excludes the center and includes each valid diagonal once.
The updated E2E suite and affected semantic lint pass.

Renderer SHA-256:
`f744130a10bc9e5934c0cac5dc19f6e6dd02ed5897ba0e56b9b2a0805d904de4`.
With the same white-circle settings and seed as the preceding change, factors
0, 0.25, 0.75, and 1 retain 189 cells and produce 151, 351, 445, and 471 total
eight-direction neighboring pairs. The new clustered image was inspected;
zero factor remains byte-identical to the user's original independent icon.

The repeated 4096-by-4096 size-1 performance check completed in 10.447 seconds
with peak child RSS of 86,116 KiB on this host. Its PNG is retained under
ignored task outputs as `cluster-maximum-diagonal.png`.

Initial strict validation rejected an unintended scenario-title rename. The
existing title was restored while preserving the revised adjacency behavior.
