## ADDED Requirements

### Requirement: Activate packaged K3s binary updates

The K3s role SHALL instruct the K3s installer to restart the service when
the role replaces the packaged binary.

#### Scenario: The packaged binary changes

- **WHEN** a role run replaces the installed K3s binary
- **THEN** the installer receives the force-restart option

#### Scenario: The packaged binary stays the same

- **WHEN** the copy task leaves the installed binary unchanged
- **THEN** the installer uses its normal service change detection
