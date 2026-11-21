## Why

Operators need a Grafana instance for dashboards on the DC1 network.

## What Changes

- Provision one Fedora VM through Xen Orchestra and deploy Grafana with Ansible.
- Serve Grafana through guest Traefik with Vault-issued HTTPS and DC1 DNS.
- Configure Vault OIDC for the existing global administrators and a local recovery login.
- Allow administrative SSH from any routed source.

## Capabilities

### New Capabilities

- `single-vm-grafana`: Single-host provisioning, service deployment, and browser sign-in.
- `grafana-identity`: Component-owned Vault access, certificate signing, and OIDC.

### Modified Capabilities

None.

## Impact

Adds a project under `infra/grafana` with linked Vault and Xen Orchestra registration.
