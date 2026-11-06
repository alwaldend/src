## Why

Existing owned teams cannot be adopted because the provider couples their
identity to credentials generated during creation. The account API can identify
owned teams without recovering or changing those credentials.

## What Changes

- Import existing teams by canonical numeric account ID using read-only metadata.
- Leave unrecoverable generated credentials null and refuse deletion of imported teams.
- Verify successful import, unchanged plans, failure preservation, and zero remote mutations through Terraform CLI fixtures.

## Capabilities

### New Capabilities

- `team-import`: Adopt owned team identities without generating credentials or changing account context.

### Modified Capabilities

None.

## Impact

The provider team resource, acceptance fixtures, and documentation change.
Created-team lifecycle and cleanup remain unchanged. infra/vastai consumes
the new importer through its own reconciliation change. No new dependencies
or live credential generation are required.
