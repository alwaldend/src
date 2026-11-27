## Why

Flux needs a host listener on TCP 443 while Traefik remains non-root and
continues using its Kubernetes ServiceAccount.

## What Changes

- Activate a systemd socket on TCP 443 and pass TLS through to loopback Traefik.
- Run the Helm-managed Traefik pod non-root on TCP 8443.
- Enforce allowed source networks at the host firewall.

## Capabilities

### Modified Capabilities

- `infra-flux` assigns the HTTPS host socket to systemd.

## Impact

Terraform moves Traefik to loopback before the Flux playbook claims
TCP 443. Shared ingress continues to verify HTTPS to the Flux backend.
