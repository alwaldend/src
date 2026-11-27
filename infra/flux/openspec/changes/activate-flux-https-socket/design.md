## Decisions

- Systemd binds TCP 443 and starts its socket proxy on demand. The proxy
  passes bytes to Traefik over IPv4 loopback. Traefik retains its existing
  Kubernetes identity and terminates TLS using the cert-manager Secret.
- Keep both ports in the Flux ingress access declaration. The workload
  apply moves the old listener away from 443 before the operator runs the Flux playbook.
- Restrict TCP 443 to the management LAN, workload LAN, and shared ingress
  peer at the guest firewall. Traefik sees the proxy's loopback address.
  Vault OIDC and the shared edge's client certificate policy still apply.
