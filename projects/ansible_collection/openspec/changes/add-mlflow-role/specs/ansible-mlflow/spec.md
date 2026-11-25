## ADDED Requirements

### Requirement: MLflow service role

The collection SHALL provide a role that installs a caller-supplied MLflow
runtime and manages its tracking service and Vast.ai coordinator on systemd.

#### Scenario: Deploying a runtime revision

- **WHEN** a caller provides the runtime archive and coordinator configuration
- **THEN** the role installs the revision and enables both loopback services
  with persistent state

#### Scenario: Dedicated storage is required

- **WHEN** the caller enables required data mounts
- **THEN** tracking and the coordinator fail to start if any of their expected
  mount points is absent

#### Scenario: The runtime uses OIDC

- **WHEN** a caller supplies an OIDC provider registry and a workload token
- **THEN** the role starts tracking with the OIDC app and installs the
  coordinator's token as a private file
