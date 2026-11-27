## Decisions

- Provision the VM on XCP-ng through the existing Xen Orchestra provider. The
  Flux AppRole uses the shared Vault OIDC login plugin and its own resource set.
  The VM root uses a dedicated XCP-ng state path.
- Reuse the shared Xen cloud-init configuration and Ansible host and K3s roles.
- Configure the API server to trust the existing Flux Vault OIDC issuer and
  add a single ClusterRoleBinding for the Flux administrator group. This
  binding must exist before Terraform can install charts with cluster-scoped
  resources. An AL resource opens the management tunnel with Vault-signed SSH
  certificates before each Terraform invocation.
  Ansible publishes the K3s public server CA to the Flux Vault namespace.
  AL provides Helm with that CA and an OIDC ID token for each invocation.
- Use the pinned HashiCorp Helm provider for release lifecycle and dependency
  ordering. Install the mirrored operator chart first, then a project-local
  chart containing the FluxInstance.
- Use the distribution data embedded in the operator and pull workload images
  from Forgejo. The Flux root uses a separate Vault state backend.
