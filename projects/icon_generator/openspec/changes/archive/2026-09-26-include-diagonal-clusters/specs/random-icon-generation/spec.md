## MODIFIED Requirements

### Requirement: Configure pixel clusterization

The CLI SHALL accept a finite `--clusterization` factor from 0 through 1,
defaulting to 0. Zero SHALL preserve independent placement and existing seeded
image content. Positive factors SHALL favor adjacency between occupied cells
in all eight surrounding positions: horizontal, vertical, and diagonal.
Adjacency SHALL NOT wrap across opposite canvas edges. Diagonal circle cells
SHALL count as neighbors even though the rendered circles do not touch at corners.
For the same seed and other settings, increasing the factor SHALL preserve
the occupied-cell count and SHALL NOT reduce the total number of occupied
neighbor pairs across all eight directions. Factor 1 SHALL provide the strongest
supported clustering effort, without guaranteeing a single connected group
or a change on every grid. The CLI SHALL reject invalid factors before creating
an output file and SHALL describe the option in help.

#### Scenario: Increase contact without increasing density

- **WHEN** the user increases clusterization with the seed and other options fixed
- **THEN** the number of occupied cells remains unchanged
- **AND** the total number of horizontal, vertical, and diagonal occupied neighbor pairs does not decrease

#### Scenario: Include corner neighbors

- **WHEN** two occupied cells meet only diagonally on the grid
- **THEN** they count as one neighboring pair for clustering
- **AND** square and circle shapes use the same grid adjacency definition

#### Scenario: Preserve default output

- **WHEN** the user omits clusterization or explicitly sets it to 0
- **THEN** the decoded image matches the independent-placement behavior

#### Scenario: Reject an invalid factor

- **WHEN** clusterization is negative, greater than 1, nonnumeric, NaN, or infinite
- **THEN** the command exits unsuccessfully with a diagnostic and creates no image
