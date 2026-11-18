# src_projects_al Vault project

This standalone root owns the existing src_projects_al AppRole and its local
policies and identity resources. It uses the existing Vault administrative
identity and a separate HTTP backend under
`src_infra_dc1_vault/tf_backend/approles/src_projects_al`.

The shared mounts and global membership groups remain in the core root. It
looks up this project's existing identity IDs by name. No remote state is read.
Bootstrap shared mounts first, then this project, then core memberships.

Use `bazel run //infra/vault/approles/src_projects_al:tf.plan` or `:tf.apply`.
Migration from core state preserves live IDs and permissions; do not recreate
resources that are already deployed. The test-only `:offline` target supports
backend-free validation and local state migration with the pinned Terraform.
