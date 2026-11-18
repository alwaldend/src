---
title: Tf
description: Terraform config
tags:
  - terraform
---

This package configures Forgejo organizations, repositories, and access from
the shared [repository catalog](../../repos/README.md). The catalog owns named
administrators, developers, repository identities, and destination names.
Its [Terraform module](../../repos/tf/README.md) supplies the consumer's
repository projection. First-party names remain stable across forges;
external fork and mirror names retain their original upstream identity.

Repository resources preselect deletion of eligible pull request source
branches in Forgejo's web merge form. The pinned provider sends the setting
on repository creation and update and does not detect out-of-band changes to
this field during refresh.

Before applying, provision the Vault OIDC authentication source with Ansible and set
`TF_VAR_vault_oauth_source_id` to its verified positive numeric ID. Obtain the
ID from `forgejo admin auth list` on the instance using the service's config
and work path; confirm that the `vault` source is active and uses OpenID
Connect with the expected Vault issuer. Do not assume source IDs survive
instance recreation.

Terraform discovers users from the authoritative Vault
`src_infra_dc1_forgejo1_users` login group, then looks up their entities,
following the PVE Terraform pattern. Discovery includes direct members and
two levels of nested groups, covering the current administrator membership.
Deeper nesting fails validation rather than silently omitting accounts.
Each account uses its Vault entity UUID as the OIDC login name and its
canonical username/email metadata. No local password or user impersonation is
used. The Forgejo AppRole receives the existing group/entity lookup permissions,
as PVE does. Apply those Vault group memberships before running this package.

The accounts are protected from Terraform deletion. OIDC group claims retain
ownership of administrator/restricted status after login. If a user already
exists outside this state, verify its source and login-name mapping before
importing it; never replace an existing account to resolve a name collision.
Resource addresses remain keyed by Vault entity name, preserving existing
managed users. The Forgejo service identity may already exist from the first
OIDC login: verify and import it before applying its newly discovered resource.

Named organization administrator and developer assignments come from the
catalog and must resolve to the discovered Vault login population. The
Developer team can read repositories, write feature branches, and open pull
requests. Its members receive no default-branch push or merge bypass.

Existing service administration remains sourced from
`src_infra_dc1_forgejo1_admins` and its direct child groups. Those grants are
retained for entities with a Vault AppRole alias, alongside catalog
administrators. The separate instance-wide OIDC administrator claim remains
owned by Vault. Validation rejects a catalog
developer who also receives administrator access through those Vault groups.
Package writers and the repository automation writer retain the grants from
`src_infra_dc1_forgejo1_package_writers` and
`src_infra_dc1_forgejo1_src_writers`. These dedicated groups avoid granting
service permissions to AppRole credential issuers. The existing singleton
collaborator still requires exactly one discovered writer. Access-group
members outside the login population also fail validation.

Repository adoption preserves remote IDs. The existing singleton repository
address has an explicit move to its catalog key, while account and service
access addresses remain stable. Import an existing untracked destination
before applying; never recreate it to resolve a name collision. Review the
full plan before execution and obtain separate approval for any deletion or
replacement. The retained branch rules continue to own administrator and
automation-writer exceptions.

```sh
bazel run //infra/forgejo/tf:tf.plan
bazel run //infra/forgejo/tf:tf.apply
```

## Runner credentials

`terraform_data.runners["secure1"]` bootstraps the repository-scoped
`secure1.dev.alwaldend.com` runner for `alwaldend/src` through the packaged
[local tool](../cmd/runner_registration/README.md). Its credentials are cached
in KV v2 mount `secrets` at
`alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/runners/secure1`, within
the Forgejo AppRole's existing namespace. The worker receives no controller
API credential. This declaration creates credentials; host configuration and
consumer read access remain owned by their deployment and Vault policies.

The provisioner runs only when its stable Terraform resource is created. It
has no timestamp, executable hash or replacement trigger. Normal applies and
tool changes do not rerun registration. An explicit replacement reuses the
Vault cache. The token is written directly to Vault and is never a Terraform
attribute. Removing the resource preserves the runner and cached credentials.
The local tool's README owns its failure and reconciliation contract.

Runner registration consumes an existing repository and uses its identity from
the shared catalog. Provision the repository first on a fresh service. A
targeted runner bootstrap does not reconcile the service's repository resources.
