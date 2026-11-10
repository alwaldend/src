## ADDED Requirements

### Requirement: Approved Jsonnet tooling

Bazel SHALL supply Google go-jsonnet v0.22.0 and its formatter through the existing
binary toolchain mechanism with a pinned upstream archive and verified integrity.
It SHALL NOT add a runtime dependency to production applications or rentals.

#### Scenario: Build execution

- **WHEN** a declared Jsonnet entrypoint and imports are compiled
- **THEN** a sandboxed action SHALL produce only declared consumer-owned JSON outputs without network access

#### Scenario: Invalid source

- **WHEN** the entrypoint contains invalid Jsonnet or references an undefined symbol
- **THEN** compilation SHALL fail and no successful output SHALL be reported
