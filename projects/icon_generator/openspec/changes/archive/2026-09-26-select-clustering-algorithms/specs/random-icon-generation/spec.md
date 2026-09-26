## MODIFIED Requirements

### Requirement: Configure occupancy density

The CLI SHALL accept a finite `--density` from 0 through 1 inclusive.
The initial placement SHALL independently occupy each cell with that
probability. Compact and strand clustering SHALL rearrange that sample without
changing its occupied-cell count. Replication SHALL use the sample as its
initial population and MAY increase the count through reproduction. Density
SHALL represent an initial occupancy probability, not an exact occupied-cell
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
`replication`, defaulting to `compact`. Zero factor SHALL preserve independent
placement and existing seeded image content for every algorithm. Explicit
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

- **WHEN** the user omits clusterization or explicitly sets it to 0
- **THEN** the decoded image matches the independent-placement behavior
- **AND** omitted algorithm selection remains equivalent to explicit compact for positive factors

#### Scenario: Reject an invalid factor

- **WHEN** clusterization is negative, greater than 1, nonnumeric, NaN, or infinite
- **THEN** the command exits unsuccessfully with a diagnostic and creates no image

#### Scenario: Reject an unknown algorithm

- **WHEN** the algorithm name is empty or unsupported, even with zero factor
- **THEN** the command reports the invalid selector and creates no image

## ADDED Requirements

### Requirement: Reproduce cells into adjacent vacancies

In replication mode, each initial and newborn occupied cell SHALL receive one
reproduction decision with probability equal to the clusterization factor.
Reproduction SHALL add exactly two distinct children chosen uniformly from the
parent's currently empty eight-neighbor positions. If fewer than two positions
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
