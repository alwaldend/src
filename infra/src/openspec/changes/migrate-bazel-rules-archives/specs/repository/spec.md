## ADDED Requirements

### Requirement: Root-owned ruleset development

The repository SHALL build and test its first-party rulesets from the root
Bazel workspace and provide standalone archives with their public rule labels.

#### Scenario: A consumer loads a published archive

- **WHEN** a standalone workspace resolves one of the archives through a registry
- **THEN** it can load the module's existing public rules

#### Scenario: CI runs the repository checks

- **WHEN** repository CI builds and tests the source checkout
- **THEN** the root wildcard selects every migrated ruleset under its
  repository-relative package path
