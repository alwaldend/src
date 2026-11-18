## Why

OpenCanary needs an independently managed Vault AppRole.

## What Changes

Move the OpenCanary AppRole and its SSH signing policy to a standalone Terraform root.

## Capabilities

### New Capabilities

- `standalone-opencanary-approle`: Independent AppRole provisioning.

### Modified Capabilities

None.

## Impact

The shared Vault root discovers the project's Vault identities by name.
