## Why

The live Vast.ai account contains configuration outside the bootstrap catalog.
Reconcile recoverable resources without replacing existing teams or credentials.

## What Changes

- Identify missing personal-account declarations and import supported resources
  when their desired configuration is recoverable.
- Expose the root's authenticated Terraform import command.
- Declare and import ENABLE_HTTPS with the user-selected value `true`.
- Declare and import the three existing owned teams without generating credentials.
- Revoke only the unmanaged team API key outside IaC, as explicitly requested,
  while preserving account keys and all Terraform-owned credentials.
- Remove import blocks after adoption and record inventory coverage limits.

## Capabilities

### New Capabilities

- `account-import`: Import supported existing resources through the authenticated infrastructure root.

### Modified Capabilities

None.

## Impact

Changes are confined to infra/vastai. Existing team and token addresses remain
unchanged. No dependencies, compute rentals, credit transfers, credential
rotations, or team recreation are required. Provider team import support is
owned by projects/vastai_terraform_provider/openspec/changes/import-existing-teams.
Account keys remain operator-managed under the user's explicit exception.
The API masks ENABLE_HTTPS; the user selected `true` as its desired value.
