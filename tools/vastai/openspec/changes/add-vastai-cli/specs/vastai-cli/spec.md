## Purpose

Provide repository users with reproducible access to the upstream Vast.ai CLI through Bazel.

## ADDED Requirements

### Requirement: Pinned CLI entry point

The repository SHALL expose a pinned Vast.ai CLI at `//tools/vastai` with repository-internal visibility.

#### Scenario: Discover commands

- **WHEN** a user runs `bazel run //tools/vastai -- --help`
- **THEN** the upstream command help is displayed successfully without requiring API credentials or provisioning resources.
