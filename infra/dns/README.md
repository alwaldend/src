---
title: Dns
description: Project-owned DNS declarations for alwaldend.com
tags:
  - dns
  - terraform
---

DNS records live in each owner's `dns/dnsconfig.json`. Each owner manages its records
through the reusable [Terraform DNS module](../../projects/tf_modules/dns_records/)
from a dedicated DNS root with independent state and DNS-only credentials.
DNS-only roots moved from `tf/` to `dns/` keep their existing backends.
This component's
[dnsconfig.json](dns/dnsconfig.json) owns shared apex and mail records.

The public apex follows `download.alwaldend.com` to the Yandex download host.
The dc1 apex continues to use the local download VM. The
[activation record](openspec/changes/activate-public-download-site/evidence.md)
documents the authorized cutover and HTTPS acceptance.

## Architecture

This diagram describes the maintained architecture source, not live inventory or service health.
The [Mermaid source](docs/architecture.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![DNS architecture](docs/assets/architecture.svg)](docs/assets/architecture.svg)

## Declaration pages

[Cloudflare DNS](cloudflare_dns.md) and [Mikrotik DNS](mikrotik_dns.md) list the
records this repository declares for each destination view, with the owning
declaration for every record. They are manually refreshed snapshots of
`dnsconfig.json` and may lag behind the declarations. Regenerate them explicitly
when a documentation refresh is wanted, for example after parallel DNS changes
have merged:

```sh
bazel run //infra/dns/cmd/dump -- --write
```

Ordinary validation checks declarations and ownership without requiring or
rewriting these snapshots. A DNS change does not require a documentation
refresh in the same PR. To explicitly check snapshot freshness, run
`bazel run //infra/dns/cmd/dump -- --check`.

## Inspect and lint declarations

```sh
bazel run //infra/dns:lint
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

`bazel test //infra/dns:config_test` runs the linter's fixture tests
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
backend through the existing AL configuration. Standalone roots can use the
[`infra_terraform_dns` macro](terraform/README.md) to package their Terraform
configuration and Bazel commands. The shared module resolves the
Cloudflare zone by name when no optional zone ID is provided. DNS and ingress
reuse the RouterOS endpoint declared in the [shared AL configuration](../al_config.bzl),
while their separate credentials remain in Vault. Neither integration requires
adding metadata fields to the existing credential entries.
Validate their source with the
owning package's offline checks. A live plan, import, apply, or state operation
requires separate authorization for that owner and operation.

Run `tf.plan` and interactive `tf.apply` from the owner's DNS package.
Service commands retain their own AL settings, providers, and backends.

## Transfer from service state

The [consolidation change](openspec/changes/consolidate-terraform-dns/proposal.md)
introduces independent roots for the owners below. Each new root uses the
same AppRole and a distinct `tf_backend/dns` path under that AppRole.
Freeze both writers before using the new revision. Complete the state transfer
before running ordinary plan or apply from this revision. A DNS plan against a
fresh backend can propose creation. A service plan against its old state can
fail or propose removal. Use the prior revision to inspect and relinquish the
old bindings.

| Owner         | Previous root                 |
| ------------- | ----------------------------- |
| Flux          | `//infra/flux/tf_setup`       |
| Forgejo       | `//infra/forgejo/tf_setup`    |
| Harbor        | `//infra/harbor/tf_setup`     |
| Ingress       | `//infra/ingress/tf`          |
| OpenHands     | `//infra/openhands/tf_setup`  |
| Proxmox       | `//infra/pve/tf`              |
| 3x-ui         | `//infra/threexui/tf_setup`   |
| Vault         | `//infra/vault/tf_setup`      |
| Xen Orchestra | `//infra/xcp_ng/tf`           |
| alwaldend.com | `//projects/alwaldend.com/tf` |

1. Obtain explicit authorization for the source state removal and the
   destination imports for one owner. Freeze all writers to both backends,
   record the prior deployed revision, and verify an independent recovery
   path for infrastructure endpoints.
2. Inspect both states using their owning wrappers. Verify the destination
   is absent or empty, and check each managed source DNS address and exact
   provider ID against the declaration and provider inventory. Cloudflare
   imports use `zone_id/record_id`. RouterOS imports use the static row ID.
   Keep that inventory private and preserve a recoverable KV version of the
   original state.
3. From the prior revision, use the previous root's state target to review
   `rm -dry-run module.dns`. Require every address to match the reviewed DNS
   inventory or a DNS data source under the same module. Run the separately
   authorized state removal and verify the scope. The previous revision retains
   the provider configuration needed to inspect the original bindings.
4. In the new `dns/` root, import each reviewed provider ID at the same
   Terraform address with `tf.import -- <address> <id>`. If an import fails,
   inspect both states before resuming and import only the missing bindings.
5. Review a fresh DNS `tf.plan` and a service plan
   against the new revision. Require zero DNS additions, changes, replacements,
   or deletions in both. Verify provider IDs and authoritative DNS answers,
   then release the writer freeze.

If a transfer cannot finish, keep both writers frozen. Reconcile back to a
single owner using the reviewed state inventory and the prior revision before
resuming normal operations.

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

Deployment snapshots are not maintained here. The declaration pages above
reflect their last manual refresh; each owner's declarations remain the
authoritative inventory.
