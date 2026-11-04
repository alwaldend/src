## Why

The new Vast.ai infrastructure project needs its own Vault identity and
secret namespace rather than sharing another component's credentials.

## What Changes

- Declare `src_infra_vastai` through the existing AppRole module.
- Grant the existing operator entity bootstrap access.
- Disable the unrelated Yandex Cloud policy and omit broad infrastructure
  group membership.

## Capabilities

### New Capabilities

- `vastai-identity`: Dedicated Vast.ai identity and operator bootstrap access.

## Impact

Changes only the checked-in Vault Terraform configuration. Consumer wiring
belongs to `infra/vastai/openspec/changes/add-vastai-project`.
An independently authorized Vault apply is required to deploy the role.
