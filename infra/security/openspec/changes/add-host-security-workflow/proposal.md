## Why

Repeatable fleet security maintenance must reuse owner inventories and preserve audit evidence without rebooting hosts.

## What Changes

- Add infra/security with scoped SSH authentication and composed owner inventories.
- Run pinned Lynis, update installed packages, and always enable native automatic updates in the full playbook.
- Preserve private reports and verify unchanged boot identities.

## Capabilities

### New Capabilities

- `host-maintenance`: Reusable fleet audit and no-reboot package maintenance.

## Impact

Uses the shared Lynis and package updater roles and existing operator Vault authentication. No new credentials or copied host inventory.
