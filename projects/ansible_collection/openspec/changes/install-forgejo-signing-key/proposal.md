## Why

The Forgejo deployment needs an optional dedicated SSH commit-signing pair
installed before service activation, with protection against implicit rotation.

## What Changes

- Add optional signing-pair inputs and guarded installation under the Forgejo root.
- Verify the deployed pair through actual instance signing in the existing fixture.

## Capabilities

### New Capabilities

- `forgejo-signing-key`: Optional guarded installation of instance signing files.

### Modified Capabilities

None.

## Impact

Forgejo role and its existing upgrade acceptance fixture. The service-specific
Vault lifecycle remains owned by infra/forgejo's add-commit-signing-key change.
