## Purpose

Give contributors a short, verifiable path from a supported host to their
first repository build while keeping host trust configuration explicit.

Baseline: the [shared truststore option at the reviewed revision](https://github.com/alwaldend/src/blob/e0d24cb3a9c250fc92665395d953f50aab2e5a6a/tools/bazelrc/project.bazelrc).
These requirements describe deferred behavior, not an already verified host matrix.

## ADDED Requirements

### Requirement: Document a supported first build

The repository entry point SHALL identify the supported host prerequisites
and link one authoritative setup procedure with a minimal public smoke build
and its expected output. Specialized component and agent procedures SHALL
remain linked from their owners instead of being duplicated in that path.

#### Scenario: A contributor starts from a fresh supported host

- **WHEN** a contributor follows the documented setup and smoke-build steps
- **THEN** the selected target produces the documented verifiable output
- **AND** the procedure identifies required tools and host assumptions before use

### Requirement: Select host trust configuration explicitly

Shared Bazel startup configuration SHALL work on documented supported hosts
without unconditionally requiring the Fedora-style Java truststore pathname.
Custom certificate authorities SHALL remain supported through explicit
host configuration, with certificate validation enabled.

#### Scenario: The current hardcoded truststore path is absent

- **WHEN** a supported clean host lacks that path and uses standard public trust
- **THEN** the documented smoke build resolves its public dependencies using
  the supported trust configuration
- **AND** TLS verification remains enabled

#### Scenario: A supported host uses the repository's private CA

- **WHEN** the documented custom-trust setup is selected on that host
- **THEN** the required private-CA endpoint validates successfully
- **AND** unrelated hosts do not inherit that host's filesystem assumptions
