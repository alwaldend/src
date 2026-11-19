---
title: Tf
description: Terraform config
tags:
  - terraform
---

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
the [AppRole bootstrap procedure](../../../projects/agents/skills/repo-infra/references/vault.md)
when distinguishing missing deployment from insufficient operator access.
