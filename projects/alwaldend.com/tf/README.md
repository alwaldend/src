---
title: Tf
description: Terraform
tags:
  - terraform
  - pve
---

This service stage consumes the [project DNS declaration](../dnsconfig.json)
through the [shared global DNS module](../../tf_modules/dns_records/global/README.md).
The DNS workflow uses the project's AppRole, HTTP state backend, and Cloudflare
credentials. This public-only module does not require RouterOS. Ordinary
service commands retain their `tf=main` authentication, including PVE login.

Operational commands require the project AppRole and Vault provider fields in
the [migration prerequisites](../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner).
`dns_enabled` defaults to `true` after verified adoption. Keep it enabled to
retain existing records; disabling it would propose deletion. The zone input
remains optional.

The `dns.plan`, `dns.show`, and `dns.apply` targets select only the DNS
workflow in this same root and backend. They use `dns=1` to load DNS credentials
without starting PVE login. Apply requires a reviewed saved plan file. Targeted
DNS planning does not establish PVE or other service health.

Follow the [DNS migration workflow](../../../infra/dns/README.md) for continued
reconciliation and recovery. The
[adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains the earlier blocked attempt and completed scoped DNS verification.

Offline formatting is available through
`//projects/alwaldend.com/tf:tf_tests.fmt_test`.
