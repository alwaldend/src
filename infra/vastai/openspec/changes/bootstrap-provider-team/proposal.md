## Why

The provider's test project needs a Terraform-generated team and project token,
using only a manually supplied personal bootstrap token in this component.

## What Changes

- Create a zero-credit test_ team and generated project API key with the provider.
- Write the generated project key directly to its own Vault credential entry.
- Keep bootstrap state separate from the provider project's test state.

## Capabilities

### New Capabilities

- `provider-team-bootstrap`: Terraform-managed team and credential delivery.

## Impact

Uses the source-built provider and existing Vault provider. Vault owns the
dedicated project AppRole and its existing secret namespace permissions.
No compute, credit transfer, host configuration, or manual project token is used.
