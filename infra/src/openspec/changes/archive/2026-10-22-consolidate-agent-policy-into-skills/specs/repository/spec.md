## ADDED Requirements

### Requirement: Compose repository policy with skill routing

The root `AGENTS.md` SHALL keep repository-wide constraints visible: fact
ownership, infrastructure and approval authority, mandatory delivery,
worktree isolation, protection of unrelated work, confidentiality, and naming
and communication standards. Execution detail SHALL live with the skill that
owns it, and the root document SHALL route to that skill through one compact
"When to load a skill" table ordered by task phase. Each routing entry SHALL
name the trigger and the skill without restating the skill's procedure.

Initial reads SHALL depend on the task: applicable policy and the owning
README first, and BUILD or MODULE declarations only when implementation or
validation depends on them.

#### Scenario: An agent starts an unfamiliar task

- **WHEN** an agent begins work in a subtree
- **THEN** it reads the applicable `AGENTS.md` chain and the owning `README.md`
- **AND** it inspects BUILD or MODULE declarations only when the task depends
  on them
- **AND** it loads the skill whose routing trigger matches the current phase

#### Scenario: A procedure needs to change

- **WHEN** a repository-wide procedure changes
- **THEN** the change belongs with the skill that owns it, not in the root
  policy
- **AND** the root document retains only the routing entry and any
  repository-wide constraint

### Requirement: One authoritative source per fact

Each fact SHALL have one authoritative source. Independently maintained
duplicates SHALL be treated as defects even when the original is hard to
reach; the owner SHALL be extended or parameterized instead of copied.
Generated projections, fixtures, and concise summaries that identify their
source are projections rather than duplicates.

#### Scenario: A generated file restates an owned fact

- **WHEN** a checked-in generated file, fixture, or summary restates a fact
  owned elsewhere
- **THEN** it identifies its source and regeneration workflow
- **AND** it is not treated as an independently maintained duplicate

## MODIFIED Requirements

### Requirement: Isolated implementation and evidence

Repository modifications MUST use a dedicated feature branch and linked
worktree. The repository SHALL own that procedure in a `repo-workspace` skill
that runs before the first mutation or task-scratch write, so isolation and
scratch placement are established before `repo-delivery` is loaded. Task
scratch MUST remain under ignored `out/<task>/` in the applicable workspace,
and configurable temporary and cache locations SHALL point there. Delivery
SHALL exclude disposable build outputs and include required
generator-maintained files updated through their owning workflow. Acceptance
SHALL identify the candidate and relevant checks and representative output
before delivery.

#### Scenario: Finish a repository change

- **WHEN** an authorized implementation is ready for delivery
- **THEN** task-owned changes pass applicable formatting, quality and semantic
  checks and are committed and published through the repository procedure
- **AND** unrelated edits, scratch and protected data are excluded

#### Scenario: Begin a task that writes

- **WHEN** a task will modify repository files or write task-owned scratch
- **THEN** the agent verifies the feature worktree and branch before the first
  write
- **AND** task-owned scratch and configurable temporary locations stay under
  the ignored `out/<task>/` directory

#### Scenario: Prepare the delivery candidate

- **WHEN** the delivery procedure stages the candidate
- **THEN** disposable build outputs are excluded
- **AND** required generator-maintained files are included after their owning
  workflow regenerates them
