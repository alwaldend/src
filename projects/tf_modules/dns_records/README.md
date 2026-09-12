---
title: DNS records
description: Canonical DNS declarations for Cloudflare and RouterOS
---

This reusable module translates an owner's decoded `dnsconfig.json` through
its provider-free [normalizer](normalize/README.md), then manages individual
global Cloudflare and dc1 RouterOS records with caller-provided providers.
Owners with only global records use the [global entrypoint](global/README.md),
which needs only Cloudflare. The combined entrypoint delegates Cloudflare
resources to that same module and consumes its normalization for RouterOS.
Repository infrastructure roots may consume its Bazel source filegroups;
the child normalizer also supports provider-free inspection.

```hcl
module "dns" {
  source             = "../../../projects/tf_modules/dns_records"
  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
}
```

`zone` defaults to `alwaldend.com`. An explicit `cloudflare_zone_id` is optional;
when it is null or empty, enabled global records resolve exactly one Cloudflare
zone matching `zone`. The token must permit zone listing and reading. Disabled
ownership and declarations without global records skip this lookup.
`enabled` defaults to false so a prepared root owns
no provider records. Keep it true after adoption: setting it false in a state
that already owns records plans deletion and is not a rollback procedure.
Disabling ownership does not suppress provider initialization. Operational
Terraform commands still require the caller's real provider configuration
and the owning Vault/AppRole prerequisites. In particular, the pinned
RouterOS provider probes its configured API even when no DNS records are
enabled. Use the provider-free normalizer for offline inspection.
The [migration runbook](../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/design.md)
owns cutover, import, and rollback ordering. This module supplies no provider
credentials or state backend.

The input contains exactly one `records` object. Each stable logical key
contains one or more scalar type members and a nonempty `dsp` list. Supported
destinations are `global`, `dc1`, and `all`; repeated destinations are deduplicated.

| Type  | Required member fields       | Default TTL |
| ----- | ---------------------------- | ----------- |
| A     | `name`, `address`            | 300         |
| AAAA  | `name`, `address`            | 600         |
| CNAME | `name`, `target`             | 600         |
| NS    | `name`, `address`            | 300         |
| MX    | `name`, `target`, `priority` | 300         |
| TXT   | `name`, `content`            | 300         |

Every member accepts an integer `ttl` from 60 through 86400 seconds. MX
priorities range from 0 through 65535. Multiple scalar values at the same
name and type use distinct logical keys; changing values never changes
resource identity. Unsupported fields and malformed members are rejected.

Names accept `@`, relative names, or zone-qualified names. Absolute owner
names ending with a dot must be inside the zone. CNAME, MX, and NS targets
accept `@`, relative names, zone-qualified names, or absolute external names
ending with a dot. Normalization lowercases domain names and omits final dots;
address and TXT bytes are preserved. `relative_name` is `@` at the apex.

`normalized_records` maps `logical_key/type/view` to `key`, `name`,
`relative_name`, `type`, `value`, `priority` (null except MX), `ttl`, `view`,
and `proxied` (always false). This output remains available while disabled.
The same keys index `module.global.cloudflare_dns_record.records` and
`module.dc1[0].routeros_ip_dns_record.records` in the combined entrypoint; the global-only
entrypoint uses `cloudflare_dns_record.records`. `import_addresses` exposes
these addresses relative to the selected module, including while staged.
The pinned providers import Cloudflare
records by `zone_id/record_id` and RouterOS records by static IDs such as `*1A`.

Terraform tests use mocked providers and pinned local provider packages.
The global-only test root has no RouterOS provider dependency. Normalization
tests require no providers. Both run without network access;
they establish declaration and resource mapping behavior, not live adoption.
