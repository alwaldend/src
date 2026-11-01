## Why

Appliances can use the shared native Lynis role with their already installed
Ansible interpreters. Historical Python probes and collectors duplicate runtime
and evidence ownership without improving the requested maintenance behavior.

## What Changes

- Use the owner inventory, existing authentication, and shared Lynis role.
- Remove version probes and the optional Python 3.6 collector/runtime survey.
- Let Lynis own private audit reports; retain boot verification and separate
  native maintenance evidence.
- Keep explicit audit, XOA OS-update, and XCP-ng vendor-update targets without
  a redundant appliance_updates_approved input.
- Preserve XCP-ng vendor prerequisites, raw status handling, placement recovery,
  and no-reboot boundaries; no new dependency or runtime installation is added.

## Capabilities

### New Capabilities

- `appliance-security`: Native appliance audit and separately selected
  no-reboot maintenance with private evidence.

### Modified Capabilities

None.

## Impact

Source remains under infra/security/appliances. Inventory and credential
ownership remain with infra/xcp_ng; controller host trust is unchanged.
The parent owns all live authorization and validation. This revision changes
source only and leaves captured scan runtimes and evidence untouched.
