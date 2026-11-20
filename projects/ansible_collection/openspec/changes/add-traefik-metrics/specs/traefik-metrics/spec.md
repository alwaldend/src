## ADDED Requirements

### Requirement: Shared Traefik metrics port

The role SHALL supply the default metrics port for deployment-owned Traefik
configuration. Deployments SHALL select which hosts expose the listener.

#### Scenario: Render an opted-in deployment

- **WHEN** the deployment configures an enabled Prometheus entry point
- **THEN** its listener uses the same port provided to the scrape targets and firewall.
