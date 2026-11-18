## Context

Secure1 has a boot disk and a cloud-init drive, so the cache device name depends
on its actual attachment position. NAS owns storage declarations; this owner
owns the VM, disk attachment and guest mount.

## Decisions

Keep the cache declaration in infra/nas/approles/src_infra_dev/runner_cache.tf.
That AppRole project owns its namespace, ACLs, shares and cache in a standalone
Terraform root with independent state and plan/apply targets. It uses the NAS
operator identity, with no input variables or remote-state reads. The shared NAS
root owns accounts and appliance-wide settings only. The VM root consumes the provisioned XO repository and
owns its disk/metadata only. It declares no TrueNAS provider and requests no
NAS AppRole or administrator credential. The namespace's filesystem ACLs scope
data access; they do not scope TrueNAS API provisioning authority.

Terraform publishes the actual cache disk device and capacity to its own Vault
KV path, alwaldend.com/vault1/approles/src_infra_dev/runner_cache. This record is
a generated projection of the VM resource, not a second configuration source.
The existing AL injector supplies it to secure1 group variables automatically.
Ansible consumes that metadata and performs guest-only checks before
formatting and mounting by UUID; it has no XO credentials, helper or API call.
Reject missing metadata, boot/invalid positions, mismatched capacity, partitions,
incompatible filesystems and mounts used outside the cache path.

## Migration Plan

Apply Vault groups/grants, the shared NAS root, the standalone development NAS
project (namespace and runner cache), XCP-ng storage, the VM root and Ansible.
Each uses its owning state and backend lock. Initial VM attachment precedes
metadata publication as described below. The authorized rollout used a scoped
XCP-ng apply to leave unrelated OIDC and AppRole changes unapplied; deployment.md
records the exact scope and final verification.

## Metadata handoff failure cases

Before implementing the revised handoff, verify that a valid injected record
renders without caller arguments and that missing metadata, invalid JSON, a
boot/invalid device and non-positive capacity fail before guest disk operations.
Terraform must not publish an unattached or invalid disk position. The existing
development AppRole policy owns reads and writes; no broader credential grant
or live Vault write is needed to validate source. Offline playbook syntax must
remain independent of Vault, and user-local fj installation must not fetch it.

The pinned XO provider leaves a new disk's computed position null in its plan
and changes it during apply. Initialize the attachment through the owning VM
apply with -target=xenorchestra_vm.secure1, then run the complete VM apply to
publish its actual device and capacity. A built-in data projection cannot repair
that provider inconsistency. Never guess a guest device or publish null positions;
a repeated complete plan must leave the VM and Vault record unchanged.

The cache module grants XCP-ng inherited MODIFY access. The module owns the
filesystem/share ACL implementation; the AppRole project owns the required parent traversal.

## Provisioning boundary checks

Verify that AppRole storage has no VM, XO or remote-state dependency; its package
includes both namespace and cache while the shared NAS package excludes them.
The selected state path must be independent and use NAS operator authentication.
VM validation must work without TrueNAS; dev:al must contain no NAS authentication.
Imported NAS appliance addresses and the existing VM backend stay unchanged.
