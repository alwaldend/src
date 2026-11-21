---
title: Grafana Vault AppRole
---

This root owns the Grafana AppRole, DC1 DNS credential access, guest SSH
host signing, Traefik's ACME enrollment, and the application OIDC client.
The OIDC client allows sign-in for `global_admins` and exposes its client
credentials to the deployment AppRole. The AppRole also has read access to
the VictoriaMetrics web credential used by the Grafana Terraform data source.

Provision the shared Vault mounts, administrator group, and OIDC scopes
first. Review this root using
`bazel run //infra/grafana/approle:tf.plan`. Apply the root
before reconciling shared Vault memberships and the Xen Orchestra resource
set. The [Grafana runbook](../README.md) gives the deployment
sequence.
