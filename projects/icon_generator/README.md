---
title: Icon generator
description: Generate random patterns of colored squares or circles as PNG icons
statuses:
  - active
languages:
  - go
tags:
  - cli
  - images
---

Icon generator is a Go CLI that creates random grid patterns as PNG images.
Choose the canvas dimensions, background, palette, cell size, density,
clustering algorithm and strength, and square or circular shapes. Rendering and PNG encoding use
only Go's standard library; no external graphics programs or libraries are required.

## Usage

From the repository root, create an output directory and generate an icon:

```sh
mkdir -p out/icon_generator
bazel_agent bazel run //projects/icon_generator -- \
  --width 256 --height 256 \
  --background '#f5f1e8' \
  --colors '#223843,#d77a61,#e3b23c' \
  --pixel-size 16 --pixel-shape square \
  --density 0.4 --seed 42 \
  --output "$PWD/out/icon_generator/squares.png"
```

For circles on a transparent background:

```sh
bazel_agent bazel run //projects/icon_generator -- \
  --background transparent \
  --colors '#223843,#d77a61,#e3b23c' \
  --pixel-shape circle --seed 42 \
  --output "$PWD/out/icon_generator/circles.png"
```

Use an absolute output path with `bazel run`. The output directory must
already exist. Existing files and symlinks are never overwritten; choose a
new path when repeating a command.

Use `--clusterization 0.75` for compact groups. To grow winding, branching strands
while keeping the same number of occupied cells, also select
`--clusterization-algorithm strands`. Sparse densities leave more room for open gaps:

```sh
bazel_agent bazel run //projects/icon_generator -- \
  --width 512 --height 512 \
  --background '#212121' --colors '#ffffff' \
  --pixel-size 16 --pixel-shape circle \
  --density 0.2 --clusterization 0.75 --clusterization-algorithm strands --seed 42 \
  --output "$PWD/out/icon_generator/clustered.png"
```

For replication from a sparse initial population of individual white pixels:

```sh
bazel_agent bazel run //projects/icon_generator -- \
  --width 1024 --height 1024 \
  --background '#212121' --colors '#ffffff' \
  --pixel-size 1 --pixel-shape square \
  --density 0.03 --clusterization 0.4 --clusterization-algorithm replication \
  --replication-inheritance 0.9 --replication-forward-bias 0.8 \
  --replication-crowding 0.75 --replication-branching 0.15 \
  --replication-field-scale 128 \
  --replication-growth-variation 0.08 --replication-seeding-variation 0.75 \
  --seed 42 --output "$PWD/out/icon_generator/replication.png"
```

Build the executable with `bazel_agent bazel build //projects/icon_generator`.
It can also run directly, without Bazel or other programs on `PATH`:

```sh
bazel-bin/projects/icon_generator/cmd/icon_generator/icon_generator_/icon_generator --help
```

## Options

| Option                            | Default           | Meaning                                                            |
| --------------------------------- | ----------------- | ------------------------------------------------------------------ |
| `--width`                         | `256`             | Output width, 1 through 4096 pixels                                |
| `--height`                        | `256`             | Output height, 1 through 4096 pixels                               |
| `--background`                    | `#ffffff`         | Opaque `#RRGGBB` or `transparent`                                  |
| `--colors`                        | `#000000`         | Nonempty comma-separated `#RRGGBB` palette                         |
| `--pixel-size`                    | `16`              | Square side or circle diameter, 1 through 4096 pixels              |
| `--pixel-shape`                   | `square`          | `square` or `circle`                                               |
| `--density`                       | `0.4`             | Average initial cell occupancy probability, from 0 through 1       |
| `--clusterization`                | `0`               | Clustering strength or reproduction probability, from 0 through 1  |
| `--clusterization-algorithm`      | `compact`         | `compact`, `strands`, or `replication`                             |
| `--replication-inheritance`       | `0`               | Probability a child keeps its parent-to-child heading, 0 through 1 |
| `--replication-forward-bias`      | `0`               | Strength of forward and gentle-turn preference, 0 through 1        |
| `--replication-crowding`          | `0`               | Strength of preference for less crowded destinations, 0 through 1  |
| `--replication-branching`         | `0`               | Probability the second child favors a sideways fork, 0 through 1   |
| `--replication-field-scale`       | `128`             | Approximate regional variation scale, 1 through 4096 logical cells |
| `--replication-growth-variation`  | `0`               | Maximum local offset to reproduction probability, 0 through 1      |
| `--replication-seeding-variation` | `0`               | Regional starting-population contrast, 0 through 1                 |
| `--seed`                          | Fresh random seed | Unsigned 64-bit seed; zero is valid                                |
| `--output`                        | `icon.png`        | Destination PNG file                                               |
| `--help`                          |                   | Print usage without generating an image                            |

Hex colors are case-insensitive. Quote colors so the shell does not interpret
`#` as a comment. Palette colors are opaque and keep their exact RGB values.

## Rendering and reproducibility

Canvas dimensions count output image pixels. Logical pixels occupy a grid
starting at the top-left corner. A 256 by 256 canvas with size 16 contains
16 by 16 cells. Squares fill their cell; circles are centered within it and
leave the background visible at the corners. Shapes have hard rasterized
edges, so very small circles can resemble squares. Partial cells at the
right and bottom are clipped without resizing the shape or padding the image.

