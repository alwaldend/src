## MODIFIED Requirements

### Requirement: OpenSpec coverage and work state

Each direct project in `projects/` and `infra/` SHALL keep its baseline
specifications and maintained changes in its own `openspec/` workspace.
`infra/src/openspec` SHALL describe evolution of the repository itself,
including shared structure, build system and development workflows. The
repository root SHALL NOT collect component specifications or change records.
Maintained work SHALL preserve outcome, decisions, tasks and requirement
deltas. Existing task authority MUST survive workflow transitions.

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
