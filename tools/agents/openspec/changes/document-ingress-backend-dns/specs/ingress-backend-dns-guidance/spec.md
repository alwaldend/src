## ADDED Requirements

### Requirement: Document the ingress backend DNS convention

The infrastructure skill SHALL document the public ingress name and dedicated
backend name when ingress must resolve a private upstream.

#### Scenario: Add a private backend route

- **WHEN** an infrastructure change routes a public hostname through ingress
- **THEN** the skill specifies that the public name routes to ingress and the dedicated backend name resolves to the private address from ingress
