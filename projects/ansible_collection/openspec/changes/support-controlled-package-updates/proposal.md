## Why

Package maintenance needs one owner for immediate updates and native automatic
updates across mutable and OSTree hosts, with reboots disabled by default.

## What Changes

- Keep immediate and automatic updates in `update_all_packages`.
- Default `update_all_packages_automatic_enabled` to true and immediate reboots
  to false. Provide a self-contained automatic-only task entrypoint and tag.
- Preserve native security-update policy, OSTree staging, package-manager
  restrictions, effective no-reboot checks, and explicit failure reporting.

## Capabilities

### New Capabilities

- `controlled-package-updates`: Unified immediate and scheduled updates with
  backend-aware execution, native timers, and explicit reboot policy.

### Modified Capabilities

None.

## Impact

The shared updater and Host Bot's automatic-only caller change. Existing native
updater packages remain the only dependencies; no custom timer is introduced.
This source revision performs no live operations. PR202 remains separate and
will adopt the canonical updater when its owner rebases it.
