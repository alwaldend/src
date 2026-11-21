## ADDED Requirements

### Requirement: Grafana service role

The Ansible collection SHALL provide a role that installs Fedora's Grafana
package, configures loopback HTTP and OIDC, and manages the systemd service.
It SHALL install caller-supplied credential files with restricted permissions.

#### Scenario: Deploy through the role

- **WHEN** a caller supplies its hostname, OIDC configuration, and credential files
- **THEN** Grafana starts and the role checks its loopback database health endpoint
