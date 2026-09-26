# random-icon-generation Specification

## Purpose

Generate PNG icons from random arrangements of configurable colored shapes,
with control over the canvas, background, palette, density, and seed.

## Requirements

### Requirement: Generate a PNG with configurable canvas dimensions

The `icon_generator` CLI SHALL accept `--width`, `--height`, and `--output`
and write one PNG with exactly the requested dimensions. Width and height
SHALL be integers from 1 through 4096, measured in output image pixels.
Generation SHALL work without installed image-processing programs or
external runtime libraries.

#### Scenario: Generate a rectangular image

- **WHEN** the user requests width 120, height 80, and an available output
  path with otherwise valid options
- **THEN** the command exits successfully and writes a decodable 120 by 80
  PNG to that path
- **AND** no image-processing executable is required on `PATH`

### Requirement: Configure background and allowed shape colors

The CLI SHALL accept `--background` as `#RRGGBB` or `transparent`, and
`--colors` as a comma-separated nonempty list of `#RRGGBB` colors. Hex digits
SHALL be case-insensitive. An occupied cell SHALL choose uniformly from
the supplied palette entries and use that opaque color for its shape.
Unoccupied cells and pixels outside shapes SHALL retain the background.
Rendering SHALL preserve palette colors without blended edge colors.

#### Scenario: Use a single foreground color

- **WHEN** the user selects a white background and a palette containing
  only `#123abc`
- **THEN** every rendered shape pixel has opaque RGB color `#123abc`
- **AND** every remaining pixel is opaque white

#### Scenario: Preserve transparent background

- **WHEN** the user selects `--background transparent`
- **THEN** pixels not covered by a shape have alpha zero
- **AND** shape pixels have alpha 255

### Requirement: Configure grid size and shape

The CLI SHALL accept `--pixel-size` as an integer from 1 through 4096 and
`--pixel-shape` as `square` or `circle`. The grid SHALL start at the top-left
canvas corner with cells of the requested size. A square SHALL fill its
cell; a circle SHALL be centered in its cell with a diameter equal to that
size. A circle SHALL include an output pixel when that pixel's center is
on or inside the circle. Shapes SHALL have hard rasterized boundaries.
The rightmost and bottommost cells SHALL be clipped at the canvas bounds
without changing their original shape size or center.

#### Scenario: Render square cells

- **WHEN** the user requests a 32 by 16 canvas, size 16, square shapes,
  and density 1
- **THEN** two adjacent 16 by 16 squares cover the canvas

#### Scenario: Render a circle within a cell

- **WHEN** the user requests a 16 by 16 canvas, size 16, circle shapes,
  and density 1
- **THEN** the cell contains a centered circle of diameter 16
- **AND** the canvas corners retain the background

#### Scenario: Clip incomplete edge cells

- **WHEN** the user requests a 25 by 18 canvas with size 16
- **THEN** the grid contains two columns and two rows
- **AND** shapes in the last column and row are clipped to the 25 by 18
  output bounds without padding or shrinking the shapes

#### Scenario: Render the smallest logical pixel

- **WHEN** size is 1, density is 1, and the shape is either supported value
- **THEN** every output pixel receives a palette color

### Requirement: Configure occupancy density

The CLI SHALL accept a finite `--density` from 0 through 1 inclusive.
Initial placement SHALL independently sample each cell. The per-cell
probability SHALL equal density unless regional replication seeding is enabled,
in which case density SHALL equal the mean of the local probabilities. Compact and strand clustering SHALL rearrange that sample without
changing its occupied-cell count. Replication SHALL use the sample as its
initial population and MAY increase the count through reproduction. Density
SHALL represent an average initial occupancy probability, not an exact occupied-cell
count or the percentage of output pixels covered by circles.

#### Scenario: Generate only background

- **WHEN** density is 0, with any valid clusterization and algorithm
- **THEN** every output pixel retains the configured background
- **AND** replication does not introduce a population spontaneously

