## ADDED Requirements

### Requirement: Route the Flux Operator UI

Shared ingress SHALL route the Flux Operator hostname to the Flux-owned
private HTTPS backend using the standard client certificate and backend TLS
verification settings.

#### Scenario: Render the Operator route

- **WHEN** the ingress playbook renders Traefik configuration
- **THEN** the Operator route uses the public Operator hostname as its match
  and backend TLS server name
- **AND** the backend address comes from the Flux DNS declaration
