## Context

The owner's `tf_setup` root packages its canonical declaration and shared
DNS module. Its adopted `dns_enabled` default remains true. See
[proposal.md](proposal.md) for the ownership change.

## Deployment boundary

Use `//infra/threexui/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply` with
`dns=1` and the existing `src_infra_threexui` AppRole and backend. The wrappers
target `module.dns` and preserve ordinary service authentication. Their plans
validate DNS and its dependencies; they do not establish overall service health.

The shared
[cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns writer coordination, adoption ordering, credential handling, and recovery.
Use checked-in import blocks and exact provider IDs in a reviewed saved plan,
then retain enabled ownership and verify a no-change follow-up plan.

## Risks and acceptance

An adoption plan must preserve provider identities and propose zero record
additions, changes, replacements, or deletions and no non-DNS managed actions.
Verify the declared DNS views and compare the complete provider inventories
before and after the reviewed apply.

## Verified adoption

Adoption completed at `2026-09-13T03:00:26.400156Z` from source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` and the receipt's source hashes.
The retained enabled `tf_setup/dns.tf` SHA256 is
`e089f896a935fe12cfd6a6e8a0bda4e465e54b67443d417329a03c1ba8ed03e7`;
the declaration SHA256 is
`1fe9f7dd3e39bbbbbf6fc3080a374c266c8804e00f30d9d4fbf3f276cb1b77ed`.
Authenticated operations verified access through the existing owner role.

The saved plan imported all 12 existing records, five Cloudflare and seven
RouterOS, with zero additions, changes, or deletions. All managed actions were
DNS-only `no-op` imports; the follow-up plan had no changes or pending imports.
Complete inventories preserved all identities and attributes: 80 Cloudflare
and 76 RouterOS records. The receipt is
`out/dns-deploy/private/infra_threexui-scoped-adopt-20260913T025849410935Z/verification.json`.
It records bindings, source hashes, and saved-plan digests; raw plans and
inventories remain in private ignored scratch.

All ten distinct declared name/type/view queries matched at
`2026-09-13T03:01:11Z`, including both NS values in each view and parsed IPv6
equivalence. `out/dns-deploy/dns-probes/infra_threexui.json` binds the answers
to the declaration and verification-receipt hashes.

The first attempt stopped before apply with no imports because two AAAA
resources proposed spelling-only updates. The declaration was normalized to
the existing provider spelling `2a0a:3840:8078:141:0:2d8e:8d85:1337`, preserving
the parsed address; the coordinator's declaration test passed. The successful
retry retained the already-enabled default and changed no source during its run.

The shared [writer audit and recovery evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns verified operational controls and their coverage limits; universal
historical writer disablement is not claimed. Enabled source shape, wrapper
builds, the declaration test, and scoped plans cover this adoption snapshot.
Final source formatting and aggregate quality gates remain with repository delivery.
