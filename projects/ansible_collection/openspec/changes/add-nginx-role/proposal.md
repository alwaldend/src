## Why

The download service needs a small native Nginx deployment behind the existing
Traefik role. A reusable collection role should own package installation,
configuration validation, and service lifecycle while consumers own routes
and storage layout.

## What Changes

- Add and package `roles/nginx` in `alwaldend.main`.
- Support consumer-provided static-server configuration and a loopback
  listener, with Nginx configuration validation before activation.
- Reload only after a valid configuration change and preserve the active
  configuration when a candidate is invalid.
- Support JSON autoindex and symlink-selected static website roots through
  ordinary Nginx configuration; leave DNS, TLS, theme, and content publication
  with their existing owners.

## Capabilities

### New Capabilities

- `ansible-nginx`: Packaged native Nginx installation, validated configuration,
  and idempotent systemd lifecycle.

### Modified Capabilities

None. The existing collection aggregation contract remains applicable.

## Impact

Changes belong in `roles/nginx` and the collection's role aggregation. The
consumer is [static hosting](../../../../../infra/download/openspec/changes/add-static-hosting/proposal.md).
The role must be tested on the selected VM distribution and exercised through
that consumer's HTTP acceptance scenarios. This plan does not run Ansible.
