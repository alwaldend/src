## Context

MLflow's Projects backend loads in the submitting CLI. The coordinator discovers
existing rentals and calls their preinstalled workers over host-verified SSH.

## Goals / Non-Goals

Provide a persistent, authenticated tracking endpoint reachable by clients and
rentals, and an installed CLI that can submit to the colocated coordinator.
Rental creation and training environment installation belong to rental owners.

## Decisions

- Use a dedicated Xen Orchestra guest, SQLite and proxied artifacts on persistent
  guest storage. This follows service ownership and avoids another database or
  object store. This is a single-host service.
- Use the existing client-certificate ingress externally and terminate guest
  HTTPS at Traefik. The MLflow OIDC plugin authenticates browser users against
  Vault and verifies separate, issuer-bound workload bearer tokens. The
  coordinator rereads its token file on each request. The AppRole root
  registers the Vault OIDC client, and `tf_setup` creates the guest and writes
  the session key into the component's Vault namespace.
- Package the plugin-enabled executable and pinned Python runfiles with Bazel.
  Installing from a mutable Python index on the guest would bypass the lock.
- Give tracking, artifacts, the coordinator, and Traefik independent data disks.
  The same ordered storage declaration drives XO disk attachment and Ansible
  mounts. Systemd requires the expected mounts before each service starts.
