## Purpose

Reduce resource contention from independent Bazel builds on host-bot through conservative per-invocation scheduling defaults.

## ADDED Requirements

### Requirement: Reduced per-invocation scheduling defaults

The managed user Bazel configuration SHALL set four concurrent jobs and a local CPU scheduling budget of `HOST_CPUS*0.25` for builds and tests. These defaults SHALL NOT be described as an aggregate host-wide resource cap.

#### Scenario: Build or test reads the managed defaults

- **WHEN** a new Bazel build or test invocation loads the managed user configuration without overriding its resource settings
- **THEN** it receives `--jobs=4` and `--local_resources=cpu=HOST_CPUS*0.25`

#### Scenario: Independent workspaces build concurrently

- **WHEN** separate Bazel servers build in multiple workspaces
- **THEN** each receives its own scheduling budget, without sharing an aggregate limit
