## Deployment boundary

Use `//infra/vault/tf_setup:dns.plan`, `dns.show`, and `dns.apply` with the
[canonical declaration](../../../../dnsconfig.json) and existing AL/Vault backend.
The [stage definitions](../../../../tf_setup/BUILD.bazel) own packaging and provider
selection. DNS access uses the configured `src_infra_dc1_vault` AppRole.
These targets select `dns=1` and `module.dns` in the same setup root and
backend. The ordinary setup and service wrappers retain their authentication
flow. Vault service configuration remains in the separate `tf` root.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, credential handling, import matching, and
recovery. Adopt only this owner's records and retain enabled defaults after
import. Require exact provider IDs and no resource changes in the saved
adoption plan, apply that reviewed file, and verify a no-change follow-up plan
and unchanged complete provider inventories.

The shared [writer audit and authenticated recovery evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
records the endpoint prerequisites and writer-control observation limits,
including external hosts, CI, and manual historical checkouts. Targeted DNS
results do not establish the health of unrelated setup or service resources.

## Adoption evidence

At `2026-09-13T03:03:58.273245+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf_setup/dns.tf` SHA256
`e089f896a935fe12cfd6a6e8a0bda4e465e54b67443d417329a03c1ba8ed03e7`.
Authenticated scoped commands used the existing AppRole and setup backend.
The reviewed plan selected five RouterOS records and two Cloudflare records
for import by exact provider ID; every managed action was no-op. Apply reported seven
imports and zero additions, changes, or deletions. The follow-up plan retained
seven no-op DNS records with no pending imports.

Complete inventories retained all 80 Cloudflare records and 76 RouterOS
records, with identical provider IDs and complete record objects. At
`2026-09-13T03:04:11.576922+00:00`, all seven DNS queries across the declared
views matched the receipt-hashed declarations and adopted records.

The baseline ownership-linter evidence covers the unchanged DNS declaration.
Current source hashes and the saved plans cover the enabled adoption source.
Final aggregate source validation remains with repository delivery. Raw
state, plans, import mappings, and provider inventories remain in private
ignored task scratch.
