# DNS migration provider feasibility

Observed 2026-09-12, checked-in source plus version-pinned public upstream source. No live provider, router, Vault, backend, import, plan, apply, or DNSControl push was contacted. Local task branch: `t3code/dns-migration-add-missing-approles`. This is capability evidence, not a provider inventory or adoption receipt.

## Verdict

Proceed with existing pins: Cloudflare `5.22.0` (`infra/dns/tf/provider.tf`) and RouterOS `1.99.1` (`infra/ingress/tf/provider.tf`). Both expose individual record resources, including every current type and multiple rows for the same owner/type. No additional provider product is necessary. Existing provider locks can be reused. Offline schema compatibility is supported; real IDs and imported-state parity remain unknown until separately authorized inventory/adoption.

## Resource mapping

| Canonical type | Cloudflare `cloudflare_dns_record`        | RouterOS `routeros_ip_dns_record`                                       |
| -------------- | ----------------------------------------- | ----------------------------------------------------------------------- |
| A / AAAA       | `content` = normalized IP                 | `address` = normalized IP                                               |
| CNAME          | `content` = absolute target               | `cname` = absolute target, no terminal dot                              |
| NS             | `content` = absolute target               | `ns` = absolute target, no terminal dot                                 |
| MX             | `content` = absolute exchange; `priority` | `mx_exchange` = absolute exchange without terminal dot; `mx_preference` |
| TXT            | `content` = text                          | `text` = joined text chunks, no presentation quotes                     |

Cloudflare uses `zone_id`, `name`, `type`, numeric `ttl`, scalar content and `proxied=false` for the current declarations. TTL `1` means automatic; ordinary TTL is 60–86400 seconds (30 minimum for Enterprise). Import each record as `<zone_id>/<dns_record_id>`. Cloudflare supports more types than the migration needs; do not silently accept types absent from the common contract. [Pinned Cloudflare schema and import](https://github.com/cloudflare/terraform-provider-cloudflare/blob/v5.22.0/docs/resources/dns_record.md).

RouterOS uses a flat `/ip/dns/static` row with `.id`. All current types are supported; field values are scalar, so apex multi-value A/AAAA/MX/TXT and delegated NS entries need one Terraform instance per provider row. `type` forces replacement on change. `ttl` is a duration string with semantic time comparison (`300s` equals `5m`). `name` and `regexp` are mutually exclusive. `match_subdomain` is optional Boolean. [Pinned resource implementation](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/routeros/resource_ip_dns_record.go).

Import RouterOS rows by their exact `.id`, e.g. `*A`. Attribute-based `name=...` imports are documented but are ambiguous for multi-value names and must not be used for those names. [Pinned import documentation](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/docs/resources/ip_dns_record.md).

## Effective TTL and normalization parity

The pin is DNSControl `4.36.1` (`tools/dnscontrol/binary_toolchain.json`). The baseline `infra/dns/dnsconfig.js` at revision `cb4b5bd3a37d29dc778ed140ab72d05efef7816a`, historical BIND snapshots, and `infra/mikrotik/router1.rsc` agree on these source semantics:

| Input                                     | Effective TTL seconds |
| ----------------------------------------- | --------------------- |
| Project JSON A                            | 300                   |
| Project JSON NS                           | 300                   |
| Project JSON AAAA / CNAME                 | 600                   |
| Shared apex A / AAAA                      | 300                   |
| Other common apex, www, Protonmail, DMARC | 300                   |
| Simplelogin records                       | 10800                 |
| Yandex MX                                 | 21600                 |
| Other Yandex records                      | 300                   |

`mods.push(A(...), TTL(600))` places TTL at domain scope; it sets the domain object's `ttl`, not `defaultTTL` and not the A record's TTL. The record factory copies `d.defaultTTL`, which is initially zero; AAAA/CNAME instead pass TTL into their record factories. [Pinned JavaScript implementation](https://github.com/StackExchange/dnscontrol/blob/v4.36.1/pkg/js/helpers.js). Zero record TTL becomes the 300-second model default during normalization. [Normalization](https://github.com/StackExchange/dnscontrol/blob/v4.36.1/pkg/normalize/validate.go), [default constant](https://github.com/StackExchange/dnscontrol/blob/v4.36.1/models/dns.go).

Do not flatten identity to only name/type/view: several apex records share those fields. Preserve the logical key/type/view in state identity, including a stable value slot if one logical entry has an array. Updating an IP must not replace its key. NS delegation has two values. Shared apex has four A values and four AAAA values in each view.

DNSControl's MikroTik conversion uses absolute names without terminal dots, canonical IP strings, strips terminal dots from CNAME/NS/MX targets, joins TXT chunks, and converts numeric TTL to duration strings. It explicitly writes `match-subdomain=no`, empty `regexp`, empty `address-list`, and empty `comment` when metadata is absent. Adding Terraform ownership comments would violate a strict unchanged adoption plan. [Pinned conversion](https://github.com/StackExchange/dnscontrol/blob/v4.36.1/providers/mikrotik/convert.go).

Legacy reads skip disabled and dynamic rows. Import only enabled, non-dynamic matching rows; do not adopt unrelated disabled entries that share the same name. [Pinned reader](https://github.com/StackExchange/dnscontrol/blob/v4.36.1/providers/mikrotik/mikrotikProvider.go). RouterOS Terraform `comment` is optional, without a preserve-on-omission diff suppressor. `disabled` has an omission-based diff suppressor; explicitly `disabled=false` and `match_subdomain=false` make the intended active exact-name semantics clear. `dynamic` is computed/read-only. [Pinned shared property schemas](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/routeros/provider_schema_helpers.go).

## Management dependency

The pinned Cloudflare `cloudflare_zone` data source supports a zone-name filter,
rejects any result count other than one, and reads the resulting zone. The
module uses it only for enabled global records without an explicit zone ID.
The List Zones endpoint requires `Zone Zone Read`; a record-edit token alone
does not establish that capability. This removes the need to add a zone ID to
the existing credential entry. [Pinned data source](https://github.com/cloudflare/terraform-provider-cloudflare/blob/v5.22.0/internal/services/zone/data_source.go),
[List Zones permissions](https://developers.cloudflare.com/api/resources/zones/methods/list/).

The baseline `infra/dns/providers.json` and ingress AL configuration address `https://router1.dc1.alwaldend.com`; the MikroTik declaration owns that same name, whose A target is `192.168.1.1`. The router snapshot assigns certificate `alwaldend.com_acme` to HTTPS/API TLS, but source inspection does not establish its SANs or expiry. Merely replacing the hostname with that IP cannot establish successful TLS verification.

The shared infra AL configuration now owns that existing endpoint value for
ingress and DNS consumers. DNS still reads only its existing username and
password from Vault; no `management_url` field or credential rewrite is needed.

The pinned RouterOS provider accepts HTTPS or TLS API endpoints, optional CA certificate, and verified TLS by default (`insecure=false`). It has no separate TLS server-name field in its schema. [Provider source](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/routeros/provider.go), [provider documentation](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/docs/index.md).

Before router-name ownership transfer, owning IaC must provide an independently reachable endpoint with a valid matching TLS identity (IP SAN or independently resolved name); independently verified Vault/backend connectivity is required for their scopes too. Do not infer that the existing certificate covers the IP, disable certificate verification, or treat source-only checks as an authenticated recovery test. This blocks live transfer of those scopes, not offline module/root preparation.
