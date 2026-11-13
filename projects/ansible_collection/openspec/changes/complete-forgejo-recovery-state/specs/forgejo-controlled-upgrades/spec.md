## MODIFIED Requirements

### Requirement: Preserve an offline recovery archive

For a changed existing installation, the role SHALL drain an active instance
using the installed binary, stop the service, and create a unique protected
archive of its local SQLite data, repositories, storage, configuration,
executable, and service unit before installing candidate inputs. The archive
SHALL remain on the managed host beneath an administrator-only directory,
outside the backup input roots. The role SHALL validate declared roots and
coverage of the managed tree, configured SQLite file and startup files before
downtime. The caller SHALL declare all external local storage and secret-file
roots and coordinate persistence-path migrations separately. The role SHALL NOT
infer backup locations from arbitrary INI settings. A backup failure SHALL NOT install candidate
inputs and SHALL recover a previously active instance with intact inputs.

#### Scenario: Drain failure

- **WHEN** the installed instance cannot drain its queues
- **THEN** candidate installation and backup do not proceed
- **AND** the installed service and inputs remain available

#### Scenario: Successful existing-instance upgrade

- **WHEN** a changed candidate is activated
- **THEN** a consistent archive predating candidate installation is retained
- **AND** the existing account and repository data survive the upgrade
