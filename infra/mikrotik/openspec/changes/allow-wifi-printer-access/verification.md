# Verification

Observed on 2026-11-07. The scoped Vault grants, both printer DNS views, and
the Wi-Fi firewall exception are deployed. The user confirmed the printer page
opens from a device connected to Wi-Fi.

## Current source

The printer record selects both dc1 and global DNS views. Cloudflare injection
is enabled and the existing DNS access policy gains its named provider-secret
read. The global A record remains unproxied and points to the private printer
address; public DNS resolution does not provide network connectivity.

The firewall lookup is inline with its resource consumer. Its exact-one-result
precondition remains. The module and mock suite introduced solely to work around
Terraform's nested-response mocking have been removed. Offline checks cover
formatting and the pinned provider schema, not live query cardinality failures.

The destination is a single IPv4 host. RouterOS returns its bare address rather
than an explicit `/32`; source uses that canonical representation to prevent
repeated cosmetic updates (`routeros-host-address-normalization`). The change
does not broaden the destination match.

The source checks comprise 35 tests covering Terraform formatting and schema
validation, the DNS suite, and repository quality, plus affected semantic lint
and strict OpenSpec validation. Exact-candidate results, local review, and
publication are tracked by the AL receipts under `out/printer-access/`.
Deployed head `7e17ea8` passed automatic review with no findings; that result
does not cover the later canonical-address correction and deployment receipt.

## Deployment receipt

The user approved local DNS, then authorized the remaining operations with
“Apply everything.” Each owning Terraform apply displayed a normal text plan
and received interactive confirmation only after review. No resource was deleted.

| UTC completion | Owning apply target and scope                                    | Applied result                                                                                                        |
| -------------- | ---------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| 16:28:55       | `//infra/mikrotik/tf:tf.apply`                                   | One local DNS record, RouterOS ID `*8E`.                                                                              |
| 16:57:36       | `//infra/vault/tf:tf.apply -- -target=module.src_infra_mikrotik` | One exact-secret firewall policy; two in-place updates for AppRole policies and the named global DNS credential read. |
| 16:58:14       | `//infra/mikrotik/tf:tf.apply`                                   | One unproxied global A record, TTL 300; existing local record unchanged.                                              |
| 16:58:51       | `//infra/mikrotik/firewall:tf.apply`                             | One forwarding exception, RouterOS ID `*29`, inserted before drop `*16`.                                              |

Post-apply owning plans report **No changes**: firewall at 17:01:41 UTC after
the canonical-address correction, DNS at 17:02:13 UTC, and targeted Vault at
17:05:34 UTC. The targeted Vault result covers the MikroTik module and its
dependencies, not unrelated Vault resources.

## Reachability and rule verification

Cordis DNS lookups through public resolvers `1.1.1.1` and `8.8.8.8`, and local
resolvers `192.168.1.1` and `192.168.2.1`, all returned
`printer1.dc1.alwaldend.com A = 192.168.1.221`. Public resolution supplies the
private address; network access still requires a route to the site.

HTTP GET by name followed the printer's redirect to `/general/status.html`
and returned HTTP 200. TCP port 631 also accepted a connection. These probes
came from the wired agent host; response bodies were not retained.

A read-only RouterOS API check through the owning injector at 17:07:04 UTC
confirmed exactly one printer rule `*29` with the intended enabled accept
action, wireless source `192.168.2.0/24` on `bridge2`, wired output `bridge1`,
and destination `192.168.1.221`. It precedes the enabled forward drop `*16`.
Only the verification summary was printed, with no credential values.

The user also confirmed that opening `http://printer1.dc1.alwaldend.com/`
from a Wi-Fi device displays the printer page. No print job was submitted or
verified; this task establishes DNS and device connectivity. Routed access
does not provide cross-subnet mDNS discovery.

## Session ergonomics

A production module existed solely to support provider mocks. It was removed
after the user questioned the boundary; production structure now follows its
single consumer. The precondition remains and the reduced offline coverage is
explicit. Other resolved workflow issues were a launcher generated with bound
arguments, sandbox-only GitHub credential lookup, and unsupported Vault metadata
field selection. Supported unbound launchers, approved keyring access, and
read-only key listing resolved them without changing shared configuration.

No raw plan, state, credential value, or personal data is retained in this record.
