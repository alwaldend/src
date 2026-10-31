## Purpose

Provide one shared role for package maintenance and native automatic updates,
with explicit enrollment and reboot policies across supported host backends.

## ADDED Requirements

### Requirement: One package maintenance owner

The shared host configuration SHALL invoke one canonical package-update role.
The normal entry point SHALL select the backend once. The configuration-only
tag SHALL configure automation without performing an immediate upgrade;
internal task files SHALL NOT be standalone public entry points.
The normal role SHALL also support `update_all_packages_upgrade_enabled: false`,
which skips immediate upgrades, reboot detection, and reboot while retaining
native automatic-update configuration. The control SHALL default to true.

#### Scenario: Use the normal entrypoint for configuration only

- **WHEN** the caller disables immediate upgrades through the role variable
- **THEN** native automatic-update policy and timers are reconciled
- **AND** immediate upgrades, reboot detection, and reboot are skipped even if
  the separate reboot control is enabled

#### Scenario: Configure automatic updates independently

- **WHEN** a caller selects the `automatic_updates` tag through the normal role
- **THEN** native updater configuration and scheduling are reconciled
- **AND** no immediate package upgrade or reboot runs

### Requirement: Supported immediate updates

Immediate updates SHALL refresh mutable-host metadata and update all installed
packages allowed by enabled repositories and package-manager policy. An OSTree
boot marker SHALL select deployment upgrades even when DNF is installed.
Unsupported appliances and package-manager failures SHALL fail explicitly.

#### Scenario: Update an immutable Fedora host

- **WHEN** an OSTree-booted host requests an immediate update
- **THEN** an updated deployment is prepared for a later boot
- **AND** no mutable package upgrade or live application runs

### Requirement: Explicit immediate reboot policy

Immediate reboots SHALL default to disabled and SHALL require both an enabled
policy and a positive reboot requirement. Reboot detection SHALL distinguish
errors from negative results during immediate maintenance.

#### Scenario: Updates require a reboot during no-reboot maintenance

- **WHEN** the backend reports reboot required and the reboot policy is false
- **THEN** the role reports the requirement without initiating or scheduling a reboot

### Requirement: Native automatic updates without automatic reboots

Automatic-update enrollment SHALL default to enabled through the single
`update_all_packages_automatic_enabled` control. Enabled mutable hosts
SHALL reuse the existing distribution updater packages and security-update
scope. Enabled OSTree hosts SHALL use native staged updates. No custom timer
SHALL be introduced, and no managed updater SHALL request automatic reboot.
Unrelated vendor configuration SHALL remain intact.

#### Scenario: Use the default maintenance policy

- **WHEN** a caller supplies neither automatic-update control nor a reboot control
- **THEN** native automatic updates are enabled
- **AND** immediate and automatic reboots remain disabled

#### Scenario: Enable mutable-host automatic updates

- **WHEN** automatic updates are enabled on a supported mutable host
- **THEN** the existing native updater package and timer are configured
- **AND** security-update policy and explicit no-automatic-reboot settings apply

#### Scenario: Disable automatic updates

- **WHEN** `update_all_packages_automatic_enabled: false` is supplied
- **THEN** an installed managed updater's installation timer is stopped and disabled
- **AND** no missing updater package is installed

#### Scenario: Stage automatic OSTree updates

- **WHEN** automatic updates are enabled on an OSTree host
- **THEN** the native automatic timer uses the staging policy after daemon reload
- **AND** automatic application and reboot policies are not selected

#### Scenario: A native updater has a separate reboot hook

- **WHEN** inspection finds a known automatic reboot command or systemd action
- **THEN** configuration fails with an actionable error before enabling that updater
- **AND** unrelated administrator overrides are not overwritten
