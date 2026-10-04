## ADDED Requirements

### Requirement: Typed Python development guidance

The discoverable repo-python skill SHALL require type hints throughout Python
code and identify the owning dependency manifest, shared tool configuration,
generated files, and pinned validation workflows.

#### Scenario: Add Python implementation

- **WHEN** an agent selects repo-python for a Python change
- **THEN** it is instructed to annotate parameters, return types, and variables,
  preserve component support constraints, and use owning Bazel workflows
- **AND** project support metadata and shared tool settings are distinguished

### Requirement: Evidence-based change review

The discoverable change-review skill SHALL review a bounded change set for
bugs and violations of applicable repository standards. It SHALL identify the
candidate and comparison baseline, evaluate affected behavior and consumers,
and report actionable findings with evidence and consequence. It SHALL
preserve the distinction between reviewing and authorizing modifications or
external review messages.

#### Scenario: Review a pull request

- **WHEN** the requested change set is a pull request
- **THEN** the review identifies its exact head and base, inspects the applicable
  policy and affected owners, and reports evidenced findings with file locations
- **AND** it does not post comments, fix source, or merge without authorization

#### Scenario: Review working-tree changes

- **WHEN** the requested change set is staged or uncommitted source
- **THEN** the review identifies the baseline and the included index, worktree,
  and untracked content, and bounds claims to that observed candidate
- **AND** it investigates relevant consumers beyond changed lines as needed

#### Scenario: Insufficient evidence or no actionable finding

- **WHEN** inspection cannot prove a suspected defect or finds no actionable issue
- **THEN** the response distinguishes uncertainty and validation limits from
  confirmed bugs without inventing findings or treating passing checks as proof

### Requirement: Select maintained work before implementation

The repo-workspace and openspec procedures SHALL apply the repository's
AGENTS.md requirement to select or create an owner-local OpenSpec change for
all repository modifications, including trivial changes. They SHALL preserve
existing authority and keep artifacts proportional to the requested work.

#### Scenario: Correct a single documentation typo

- **WHEN** an agent begins a one-line documentation correction
- **THEN** it verifies isolation and selects a matching owner-local OpenSpec change
  before editing, reusing an active matching record where possible
- **AND** it uses skip_specs for work without a behavioral delta rather than
  inventing requirements or requesting redundant approval
