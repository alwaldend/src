---
title: Terraform DNS
description: Project-owned DNS configuration
---

This service stage consumes the [project DNS declaration](../dnsconfig.json)
through the [shared global DNS module](../../tf_modules/dns_records/global/README.md).
AL selects the project's AppRole, HTTP state backend, and Cloudflare provider
injection through the stage's `tf=main` label. This public-only module does not
require RouterOS.

Operational commands require the project AppRole and Vault provider fields in
the [migration prerequisites](../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner).
`dns_enabled` defaults to `true` after verified adoption. Keep it enabled to
retain existing records; disabling it would propose deletion. The zone input
remains optional.

Follow the [DNS migration workflow](../../../infra/dns/README.md) for continued
reconciliation and recovery. The project's adoption change records its import,
no-change plan, and verification evidence.

Offline formatting is available through
`//projects/autoscroll/tf:tf_tests.fmt_test`.
