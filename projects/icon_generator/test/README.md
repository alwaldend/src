# Icon generator acceptance

The command-level checks are written before the implementation. They run the
built executable with an empty `PATH` and decode its PNGs using Go's standard
library. No image-processing program is involved.

Failure cases:

- Invalid regional field scale or variation values, even for inactive algorithms.
- Regional seeding changes average occupancy probability or fails to create broad
  spatial contrast; clipped cells or degenerate fields break normalization.
- Zero controls change old output, growth variation creates births at zero base
  probability, or an all-empty/full regional sample falls back to uniform sampling.
- Regional growth loses parents, recolors them, adds an odd number of children,
  fails replay, or mishandles narrow, partial-cell, empty, and full grids.

- Missing flag values, unknown flags, and unexpected positional arguments.
- Zero, negative, excessive, or noninteger dimensions and logical pixel size.
- Negative, excessive, NaN, or infinite density.
- Negative, excessive, nonnumeric, NaN, or infinite clusterization.
- Empty or unknown algorithm selection, including when clusterization is zero.
- Omitted algorithm selection differs from compact for positive factors.
- Nonfinite or out-of-range replication growth controls are accepted.
- Zero growth controls change legacy output; inheritance without forward bias
  changes placement; growth controls affect inactive replication or other algorithms.
- Growth preferences prevent a valid two-child birth or remove/recolor parents.
- Crowding preference fails to reduce contacts per cell across representative
  sparse fixed-seed populations; controls fail deterministic replay.
- Selecting strands changes current images, or zero factor depends on the algorithm.
- Compact placement changes the sampled count or reduces eight-neighbor adjacency.
- Replication removes parents, creates single/duplicate/nonadjacent children,
  reseeds after extinction, skips descendants, wraps boundaries, or fails to terminate.
- Strand clustering changes the occupied-cell count, produces dense blobs at sparse density,
  wraps adjacency across canvas edges, or changes default independent placement.
- Clustering fails on single-cell, single-row, single-column, empty, or full grids.
- Empty palettes or palette entries, malformed hex, and unsupported shapes.
- Negative, overflowing, or malformed seeds; explicit zero must remain valid.
- Existing output files, directories, and symlinks must remain unchanged.
- Missing or unwritable output directories must fail without creating parents.
- A write failure must remove the incomplete file. The Linux E2E harness
  enforces a zero file-size limit in a subprocess to exercise this path.
- Entropy acquisition and close errors must retain their operation context;
  these platform failures are reviewed in source, without adding production
  fault-injection switches solely for testing.

Geometry checks use small hand-derived masks for odd and even circles and
partial cells. Other checks cover dimensions, defaults, help, exact palette
membership, alpha, density endpoints, and explicit or reported seed replay.
Randomness checks avoid assertions that can fail by chance.

Strand comparisons use fixed seeds and count occupied cells and horizontal,
vertical, and diagonal neighbor pairs in decoded output. They require unchanged
counts, few isolated cells, few crowded cells, and few filled two-by-two blocks
on representative sparse grids. Seed replay and palette checks also run with
clustering enabled. Visual inspection complements these morphology checks.

Algorithm checks compare decoded images, population counts, replay, and compact
adjacency. Hand-derived tiny-grid replication cases cover exactly two births,
no birth when only one space remains, and extinction with no initial population.
Larger replication cases check parent retention, offspring parity, and growth
beyond the first generation. All three algorithms retain repeatable PNG examples.

Replication growth checks exercise zero, intermediate, and maximum controls,
including fully constrained tiny grids that must still produce paired births.
Sparse crowding comparisons use four fixed seeds and compare neighbor contacts
per occupied cell. Direction and fork appearance are also inspected in rendered
comparisons; these controls do not guarantee a morphology for every seed.

`TestExamples` retains PNGs and `manifest.json` in Bazel's undeclared test
outputs. The manifest records executable arguments, seeds, and checksums of
decoded dimensions and RGBA pixels. Run the recorded commands in a fresh
output directory to reproduce the images.
