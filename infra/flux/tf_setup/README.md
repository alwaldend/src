---
title: Tf setup
description: Terraform setup
tags:
  - terraform
  - pve
---

This root owns the [Flux DNS declaration](../dnsconfig.json). DNS ownership
defaults to enabled after [verified adoption](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md);
keep it enabled to retain the managed records.

Use `//infra/flux/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply` for
[scoped DNS reconciliation](../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner).
These targets retain the setup backend and `src_infra_flux` AppRole, select
`dns=1`, and target `module.dns`. Apply requires a reviewed saved plan.
Reconciliation of unchanged declarations must produce no DNS changes.

The ordinary setup targets retain their service authentication behavior.
DNS checks do not establish PVE VM or service health; current aggregate source
validation remains part of repository delivery.
