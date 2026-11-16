## Purpose

Allow deployment of a previously registered Forgejo worker without creating
another identity or exposing its credentials in deployment output.

## ADDED Requirements

### Requirement: Initial credential import

The role SHALL accept API-issued instance, name, id, uuid and token fields for
initial configuration. It MUST reject incomplete or mismatched inputs, protect
the worker file with mode 0600 and suppress secret-bearing output.

#### Scenario: Import an existing worker

- **WHEN** complete matching credentials are supplied and no worker file exists
- **THEN** the role writes the identity without contacting a registration API

#### Scenario: Reject invalid credentials

- **WHEN** supplied credentials are incomplete or mismatch the configured worker
- **THEN** deployment fails without writing a worker file or revealing secrets

### Requirement: Preserve existing registration

The role SHALL preserve existing worker files and retain its registration-token
flow when no API credentials are supplied.

#### Scenario: Repeat deployment

- **WHEN** a worker file already exists and bootstrap inputs are absent
- **THEN** deployment preserves that identity without requiring those inputs

### Requirement: Prerequisites precede service activation

The role SHALL expose installation and service phases. Installing a fresh
worker MUST defer service activation until the caller completes prerequisites,
including when a failed play forces handlers to run.

#### Scenario: Prerequisite setup fails

- **WHEN** the caller's prerequisite setup fails between installation and service
- **THEN** pending restart handlers leave the fresh worker stopped
