# Verification evidence

Observed 2026-09-27 in the isolated icon-generator feature worktree.

## Behavior and compatibility

E2E cases preceded implementation and failed on the missing regional flags.
The implemented controls pass project E2E checks for invalid values, inactive
algorithms, zero-strength compatibility, average initial density, spatial
contrast, parent/color retention, paired births, replay, one-cell and narrow
grids, partial cells, scale endpoints, and density endpoints. The original
forced disk-write failure case is unchanged and passes.

Eight E2E example PNGs were replayed byte-for-byte from their argument manifest.
The new baseline exactly matches the previous mixed-growth variant 10.
Source review confirms centered bounded seeding probabilities, local clamped
reproduction probabilities, no births at zero base probability, and retained
regional samples when the initial grid is empty or full. Field evaluation uses
constant storage rather than a full-resolution noise buffer.

## Comparison renders

Ten final PNGs, their initial-population PNGs, a replay/checksum manifest, and
an HTML gallery are retained under ignored
`out/icon-generator-plan/regional-replication-renders/`.
All use 1024 by 1024 white one-pixel squares on `#212121`, average initial
density 0.01, base reproduction probability 0.5, seed 5333116260189394413,
inheritance 0.65, forward bias 0.65, crowding 0.75, and branching 0.2.

| Render                    | Scale | Seeding variation | Growth variation | Initial cells | Final cells |
| ------------------------- | ----- | ----------------- | ---------------- | ------------- | ----------- |
| Prior variant 10 baseline | 128   | 0                 | 0                | 10535         | 185481      |
| Regional seeding          | 128   | 1                 | 0                | 10479         | 155697      |
| Regional growth           | 128   | 0                 | 0.04             | 10535         | 189945      |
| Subtle regions            | 128   | 0.4               | 0.04             | 10461         | 192439      |
| Balanced regions          | 128   | 0.75              | 0.08             | 10444         | 212582      |
| Broad regions             | 256   | 0.75              | 0.08             | 10564         | 219942      |
| Fine regions              | 64    | 0.75              | 0.08             | 10527         | 214539      |
| Strong regions            | 128   | 1                 | 0.15             | 10479         | 261889      |
| Very broad regions        | 384   | 1                 | 0.12             | 10698         | 275620      |
| Gentle broad regions      | 256   | 0.9               | 0.035            | 10581         | 189585      |

Visual inspection of gentle broad regions shows larger quiet areas and gradual
transitions into branching growth with internal gaps. Very broad, stronger
variation forms denser masses, demonstrating why restrained growth variation
is useful for this visual goal. All PNG dimensions and palette values were
decoded and checked; every final population retains paired-birth parity.

## Bounds and integration

A 4096-square, size-one run with density 0.03, base probability 1, field scale 1,
both variations 1, seed 42, and variant 10's four placement controls completed
in 5.03 seconds with peak RSS 148344 KiB on this host. Its decoded dimensions
and black/white palette were correct; 15983865 cells were occupied (95.27%).
These measurements are observations, not performance guarantees.

Strict OpenSpec validation and the project E2E, project OpenSpec, and shared
Hugo site tests passed. Exact-candidate formatting, quality, affected lint,
build, and publication evidence belongs to the repo-delivery receipts.