#### Scenario: Occupy every cell

- **WHEN** density is 1, with any valid clusterization and algorithm
- **THEN** every cell contains the selected shape
- **AND** circle boundaries still leave background visible outside circles

### Requirement: Reproduce a pattern from its seed

The CLI SHALL accept an optional unsigned 64-bit integer `--seed`, including
zero. For the same tool version, seed, and rendering settings, the decoded
output dimensions and pixel values SHALL be identical. If omitted, the
command SHALL obtain a fresh seed and report the effective seed on stderr
so the user can reproduce the result. Reproducibility SHALL concern image
content, without requiring byte-identical PNG encodings across toolchains.

#### Scenario: Replay an explicitly seeded image

- **WHEN** two invocations use the same explicit seed and rendering options
  but different available output paths
- **THEN** decoding both PNGs produces identical dimensions and pixels

#### Scenario: Replay an automatically seeded image

- **WHEN** an invocation omits the seed and the user repeats its rendering
  options with the reported seed
- **THEN** the repeated image has identical decoded dimensions and pixels

### Requirement: Provide usable defaults and help

The CLI SHALL default to width 256, height 256, background `#ffffff`,
palette `#000000`, pixel size 16, square shapes, density 0.4, a fresh seed,
and output `icon.png`. The CLI SHALL provide `--help` describing options,
defaults, accepted values, seed replay, and the distinction between output
pixels and logical grid cells. Help SHALL exit successfully without writing
an image.

#### Scenario: Generate an image with defaults

- **WHEN** the command is invoked without options and `icon.png` does not
  already exist
- **THEN** it writes a 256 by 256 PNG using black 16 by 16 square cells,
  a white background, and occupancy probability 0.4

#### Scenario: Inspect usage without generating an image

- **WHEN** the user invokes `--help`
- **THEN** the command describes all supported options and exits
  successfully without creating or changing an image file

### Requirement: Reject invalid input and report output failures

The CLI SHALL reject invalid dimensions, size, density, palette, background,
shape, seed, unknown flags, and unexpected positional arguments with a
nonzero exit status and a diagnostic on stderr. Input validation SHALL
complete before creating the output file. The command SHALL refuse to
overwrite an existing output file. File creation, encoding, and close
failures SHALL produce a nonzero exit status and a diagnostic; a failed
write SHALL remove only the incomplete output created by that invocation.

#### Scenario: Reject invalid rendering options without output

- **WHEN** a rendering option is invalid, such as a zero width, size 4097,
  density `NaN`, an empty palette, a malformed hex color, or shape `triangle`
- **THEN** the command reports the invalid option and exits unsuccessfully
- **AND** it does not create an output image

#### Scenario: Preserve an existing output file

- **WHEN** the requested output path already exists
- **THEN** the command exits unsuccessfully and reports the conflict
- **AND** the existing file remains unchanged

#### Scenario: Report an unavailable destination

- **WHEN** the output parent directory does not exist or cannot be written
- **THEN** the command exits unsuccessfully and identifies the output error

### Requirement: Configure pixel clusterization

The CLI SHALL accept a finite `--clusterization` factor from 0 through 1,
defaulting to 0, and `--clusterization-algorithm` as `compact`, `strands`, or
`replication`, defaulting to `compact`. Zero factor SHALL prevent reproduction. With regional seeding disabled,
it SHALL preserve independent placement and existing seeded image content for
every algorithm. Explicit
strand selection SHALL preserve the current seeded strand behavior.

Compact clustering SHALL preserve the sampled occupied-cell count and favor
dense groups. Increasing the factor with the same seed and other options SHALL
NOT reduce the total eight-neighbor occupied pair count. Strand clustering
SHALL preserve the sampled count and favor winding strands with occasional
branches and open gaps at sparse densities; higher factors SHALL favor longer
strands without a monotonic neighbor-count guarantee. Replication SHALL use the
factor as the probability of reproduction as specified by its lifecycle.

