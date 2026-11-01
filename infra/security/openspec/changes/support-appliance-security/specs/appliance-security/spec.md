## ADDED Requirements

### Requirement: Trusted appliance access

The workflow SHALL consume the owner inventory and credential selectors and
SHALL require controller-managed SSH host verification without accepting new
keys during connection.

#### Scenario: Missing controller host trust

- **WHEN** an appliance key is absent or differs
- **THEN** authentication fails without changing trust or exposing credentials

### Requirement: Shared private appliance audit

The workflow SHALL use `alwaldend.main.lynis_audit` for both appliances with
their already installed Ansible interpreters. It SHALL reuse the pinned Lynis
source and the shared role's private controller evidence lifecycle,
and clean temporary remote scratch without installing packages or applying
hardening. The XCP-ng interpreter SHALL come from its owning inventory.

#### Scenario: Supported appliance runtime

- **WHEN** the operator selects the audit entry point with a supported interpreter
- **THEN** the shared role runs Lynis and retains private evidence, and boot
  verification records unchanged boot identity

#### Scenario: Missing or unsupported interpreter

- **WHEN** the selected interpreter is absent or outside the pinned core's
  managed Python range
- **THEN** native Ansible execution fails without installing a runtime or
  falling back to a second Lynis implementation

#### Scenario: Scanner fails

- **WHEN** Lynis fails or times out
- **THEN** available evidence is retained privately before the task fails

### Requirement: Explicit no-reboot maintenance boundaries

The workflow SHALL separate updates from audits, use vendor package managers,
preserve boot identity, and refuse unsupported automatic-update paths.
Separately selected update targets SHALL NOT require an additional approval
variable; deployment authority remains with the operator. The workflow
SHALL NOT reboot, schedule reboot, run RPU, or stop/reboot guest VMs. The
XCP-ng control-plane update path SHALL preserve the original host placement
policy across transient host disable and toolstack restart when XAPI remains
responsive, and SHALL report failed restoration explicitly.

#### Scenario: XCP-ng has a modern audit interpreter

- **WHEN** the operator selects the separately authorized XCP-ng update path
- **THEN** maintenance still uses raw transport and vendor yum without assuming
  that the audit interpreter supplies vendor package-manager module bindings

#### Scenario: Unverified XCP-ng maintenance prerequisite

- **WHEN** pool state or native automatic-update support is unverified
- **THEN** the relevant operation fails with a concrete blocker before mutation

#### Scenario: Supported XOA OS maintenance

- **WHEN** the separately authorized XOA target selects a Debian host with
  existing module dependencies
- **THEN** the shared role updates packages with immediate and automatic reboot
  disabled, installing the authorized standard native component if needed

#### Scenario: XOA preparation or boot verification fails

- **WHEN** evidence preparation, host facts, maintenance, or boot comparison fails
- **THEN** an outer cleanup path still closes and removes allocated private SSH transport
- **AND** absent transport is skipped and the execution remains failed
