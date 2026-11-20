## Why

Provide a small VictoriaMetrics deployment for collecting and exploring
Prometheus-compatible metrics on one VM.

## What Changes

- Define one Fedora guest on XCP-ng and a dedicated Vault deployment identity.
- Manage VictoriaMetrics single-node with Ansible and systemd.
- Configure a loopback scrape job and persistent metric storage.
- Publish component-owned split-view DNS and expose VMUI through authenticated ingress.

## Capabilities

### New Capabilities

- `single-node-metrics`: Provisions one guest and configures a local metrics database and VMUI.

### Modified Capabilities

None.

## Impact

Adds Terraform and Ansible packages under `infra/victoria_metrics`, including
the dedicated AppRole root. The shared Vault and Xen Orchestra roots manage
the AppRole's group memberships and resource set.