All algorithms SHALL use the eight surrounding positions: horizontal, vertical,
and diagonal. Adjacency SHALL NOT wrap across opposite canvas edges. Diagonal
circle cells SHALL count as neighbors even though their shapes do not touch at
corners. The CLI SHALL reject invalid factors and unknown or empty algorithm
names before creating an output file and SHALL describe selection in help.

#### Scenario: Increase contact without increasing density

- **WHEN** the user increases compact or strand clusterization with the seed and other options fixed
- **THEN** the number of occupied cells remains unchanged
- **AND** compact mode does not reduce the occupied neighbor-pair count

#### Scenario: Include corner neighbors

- **WHEN** two occupied cells meet only diagonally on the grid
- **THEN** they count as one neighboring pair for clustering
- **AND** square and circle shapes use the same grid adjacency definition

#### Scenario: Grow sparse organic patterns

- **WHEN** the user selects strands, sparse density, and strong clusterization
- **THEN** the pattern favors narrow winding strands, occasional branches, and open gaps
- **AND** identical settings and seed reproduce the pattern

#### Scenario: Preserve default output

- **WHEN** the user omits clusterization or explicitly sets it to 0, with regional seeding disabled
- **THEN** the decoded image matches the independent-placement behavior
- **AND** omitted algorithm selection remains equivalent to explicit compact for positive factors

#### Scenario: Reject an invalid factor

- **WHEN** clusterization is negative, greater than 1, nonnumeric, NaN, or infinite
- **THEN** the command exits unsuccessfully with a diagnostic and creates no image

#### Scenario: Reject an unknown algorithm

- **WHEN** the algorithm name is empty or unsupported, even with zero factor
- **THEN** the command reports the invalid selector and creates no image

### Requirement: Reproduce cells into adjacent vacancies

In replication mode, each initial and newborn occupied cell SHALL receive one
reproduction decision with probability equal to the clusterization factor,
locally modified when regional growth variation is enabled. Zero base
clusterization SHALL disable reproduction even with regional growth variation.
Reproduction SHALL add exactly two distinct children chosen from the parent's
currently empty eight-neighbor positions, uniformly when placement controls
are zero and with configured preferences otherwise. If fewer than two positions
are available, the cell SHALL add none. A death decision SHALL add no children
and SHALL NOT erase or recolor the parent. Newborn cells SHALL receive their own
decision. Cells SHALL never be processed more than once or overwrite occupied
positions, and processing SHALL terminate when no decisions remain.

#### Scenario: Produce exactly two children

- **WHEN** a cell reproduces and at least two adjacent empty positions exist
- **THEN** exactly two distinct adjacent positions become occupied
- **AND** the parent remains occupied and its children may subsequently reproduce

#### Scenario: Stop a branch

- **WHEN** a cell dies or has fewer than two empty adjacent positions
- **THEN** it adds no children and remains visible

#### Scenario: Preserve the initial population

- **WHEN** replication completes
- **THEN** every initial occupied cell remains unchanged
- **AND** the final population equals the initial population plus an even number of newborn cells

### Requirement: Configure replication growth preferences

The CLI SHALL accept `--replication-inheritance`, `--replication-forward-bias`,
`--replication-crowding`, and `--replication-branching` as finite numbers from 0
through 1 inclusive, defaulting to zero. Help SHALL describe each option.
These controls SHALL affect only replication. Their zero defaults SHALL preserve
existing seeded image content when regional controls are also disabled, and valid values SHALL NOT alter compact or strand
output. Invalid values SHALL fail before creating an output file, even when the
selected algorithm does not use them.

Inheritance SHALL set the probability that a newborn keeps its parent-to-child
direction as its heading instead of taking a random heading. Roots SHALL start
with random headings. Forward bias SHALL favor forward and gentle turns over
sideways and backward choices relative to the heading. Inheritance alone SHALL
NOT change image content when forward bias and branching are zero.

