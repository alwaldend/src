## Purpose

Support the official standalone T3 Code distribution through the existing
shared role while retaining its established npm installation mode.

## ADDED Requirements

### Requirement: Explicit standalone installation

The role SHALL support an explicit standalone mode that verifies an archive
checksum, installs administrator-owned runtime files, and invokes the
installed executable directly. The default npm mode SHALL remain available.

#### Scenario: Start an installed standalone service

- **WHEN** the role renders a standalone service
- **THEN** its command runs the installed T3 executable without npm or npx

### Requirement: Idempotent activation and private authentication

Changed startup inputs SHALL restart the service, unchanged inputs SHALL
leave a running service alone, and authentication environment writes MUST
be protected from Ansible logging. Optional initial settings SHALL seed an
absent settings file without replacing later UI changes.

#### Scenario: Reapply an unchanged configuration

- **WHEN** runtime and service inputs are unchanged
- **THEN** the role ensures the service is running without restarting it

#### Scenario: Retry interrupted extraction

- **WHEN** extraction left an executable but incomplete supporting files
- **THEN** reapplication repairs the bundle before service activation
