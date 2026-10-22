## Purpose

Provide a reusable packaged Ansible role for native Nginx installation and
validated, idempotent service configuration on repository-managed hosts.

## ADDED Requirements

### Requirement: Packaged native Nginx service

The collection SHALL expose an `nginx` role that installs the distribution
package and manages its native service on the documented supported guest
distribution. Consumers SHALL configure listeners and static document roots
without modifying the role's implementation.

#### Scenario: Install the download consumer

- **WHEN** the consumer applies the role with a loopback listener and static roots
- **THEN** the packaged Nginx service starts with that configuration
- **AND** no package-default public listener remains unintentionally enabled

### Requirement: Validate before configuration activation

The role SHALL validate a complete candidate configuration before activating
it. An invalid candidate SHALL fail the deployment while retaining the last
valid persistent configuration and the running service.

#### Scenario: Submit an invalid virtual host

- **WHEN** a consumer supplies configuration rejected by Nginx validation
- **THEN** Ansible reports a failure
- **AND** existing HTTP content remains served using the previous configuration

### Requirement: Idempotent configuration and content separation

An unchanged role run SHALL NOT restart or reload the service. A valid changed
configuration SHALL be activated through a service reload. Applying the role
SHALL NOT upload, extract, delete, or select consumer release content.

#### Scenario: Repeat a completed deployment

- **WHEN** the same configuration is applied again
- **THEN** service lifecycle operations are not repeated
- **AND** consumer release files and selected-site links remain unchanged
