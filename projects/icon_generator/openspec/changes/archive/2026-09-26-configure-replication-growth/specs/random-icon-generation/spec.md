## MODIFIED Requirements

### Requirement: Reproduce cells into adjacent vacancies

In replication mode, each initial and newborn occupied cell SHALL receive one
reproduction decision with probability equal to the clusterization factor.
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

## ADDED Requirements

### Requirement: Configure replication growth preferences

The CLI SHALL accept `--replication-inheritance`, `--replication-forward-bias`,
`--replication-crowding`, and `--replication-branching` as finite numbers from 0
through 1 inclusive, defaulting to zero. Help SHALL describe each option.
These controls SHALL affect only replication. Their zero defaults SHALL preserve
existing seeded image content, and valid values SHALL NOT alter compact or strand
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

- **WHEN** the four controls are omitted or explicitly zero
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
