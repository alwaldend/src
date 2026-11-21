## ADDED Requirements

### Requirement: Develop the ISO ruleset in the root Bazel workspace

The repository SHALL build and test the ISO ruleset in the root Bazel workspace.

#### Scenario: CI discovers workspace boundaries

- **WHEN** repository CI discovers the declared Bazel workspaces
- **THEN** it includes the ISO ruleset in the root phases
- **AND** it runs no separate workspace phases for the ISO ruleset
