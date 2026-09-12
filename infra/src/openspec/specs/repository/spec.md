# Repository Evolution Specification

## Purpose

Define how the alwaldend/src monorepo evolves its shared structure, build
system and development workflows. Component behavior is specified in each
owner's local OpenSpec workspace. This baseline was inspected at source revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f` on 2026-09-08 and incorporates the
OpenSpec migration delivered with this specification. Changes in `infra/src`
record repository evolution; they do not collect unrelated component work or claim
runtime health or deployed infrastructure.

Sources: [repository map](../../../../../README.md), [agent policy](../../../../../AGENTS.md),
[root build](../../../../../BUILD.bazel), [project boundary](../../../../../projects/README.md),
[infrastructure boundary](../../../../README.md),
[tool boundary](../../../../../tools/README.md), and [OpenSpec workflow](../../README.md).

## Requirements

### Requirement: Component ownership and source boundaries

The repository MUST keep component purpose and publication boundaries in the
nearest owner README, executable and dependency structure in BUILD and MODULE
files, and agent policy in the applicable AGENTS.md chain. Specifications SHALL
describe component requirements with links to these owners.

#### Scenario: Select an implementation owner

- **WHEN** a change affects a project, tool or infrastructure component
- **THEN** its owning documents and Bazel workspace determine the affected
  paths and supported consumers
- **AND** a repository-wide catalog or specification does not grant broader
  publication or execution authority

### Requirement: Reproducible Bazel development

Repository agent build and test operations MUST use `bazel_agent bazel` in the
owning workspace with pinned dependencies and the repository's shared agent
configuration. Generated dependency and catalog files MUST be updated through
their owning generator.

Every nested module that declares `MODULE.bazel` MUST declare only
dependencies its own sources use, with a resolvable version or an override
that applies when that module is the root. Every nested module directory MUST
appear in the root `.bazelignore` so root target expansion does not cross the
workspace boundary. Each nested workspace MUST build and test standalone.

#### Scenario: Change a shared dependency

- **WHEN** an implementation adds an external build input
- **THEN** its owner records an immutable version and integrity information
- **AND** applicable generated locks and package checks validate the declared
  dependency through the pinned build workflow

#### Scenario: Build a nested module standalone

- **WHEN** a nested Bazel workspace is built or tested on its own
- **THEN** its module graph resolves without depending on root-only overrides
- **AND** the workspace builds and tests through its shared configuration

#### Scenario: Expand root targets across a nested boundary

- **WHEN** the root workspace expands targets beneath a nested module directory
- **THEN** the nested workspace is excluded by the root ignore list
- **AND** root expansion neither loads nor silently omits that module

### Requirement: Isolated implementation and evidence

Repository modifications MUST use a dedicated feature branch and linked
worktree. Task scratch MUST remain under ignored `out/<task>/` in the
applicable workspace. Acceptance SHALL identify the candidate and relevant
checks and representative output before delivery.

#### Scenario: Finish a repository change

- **WHEN** an authorized implementation is ready for delivery
- **THEN** task-owned changes pass applicable formatting, quality and semantic
  checks and are committed and published through the repository procedure
- **AND** unrelated edits, scratch and protected data are excluded

### Requirement: Declarative infrastructure authority

Persistent infrastructure configuration MUST be defined in checked-in owning
infrastructure code. Infrastructure operations MUST require the user's exact
operation and scope authorization; an implementation validation request SHALL
NOT authorize live deployment.

#### Scenario: Validate an infrastructure specification

- **WHEN** an agent adds or updates an infra component specification
- **THEN** it inspects source and runs applicable offline checks
- **AND** it does not infer deployed state or execute live apply or deployment
  from the existence of the specification

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

### Requirement: Legacy goal migration and deprecation

The migration SHALL preserve all tracked maintained goal history and its
recorded acceptance and execution states with source-to-change mappings and
checksums. The goal skill MUST be disabled for discovery and the goal tool
MUST identify itself as deprecated compatibility functionality.

#### Scenario: Inspect a migrated historical record

- **WHEN** a reader follows a legacy goal's migration entry
- **THEN** its OpenSpec change exposes the original history and status
- **AND** historical completion does not establish acceptance for a different
  candidate or automatically apply old specification deltas
