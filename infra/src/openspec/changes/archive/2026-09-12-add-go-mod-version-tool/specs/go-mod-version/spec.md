## ADDED Requirements

### Requirement: Configured go.mod version

The repository SHALL provide a `tools/go_mod` command that discovers every
tracked `go.mod` file and sets its `go` directive to the version configured in
the tool's BUILD file. The command MUST offer an update mode and a check mode.

#### Scenario: Update module files

- **WHEN** a user runs the update target
- **THEN** every tracked `go.mod` file whose `go` directive differs from the
  configured version is rewritten to that version
- **AND** files already on the configured version remain unchanged

#### Scenario: Check module files

- **WHEN** the check test runs
- **THEN** it fails if any tracked `go.mod` file has a `go` directive
  different from the configured version
- **AND** it passes when every tracked `go.mod` file matches

### Requirement: Repository quality integration

The go.mod version check MUST run as part of `//:repo_quality_test` so a
version mismatch fails repository quality.

#### Scenario: Quality suite covers the check

- **WHEN** a `go.mod` file's `go` directive differs from the configured
  version
- **THEN** `//:repo_quality_test` fails
