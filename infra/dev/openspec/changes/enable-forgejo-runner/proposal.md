## Why

Secure1 has a repository-scoped Forgejo registration but no deployed worker.
Enable the worker through its owning Ansible deployment at the user's request.

## What Changes

- Deploy the existing runner role on secure1 using its cached Vault credentials.
- Install the existing repository build prerequisites for the runner account.
- Add strictly verified, Vault-authenticated SSH inspection for secure1.
- Document the host executor trust boundary and current workflow activation state.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-forgejo-runner`: replace the host-only baseline with an enabled worker.

## Impact

infra/dev Ansible, its Vault credential references and documentation; the shared
runner role gains support for existing API-issued credentials. No new dependencies,
workflow activation, VM replacement, deployment identity in jobs, or host reboot.
