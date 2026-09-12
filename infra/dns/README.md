---
title: Dns
description: Project-owned DNS declarations for alwaldend.com
tags:
  - dns
  - terraform
---

DNS records live in each owner's `dnsconfig.json`. Each owner manages its records
through the reusable [Terraform DNS module](../../projects/tf_modules/dns_records/)
from its designated `tf_setup` or `tf` root. This component's
[dnsconfig.json](dnsconfig.json) owns shared apex and mail records.

## Inspect and lint declarations

```sh
bazel_agent bazel run //infra/dns:lint
```

The linter discovers `dnsconfig.json` files in the current workspace at runtime
and prints their DNS declarations as a table. It includes nested project
workspaces and excludes task scratch, tool caches, and generated directories.
There is no central source registry to maintain when a project adds a file.

A canonical domain name must belong to one source file. Multiple values, record
types, and destination views within that file are valid. Duplicate JSON object
keys and domains declared by different source files are errors. The command
reads checked-in declarations without contacting Terraform backends or DNS
providers.

`bazel_agent bazel test //infra/dns:config_test` runs the linter's fixture tests
and checks every DNS declaration in the current checkout. The repository test
resolves the root `MODULE.bazel` runfile to the checkout and runs without sandboxing
or cached test results so newly added files and nested workspaces are included
on every run. It only reads local files.

## Terraform records

The [provider-free normalization module](../../projects/tf_modules/dns_records/normalize/)
owns the Terraform input schema, relative-name handling, TTL defaults, and stable
resource keys. Each logical key can contain several record-type members; `dsp`
selects `global`, `dc1`, or `all`. Explicit member TTLs preserve shared apex and
mail settings. The normalization module remains an internal dependency of each
owner's Terraform DNS module.

Owner-local roots receive their provider credentials and Vault HTTP state
backend through the existing AL configuration. The shared module resolves the
Cloudflare zone by name when no optional zone ID is provided. DNS and ingress
reuse the RouterOS endpoint declared in the [shared AL configuration](../al_lib.lua),
while their separate credentials remain in Vault. Neither integration requires
adding metadata fields to the existing credential entries.
Validate their source with the
owning package's offline checks. A live plan, import, apply, or state operation
requires separate authorization for that owner and operation.

Service roots expose `dns.plan`, `dns.show`, and `dns.apply` alongside their
ordinary Terraform commands. They retain the owning root, AppRole and backend,
select the DNS AL calls, and plan `module.dns`. Required root inputs still load
through the existing read-only Vault injectors. Unrelated service authentication
is not selected. Apply accepts only a reviewed saved plan file and rejects
Terraform CLI argument environment overrides. See the runbook for the optional
import-ID maps and adoption checks.

## Migration and recovery

The [migration change](openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/) and
[cutover runbook](openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
describe adoption, imports, recovery, and retirement requirements. Each owner's
adoption change records its verified live transfer; source implementation alone
does not establish operational ownership. The completed coordination record
links the operational acceptance evidence and its coverage limits.

Current source no longer exposes the central DNSControl deployment commands.
Before an owner begins Terraform reconciliation, stop central deployment jobs
identified by the active and scheduled writer audit. Record coverage and
unavailable observations, and coordinate one owner at a time. Adopt existing
records into the owner's state with an exact-ID, no-change plan. Fresh complete
inventories may establish missing declarations for a separately reviewed
additions-only plan that preserves all existing records.
Rollback requires stopping the affected Terraform writers and following the
runbook's recorded prior revision and state-reconciliation procedure.

Historical `zones/*.zone` files are deployment snapshots. They are not a current
inventory or a declared-state export and must not be hand-edited.
