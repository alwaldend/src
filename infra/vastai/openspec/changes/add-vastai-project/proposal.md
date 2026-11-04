## Why

Vast.ai needs an infrastructure owner with a dedicated Vault authentication
context, following the existing infrastructure project conventions.

## What Changes

- Add the project README, AL configuration, and internal Bazel targets.
- Select `src_infra_vastai` for Vault commands and a component-local Terraform
  HTTP backend for future service configuration.
- Register this owner's OpenSpec source in the generated validation inventory.
- Inject the component secret's `api_key` field as `VAST_API_KEY` for an
  authenticated wrapper around the existing pinned CLI.

## Capabilities

### New Capabilities

- `vastai-configuration`: Project ownership, AL authentication, and scoped
  state-backend configuration, including API-key injection for the CLI.

## Impact

Adds `infra/vastai` and updates generated OpenSpec inventory metadata.
The identity declaration is tracked separately in
`infra/vault/openspec/changes/add-vastai-approle`.
The CLI wrapper reuses `tools/vastai` and the existing Vault injector. It adds
no provider dependency or resource provisioning.
