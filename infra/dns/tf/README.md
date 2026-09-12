---
title: Shared DNS Terraform
description: Owner state for common apex and mail records
tags:
  - terraform
---

This root consumes [the canonical common declarations](../dnsconfig.json)
through the reusable DNS module. It uses the `src_infra_dns` AppRole and its
existing Vault-backed HTTP state. Both Cloudflare and RouterOS credentials
are injected through the `tf=1` stage label.

Shared apex and mail records were adopted on 2026-09-13. `dns_enabled` now
defaults to `true`; reconciliation against the adopted state and unchanged
declarations must propose no record changes. The [adoption evidence](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
records 60 imports and preservation of both complete provider inventories.

Operational calls require the AppRole policy and Vault credential fields
described by [the cutover procedure](../openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md).
The module's `import_addresses` output gives addresses relative to `module.dns`;
actual provider IDs come from the authorized adoption inventory.

Build the wrappers or run `:tf_tests.fmt_test` for offline validation. Import,
plan and apply follow the separately authorized adoption procedure.
