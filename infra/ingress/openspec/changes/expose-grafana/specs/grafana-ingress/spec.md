## ADDED Requirements

### Requirement: Client-authenticated Grafana ingress

Shared ingress SHALL serve the Grafana hostname with its standard mutual TLS
policy and verify backend HTTPS using the Grafana service hostname.

#### Scenario: Connect to Grafana from outside DC1

- **WHEN** a client presents a trusted certificate to the shared ingress
- **THEN** the Grafana request reaches the private guest through its dedicated backend hostname

#### Scenario: Connect without a trusted certificate

- **WHEN** a client connects to the Grafana ingress without a trusted certificate
- **THEN** ingress rejects the TLS connection
