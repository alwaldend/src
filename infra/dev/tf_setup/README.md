# Secure1 provisioning

This root owns secure1's Xen Orchestra VM and uses the development AppRole's
own resource set and Vault backend at `src_infra_dev/tf_backend/tf_setup`.
Its single 500 GiB boot disk, 8 CPUs and 16 GiB RAM match the prior allocation.
Hostname, address and deterministic MAC derive from the canonical development
DNS declaration. Shared cloud-init installs the Ansible user and guest tools;
Ansible then configures the host-only baseline. Actions remains disabled.

Bootstrap the development OIDC identity and apply its resource-set assignment
through infra/xcp_ng before running `:tf.apply`. The XO session is scoped to
this owner and revoked at exit; no shared administrator token is injected.
For an explicitly authorized recreation, use `:tf.apply -- -replace=xenorchestra_vm.secure1`
and review its interactive replacement plan before confirming.

The existing private-key injector seeds the initial Ed25519 host key. Cloud-init
is treated as creation-only input, so subsequent invocations do not rotate the
VM key. Its private material stays in the Vault-protected Terraform state and
VM bootstrap configuration; temporary injector files are removed at shutdown.
The public `ssh_host_public_key` output establishes strict controller trust
before Ansible signs the permanent host certificate. No SSH checks are disabled.
