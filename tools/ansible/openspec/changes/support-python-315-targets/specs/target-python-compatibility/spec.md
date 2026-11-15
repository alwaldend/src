## ADDED Requirements

### Requirement: Supported Python 3.15 targets

The packaged Ansible CLI SHALL execute modules on dev1's Python 3.15 and secure1's existing interpreter using supported ansible-core. Dependency versions SHALL be resolved by the owning uv lock workflow. No separate host interpreter or unsupported third-party patch SHALL be required.

#### Scenario: Gather development host facts

- **WHEN** the packaged deployment gathers facts on secure1 and dev1
- **THEN** both module executions complete under strict SSH host verification
- **AND** dev1 does not fail because dataclasses._is_type is absent

### Requirement: Packaged native SSH askpass

Ansible wrappers SHALL package the matching ansible CLI and make it available
to native SSH askpass without requiring a host-installed sshpass executable or altering the caller PATH.

#### Scenario: Inspect a password-authenticated appliance

- **WHEN** the packaged playbook uses ssh_askpass and a Vault-injected password
- **THEN** its helper resolves from pinned runfiles and the SSH connection completes
