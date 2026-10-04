## ADDED Requirements

### Requirement: Repository OpenSpec integration

The repo-openspec skill SHALL define owner selection, pinned
CLI invocation, and maintained-work conventions, and SHALL route operation
procedures to the discoverable upstream OpenSpec skills. It SHALL NOT independently
maintain generic create, apply, specification-writing, or archive procedures.
Repository policy SHALL remain owned by AGENTS.md.

#### Scenario: Apply work in a component workspace

- **WHEN** an agent implements a maintained component change
- **THEN** the integration skill selects the owning OpenSpec workspace and pinned
  Bazel invocation, and routes implementation procedure to openspec-apply-change
- **AND** the upstream operation skill owns artifact instructions and task steps

#### Scenario: Validate before archiving

- **WHEN** an implementation is ready for acceptance and archive
- **THEN** the integration skill routes to openspec-verify-change and
  openspec-archive-change while preserving repository validation and authority

## MODIFIED Requirements

### Requirement: Select maintained work before implementation

The repo-workspace and repo-openspec procedures SHALL apply the repository's
AGENTS.md requirement to select or create an owner-local OpenSpec change for
all repository modifications, including trivial changes. They SHALL preserve
existing authority and keep artifacts proportional to the requested work.

#### Scenario: Correct a single documentation typo

- **WHEN** an agent begins a one-line documentation correction
- **THEN** it verifies isolation and selects a matching owner-local OpenSpec change
  before editing, reusing an active matching record where possible
- **AND** it uses skip_specs for work without a behavioral delta rather than
  inventing requirements or requesting redundant approval
