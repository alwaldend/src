# Agents Specification

## Purpose

Describe the reusable repository skills and their packaging and discovery
contract as owned by `projects/agents`. Component behavior, repository policy,
and requested authority remain at their own owners.

Sources: [project README](../../../README.md),
[project packaging](../../../BUILD.bazel), and
[discovery declaration](../../../../../.agents/BUILD.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08; durable-work routing includes this OpenSpec migration.

## Requirements

### Requirement: Preserve natural fact ownership

This project SHALL keep component behavior at its owning README and
implementation, repository policy at applicable `AGENTS.md` files, and
requested outcome and action authority in the user interaction. Project
documentation SHALL NOT independently authorize actions.

#### Scenario: A document describes an unimplemented interface

- **WHEN** a document describes an interface without an implemented provider
  and supporting evidence
- **THEN** it states that status explicitly rather than implying a supported
  runtime capability

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
