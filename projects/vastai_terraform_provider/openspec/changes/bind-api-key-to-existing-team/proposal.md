## Why

An existing personal issuer can create API keys, but team permission delegation
requires the owned team's context. A new credential must be bound without
changing the shared issuer or recreating the funded team.

## What Changes

- Allow an API-key resource to bind its newly generated key to an existing owned team.
- Retain the generated key in state before activation so failed binding is recoverable.
- Verify binding and revocation through Terraform CLI fixtures before live use.
- Default API-key lifecycle authentication to the configured provider; preserve only an optional same-apply bootstrap override.
- Remove redundant issuer inputs without rotating generated keys.

## Capabilities

### New Capabilities

- `existing-team-key`: Independently generated credentials for an existing owned team.

### Modified Capabilities

None.

## Impact

The API-key resource and its acceptance fixtures; no new dependency, team
creation, credit transfer, or shared-credential context change.
