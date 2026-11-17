## Decisions

The shared NAS root reads Vault's truenas_admins and truenas_user entity groups
and creates one account per distinct member. NAS groups mirror those memberships;
only truenas_admins grants FULL_ADMIN. Account names longer than 32 characters
use the existing stable shortening rule. No account receives shell or sudo rights.
Generate 64-character passwords, coordinate generation and write-only updates
with one rotation version, and publish credentials through each AppRole's
existing shared-read path. Keep account publication independent of datasets.

Each infra/nas/approles/<approle> directory is a standalone Terraform root with
its own HTTP backend and plan/apply targets. It owns its namespace, permissions,
shares and child datasets. The development root owns namespace traversal for
XCP-ng and the runner cache; the shared root has no namespace loop or storage
module call. The reusable dataset module retains filesystem/share defaults.
Namespace creation is part of that AppRole project's apply, not account creation.

Use the existing NAS operator identity and API credential in both NAS roots.
The shared root uses its short-lived injected AppRole token directly; it does
not request child-token creation rights.
Store AppRole state under src_infra_nas/tf_backend/approles/<approle>. Reuse one
connection injector with a truenas selection label and explicitly extend the
core backend call for the AppRole backend. No development AppRole grant or VM
provider is added. Separate state isolates lifecycle, not TrueNAS API authority;
its method-based RBAC cannot restrict administrator calls to dataset paths.

## Decision review

Proceed with separate states and targets, as explicitly selected by the user.
Keeping modules in one state would still couple applies. Giving the development
identity a NAS administrator credential would defeat the accepted credential
boundary. The repository's independent Vault AppRole root supplies the existing
backend/configuration pattern; no dependency or new framework is needed.

## Bootstrap and ownership transition

Apply Vault groups/grants, then the shared NAS root (accounts, shared encrypted
parent and SMB transport), then the development NAS root (namespace and cache),
then XCP-ng, VM provisioning and Ansible. Unlock the encrypted parent first.
Every stage uses its own backend lock. No remote state or parent input is
needed. The authorized rollout scoped XCP-ng to the development resource set to
leave unrelated OIDC/AppRole changes unapplied, and initialized the VM disk
before publishing metadata because the pinned provider plans a null position.

Keep the adopted pool, seven datasets and app addresses unchanged. The
authorized rollout is recorded in validation.md. If an earlier
revision was applied elsewhere, its namespaces must be adopted into the matching
AppRole states before removing their old ownership. In particular, the previous
module.approle_namespaces["src_infra_dev"] resources belong under module.namespace
in the development state, and module.runner_storage.module.runner_cache resources
belong under module.runner_cache. Reconcile existing namespaces for other roles
into their own projects too. Never apply a proposed removal as storage deletion;
live reconciliation requires its own authorization and reviewed state inventory.

## Failure cases to verify before implementation

- Both roots validate independently using the pinned providers; the AppRole root
  packages its namespace and child module without the shared root's files.
- The shared root cannot plan or initialize the development storage module.
- Each launcher selects exactly one backend: unchanged core state or the unique
  AppRole state path. Both select NAS authentication and the shared connection;
  VM provisioning selects neither NAS authentication nor its API injector.
- Namespace and child names, encrypted inheritance, private SMB ACLs, XCP-ng
  traversal/MODIFY grants and cache quota/copies remain coherent.
- The AppRole cache parent uses its namespace output to order creation. Shared
  account publication has no dependency on storage projects.
- Plan/apply launcher manifests and materialized runfiles match current source.
  Validation, specification checks and semantic lint must not contact live NAS.
