## ADDED Requirements

### Requirement: Derived repository inventories

Repository inventory projections SHALL derive specification-validation membership from owner source exports and project descriptions from owner READMEs. CI and the complete repository checker SHALL use the same discovery of module boundaries from .bazelignore.

#### Scenario: Add an owner source export

- **WHEN** an owner declares an OpenSpec source export
- **THEN** inventory regeneration includes it in repository validation
- **AND** freshness validation rejects the previous inventory

#### Scenario: Add a nested workspace

- **WHEN** .bazelignore declares a directory containing MODULE.bazel
- **THEN** both CI and the full repository checker discover it
- **AND** ignored directories without a module remain outside the build matrix

### Requirement: Complementary artifact boundary checks

Repository validation SHALL provide explicit artifact-dependency checks that reject forbidden indirect dependencies independently of package exceptions and exclude execution-tool edges. Its documentation SHALL state the audited roots and graph limitations.

#### Scenario: A runtime target consumes a forbidden artifact indirectly

- **WHEN** an audited target reaches a tools, infrastructure or user artifact through ordinary dependencies
- **THEN** its dependency-boundary check fails

#### Scenario: A build uses an execution tool

- **WHEN** an audited target uses a declared execution tool
- **THEN** that edge is excluded from artifact-dependency enforcement
