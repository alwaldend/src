# course-placeholder Specification

## Purpose

Describe the MLOps course project placeholder and its documentation contract
while distinguishing course metadata from future application implementation.

## Requirements

### Requirement: Publish truthful course placeholder documentation

The project SHALL provide a README identifying it as an MLOps course, linking
the user-supplied course page, and stating that no implementation code exists.
Its documentation SHALL be included through the repository documentation targets.
The placeholder SHALL NOT require a release deployment target.

#### Scenario: Inspect the course project

- **WHEN** a reader opens the project README
- **THEN** its description and external link identify the course
- **AND** the implementation is explicitly described as a placeholder

#### Scenario: Build project documentation

- **WHEN** repository consumers select the project documentation
- **THEN** its README and baseline specification are included
- **AND** release aggregation does not require an absent course release target
