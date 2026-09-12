## Deployment boundary

Use the owning `//infra/mikrotik/tf:tf` wrappers with the
[canonical declaration](../../../../dnsconfig.json) and existing AL/Vault backend.
The [stage definitions](../../../../tf/BUILD.bazel) own packaging and provider
selection. DNS access uses the configured `src_infra_mikrotik` AppRole.
The root owns DNS only; router export files remain documentation inputs.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, credential handling, import matching, and
recovery. Adopt only this owner's records and retain enabled defaults after
import. Require a no-change adoption plan and verify unchanged unrelated DNS
records and non-DNS resources before applying the reviewed plan.

The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
records the writer audit and authenticated recovery paths used for endpoint
adoption. Its coverage limits for external hosts, external CI, and historical
checkouts remain part of this owner's writer-control evidence.

## Adoption evidence

At `2026-09-13T02:34:28.923664+00:00`, the scoped Vault bootstrap created nine
resources for the dedicated AppRole and DNS access, with no changes or deletions.

At `2026-09-13T02:35:45.838942+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled `tf/dns.tf` SHA256
`cbe1a2154b34b96fa5957d1bc2557f6a2befb053ff47593e8275c7874d938302`.
The DNS-only root imported four existing RouterOS records by exact provider IDs
`*16`, `*23`, `*24`, and `*25`. Its full saved plan contained four managed
`no-op` actions; the saved-plan apply added, changed, and destroyed zero resources.
All 76 records in the complete RouterOS DNS inventory remained equal as JSON
objects keyed by provider ID, including unrelated records.

At `2026-09-13T02:39:55Z`, all four post-adoption A/AAAA queries through
`192.168.1.1` succeeded without truncation and matched the adopted records'
declared addresses. These checks establish DNS resolution for the adopted
records; router exports remain documentation inputs.

The unchanged declaration reuses the baseline ownership-linter evidence.
Current activation checks cover the source default and hash and the actual
full no-change plan. No new offline formatting run is claimed here. Final
aggregate source quality gates remain with repository delivery. Raw state,
plans, credentials, and provider inventories remain in private ignored scratch.
