## ADDED Requirements

### Requirement: Automatic staging without reboot

Host Bot SHALL use the canonical updater to enable automatic OSTree staging and SHALL disable automatic reboot.
It SHALL invoke the normal role entrypoint with immediate upgrades disabled in
owner variables, rather than selecting an internal task file.

#### Scenario: Reconcile automatic updates

- **WHEN** the automatic update configuration is deployed
- **THEN** the native timer SHALL be enabled with staging policy and the running boot SHALL remain unchanged

### Requirement: Preserve approved appliance SSH trust

Host Bot SHALL retain the operator-approved XCP-ng and XOA public host keys in the existing user known_hosts file without replacing unrelated entries.

#### Scenario: Reconcile appliance trust

- **WHEN** the SSH trust tasks run
- **THEN** the approved entries SHALL be present and other known hosts SHALL remain intact
