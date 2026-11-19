# Shared deployment evidence

The user authorized DNS deployment one owner at a time with existing records
and services preserved. Owner-local adoption changes own individual outcomes;
this file records shared prerequisites and their observation limits.

## Writer audit

At `2026-09-13T00:51Z`, source revision
`246c4e111b4c69f442a0530f1342fb886a7074c6` exposed no central DNSControl write
target or scheduled caller. Inspection of host-bot found no matching active
process, timer, service-unit reference, or standard cron/at configuration.
The audit covered six system and two user timers, and 95 system and two user
unit files. No active central writer was identified to stop.

Other hosts, external CI, and manual invocations from historical checkouts
remain outside that observation. Removing current entrypoints does not disable
old revisions. Whole-provider inventories are compared before and after every
adoption to detect concurrent record changes. Shared credentials remain valid
because the new owner roots use them too.

## Authenticated recovery

At `2026-09-13T01:07:24Z`, before adopting endpoint owners, authenticated probes reached Vault directly
at `192.168.1.218:8200` while verifying the TLS identity
`vault.alwaldend.com`, and RouterOS directly at `192.168.1.1:443` while verifying
`router1.dc1.alwaldend.com` against the repository's CA. Both bypassed DNS
resolution and proxies without disabling certificate validation. Vault token
lookup succeeded; RouterOS returned its authenticated system resource response
and confirmed the direct peer address. Existing networking, certificates and
credentials were sufficient; no recovery configuration was changed.

## Provider discovery and adoption gates

The existing Cloudflare credential could list and read the configured zone.
The shared module now discovers the zone when its optional ID is absent.
DNS and ingress use their existing shared RouterOS endpoint; no Vault secret
metadata fields were added or overwritten.

Initial live snapshots contained 80 Cloudflare records and 76 enabled static
RouterOS records. These are observations, not desired-state counts. The source
baseline and live inventory differ, so every owner must independently match
its records before import. A discrepancy pauses that owner for diagnosis.

Vault prerequisite plans are limited to the selected role's DNS access or
new AppRole module. Existing-record adoption imports exact provider IDs through
a reviewed no-change plan. Verified missing declarations may use an exact
additions-only plan. Every batch preserves all preexisting provider identities
and attributes, including unrelated records. Raw plans and inventories were
retained in access-restricted ignored task scratch for verification; sanitized
audit receipts retain the observations after cleanup.

## Verified completion

All 45 nonempty owners were verified by `2026-09-13T03:16:05.836216Z` from
source baseline `9ec1bf2d11863518d7a02c57693ca9c016403f65` and the per-owner
receipt hashes. They retain enabled defaults. The final audit at
`2026-09-13T03:16:25Z` verified 155 existing-record adoptions and 13 new
OpenHands records: all 168 declared bindings, comprising 83 global and 85 dc1.
The empty user declaration remains metadata only and has no DNS owner state.

Final complete inventories contain 84 Cloudflare and 85 RouterOS records.
The additional Cloudflare row predates this rollout and remains outside
declared ownership. Every preexisting row retained its identity and attributes;
only the four global and nine dc1 OpenHands records were created. Each final
owner plan is a no-op, current canonical hashes match receipts, and available
activation-file hashes match the retained source. Identical OpenHands resume
receipts were counted once.

The sanitized aggregate is `out/dns-deploy/coverage-audit.json`, SHA256
`9725182a6c14cc8f3603f5073e3405eb6c363dd89e8cad237e90d2a4f71d58e0`.
It records source scope, owner receipts, plan digests, record comparisons,
and inventory observations, with no pending owners or errors.

The last-phase resolver audit at `2026-09-13T03:17:05Z` passed all 58 queries
for its 11 owners: 17 global and 41 dc1. Its report is
`out/dns-deploy/dns-probes/summary.json`, SHA256
`a329d3c019efc84aac28c30874ef68235f202a08e32b52a147bd669d75f35a3b`.
Earlier owner-local evidence covers prior batches. The DNS handler does not
report TTL; provider and plan comparisons cover stored attributes. DNS checks
do not establish application, hypervisor, or TLS readiness. Writer-audit
coverage limits above remain applicable after completion.

Final rollout candidate formatting, quality gates, and publication remain
pending repository delivery; these operational receipts do not replace them.
