## Why

The download deployment playbook contains reusable host configuration. Extract
it into `alwaldend.main.download_host` so it can be reviewed and merged as an
independent collection change before the deployment consumes it.

## What Changes

- Move host tasks and deduplication templates into `download_host`, with
  filesystem management in `tasks/filesystem.yaml`.
- Document caller variables and the play-level `force_handlers` requirement.
- Create publication roots only; publishers own per-site directory creation.
- Keep Traefik independent of the content disk.
- Keep attachment checks, publisher setup, and daily maintenance; remove the
  reviewed capacity and play-size restrictions, service overrides, and undeployed
  migration cleanup. Terraform owns sizing; inventory owns target selection.

## Capabilities

No specification-level behavior changes; this is a refactor of the implementation
in download PR #102. `skip_specs: true` records that scope.

## Impact

The role is packaged by the existing collection aggregation and targets trunk.
The consumer conversion in PR #102 follows this prerequisite's merge. Inventory,
Vault injection, routing templates, publishing, and live deployment are outside
this role PR.
