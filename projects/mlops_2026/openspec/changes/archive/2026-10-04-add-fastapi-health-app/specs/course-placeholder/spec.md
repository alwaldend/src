## MODIFIED Requirements

### Requirement: Publish truthful course placeholder documentation

The project SHALL provide a README identifying it as an MLOps course and linking
the user-supplied course page. After implementation it SHALL describe the runnable
health service and its supported endpoints and startup command, rather than
claiming that no code exists. Its documentation SHALL remain included through
repository documentation targets. The local application SHALL NOT require a
release deployment target.

#### Scenario: Inspect the course project

- **WHEN** a reader opens the project README after implementation
- **THEN** its description and external link identify the course
- **AND** it documents the health service, OpenAPI endpoints, and local startup

#### Scenario: Build project documentation

- **WHEN** repository consumers select the project documentation
- **THEN** its README and baseline specification are included
- **AND** release aggregation does not require an absent course release target
