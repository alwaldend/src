## Deployment boundary

Use the owning `//infra/pve/tf:dns.plan`, `dns.show`, and `dns.apply` wrappers with the
[canonical declaration](../../../../dnsconfig.json) and existing AL/Vault backend.
The [stage definitions](../../../../tf/BUILD.bazel) own packaging and provider
selection. DNS access uses the configured `src_infra_dc1_pve1` AppRole.
The wrappers select `dns=1` and target `module.dns` in the same root and backend;
apply consumes the reviewed saved plan. Existing-role access was verified by
successful authenticated backend and provider operations, without a new bootstrap.
Ordinary Proxmox authentication remains separate. This workflow verifies DNS and
its dependencies; it does not establish PVE host, API, or service health.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, credential handling, import matching, and
recovery. Retain enabled defaults after import because disabling ownership would
propose deletion. The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
records the writer audit's coverage limits and independent authenticated recovery
paths. It does not establish universal disablement of historical writers.

The [earlier PVE reachability failure](../../../../../../projects/alwaldend.com/openspec/changes/archive/2026-09-13-adopt-dns-records/design.md#earlier-blocked-attempt)
remains a separate service observation. Successful DNS adoption and resolution
do not supersede that host/API result.

## Completed scoped adoption

At `2026-09-13T03:02:22.927061+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled `tf/dns.tf` SHA256
`a960168c202dac64309abdd85f27ea609b23500d9c7bc2db0fdabbf6f7bbf424`.
The reviewed scoped plan contained five exact RouterOS imports and five managed
`no-op` actions. Applying that saved plan completed all five declarative imports
with zero additions, changes, or deletions. The follow-up scoped plan contained
five DNS `no-op` actions and no pending imports. All 76 RouterOS records remained
equal as complete JSON objects keyed by provider ID, including unrelated records.

At `2026-09-13T03:02:36.959684+00:00`, all five declared `dc1` DNS queries through
`192.168.1.1` matched their expected values without errors or truncation. The DNS
evidence is bound to the adoption receipt and canonical declaration by SHA256.

The unchanged declaration reuses the baseline ownership-linter evidence.
The source default/hash and scoped plan checks cover activation; final aggregate
source quality gates remain with repository delivery. Raw plans, import maps,
inventories, and credentials remain in private ignored task scratch.
