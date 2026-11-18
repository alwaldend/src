# OpenCanary Vault AppRole

This Terraform root owns the `src_infra_opencanary` AppRole and its SSH host
signing policy in separate Vault-backed state. The shared Vault root owns the
mounts, the operator group, and the shared AppRole and Ansible memberships.

For a fresh installation, provision the shared Vault prerequisites, apply
this root, then reconcile shared group memberships before configuring Xen
Orchestra. Review with `bazel run //infra/vault/approles/src_infra_opencanary:tf.plan`.
If an earlier candidate already provisioned the role in shared state, migrate
its state to this root before applying either root.
