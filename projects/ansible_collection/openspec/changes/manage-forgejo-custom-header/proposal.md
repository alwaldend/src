## Why

Forgejo loads its custom header template at startup. A deployment that
supplies a header needs the role's guarded activation and recovery workflow.

## What Changes

- Accept an optional controller-side custom header source.
- Stage and activate changed header contents with the service startup inputs.

## Capabilities

### New Capabilities

- `forgejo-custom-header`: Optional, recoverable custom header deployment.

### Modified Capabilities

None.

## Impact

The Forgejo Ansible role and callers that supply a custom header.
