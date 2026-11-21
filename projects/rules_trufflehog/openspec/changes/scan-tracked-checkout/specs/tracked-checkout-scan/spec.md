## Purpose

Detect credentials in every tracked file in a materialized Git checkout,
including the bytes stored by Git LFS.

## ADDED Requirements

### Requirement: Scan the tracked checkout

The checkout scan SHALL submit tracked files to the filesystem scanner and
exclude Git metadata and untracked paths.

#### Scenario: An LFS payload is materialized

- **WHEN** the checkout contains a tracked LFS payload
- **THEN** TruffleHog scans its file contents

#### Scenario: An LFS file is still a pointer

- **WHEN** a tracked LFS file has not been materialized
- **THEN** the checkout scan fails
