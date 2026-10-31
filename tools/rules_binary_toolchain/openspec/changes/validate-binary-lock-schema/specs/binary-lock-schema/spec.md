## ADDED Requirements

### Requirement: Strict binary lock metadata

The parser SHALL reject unknown fields in lock, toolchain, archive, binary and action structures. Archive toolchain metadata SHALL preserve supported execution-platform constraints and reject obsolete rules fields.

#### Scenario: Obsolete constraints field

- **WHEN** an archive declares execution constraints under rules
- **THEN** parsing fails with an unknown-field diagnostic instead of generating an unconstrained toolchain

#### Scenario: Supported execution constraints

- **WHEN** an archive declares exec_compatible_with under toolchain
- **THEN** generated toolchain declarations retain those constraints
