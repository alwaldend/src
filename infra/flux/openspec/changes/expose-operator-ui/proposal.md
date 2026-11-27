## Why

Operators need browser access to the Flux status UI through the shared
client-certificate ingress.

## What Changes

- Enable the Flux Operator UI with Vault OIDC and read-only Kubernetes access.
- Configure private HTTPS and automatic certificate renewal through cert-manager.
- Route the Operator DNS name to shared ingress.

## Capabilities

### Modified Capabilities

- infra-flux adds the Terraform-managed Operator web endpoint.

## Impact

Flux AppRole, DNS, host firewall, and workload configuration change.
The shared ingress project owns the public Traefik route.
