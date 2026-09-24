## Context

See [proposal.md](proposal.md). The collection packages native service roles,
including Traefik, but the inspected baseline has no Nginx role. The consumer
will run Nginx on a VM, using existing Fedora-oriented provisioning patterns.

## Goals / Non-Goals

**Goals:** A small reusable role that installs Nginx, validates complete
configuration, and manages its systemd service predictably.

**Non-Goals:** TLS issuance, DNS, website rendering, archive extraction,
release management, or a generic reverse-proxy configuration framework.

## Decisions

- Add `roles/nginx` with normal defaults, tasks, handlers, templates, README,
  and Bazel role packaging. Extend collection aggregation once.
- Use the guest distribution package and native systemd unit. Keep install
  logic compatible with the actual selected image instead of inheriting the
  Caddy role's apt-only assumption. No container or custom Nginx build is
  needed for native JSON autoindex.
- Accept consumer-provided configuration, listener, and read-only document
  roots. The download consumer owns virtual hosts, CORS, listing behavior,
  and static-site routing. Avoid copying that configuration into role defaults.
- Validate the entire candidate configuration using Nginx before replacing
  the active configuration. A fragment-only syntax check is insufficient when
  it depends on includes or server context.
- Enable and start the service on initial deployment. Notify a reload handler
  for accepted changes; an unchanged rerun must not restart it. Invalid input
  must leave both the active service and its persistent configuration intact.
- Leave uploaded content and ownership boundaries under the consumer's
  control. The service account receives no publication privileges.

## Risks / Trade-offs

- Package defaults can add an unintended public listener -> replace or disable
  their default virtual host through declared role behavior and verify sockets.
- A valid candidate can still fail at runtime -> exercise HTTP behavior and
  document failures rather than treating a syntax check as service acceptance.
- Broad portability increases role complexity -> support the selected guest
  image first and document the actual distribution coverage.

## Migration Plan

This is a new optional role, so existing consumers do not change until they
select it. Package and exercise it in an isolated consumer first; the linked
download change owns any live Ansible execution.
