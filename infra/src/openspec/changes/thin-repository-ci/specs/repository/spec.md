## ADDED Requirements

### Requirement: Thin CI workflow wrappers

Repository CI workflows MUST be maximally thin wrappers. Executable logic SHALL
be implemented in `bazel run` targets or committed scripts; workflow YAML SHALL
contain only declarative orchestration and short command invocations. The
discoverable `repo-ci` skill SHALL document this rule and link owning sources.

#### Scenario: A workflow needs executable job logic

- **WHEN** a workflow requires resource checks, authentication, parsing, loops,
  or other executable behavior
- **THEN** that behavior resides in a committed implementation invoked by the
  workflow, preferably through `bazel run`
- **AND** inline programs, heredocs, and composite-action YAML do not substitute
  for the implementation

#### Scenario: An agent changes repository CI

- **WHEN** an agent discovers the repository skills for CI work
- **THEN** `repo-ci` explains the workflow and command boundaries, authentication
  constraints, and applicable offline and live validation entry points
