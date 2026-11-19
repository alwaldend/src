## Deployment boundary

Use the owning `//infra/ingress/tf:dns.plan`, `:dns.show`, and `:dns.apply` wrappers with the
[canonical declaration](../../../../dnsconfig.json) and existing AL/Vault backend.
The [stage definitions](../../../../tf/BUILD.bazel) own packaging and provider
selection. DNS access selects `dns=1` with the existing `src_infra_ingress`
AppRole and backend. The wrappers target `module.dns`; their plans validate
DNS and its dependencies, not overall ingress service health. Ordinary service
authentication and management remain separate.

## Adoption and recovery

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, credential handling, import matching, and
recovery. Checked-in import blocks bind exact provider IDs through a reviewed
saved plan. Require no record changes, no non-DNS managed actions, a no-change
follow-up plan, and complete provider-inventory equality. Retain enabled defaults.

The owner manages infrastructure endpoint names. Verify the applicable independent
authenticated recovery paths before its batch, as required by the shared procedure.

## Verified adoption

Adoption completed at `2026-09-13T02:55:17.910904Z` from source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` and the receipt's source hashes.
The retained enabled `tf/dns.tf` SHA256 is
`e089f896a935fe12cfd6a6e8a0bda4e465e54b67443d417329a03c1ba8ed03e7`.
Authenticated operations verified access through the existing owner role.

The scoped saved plan imported all four records, two in each view, with
zero additions, changes, or deletions. Its managed actions were DNS-only
`no-op` imports; the follow-up plan had no changes or pending imports.
Complete inventories retained all identities and attributes: 80 Cloudflare
and 76 RouterOS records, including unrelated owners. The receipt is
`out/dns-deploy/private/infra_ingress-scoped-adopt-20260913T025352652207Z/verification.json`.
It records bindings, source hashes, and saved-plan digests; raw plans and
inventories remain in private ignored scratch.

All four declared A queries for `host1.ingress.alwaldend.com` and
`ingress.alwaldend.com` matched `81.26.185.118` through the public and dc1
resolvers at `2026-09-13T02:56:07Z`. The report is
`out/dns-deploy/dns-probes/infra_ingress.json` and binds its declaration and
verification-receipt hashes.

The shared [writer audit and authenticated recovery evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns the verified controls and their coverage limits; universal historical
writer disablement is not claimed. Unchanged declarations reuse the
coordinator's prerequisite linter evidence. Enabled source shape, wrapper builds,
and the scoped plans cover this adoption snapshot; final source formatting
and aggregate quality gates remain with repository delivery.
