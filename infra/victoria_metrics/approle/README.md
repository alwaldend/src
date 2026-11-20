---
title: Vault AppRole
---

This Terraform root owns the AppRole for the
[VictoriaMetrics deployment](../README.md). Provision
the shared Vault prerequisites, apply this root, then reconcile the shared
AppRole and Ansible memberships before synchronizing its Xen Orchestra
identity.

Use `bazel run //infra/victoria_metrics/approle:tf.plan` to review the
AppRole configuration. The root uses the Vault infrastructure credentials
and its existing Vault HTTP state path.

The AppRole can enroll Traefik through the component's Vault ACME directory.
The PKI role allows the VictoriaMetrics service hostname.

This root generates the web login and stores it in Vault. The
`global_admins` group inherits read access to the login through the
component's readers group. The deployment AppRole reads a separate bcrypt
verifier. Terraform's Vault-backed state contains the generated password,
so preserve its existing access restrictions. Reapply the guest playbook
after changing the credentials to activate the new verifier.