Initial placement independently samples each cell. Its probability is the requested
density, or a local probability when regional seeding is enabled. These local
probabilities average to density. Each shape uses a randomly selected palette entry.
Density is a probability, not an
exact shape count: zero draws only the background and one occupies every
cell. Circle corners remain background even at density one.

The algorithm determines how positive clusterization changes placement:

- `compact` (default) moves occupied cells to increase neighboring contacts.
  Higher factors make more attempts and cannot decrease the number of occupied
  neighbor pairs with the same seed and other settings. Cell count is preserved.
- `strands` grows winding strands with occasional branches. Higher factors favor
  longer strands; placement prefers uncrowded cells to leave open gaps. Cell count
  is preserved, but the neighbor count can rise or fall. Dense or tiny canvases
  have less room for open strands.
- `replication` uses density for the initial population and clusterization as the
  base probability that each cell reproduces, optionally varied across regions.
  Reproduction adds exactly two distinct
  children chosen randomly from the parent's empty adjacent cells. A cell with
  fewer than two vacancies adds none. A cell that dies adds none and remains
  visible. Children get the same one-time decision, so final coverage can exceed
  the initial density; high reproduction probabilities can fill most of the image.
  Initial cells are processed in row order, followed by children in birth order.
  Each cell is processed once, and growth stops when no decisions remain.

Growth uses all eight neighbors: horizontal, vertical, and diagonal. Opposite
canvas edges are not neighbors. Diagonally adjacent circles count as neighbors
even though their shapes do not physically touch at the corners. The result can
contain several groups and isolated cells; no factor promises one connected group.
Zero disables reproduction and preserves the original independent pattern when
regional seeding is off. Explicit regional seeding can change that initial pattern. Compact
and strands can change positions and colors within the configured palette.
Replication retains each initial cell and its color; children use the palette.

### Replication growth controls

The four placement controls default to zero, which reproduces the original
uniform-neighbor rule when regional variation is also off. They are validated for every algorithm but affect only
replication. They change placement, while `--clusterization` still controls
reproduction probability and `--density` still sets the initial population.

Inheritance gives a newborn its parent-to-child direction with the configured
probability; otherwise it gets a random direction. Roots start with random
directions. Forward bias uses those directions to favor continuing forward or
turning gently, with weaker preferences for sideways or backward growth.
Use inheritance together with forward bias to encourage persistent paths.
Inheritance alone has no visible effect when forward bias is zero.

Crowding avoidance lowers the weight of destinations surrounded by more occupied
neighbors. It excludes the parent and counts the first child when placing the
second. This is a soft preference: crowded vacancies remain eligible.

Branching controls how often the second child favors a direction 90 degrees to
either side of the first child's direction. It controls the fork angle, not the
number of children. Every successful reproduction still adds exactly two children;
if fewer than two vacancies exist, it adds none. Children can continue reproducing.

Controls do not promise thin or connected strands. High initial density or
reproduction probability can still fill gaps. Different placements can change
later vacancy availability and final coverage, so comparisons should keep the
seed, density, and reproduction probability fixed and report final coverage.

### Regional variation

Seeding and growth variation use the same seed-derived smooth field, so fertile
and quiet areas line up. `--replication-field-scale` sets their approximate size
in logical cells: 128 means about 128 output pixels with pixel-size 1, or 1024
output pixels with pixel-size 8. It is not an exact region diameter. Smaller
values vary more rapidly; larger values form broader regions.

`--replication-growth-variation` bounds the local additive change to reproduction
probability. For example, base clusterization 0.5 and variation 0.1 keep local
probabilities within 0.4..0.6. Values are clamped to 0..1, so their mean can shift
near those limits. Base clusterization zero always disables reproduction.

`--replication-seeding-variation` controls the contrast of the starting population:
zero is uniform, while one uses the strongest regional contrast that keeps all
initial probabilities within 0..1. Their average remains the requested density,
including partially clipped cells, but the sampled cell count is still random.
Seeding variation works with clusterization zero, allowing the starting population
to be inspected separately. Empty/full density remains empty/full, and a constant
field uses uniform probabilities.

Both variation strengths default to zero, preserving earlier seeded images.
Regional options affect replication only and add a normalization pass over the
logical grid without allocating a full-canvas field buffer. Strong growth
variation can fill favorable areas; moderate values preserve more internal gaps.

Successful invocations print `seed: NUMBER` on stderr. Reuse that seed with
the same settings and generator version to reproduce the decoded pixels.
PNG compression bytes may change between Go toolchain versions. Different
seeds are not guaranteed to yield different images.

Invalid options fail before output creation. Output errors return a nonzero
status with a diagnostic, and a failed write removes the incomplete file
created by the invocation. The dimension limits bound the main image buffer
at 64 MiB. Clustering uses at most another 16 MiB for the cell grid. Replication
also uses a queue of up to 64 MiB. Encoding and the Go runtime use additional memory.

## Verification

```sh
bazel_agent bazel test //projects/icon_generator:test
```

The end-to-end suite runs the built command with an empty `PATH`, decodes
its images, checks geometry and color behavior, and exercises failure paths.
Its Bazel undeclared outputs contain `squares.png`, `circles.png`,
`transparent.png`, `clustered.png` (strands), `compact.png`, `replication.png`,
`replication-growth.png`, `replication-regions.png`,
and a `manifest.json` with replay arguments,
seeds, and checksums of decoded dimensions and RGBA pixels. They can be
inspected and recreated in a fresh directory.

[Project page](https://alwaldend.com/projects/icon_generator/)