Crowding SHALL softly favor destinations with fewer occupied eight-neighbors,
excluding the parent and including the first child when selecting the second.
Branching SHALL set the probability that the second child favors a heading
90 degrees to either side of the first child's direction. All valid vacancies
SHALL retain positive selection weights; these controls SHALL NOT add a new
reproduction rejection or permit partial births.

#### Scenario: Reproduce an existing image

- **WHEN** the four controls are omitted or explicitly zero and regional controls are disabled
- **THEN** replication produces the existing seeded image
- **AND** using the same controls and seed replays the same decoded pixels

#### Scenario: Grow with inherited direction

- **WHEN** inheritance and forward bias are enabled
- **THEN** children can continue outward using their inherited heading
- **AND** larger forward bias strengthens the forward preference

#### Scenario: Avoid crowded destinations

- **WHEN** crowding avoidance is enabled
- **THEN** more crowded destinations have lower relative weights than equally directed less crowded destinations
- **AND** a crowded destination remains eligible if it is vacant

#### Scenario: Favor sideways forks

- **WHEN** a reproduction decision selects sideways branching
- **THEN** the second child's placement favors a quarter-turn from the first child's direction
- **AND** two distinct children are still placed whenever two vacancies exist

#### Scenario: Validate controls independently of selection

- **WHEN** any replication control is negative, above one, nonnumeric, NaN, or infinite
- **THEN** the command fails with a diagnostic and creates no output for any algorithm

### Requirement: Configure smooth replication regions

The CLI SHALL accept `--replication-field-scale` as an integer from 1 through
4096 logical cells, default 128, and `--replication-growth-variation` and
`--replication-seeding-variation` as finite values from 0 through 1, default 0.
All values SHALL be validated before output creation for every algorithm.
Only replication SHALL use regional controls. Both variation strengths zero
SHALL preserve all existing seeded image content regardless of field scale.

A seed-derived smooth field SHALL influence seeding and growth in the same
spatial regions. Scale SHALL determine the approximate distance over which
values vary, measured in logical cells. Growth variation SHALL bound the
additive deviation from the base reproduction probability before clamping to
[0,1]. Clamping MAY change the mean reproduction probability. Zero base
probability SHALL prevent all reproduction.

Seeding variation SHALL interpolate between uniform sampling and regionally
contrasted probabilities. The mean initial occupancy probability over the
logical grid SHALL remain density within floating-point precision. This SHALL
NOT require an exact occupied-cell count. Density 0 and 1 SHALL retain their
existing behavior. Constant fields SHALL fall back to uniform probabilities.
Regional seeding SHALL remain active when reproduction is disabled, and its
sample SHALL be retained even when all cells happen to be occupied or vacant.
Every logical cell, including a clipped edge cell, SHALL have equal weight in
the mean. The field SHALL be evaluated without a full-resolution field buffer.

#### Scenario: Produce regional variation

- **WHEN** the user enables seeding and growth variation
- **THEN** starting populations and reproduction probabilities vary gradually across the same regions
- **AND** repeated invocations with the same seed and settings reproduce the image

#### Scenario: Preserve average initial density

- **WHEN** regional seeding is enabled with reproduction disabled
- **THEN** the local initial probabilities average to density
- **AND** stronger seeding variation creates greater spatial contrast without promising an exact cell count

#### Scenario: Inspect the starting population

- **WHEN** clusterization is zero with regional seeding enabled
- **THEN** the output contains only the sampled regional initial population
- **AND** regional growth variation does not introduce births

#### Scenario: Handle narrow and degenerate fields

- **WHEN** the grid has one cell, one row, one column, clipped cells, or a field scale larger than its dimensions
- **THEN** generation remains finite and deterministic with probabilities within [0,1]
- **AND** density endpoints and paired births remain valid

#### Scenario: Reject invalid regional options

- **WHEN** field scale is noninteger or outside 1..4096, or a variation is nonfinite or outside 0..1
- **THEN** the command reports the invalid option and creates no image
