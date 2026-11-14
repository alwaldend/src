## Purpose

Preserve recoverable local Forgejo state while activating and verifying a
changed service executable or configuration through the Ansible role.

## ADDED Requirements

### Requirement: Detect and validate a candidate before downtime

The Forgejo role SHALL compare installed and candidate executable,
configuration, and unit content. It SHALL reject an older version, incomplete
installation, unsupported backup layout, and an invalid candidate before
stopping a running service or replacing installed startup inputs.
It SHALL reject removing or changing the configured database encryption key.

#### Scenario: Unchanged deployment

- **WHEN** installed startup inputs equal the candidate
- **THEN** the role ensures the service is started without draining queues,
  archiving data, reloading, or restarting the active process

#### Scenario: Reject a downgrade

- **WHEN** the candidate version is older than the installed version
- **THEN** the role fails before modifying installed startup inputs or stopping
  the running instance

### Requirement: Preserve an offline recovery archive

For a changed existing installation, the role SHALL drain an active instance
using the installed binary, stop the service, and create a unique protected
archive of its local SQLite data, repositories, storage, configuration,
executable, and service unit before installing candidate inputs. The archive
SHALL remain on the managed host beneath an administrator-only directory,
outside the backup input roots. Unsupported external persistence or uncovered
paths SHALL fail before downtime. A backup failure SHALL NOT install candidate
inputs and SHALL recover a previously active instance with intact inputs.

#### Scenario: Drain failure

- **WHEN** the installed instance cannot drain its queues
- **THEN** candidate installation and backup do not proceed
- **AND** the installed service and inputs remain available

#### Scenario: Successful existing-instance upgrade

- **WHEN** a changed candidate is activated
- **THEN** a consistent archive predating candidate installation is retained
- **AND** the existing account and repository data survive the upgrade

### Requirement: Verify activation without unsafe automatic rollback

The role SHALL start a changed candidate with systemd and require successful
HTTP health, expected running version, and non-repairing doctor checks. If
installation or activation fails, it SHALL stop the candidate and report the
retained recovery location. It SHALL NOT start the prior binary against data
that the candidate may have migrated. Fresh installations SHALL use the same
acceptance checks without claiming a previous recovery archive exists.
An unfinished activation SHALL block subsequent runs until an operator has
deliberately resolved its recovery or retry state.

#### Scenario: Candidate fails after migration may have started

- **WHEN** startup or a subsequent acceptance check fails after installation
- **THEN** the candidate is stopped and recovery files remain available
- **AND** recovery requires explicit restoration of compatible data and inputs

#### Scenario: Rehearse recovery

- **WHEN** the local integration fixture restores the complete archive
- **THEN** the prior binary starts with its matching database and the existing
  account and repository remain accessible
