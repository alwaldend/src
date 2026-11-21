## ADDED Requirements

### Requirement: Forgejo repository CI reporting

Repository CI running on Forgejo SHALL provide a collapsible log group for each
Bazel phase and a job summary with the result and elapsed duration of each phase
in each workspace.

#### Scenario: One phase fails

- **WHEN** a Bazel phase fails in Forgejo
- **THEN** its diagnostics remain available in the workspace and phase group
- **AND** the summary marks the phase as failed and lists the results of the
  remaining phases
- **AND** the action reports failure

#### Scenario: Run the command locally

- **WHEN** the command runs outside Forgejo
- **THEN** it streams the phase logs as plain text
