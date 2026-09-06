---
title: Tf
description: Terraform config
tags:
  - terraform
---

This package configures Forgejo organizations, repositories, and access. Before
applying, provision the Vault OIDC authentication source with Ansible and set
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

Organization administrators come from `src_infra_dc1_forgejo1_admins` and its
direct child groups. Package writers and the repository automation writer
come from the corresponding `src_infra_dc1_forgejo1_package_writers` and
`src_infra_dc1_forgejo1_src_writers` Vault groups. These dedicated groups avoid
granting service permissions to AppRole credential issuers. The repository's
existing singleton collaborator requires exactly one discovered writer.
Access-group members outside the login population also fail validation.

```sh
bazel_agent bazel run //infra/forgejo/tf:tf.plan
bazel_agent bazel run //infra/forgejo/tf:tf.apply
```
