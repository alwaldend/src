## Purpose

Define reusable, data-preserving management of a native PostgreSQL service
and explicitly declared application access through the Ansible collection.

## ADDED Requirements

### Requirement: Select the service compatibly

The role SHALL install its configured native packages and enable and start
the configured PostgreSQL service. It SHALL support the documented
`postgresql_service` input and provide a documented compatibility path for
inventories using the previously referenced `postgresql` variable. Conflicting
explicit service selections SHALL fail before changing a service.

#### Scenario: Use the documented service input

- **WHEN** an initialized host supplies `postgresql_service`
- **THEN** the role manages that service without requiring an undefined alias

#### Scenario: Preserve an existing service override

- **WHEN** an existing inventory supplies only the legacy `postgresql` override
- **THEN** the role preserves that service selection and identifies its
  migration path without switching to a different service

### Requirement: Initialize only explicitly managed empty storage

Cluster initialization SHALL require explicit management inputs and validated
storage. It SHALL initialize only an empty data directory, reuse a compatible
initialized cluster, and fail safely on nonempty unrecognized data, missing
required storage, or an incompatible cluster version. Routine convergence
SHALL NOT reinitialize, replace, migrate, or delete existing database data.

#### Scenario: Initialize a fresh managed host

- **WHEN** initialization is enabled and required storage is present and empty
- **THEN** the role initializes the cluster once before starting the service

#### Scenario: Rerun against an initialized cluster

- **WHEN** a compatible cluster already exists with unchanged inputs
- **THEN** the role preserves its contents and avoids initialization and
  unnecessary service restarts

#### Scenario: Refuse unsafe initialization

- **WHEN** storage is missing, nonempty but unrecognized, or version-incompatible
- **THEN** the role fails before initializing data or starting against a
  fallback directory

### Requirement: Configure application access explicitly

Listener and authentication management SHALL be opt-in. When enabled, the
role SHALL converge the declared listener and ordered authentication rules
without adding broad trust access. Credentials SHALL be supplied by the
consumer and excluded from ordinary logs, diffs, and tracked fixtures.

#### Scenario: Restrict an application to local authenticated access

- **WHEN** a consumer declares loopback listeners and password-authenticated
  access for its application database and role
- **THEN** that access succeeds with the supplied credential, invalid
  credentials fail, and remote application access is not opened

### Requirement: Reconcile declared databases and extensions

The role SHALL manage explicitly declared application roles, database
ownership, and per-database extensions idempotently. Application roles SHALL
NOT receive superuser privileges by default. Privileged setup SHALL use a
separate local administrative context. Unrelated roles, databases, and
extensions SHALL remain unchanged; omission SHALL NOT imply deletion.

#### Scenario: Prepare an owned application database

- **WHEN** a consumer declares an application role, its owned database, and
  an available extension
- **THEN** the role creates or reconciles those objects, installs the extension
  in the requested database, and verifies the application can connect

#### Scenario: Preserve unrelated database objects

- **WHEN** managed objects already match and the cluster contains other objects
- **THEN** a repeat run preserves both managed and unrelated data and makes
  no unnecessary object changes

#### Scenario: Required extension is unavailable

- **WHEN** a declared extension cannot be installed for the selected version
- **THEN** setup fails with a redacted diagnostic before declaring the
  application database ready

### Requirement: Preserve consumers that do not opt into management

Existing initialized-host consumers SHALL retain package and service
management without silently acquiring application databases, authentication
replacement, listener changes, or cluster initialization. New management
inputs SHALL default to disabled or empty and have documented migration and
validation procedures.

#### Scenario: Run with existing minimal inputs

- **WHEN** an initialized host supplies only existing package and service inputs
- **THEN** the role manages those packages and service while preserving its
  cluster, authentication rules, listeners, and database objects
