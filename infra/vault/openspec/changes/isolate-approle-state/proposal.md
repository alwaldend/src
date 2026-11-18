## Why

The compact central Vault Terraform state grew above 1 MiB while applying
NAS prerequisites. The six changes completed but the new state could not be
saved. Recover its protected state and move selected AppRole projects into
independent states without restarting Vault or recreating live resources.

## What Changes

- Give src_infra_nas, src_infra_dev and src_projects_al standalone Vault roots,
  separate backend paths and plan/apply/state targets.
- Preserve their existing account, DNS, SSH and PKI behavior; shared groups
  discover their existing entity/group IDs by name.
- Split the protected recovery state with native Terraform state moves, push
  without force, and verify no-change resource plans before continuing deployment.

## Capabilities

### New Capabilities

- `approle-state`: Independently managed Vault AppRole projects and state recovery.

### Modified Capabilities

None.

## Impact

Uses the existing Vault provider and independent download AppRole pattern.
No dependency, Vault restart, storage-format change or server-limit change.
Live migration is explicitly authorized; secret-bearing recovery files remain
private and are removed only after all remote states are verified.
