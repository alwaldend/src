## MODIFIED Requirements

### Requirement: Configure pixel clusterization

The CLI SHALL accept a finite `--clusterization` factor from 0 through 1,
defaulting to 0. Zero SHALL preserve independent placement and existing seeded
image content. Positive factors SHALL favor winding strands with occasional
branches and open gaps at sparse densities. Higher factors SHALL favor longer
strands while preserving the occupied-cell count for the same seed and other
settings. Total neighbor count SHALL NOT be a monotonicity guarantee.

Clustering SHALL use all eight surrounding positions: horizontal, vertical,
and diagonal. Adjacency SHALL NOT wrap across opposite canvas edges. Diagonal
circle cells SHALL count as neighbors even though their shapes do not touch
at corners. Dense and tiny grids SHALL still preserve their sampled cell count,
even when the grid cannot accommodate open strands. The CLI SHALL reject invalid
factors before creating an output file and SHALL describe the option in help.

#### Scenario: Increase contact without increasing density

- **WHEN** the user increases clusterization with the seed and other options fixed
- **THEN** the number of occupied cells remains unchanged
- **AND** placement favors longer strands without promising a monotonic total neighbor count

#### Scenario: Include corner neighbors

- **WHEN** two occupied cells meet only diagonally on the grid
- **THEN** they count as one neighboring pair for clustering
- **AND** square and circle shapes use the same grid adjacency definition

#### Scenario: Grow sparse organic patterns

- **WHEN** the user selects sparse density and strong clusterization
- **THEN** the pattern favors narrow winding strands, occasional branches, and open gaps
- **AND** identical settings and seed reproduce the pattern

#### Scenario: Preserve default output

- **WHEN** the user omits clusterization or explicitly sets it to 0
- **THEN** the decoded image matches the independent-placement behavior

#### Scenario: Reject an invalid factor

- **WHEN** clusterization is negative, greater than 1, nonnumeric, NaN, or infinite
- **THEN** the command exits unsuccessfully with a diagnostic and creates no image
