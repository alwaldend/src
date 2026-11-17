## Why

The NAS owner needs Terraform declarations for its existing storage and
Nextcloud lifecycle, independent of the existing DNS state.

## What Changes

- Pin the approved official TrueNAS provider and reuse the owner's Vault identity.
- Declare the existing pool, seven datasets, and Nextcloud lifecycle directly in HCL.
- Preserve imported addresses, deletion protection, encryption inheritance, and a no-change live plan.
- Keep only formatting and provider-schema checks; remove the generic variable interface, JSON inventory, unused share/snapshot resources, and their mock tests.

## Capabilities

### New Capabilities

- `truenas-services`: Manage the existing NAS storage and app lifecycle with independent state.

### Modified Capabilities

None. DNS keeps its existing root and state.

## Impact

Source belongs to `infra/nas`, with the approved provider pin under
`third_party/terraform`. Existing resources are already imported. Simplification
changes source only; it requires no state migration or NAS apply.
