---
title: Tf
description: Terraform config
tags:
  - terraform
---

## Public certificate snapshots

This root manages Vault's PKI and SSH issuers. Public snapshots are owned and
refreshed independently by [data/ssl/alwaldend.com/tf](../../../../data/ssl/alwaldend.com/tf/README.md)
and [data/ssh/tf](../../../../data/ssh/tf/README.md); this root does not write them.

The `removed` blocks in `outputs.tf` retire the four former `local_file`
exports with `destroy = false`. A later authorized apply forgets those state
entries without deleting their files. No state migration or Vault operation is
performed merely by moving the checked-in snapshots.

## DNS identities

[Missing DNS-owner AppRoles](approles_dns.tf) compose the existing reusable
AppRole module with an owner-specific DNS policy in `approles/<name>/`.
[DNS access](dns_access/main.tf) grants read-only access to the existing
provider credentials for the owner's views. Existing component identities
receive these policies through their existing module declarations.

[DNS-only group membership](group_dns_approles.tf) adds no policies. These
identities retain their own state and named shared-secret access without
joining the general infrastructure AppRole group or receiving cloud
provisioning, SSH, or PKI permissions.

The operator bootstrap membership follows the existing component pattern.
An authorized Vault apply must create the roles and policies before their DNS
stages can authenticate; source validation does not establish deployed access.
Follow the [DNS migration](../../dns/README.md) for the coordinated cutover and
the [AppRole bootstrap procedure](../../../tools/agents/skills/repo-infra/references/vault.md)
when distinguishing missing deployment from insufficient operator access.

## Vast team token delivery

[The Vast infrastructure role](approle_src_infra_vastai.tf) declares access to
only `src_infra_image_gen/vastai` in addition to its existing namespace. The
policy grants data lifecycle operations and exact metadata read/delete access,
with no sibling or component-state access. It lets the team root deliver and
later remove the generated key using its own identity.

[The image-generation role](approle_src_infra_image_gen.tf) owns its own KV and
state subtree. Its operator member can bootstrap the AppRole and browse that
namespace in Vault. List-only permissions on exact metadata ancestors support
UI navigation without access to other secret values. An additional exact-path
metadata read/delete policy lets the component manage and remove all versions
of its web-login entry. Its explicit read permission retains metadata access
when this exact-path policy takes precedence over the namespace wildcard. It does not join the broad infrastructure AppRole group.

The legacy personal image-generation grant was removed through a reviewed
scoped apply on 2026-11-10. The replacement role and delivery grant were then
deployed under the user's credential-only authorization; no rental was created.
The [component receipt](../../image_gen/openspec/changes/add-imagegen-rental/design.md#credential-only-deployment)
records the operation. The existing personal user role and its browser access are independent
and remain unchanged.
