---
title: Tf
description: Terraform config
tags:
  - terraform
---

Use this package's `dns.plan`, `dns.show`, and `dns.apply` targets for the
[scoped DNS adoption workflow](../../dns/README.md).
They select `dns=1` and target `module.dns` in this root, retaining the
`src_infra_dc1_pve1` AppRole and existing backend. Use the
[owner preparation and import procedure](../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner)
for declarative import maps and reviewed saved-plan apply.

`dns_enabled` defaults to `true` after the
[completed adoption](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md).
Reconciliation against adopted state and unchanged declarations produces no DNS
changes. DNS targeting covers DNS and its dependencies; it does not establish
PVE host, API, or service health. Ordinary Terraform commands retain their
existing Proxmox authentication.
