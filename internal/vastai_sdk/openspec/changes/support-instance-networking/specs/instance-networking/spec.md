## ADDED Requirements

### Requirement: Instance networking

The implementation SHALL provide the following behavior: Pass optional Docker options through the generated create protocol. Read typed port bindings without retaining API-returned credentials; preserve omitted-field compatibility.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets
