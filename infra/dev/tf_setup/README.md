# Secure1 provisioning

This root owns secure1's Xen Orchestra VM and uses the development AppRole's
own resource set and Vault backend at `src_infra_dev/tf_backend/tf_setup`.
Its 500 GiB boot disk, 8 CPUs and 16 GiB RAM match the prior allocation. A separate
800-GiB cache disk uses the `Runner cache` NAS-backed storage repository.
The [NAS AppRole project](../../nas/approles/src_infra_dev/README.md) owns the
namespace and cache in separate state using NAS operator authentication.
This VM root has no TrueNAS provider or NAS credentials;
it consumes the provisioned XO repository and owns the disk attachment.
Appending the disk uses the pinned provider's in-place update support.
Hostname, address and deterministic MAC derive from the canonical development
DNS declaration. Shared cloud-init installs the Ansible user and guest tools;
Ansible then configures the runner and its UUID-mounted cache filesystem.

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

## Initial cache provisioning

After separate authorization, apply Vault groups/grants, then the complete
`//infra/nas/truenas:tf.apply` root for accounts and shared settings. Apply
`//infra/nas/approles/src_infra_dev:tf.apply` for its namespace and runner cache
dataset/share/ACLs. Then apply `//infra/xcp_ng/tf:tf.apply`
to create and delegate the storage repository, and `//infra/dev/tf_setup:tf.apply`
to attach the cache disk. Each stage uses a complete apply in its own state.

The complete VM apply publishes the attached cache device and capacity to
`secrets/alwaldend.com/vault1/approles/src_infra_dev/runner_cache`, using the
existing development AppRole permissions. Terraform owns this generated record;
do not edit its values manually. The normal Ansible wrapper loads it from Vault,
so deployment requires no Terraform output command or extra-vars handoff.
See the [playbook guide](../ansible/README.md). No live apply ran for this change.
