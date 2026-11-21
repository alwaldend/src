## ADDED Requirements

### Requirement: Bound Bazel server lifetime in Forgejo CI

The shared CI command SHALL shut down each workspace's Bazel server after
attempting its build and test phases on Forgejo, before starting the next
workspace.

#### Scenario: A phase fails on Forgejo

- **WHEN** a workspace's build or test fails
- **THEN** the CI command attempts to shut down the workspace's Bazel server
- **AND** the command still attempts the remaining workspaces

#### Scenario: Run shared CI locally

- **WHEN** the CI command runs outside Forgejo
- **THEN** it leaves Bazel server lifecycle to the local user
