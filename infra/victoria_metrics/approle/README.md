---
title: Vault AppRole
---

This Terraform root owns the AppRole for the
[VictoriaMetrics deployment](../README.md). Provision
the shared Vault prerequisites, apply this root, then reconcile the shared
AppRole and Ansible memberships before synchronizing its Xen Orchestra
identity.

Use `bazel run //infra/victoria_metrics/approle:tf.plan` to review the
AppRole and `tf.apply` to deploy it. The root uses the Vault infrastructure
credentials for its policies and HTTP state.

The AppRole can enroll Traefik through the component's Vault ACME directory.
The PKI role allows the VictoriaMetrics service hostname.

The `global_admins` group inherits read access to the login through the
component's readers group. The component AppRole manages the `web_auth`
entry in its own subtree. The separate
[`tf_setup` package](../tf_setup/web_auth/README.md) generates and stores the
credentials using these permissions.
