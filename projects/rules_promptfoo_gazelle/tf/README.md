---
title: Project DNS Terraform
description: Project-owned DNS root with repository operational packaging
---

This root reads the canonical [DNS declarations](../dnsconfig.json) and calls
the shared [DNS records module](../../tf_modules/dns_records/global/README.md).
Its [AL configuration](../al.lua) owns its Vault identity and HTTP state
backend.

Run the root-workspace commands
`//infra/dns/projects:rules_promptfoo_gazelle_tf.<operation>` through `bazel_agent`.
The offline formatting check is
`//infra/dns/projects:rules_promptfoo_gazelle_tf_tests.fmt_test`.
The [operational adapter](../../../infra/dns/projects/README.md) imports this
project's source filegroup; the reusable Bazel module does not depend on the
monorepo.

The operational root requires its deployed AppRole and DNS credential access.
`dns_enabled` defaults to `true` after verified adoption. Keep it enabled to
retain existing records; disabling it would propose deletion. The
[adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains import and verification evidence. Use the shared
[DNS workflow](../../../infra/dns/README.md) for reconciliation and recovery,
including its [access prerequisites](../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner).
