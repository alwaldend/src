## MODIFIED Requirements

### Requirement: Configure occupancy density

The CLI SHALL accept a finite `--density` from 0 through 1 inclusive.
Initial placement SHALL independently sample each cell. The per-cell
probability SHALL equal density unless regional replication seeding is enabled,
in which case density SHALL equal the mean of the local probabilities. Compact and strand clustering SHALL rearrange that sample without
changing its occupied-cell count. Replication and tips SHALL use the sample as their
initial population and MAY increase the count through growth. Density
SHALL represent an average initial occupancy probability, not an exact occupied-cell
count or the percentage of output pixels covered by circles.

#### Scenario: Generate only background

- **WHEN** density is 0, with any valid clusterization and algorithm
- **THEN** every output pixel retains the configured background
- **AND** growth does not introduce a population spontaneously

#### Scenario: Occupy every cell

- **WHEN** density is 1, with any valid clusterization and algorithm
- **THEN** every cell contains the selected shape
- **AND** circle boundaries still leave background visible outside circles

### Requirement: Configure pixel clusterization

The CLI SHALL accept a finite `--clusterization` factor from 0 through 1,
defaulting to 0, and `--clusterization-algorithm` as `compact`, `strands`,
`replication`, or `tips`, defaulting to `compact`. Zero factor SHALL prevent reproduction. With regional seeding disabled,
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
Tips SHALL use the factor to scale growth length as specified by its lifecycle.

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

## ADDED Requirements

### Requirement: Grow persistent tips with asymmetric branches

The `tips` algorithm SHALL retain initial occupied cells and their colors and
extend eight-connected paths into vacant cells. A tip SHALL ordinarily add one
cell along a persistent, gradually turning heading. Occasional forks SHALL
create a side tip with a smaller remaining growth budget than the continuing
main tip. Tips SHALL stop on exhausted budgets, canvas boundaries, occupied
destinations, or contact with existing growth outside their recent trail and
local fork junction. Growth SHALL be finite without wraparound or overwrites.

Positive clusterization SHALL scale the available growth length; zero SHALL
preserve independent placement. Final coverage SHALL NOT be constrained to the
initial density or guaranteed monotonic in clusterization. Sparse starting
populations SHALL favor narrow winding paths with open gaps. Replication-specific
options SHALL NOT affect tips. The CLI SHALL describe tips and its density and
factor semantics. Existing algorithms SHALL retain their seeded behavior.

#### Scenario: Extend and branch from sparse roots

- **WHEN** tips is selected with sparse density and positive clusterization
- **THEN** paths grow from the sampled initial cells and can add a single child
- **AND** occasional forks allocate less growth to side branches than main tips
- **AND** every resulting connected component contains an initial cell

#### Scenario: Retain existing cells and replay growth

- **WHEN** a tips image is generated and then replayed with the same settings and seed
- **THEN** both images have identical decoded pixels
- **AND** every initial cell retains its palette color

#### Scenario: Stop in constrained space

- **WHEN** a tip meets existing growth outside its local trail or junction, or a canvas edge
- **THEN** it stops without overwriting cells or wrapping across the canvas
- **AND** single-cell, narrow, clipped, empty, and full grids terminate correctly

#### Scenario: Keep controls scoped to replication

- **WHEN** valid replication controls change while tips settings and seed remain fixed
- **THEN** the tips image remains unchanged
