## Deployment boundary

Use the owning `//infra/nas/tf:tf` wrappers with the
[canonical declaration](../../../../dnsconfig.json) and existing AL/Vault backend.
The [stage definitions](../../../../tf/BUILD.bazel) own packaging and provider
selection. DNS access uses the configured `src_infra_nas` AppRole.
The root owns DNS only; NAS service provisioning remains outside it.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, credential handling, import matching, and
recovery. Adopt only this owner's records and retain enabled defaults after
import. Require a no-change adoption plan and verify unchanged unrelated DNS
records and non-DNS resources before applying the reviewed plan.

## Adoption evidence

The scoped Vault bootstrap created nine AppRole resources at
`2026-09-13T02:06:22.046516+00:00`. Adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled `tf/dns.tf` SHA256
`cbe1a2154b34b96fa5957d1bc2557f6a2befb053ff47593e8275c7874d938302`.

An earlier attempt imported seven existing records, then verification rejected
omitted RouterOS provider-default fields. The verifier was corrected to accept
null `comment` and `match_subdomain` only for a no-op plan; live matching still
rejects nonempty comments and `match_subdomain=true`. The successful resume at
`2026-09-13T02:11:45.684054+00:00` reused all seven bindings and imported no new
records: five RouterOS records and two Cloudflare records. Its full saved plan
and apply changed no managed resources. All 76 RouterOS and 80 Cloudflare records
remained equal as complete JSON objects keyed by provider ID.

At `2026-09-13T02:15:04Z`, all seven declared DNS results matched: five site-local
records through `192.168.1.1` and two public records through `1.1.1.1`. The checks
covered the declared A and CNAME records in both views without truncated results.

The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns authenticated recovery results and the writer audit's coverage limits.
It does not establish universal disablement of historical checkouts.

Baseline source and ownership-linter checks cover unchanged canonical
declarations. Successful stage builds, source shape/hash checks, and the full
adoption plan cover activation. Final aggregate source validation remains with
repository delivery. Raw state, plans, inventories, and credentials remain in
private ignored task scratch.
