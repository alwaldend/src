## ADDED Requirements

### Requirement: Managed MLflow tracking

The infrastructure SHALL declare VM provisioning, credentials, and deployment
for an authenticated HTTPS MLflow service with persistent tracking and artifact
storage.

#### Scenario: Tracking data survives a service restart

- **WHEN** an authenticated client logs a run and uploads an artifact
- **THEN** the run and artifact remain accessible after the service restarts

### Requirement: Vast.ai Projects execution

The deployed runtime SHALL include the MLflow Vast.ai project backend and a
loopback coordinator for configured existing rentals.

#### Scenario: A project targets a configured worker

- **WHEN** the plugin submits a project to the coordinator for an allowed worker
- **THEN** that worker can retrieve the project artifact and report the outcome
  against the same MLflow run

### Requirement: Dedicated service storage

The deployment SHALL attach separate persistent disks for MLflow tracking,
MLflow artifacts, the Vast.ai coordinator, and Traefik. Each service SHALL
require its configured data mounts before starting.

#### Scenario: A service data disk is missing

- **WHEN** a required data mount is absent
- **THEN** the service does not start with an empty directory on the boot disk

### Requirement: Tracking authentication

The tracking service SHALL use Vault OIDC for browser sessions and verify issuer-
and audience-bound bearer tokens for registered workload accounts.

#### Scenario: A workload rotates its tracking token

- **WHEN** the coordinator's token file is replaced with a valid bearer token
- **THEN** subsequent tracking requests use the replacement without a restart
