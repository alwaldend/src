## Context

Shared ingress verifies HTTPS to its upstream and requires a browser client
certificate. The Flux Operator chart has a separate Vault OIDC client.

## Goals / Non-Goals

- Keep the Operator application objects in the Flux Terraform workload root.
- Renew backend certificates without storing a Vault login secret on the cluster.

## Decisions

- Bind Traefik to TCP 443 on the host network. Keep its health and metrics
  listeners on loopback. The host firewall and Traefik middleware admit the
  DC1 LANs and shared ingress WireGuard peers.
  Traefik runs as UID 0 with only `NET_BIND_SERVICE` so it can bind port 443
  on K3s. Its root filesystem is read-only and privilege escalation is disabled.
- Resolve the Operator hostname to the host from DC1 DNS and to shared
  ingress from global DNS. Vault OIDC applies to both paths.
- Publish the private backend address in global DNS so the ingress VM can
  resolve it using its cloud resolver.
- cert-manager authenticates to Vault with a short-lived service account token.
  Vault reaches the K3s API over the private LAN and verifies the API against
  the public CA already published by bootstrap. Only the cert-manager issuer
  service account can use the signing role.
  Its certificate request includes the Operator hostname as the common
  name and a DNS subject alternative name.
- Use the standard Kubernetes view role, which excludes Secret access, and
  prefix UI identities independently of the cluster API administrator group.

## Risks / Trade-offs

Vault certificate issuance requires TCP 6443 from the Vault host. The host
firewall scopes this exception to Vault. The workload root manages the Vault
auth endpoint configuration after the API is available.
