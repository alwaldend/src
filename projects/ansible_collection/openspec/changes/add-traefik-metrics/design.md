## Context

Traefik already ships a Prometheus handler. Managed hosts use TOML
configuration and either UFW or a managed firewalld zone.

## Goals / Non-Goals

Collect metrics privately without adding an exporter or creating public routes.

## Decisions

Define the handler on a dedicated listener in each deployment's static
configuration. Gate mixed inventories by group. Bazel passes the DNS-derived collector address and the port from the
Traefik role defaults directly to Ansible as an extra-vars file. The
package wrapper accepts multiple files in precedence order.
Keep collector-only access in the deployment firewall rules or the host's
existing firewalld zone.

## Risks / Trade-offs

A host needs a reachable private route for collection. Public-only hosts
require a different transport before they can be monitored.
