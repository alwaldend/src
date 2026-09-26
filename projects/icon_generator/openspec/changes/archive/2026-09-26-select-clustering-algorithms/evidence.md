# Verification evidence

Observed 2026-09-26 in the isolated icon-generator feature worktree.

The user confirmed initial-population density for replication and selected
compact as the default. Help, defaults, E2E assertions, and the maintained
requirements reflect those choices.

## Executable behavior

E2E cases were written before implementation. The initial run failed because
`--clusterization-algorithm` was undefined. After implementation,
`bazel_agent bazel test //projects/icon_generator:test` passed. The suite
checks compact default selection, zero-factor equivalence, deterministic replay,
compact neighbor-pair monotonicity, strand morphology, parent/color retention,
paired births, insufficient vacancies, descendants, and empty/full grids.
The existing forced disk-write failure test remains unchanged and passes.

The six example commands in the E2E manifest were replayed into a fresh directory;
all six PNGs matched byte-for-byte. Source review confirmed that each newborn
occupies a distinct vacancy, only eight bounded neighbors are considered, and
each cell enters the finite FIFO queue once. No new dependency was added.

## Prior output compatibility

Regenerated 512 by 512 white circles on `#212121`, size 16, density 0.2,
factor 0.75, seed 5333116260189394413, selecting each algorithm explicitly:

| Algorithm | Prior image                   | PNG SHA-256, unchanged                                             |
| --------- | ----------------------------- | ------------------------------------------------------------------ |
| compact   | Eight-neighbor compact output | `8bf0d122064ace1cf661bec74723e6ac066872f5ab61eea80fda9edcd8734b60` |
| strands   | Organic strand output         | `b28d0b8bad8503098f85f9378fdadae1a80ade23c14f2f735e51581b5aa44cf8` |

## Replication previews and bounds

Generated 1024 by 1024 previews with individual white square pixels on `#212121`,
density 0.03, and seed 5333116260189394413. The sampled initial count is 31,600.
Compact and strands preserve it. Replication produces these populations:

| Reproduction probability | Final cells | Final coverage |
| ------------------------ | ----------- | -------------- |
| 0                        | 31,600      | 3.01%          |
| 0.25                     | 62,298      | 5.94%          |
| 0.4                      | 143,876     | 13.72%         |
| 0.55                     | 616,672     | 58.81%         |

Visually inspected the 0.4 preview: discrete branching groups with dark gaps,
using the requested one-pixel geometry. Higher reproduction probability can
consume those gaps, consistent with the documented initial-density semantics.
Preview PNGs and their replay manifest are retained under ignored
`out/icon-generator-plan/algorithm-examples/`.

A 4096 by 4096, size-one replication run with density 0.03, probability 1,
and seed 42 completed in 1.46 seconds with peak RSS 150,760 KiB on this host.
Decoded dimensions and black/white palette were correct; it occupied 16,544,114
cells (98.61%). These are observations, not performance guarantees.

Strict change validation, the project E2E target, project spec validation, and
the shared Hugo site test passed. Repository formatting completed, and the
resulting diff contains only task-owned changes. Exact-candidate repository quality, affected
lint, project/spec checks, site build, and publication are recorded in the
repo-delivery receipts; disposable logs and images remain outside Git.
