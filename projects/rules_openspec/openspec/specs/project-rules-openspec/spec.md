# project-rules-openspec Specification

## Purpose

Define reusable Bazel validation and execution behavior for owner-local OpenSpec workspaces.

## Requirements

### OpenSpec validation

The system SHALL validate one owner-local OpenSpec workspace in a sandbox using the pinned OpenSpec CLI.

#### Scenario: current validation

- **WHEN** a caller runs the validation target
- **THEN** OpenSpec validates current specifications and active changes with strict mode

#### Scenario: archived validation

- **WHEN** a caller requests archived validation
- **THEN** OpenSpec validates completed archived changes with strict mode
