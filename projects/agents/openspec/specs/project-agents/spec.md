# Agents Specification

## Purpose

Describe the repository agent-system contract, shared procedures, and evidence
boundaries owned by `projects/agents`. The architecture defines intended
composition; the current-state guide identifies implemented capabilities.

Sources: [project README](../../../README.md),
[current state](../../../docs/current-state.md),
[architecture](../../../docs/architecture.md), and
[documentation packaging](../../../BUILD.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08; durable-work routing includes this OpenSpec migration.

## Requirements

### Requirement: Preserve natural fact ownership

The agent-system documentation SHALL keep component behavior at its owning
README and implementation, repository policy at applicable `AGENTS.md` files,
and requested outcome and action authority in the user interaction. Architecture
and derived catalogs SHALL NOT independently authorize actions.

#### Scenario: A design document names a future interface

- **WHEN** an interface is described by the architecture or roadmap without an
  implemented provider and supporting evidence
- **THEN** the current-state guide treats it as intended composition rather than
  a supported runtime capability.

### Requirement: Package canonical skills separately from discovery

Enabled repository skills SHALL have canonical project-owned directories and
`skill_library` targets. `.agents/skills/` SHALL expose relative symlinks to
those directories, while Bazel builds the canonical targets. The deprecated
goal skill SHALL be excluded from enabled discovery.

#### Scenario: An agent discovers a reusable procedure

- **WHEN** an enabled skill is selected through `.agents/skills/`
- **THEN** its content comes from its canonical owning project without creating
  a second Bazel package for the discovery link.

### Requirement: Use OpenSpec for maintained work

Maintained agent-system work SHALL use this project's OpenSpec changes for its
proposal, acceptance requirements, tasks, and continuation evidence. Migrated
goal history SHALL remain available as historical evidence with its original
acceptance and observation limits.

#### Scenario: Work resumes after migration

- **WHEN** an agent continues maintained repository-agent work
- **THEN** it follows the corresponding OpenSpec change and retained migration
  history instead of creating a new legacy goal record.

### Requirement: Bound claims made from evaluations and catalogs

Documentation SHALL distinguish offline eval-configuration validation from
behavioral correctness, configured capability declarations from runtime
availability, and stored observations from current provider health.

#### Scenario: An offline skill evaluation configuration passes

- **WHEN** the configured Promptfoo validation target succeeds
- **THEN** the result establishes validity of the evaluation harness without
  claiming successful skill routing or real task completion.
