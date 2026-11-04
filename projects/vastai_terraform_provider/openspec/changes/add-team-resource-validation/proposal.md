## Why

The provider needs account and team resource lifecycles that can be verified
without renting paid compute or configuring live credentials.

## What Changes

- Add SSH key, environment variable, team role, team, and API-key resources.
- Consume the shared generated SDK through its Options constructor.
- Verify lifecycle, imports, failed creation, external deletion, and key cleanup
  through the real Terraform CLI against loopback API fixtures.
- Deliver the provider implementation separately from Terraform configurations,
  AppRole provisioning, Vault policies, credential injection, and live state.

## Capabilities

### New Capabilities

- `team-resource-validation`: Non-compute resource lifecycle and team ownership.

## Impact

Extends the provider using the approved Framework dependency graph and existing
shared SDK. PR #249 retains the Terraform and credential integration. No live
API calls, infrastructure deployment, or Registry publication are included.
