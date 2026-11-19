---
title: Tf setup
description: Terraform setup
tags:
  - terraform
  - pve
---

This root owns Vault's canonical [DNS declaration](../dnsconfig.json) and
retains its existing `src_infra_dc1_vault` AppRole and setup state backend.
Use `//infra/vault/tf_setup:dns.plan`, `dns.show`, and `dns.apply` for the
[scoped DNS workflow](../../dns/README.md). They select `dns=1` and
`module.dns`; apply requires a reviewed saved plan. Ordinary setup and service
wrappers retain their existing authentication flow.

`dns_enabled` defaults to `true` after verified adoption. Keep it enabled to
retain existing records; disabling it would propose deletion. The
[adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains the import and verification evidence. Targeted DNS checks do not
establish the health of other resources in the root.

Offline formatting is available through
`//infra/vault/tf_setup:tf_tests.fmt_test`.
