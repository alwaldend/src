## ADDED Requirements

### Requirement: Instance networking

The implementation SHALL provide the following behavior: Add replacement-only sensitive docker_options with an empty compatible default and computed mapped_ports. Preserve imported and legacy instances. Exercise actual Terraform lifecycle fixtures.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets
