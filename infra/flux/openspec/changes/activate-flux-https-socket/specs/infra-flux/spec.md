## ADDED Requirements

### Requirement: Socket-activated HTTPS

Flux SHALL expose the Operator on TCP 443 using a systemd socket that
passes TLS to non-root Traefik bound to loopback inside the host network.

#### Scenario: Connect over HTTPS

- **WHEN** a client allowed by the host firewall connects to TCP 443
- **THEN** systemd activates its TCP proxy to the loopback listener
- **AND** Traefik terminates TLS and routes to the Operator in Kubernetes
