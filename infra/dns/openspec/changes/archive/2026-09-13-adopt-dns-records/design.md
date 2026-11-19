## Deployment boundary

Use `//infra/dns/tf:tf.import`, `:tf.plan`, and `:tf.apply` with the
[canonical shared declaration](../../../../dnsconfig.json), existing Vault HTTP
backend, and `src_infra_dns` AppRole. The [stage definition](../../../../tf/BUILD.bazel)
owns provider injection for both declared views. Adoption covers only this
owner's apex and mail records.

## Adoption and recovery

The shared [cutover procedure](../2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns prerequisites, ordering, writer coordination, credential handling, record
matching, and recovery. Retain enabled source defaults after exact-ID imports.
Require zero record additions, changes, replacements, or deletions in the
adoption plan, then verify both views and preservation of unrelated records.

## Verified adoption

Adoption completed at `2026-09-13T02:42:01.133113Z` from source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` and the receipt's source hashes.
The retained enabled `tf/dns.tf` SHA256 is
`8a13afc0fc4da54f69e82221feaa2928bf7480c6b2e69620c91dfa5a10545a2b`.
The existing `src_infra_dns` AppRole successfully accessed its owning backend
and both providers. No new owner role was needed for this batch.

All 60 existing records were imported by provider ID: 30 Cloudflare and
30 RouterOS. The full saved plan contained 60 managed `no-op` actions;
its apply reported zero additions, changes, or deletions. The complete before
and after inventories were identical by provider identity and attributes:
80 Cloudflare records and 76 RouterOS records, including unrelated owners.
Provider bindings, source hashes, plan digest, and comparisons are recorded in
`out/dns-deploy/private/infra_dns-adopt-20260913T023748162472Z/verification.json`.
Raw plans and inventories remain in private ignored scratch.

At `2026-09-13T02:42:38Z`, eight representative apex A, AAAA, MX, and TXT
queries through `1.1.1.1` and `192.168.1.1` matched the pre-adoption answers.
The reports are `out/dns-deploy/apex-before-dns-verification.json` and
`out/dns-deploy/apex-post-dns-verification.json`; all responses were complete.
These queries supplement the complete provider-inventory comparison.

The shared [writer audit and recovery evidence](../2026-09-13-migrate-project-dns-to-terraform/operations.md)
records the verified operational controls and their coverage limits; it does
not prove every historical checkout or external writer disabled. The unchanged
declarations reuse the coordinator's prerequisite linter evidence. Wrapper
builds, enabled source shape, and the successful live plan cover this adoption
snapshot; final source formatting and aggregate quality gates remain with
repository delivery. The shared migration stays active for pending owners.
