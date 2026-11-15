## Why

The user requested recreation of secure1 under infra/dev instead of migrating
its old Terraform state. Its development identity needs an XO resource set.

## What Changes

- Assign the existing Fedora cloud template, local storage and wired network to src_infra_dev.
- Use the current per-component OIDC and resource-set authorization workflow.

## Capabilities

### New Capabilities

- `development-resource-set`: Scoped virtualization allocation for the development owner.

### Modified Capabilities

None.

## Impact

infra/xcp_ng Terraform inventory and the linked infra/dev and Vault changes.
No new provider, shared administrator credential or pool-wide permission is added.
