## Why

Limit the Flux infrastructure project to its DNS records and the existing
Vault AppRole used to manage them.

## What Changes

- **BREAKING** Remove the Flux VM, Ansible, and Kubernetes packages, including
  the Flux-hosted reconciliation entry point for Harbor.
- Preserve the standalone DNS package and move the Flux-specific AppRoles,
  identity groups, PKI, OIDC, and Transit configuration into a Flux-owned
  Terraform root.
- Remove the deleted inventory from fleet security and point Harbor's
  encryption instructions at its own command.

## Capabilities

### Modified Capabilities

- `infra-flux` retains DNS and AppRole integration and retires cluster
  provisioning, configuration, and reconciliation.

## Impact

The [state transfer record](design.md#state-transfer) documents how the
existing Vault objects were adopted and the shared state bindings relinquished.
