# Verification evidence

Observed 2026-09-26 in the isolated icon-generator feature worktree.

## Behavior and compatibility

E2E additions preceded implementation and failed on the missing replication
control flags. The implemented controls pass the project E2E suite, including
invalid values, unchanged inactive algorithms, retained parents and colors,
paired births in constrained grids, replay, and sparse crowding comparisons.
The original forced disk-write failure test remains unchanged and passes.

Seven E2E example PNGs replayed byte-for-byte from their argument manifest.
Zero controls also exactly reproduce the prior 1024-square replication PNG
(density 0.03, probability 0.4, seed 5333116260189394413). Source review confirms
positive vacancy weights, sibling contact accounting, eight bounded neighbors,
and a single queue entry per occupied cell. All four controls default to zero;
compact remains the default algorithm.

## Comparison renders

Ten PNGs and an argument/checksum manifest are retained under ignored
`out/icon-generator-plan/replication-growth-renders/`, alongside an HTML gallery.
All use 1024 by 1024 output, individual white square pixels on `#212121`,
initial density 0.01, reproduction probability 0.5, and seed 5333116260189394413.

| Render                   | Inheritance | Forward bias | Crowding | Branching | Final cells |
| ------------------------ | ----------- | ------------ | -------- | --------- | ----------- |
| Baseline                 | 0           | 0            | 0        | 0         | 161389      |
| Forward, random headings | 0           | 1            | 0        | 0         | 157303      |
| Inherited, gentle        | 1           | 0.5          | 0        | 0         | 171615      |
| Inherited, strong        | 1           | 1            | 0        | 0         | 173107      |
| Crowding only            | 0           | 0            | 1        | 0         | 192047      |
| Sideways forks only      | 0           | 0            | 0        | 0.35      | 158421      |
| Direction and space      | 1           | 0.8          | 1        | 0         | 195301      |
| Occasional forks         | 1           | 0.8          | 1        | 0.1       | 195903      |
| Frequent forks           | 1           | 0.8          | 1        | 0.5       | 185581      |
| Mixed growth             | 0.65        | 0.65         | 0.75     | 0.2       | 185481      |

Visual comparison of baseline and occasional forks shows more outward,
branching growth with the combined controls. Dense patches remain, consistent
with two-child reproduction and the documented absence of a thin-strand
guarantee. All outputs retain exactly the requested dimensions and two colors.

## Bounds and integration

A 4096 by 4096 size-one run with density 0.03, probability 1, inheritance 1,
forward bias 1, crowding 1, branching 0.2, and seed 42 completed in 3.33 seconds,
with peak RSS 149872 KiB. Decoding confirmed the expected dimensions and palette;
16264554 cells were occupied (96.94%). This is an observation on this host,
not a performance guarantee. Packed headings retain the 64 MiB maximum queue.

Strict change validation and project E2E, project OpenSpec, and shared Hugo
site tests passed. Repository formatting completed and its diff was inspected.
Exact-candidate quality, lint, build, and publication evidence belongs to the
repo-delivery receipts in task scratch.
