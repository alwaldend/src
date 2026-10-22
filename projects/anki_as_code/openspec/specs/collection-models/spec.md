# collection-models Specification

## Purpose

Define first-party collection configuration, stored state and execution data through generated protobuf contracts.

## Requirements

### Requirement: Generated collection contracts

The API SHALL define editable collection and plan data in a first-party protobuf schema, using existing repository generation rules. It SHALL NOT copy upstream Anki schemas or check in generated language sources.

#### Scenario: Independent generation

- **WHEN** the collection API Bazel targets are built
- **THEN** schema generation and Go compilation succeed without the collection implementation
