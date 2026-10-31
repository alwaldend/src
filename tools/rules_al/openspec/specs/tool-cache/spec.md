# Explicit tool installation

## Purpose

Record AL's explicit source-built installation boundary after removal of the
private executable cache.

## Requirements

### Requirement: Use explicit source-built installation

AL SHALL be built by Bazel and installed only through an explicitly requested
installation target. It SHALL NOT expose `al tool` or maintain a private executable
cache. Updating source SHALL NOT change an installed executable automatically.

#### Scenario: Update an installed command

- **WHEN** an operator requests the repository's AL install target
- **THEN** the installer atomically copies the current source-built executable.

#### Scenario: Invoke the removed tool subcommand

- **WHEN** a caller invokes `al tool`
- **THEN** AL reports usage failure without inspecting a tool cache.
