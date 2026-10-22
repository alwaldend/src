## ADDED Requirements

### Requirement: Restrict deployment identity to its authorized hosts

The checked-in deployment signing policies and host authorization configuration
SHALL jointly restrict each component identity to its declared host scope. A
component identity MUST NOT obtain a certificate accepted for another
component's administration unless that shared scope is explicitly declared.

#### Scenario: A component deploys to its own host

- **WHEN** a deployment certificate is issued to a component identity
- **THEN** the declared target host accepts its authorized principal
- **AND** an unrelated component's host rejects that certificate

#### Scenario: An identity requests another component's principal

- **WHEN** a component identity requests a principal outside its declared scope
- **THEN** the signing policy refuses the request

### Requirement: Bound deployment certificate authority

Deployment certificates SHALL have an explicit maximum lifetime and allowed
extensions justified by their deployment workflow. Migration SHALL preserve a
documented authorized administration path while removing superseded broad
deployment grants after the restricted path has been verified.

#### Scenario: Validate a deployment identity migration

- **WHEN** the proposed desired state is exercised in an isolated environment
- **THEN** its issued certificate has the documented lifetime and extensions
- **AND** permitted access succeeds and cross-component access fails
- **AND** removing the old grant leaves the authorized administration path usable
