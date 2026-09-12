## Deployment boundary

The owning `infra/xcp_ng/tf` root consumes the
[canonical declaration](../../../../dnsconfig.json). Its
[stage definitions](../../../../tf/BUILD.bazel) expose `dns.plan`, `dns.show`,
and saved-plan `dns.apply` commands selecting `dns=1` and targeting
`module.dns`. The existing root, HTTP backend, and `src_infra_xcp_ng` AppRole
through named `xcp_ng` Vault authentication remain in use. Ordinary service
wrappers retain their `xcp_ng_tf`, `xoa`, and Vault environment labels and
authentication.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, credential handling, import matching, and recovery. Its
[operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns writer-control observations, their coverage limits, and independent
authenticated recovery. Adoption does not establish universal disablement of
historical writers.

## Verified adoption

Adoption was verified at `2026-09-13T03:05:42.186806Z` against the source
baseline `9ec1bf2d11863518d7a02c57693ca9c016403f65` with the enabled source
captured by content hashes. The current owner source matches all 12 recorded
post-adoption hashes. The canonical declaration SHA256 is
`2a19d3febeddcadce6902fa48ee75a869bbd9dd8f8089bca0c46c4a3d0ac5abc`;
the enabled `tf/dns.tf` SHA256 is
`a960168c202dac64309abdd85f27ea609b23500d9c7bc2db0fdabbf6f7bbf424`.
Its checked-in `dns_enabled` default is true.

The existing AppRole successfully authenticated the scoped backend and DNS
provider flow. The saved adoption plan contained exactly four managed DNS
resources, each with a matching declarative import and a `no-op` action:

| Normalized record key | RouterOS ID | Name                             |
| --------------------- | ----------- | -------------------------------- |
| `host1/A/dc1`         | `*77`       | `host1.xcp-ng.alwaldend.com`     |
| `root/A/dc1`          | `*76`       | `xcp-ng.alwaldend.com`           |
| `xoa_host1/A/dc1`     | `*79`       | `host1.xoa.xcp-ng.alwaldend.com` |
| `xoa_root/A/dc1`      | `*78`       | `xoa.xcp-ng.alwaldend.com`       |

The owning address prefix is
`module.dns.module.dc1[0].routeros_ip_dns_record.records`. Actual apply output
reported four imports and zero additions, changes, or deletions. The saved
follow-up plan contained the same four DNS resources, all `no-op`, with no
remaining imports. Both plans reported no errors and were scoped plans
(`complete=false`).

The complete RouterOS inventories contain 76 records before and after;
their full record objects are equal when keyed by unique provider ID. The
21 non-DNS managed resource objects in the adoption and follow-up prior states
are equal by address, and neither plan contains non-DNS managed actions.
This owner declares only the `dc1` view; Cloudflare inventory verification is
not applicable to this adoption.

| Evidence artifact             | Verified SHA256                                                    |
| ----------------------------- | ------------------------------------------------------------------ |
| Saved adoption plan           | `b9d1244ca80df0a3c90df9480458b7575b5d51c0bb48ada1f886be20e633a099` |
| Saved follow-up plan          | `abf17f2c74c0b2213ca245ef571e0f531934fa907d50eeba9c9d12f2d06100d0` |
| Declarative import map        | `8a5712053ed5bdab799a2b1b6699f1cb2714aea6807d0d60951f7f453104e00b` |
| Adoption verification receipt | `28c29f570b3b11420abff15588b09a9acd69fcd4a7b22c76e07f7e37270e98ab` |
| Post-adoption DNS receipt     | `88a4413e56eada020f1f201d37fdca8980c4d49a75a8a986f4d0f6a510bc606b` |

## DNS verification and limits

At `2026-09-13T03:06:51.262515Z`, the post-adoption probe bound its results to
the declaration and verification-receipt hashes above. All four `A` queries
through resolver `192.168.1.1` returned the declared `dc1` answers:
`host1.xcp-ng.alwaldend.com` and `xcp-ng.alwaldend.com` returned
`192.168.1.213`; `host1.xoa.xcp-ng.alwaldend.com` and
`xoa.xcp-ng.alwaldend.com` returned `192.168.1.206`. The receipt records no
truncation or query errors.

The coordinator authorized reuse of unchanged-declaration linter evidence at
the baseline revision above. The prerequisite batch passed nine tests, and the
SRI activation batch passed four, including shared DNS linter and repository
checks. These are reused source checks, not an XCP-ng formatting-test run.
Current aggregate source quality gates remain with the delivery coordinator.

This evidence verifies scoped DNS adoption, DNS answers, provider inventory
preservation, and recorded non-DNS state preservation. It does not establish
full Xen Orchestra or XCP-ng service health. Keep ownership enabled after
import; disabling it would propose deletion. Raw plans, state, inventories,
and credentials remain in access-restricted ignored task scratch and are not
publication artifacts.
