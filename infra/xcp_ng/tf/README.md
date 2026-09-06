---
title: Terraform
description: AppRole resource sets and Vault-backed Xen Orchestra authentication
---

This configuration manages one Xen Orchestra resource set per Vault AppRole
and configures the installed XO OIDC plugin through a packaged JSON-RPC helper.
See the [project documentation](../README.md) for requirements and checks.

Resource-set subjects reference exact synchronized OIDC users whose OIDC
subject identifiers match Vault AppRole entity UUIDs. AppRole groups can include other
service entities, so using those groups would grant access across deployments.
Bootstrap XO's OIDC configuration first, then log in through OIDC as each
AppRole and run the full Terraform apply. Identity discovery runs at plan
time and never creates users. AppRoles that have not logged in are listed in
`approles_pending_oidc_login`; their sets remain without subjects until the
next apply after login. Groups remain responsible for login entitlement and
administrative access, separately from these per-AppRole grants.

`resource_set_inventory` supplies `pool`, `template`, `storage` and `network`
names per AppRole. Native data sources resolve a unique pool and unique
objects within that pool. Unassigned resource sets remain empty. These inputs
have no UUID defaults; resolved IDs remain necessary for XO API operations.

`existing_vms` maps hostnames to an owning AppRole and pool name. Discovery
requires exactly one VM with both the hostname and AppRole tag in that pool.
The owner receives an explicit `admin` ACL. Resource-set membership permits
provisioning and does not transfer existing VM ownership. Missing or ambiguous
VMs, unknown AppRoles and unsynchronized OIDC users fail planning. ACL resource
keys retain the AppRole and resolved VM ID, preserving their addresses when
the same existing VMs are found through names.

Native `xenorchestra_acl.pool_administrator` resources manage the synchronized
administrators group's permissions on connected pools. The read-only helper
identifies existing matching ACLs, and an `import` block adopts them without
recreating access. Fresh installations must bootstrap `terraform_data.oidc`
before the full plan so the synchronized groups exist. There is no shared
pool-view grant and the helper performs no ACL writes.

The pinned `hashicorp/external` provider calls the packaged helper's read-only
`oidc-identities` command. XO exposes login names as `email` and authentication
providers in `authProviders`; the helper accepts only OIDC-authenticated
identities from the configured Vault issuer and indexes them by
`authProviders["oidc:" + issuer].id`, the Vault entity UUID.
Duplicate OIDC subjects are an error. Mutable login names and inherited group
membership do not select the ACL subject. Group discovery selects only
synchronized OIDC groups and rejects duplicate names. The provider has no
group datasource and its user datasource does not expose immutable external
subjects. API credentials remain in the injected environment.
