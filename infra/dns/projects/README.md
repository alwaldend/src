---
title: Nested project DNS commands
description: Repository operational wrappers for project-owned Terraform roots
---

This package exposes independent Terraform commands for DNS roots owned by
nested Bazel projects. Each project's `tf/` directory owns its configuration,
provider lockfile, and state boundary; its `al.lua` selects its own Vault
AppRole and HTTP backend. Its existing `dnsconfig.json` remains the record
source.

For example, run `//infra/dns/projects:rules_docs_tf.plan` from the monorepo
workspace when a live plan is authorized. The offline formatting check is
`//infra/dns/projects:rules_docs_tf_tests.fmt_test`. Other owners use the
same target naming with their module name.

The nested modules export source filegroups. This adapter aliases those
declared files into runfiles at their checkout paths and declares the one
shared `//projects/tf_modules/dns_records/global:global` module dependency.
Terraform's relative module and JSON paths therefore resolve in both the
checkout and Bazel runfiles. The adapter does not copy source or introduce
a dependency from reusable Bazel modules back to the monorepo.

Each operational command selects the ordinary `tf=main` pipeline, including
DNS credential injection. Deploy the owner's AppRole and credential access
before running its operational root. The initial `dns_enabled = false`
suppresses record creation while preparing adoption.

Stop central deployments and prevent older central candidates from running
before importing existing records or enabling owner reconciliation. Follow
the [DNS migration runbook](../README.md) for adoption and recovery.
