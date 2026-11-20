## ADDED Requirements

### Requirement: Private Prometheus metrics

Ingress SHALL serve Traefik metrics on a dedicated port whose incoming
firewall rule is restricted to the VictoriaMetrics guest.

#### Scenario: Scrape an active ingress host

- **WHEN** VictoriaMetrics requests metrics through the ingress WireGuard address
- **THEN** Traefik responds with Prometheus metrics.
