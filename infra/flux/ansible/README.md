---
title: Flux K3s bootstrap
---

After provisioning the VM, wait for cloud-init and verify the initial SSH host
key out of band. Run `bazel run //infra/flux/ansible`. AL uses the Flux AppRole
to obtain a short-lived Vault-signed SSH client certificate for the `ansible`
user, which cloud-init configures on the guest. Ansible installs a Vault-signed
host certificate during bootstrap.

The shared K3s role installs the repository-pinned binary and writes the server
configuration. The API server trusts the existing Flux Vault OIDC issuer.
The only project manifest in the auto-deploy directory grants Kubernetes
cluster administration to the Vault group `src_infra_flux_admins`, which
includes the deployment AppRole. Terraform uses that grant to install charts
with cluster-scoped resources.

After K3s starts, the shared role publishes the public API server CA at the
Vault path configured by Flux. It updates that entry if the CA has changed.

The host firewall accepts SSH and lets cluster pods reach the API.
The [Terraform workload root](../tf/README.md) opens its SSH forward
automatically and authenticates with a Vault OIDC token.
