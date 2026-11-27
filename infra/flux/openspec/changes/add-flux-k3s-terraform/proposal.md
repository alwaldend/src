## Why

Provide a minimal K3s host for Flux and make Terraform responsible for the
Flux application lifecycle.

## What Changes

- Provision the Flux VM on XCP-ng through a Vault-backed Xen Orchestra session
  and bootstrap K3s using repository-pinned artifacts.
- Bootstrap the Vault OIDC authenticator and its administration binding.
  Leave the K3s auto-deploy directory free of project workloads.
- Deploy the Flux Operator and FluxInstance from a separate Terraform root
  authenticated with a Vault OIDC token over the SSH tunnel.

## Capabilities

### Modified Capabilities

- `infra-flux` defines independent host bootstrap and workload lifecycle.
