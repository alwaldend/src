## Why

Keeping shared settings with their inventory group makes the inventory easier
to read and avoids repeating host variables.

## What Changes

Update `repo-infra` to prefer `group_vars` for settings shared by an inventory
group and `host_vars` for host-specific overrides. Add an evaluation case.

## Impact

Changes the canonical Ansible reference and its offline evaluation fixture.
