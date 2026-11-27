## Why

Flux workloads need ingress and certificate controllers installed through the
workload Terraform root.

## What Changes

- Add Terraform Helm releases for Traefik and cert-manager.
- Mirror the pinned controller charts and image indexes in Forgejo.

## Capabilities

### Modified Capabilities

- `infra-flux` manages Traefik and cert-manager in the workload root.
