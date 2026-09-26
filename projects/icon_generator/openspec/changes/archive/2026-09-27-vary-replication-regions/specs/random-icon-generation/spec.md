## MODIFIED Requirements

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

## ADDED Requirements

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
