## Why

Flux and Harbor package K3s v1.36.2+k3s1 through the shared Ansible role.
Update both release artifacts to the current stable channel version. The
role copies the binary before invoking the installer, so the installer can
miss a binary change and leave the old process running.

## What Changes

- Pin the upstream Linux amd64 binary and installer source to v1.37.2+k3s1.
- Tell the installer to restart K3s when Ansible replaces the binary.

## Capabilities

### New Capabilities

- `ansible-k3s` describes how the role activates a packaged binary upgrade.

### Modified Capabilities

None.

## Impact

The shared third-party dependency and Ansible role affect the Flux and Harbor
deployment packages. Acceptance covers artifact integrity, the reported binary
version, the installer's restart input, and both playbooks' syntax checks.
