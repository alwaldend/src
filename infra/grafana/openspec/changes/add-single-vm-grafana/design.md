## Context

The existing VictoriaMetrics deployment provisions a Fedora guest in DC1
through Xen Orchestra and uses Vault for deployment authentication and TLS.

## Goals / Non-Goals

Provide one persistent Grafana service with the usual infrastructure targets,
DC1-local browser access, and Vault OIDC. Dashboards and data sources are
configured by operators.

## Decisions

- Use the Fedora RPM, its systemd service, and SQLite on the persistent VM disk.
- Read the canonical DNS declaration directly in Terraform and both Ansible
  packages for the static IPv4 address and HTTPS hostname.
  Guest Traefik handles redirects, certificates, and proxying to loopback.
- Deploy Grafana through the collection role. Use Vault's UI authorization
  endpoint for browsers and its API endpoints for token exchange and UserInfo.
- Bind Traefik metrics to the guest address and allow the VictoriaMetrics
  collector to reach it.
- Use a dedicated Vault client assigned to `global_admins`. Grafana requires
  group membership and maps the group to the Grafana server administrator role.
  A separate local admin credential allows recovery from OIDC outages.
- Use separate Vault HTTP state for VM and DNS operations and a standalone
  Terraform root under `infra/grafana/approle` for Vault identity and
  certificate access. Use Vault administrative credentials and preserve the
  identity root's existing backend path.

## Risks / Trade-offs

- Guest or disk failure can interrupt service and lose dashboard changes.
- Confirm the configured address is available before provisioning. Vault must
  resolve the DC1 record and reach guest TCP port 80 for ACME.
