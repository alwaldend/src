## ADDED Requirements

### Requirement: Scoped Grafana identity

Vault SHALL provide a Grafana-specific AppRole with DC1 DNS credential access,
SSH host signing, ACME enrollment, and its own OIDC client credentials.

#### Scenario: Grafana OIDC login

- **WHEN** an entity requests Grafana sign-in from Vault
- **THEN** Vault authorizes only the existing global administrators group and uses the Grafana HTTPS callback
