## ADDED Requirements

### Requirement: Configure pixel clusterization

The CLI SHALL accept a finite `--clusterization` factor from 0 through 1,
defaulting to 0. Zero SHALL preserve independent placement and existing seeded
image content. Positive factors SHALL favor horizontal and vertical contacts
between occupied cells, with no adjacency across opposite canvas edges.
For the same seed and other settings, increasing the factor SHALL preserve
the occupied-cell count and SHALL NOT reduce the number of occupied neighbor
pairs. Factor 1 SHALL provide the strongest supported clustering effort,
without guaranteeing a single connected group or a change on every grid.
The CLI SHALL reject invalid factors before creating an output file and SHALL
describe the option in help.

#### Scenario: Increase contact without increasing density

- **WHEN** the user increases clusterization with the seed and other options fixed
- **THEN** the number of occupied cells remains unchanged
- **AND** the number of horizontal and vertical occupied neighbor pairs does not decrease

#### Scenario: Preserve default output

- **WHEN** the user omits clusterization or explicitly sets it to 0
- **THEN** the decoded image matches the independent-placement behavior

#### Scenario: Reject an invalid factor

- **WHEN** clusterization is negative, greater than 1, nonnumeric, NaN, or infinite
- **THEN** the command exits unsuccessfully with a diagnostic and creates no image

## MODIFIED Requirements

### Requirement: Configure occupancy density

The CLI SHALL accept a finite `--density` from 0 through 1 inclusive.
The initial placement SHALL independently occupy each cell with that
probability. Clusterization SHALL rearrange that sample without changing its
occupied-cell count. Density SHALL represent an initial occupancy probability,
not an exact occupied-cell count or the percentage of output pixels covered
by circles.

#### Scenario: Generate only background

- **WHEN** density is 0, with any valid clusterization
- **THEN** every output pixel retains the configured background

#### Scenario: Occupy every cell

- **WHEN** density is 1, with any valid clusterization
- **THEN** every cell contains the selected shape
- **AND** circle boundaries still leave background visible outside circles
