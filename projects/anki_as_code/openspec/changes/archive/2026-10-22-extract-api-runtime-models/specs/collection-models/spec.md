## ADDED Requirements

### Requirement: Generated collection contracts

The API SHALL define editable collection and plan data in a first-party protobuf schema, using existing repository generation rules. It SHALL NOT copy upstream Anki schemas or check in generated language sources.

#### Scenario: Independent generation

- **WHEN** the extracted API Bazel targets are built against master
- **THEN** schema generation and Go compilation succeed without the collection implementation

### Requirement: Runtime models use generated messages

The model package SHALL contain only shared runtime data types. It SHALL use generated protobuf messages for serialized data and SHALL NOT duplicate their fields or contain business logic.

#### Scenario: Independent model compilation

- **WHEN** the model Bazel target is built against master
- **THEN** it compiles using only the collection API and existing upstream Anki note-type definitions
