## MODIFIED Requirements

### Requirement: OpenSpec coverage and work state

Each direct project in `projects/` and `infra/` SHALL keep its baseline
specifications and maintained changes in its own `openspec/` workspace.
`infra/src/openspec` SHALL describe evolution of the repository itself,
including shared structure, build system and development workflows. The
repository root SHALL NOT collect component specifications or change records.
Every repository change MUST use an OpenSpec change in its owning workspace,
including trivial edits, documentation, configuration, and dependency changes.
A matching active change SHALL be reused when available, and artifacts SHALL
remain proportional to scope. Read-only questions and reviews do not require
a change record. Maintained work SHALL preserve outcome, decisions, tasks and
requirement deltas. Existing task authority MUST survive workflow transitions.

#### Scenario: Change a component contract

- **WHEN** a change affects a project's supported behavior
- **THEN** its specifications and change artifacts remain in that project's
  `openspec/` directory
- **AND** the pinned CLI and context routing select that owner workspace

#### Scenario: Evolve the repository itself

- **WHEN** a change affects the monorepo's shared structure or development workflow
- **THEN** `infra/src/openspec` records its repository requirements and evolution
- **AND** component contracts remain with their respective owners

#### Scenario: Resume unfinished maintained work

- **WHEN** an agent resumes a named OpenSpec change
- **THEN** it reads that owner's change artifacts, candidate, evidence and next action
- **AND** unfinished or blocked acceptance remains explicit until evidence
  supports completion

#### Scenario: Validate a standalone project

- **WHEN** a project has its own Bazel module
- **THEN** repository OpenSpec validation still checks its declared local workspace
- **AND** a successful repository check does not omit that project's requirements

#### Scenario: Make a trivial change

- **WHEN** an agent makes a one-line repository correction
- **THEN** it selects or creates a matching owner-local OpenSpec change before
  editing, with small artifacts and skip_specs if behavior is unchanged
- **AND** the record adds no approval gate to already authorized work

#### Scenario: Conduct a read-only review

- **WHEN** the task only inspects a bounded change set and reports findings
- **THEN** no new OpenSpec change is required for that review
- **AND** the review itself does not authorize implementation or publication
