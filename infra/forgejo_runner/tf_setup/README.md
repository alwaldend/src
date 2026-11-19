---
title: Terraform setup
description: Forgejo runner VM
tags:
  - terraform
  - pve
---

Use this package's `dns.plan`, `dns.show`, and `dns.apply` targets for the
[scoped DNS adoption workflow](../../dns/README.md).

These targets select `dns=1` and target `module.dns` with the existing
`src_infra_forgejo_runner` AppRole and setup backend. Plan declared imports,
inspect the saved plan with `dns.show`, and apply only the reviewed file with
`dns.apply`. Ordinary setup commands retain their `tf=setup` flow.

`dns_enabled` defaults to `true` after verified adoption. Keep it enabled to
retain existing records; disabling it would propose deletion. The
[adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains the import, follow-up plan, inventory, and DNS verification evidence.
The shared [cutover procedure](../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner)
owns prerequisites and recovery. Scoped DNS checks do not establish VM or runner
service health.
