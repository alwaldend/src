## Context

The existing [`tf_setup` root](../../../../tf_setup/README.md) consumes the owner's
canonical DNS declaration and retains enabled DNS ownership after provisioning.
The scoped `dns.plan`, `dns.show`, and saved-plan `dns.apply` wrappers select
`dns=1`, the existing `src_infra_openhands` AppRole, and the same backend.
See [the proposal](proposal.md) for the adoption scope.

## Deployment boundary

Use `//infra/openhands/tf_setup:dns.plan`, `dns.show`, and `dns.apply` for
`module.dns` and its dependencies. The root also manages VMs; ordinary setup
retains `tf=setup` and `xoa_login=1`, while DNS operations do not start XO login.
The saved DNS plan contains every declared DNS resource and no other managed
resource change. Successful scoped operations establish existing-role DNS and
backend access; this batch required no new AppRole bootstrap.

Fresh complete inventories established that all 13 declared records were
missing and conflict-free. The authorized batch created four Cloudflare and
nine RouterOS records; it imported no existing records. Keep ownership enabled
so reconciliation retains those bindings and preserves existing provider records.

The shared [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns prerequisites, writer coordination, import ordering, credential handling,
plan acceptance, and recovery. The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns writer-control observations and their coverage limits; this adoption does
not establish universal disablement of historical writers.

## Risks and acceptance

- A full-root apply could create or change VMs. The saved DNS-scoped plan must
  contain exactly four Cloudflare and nine RouterOS creates, with no other
  managed changes, updates, replacements, or deletions.
- Incorrect record matching can change a live identity or view. Compare exact
  provider IDs and declared attributes. After initial provisioning, require
  zero DNS additions, changes, replacements, or deletions in a follow-up plan.
- Overlapping writers or disabling management after provisioning can damage managed
  records. Follow the shared cutover and recovery procedure.

## Completed provisioning and verification recovery

The adoption source baseline was `9ec1bf2d11863518d7a02c57693ca9c016403f65`, with
enabled `tf_setup/dns.tf` SHA256
`e089f896a935fe12cfd6a6e8a0bda4e465e54b67443d417329a03c1ba8ed03e7`.
At `2026-09-13T03:08:43Z`, applying the reviewed saved plan created exactly
13 records, with zero imports, updates, replacements, or deletions. The initial
post-apply verifier halted with `cloudflare_inventory_truncated` because the
Cloudflare inventory pagination metadata was inconsistent.

At `2026-09-13T03:16:05.836216+00:00`, a verification-only resume completed;
it performed no second apply, import, or source activation. Complete inventories
contained 84 Cloudflare and 85 RouterOS records. Every previous provider ID and
complete record object remained unchanged, with exactly four Cloudflare and
nine RouterOS additions matching the declarations. The follow-up scoped plan
contained 13 DNS `no-op` actions and no pending imports.

At `2026-09-13T03:16:28.012974+00:00`, all 13 declared record queries matched:
four global queries through `1.1.1.1` and nine `dc1` queries through
`192.168.1.1`, without errors or truncation. The DNS report binds its results
to the declaration and updated verification receipt by SHA256.

The unchanged declaration reuses the baseline ownership-linter evidence.
Source hashes and the saved plans cover the enabled adoption source; aggregate
source quality and delivery gates remain pending root validation. Sanitized
receipts retain created provider IDs and artifact digests. Raw inventories,
state, plans, and credentials remain in private ignored task scratch.

## Separate service readiness

At approximately `2026-09-13T01:55Z`, the three declared service IPs were
unreachable on TCP 22 and 443. Forcing the Canvas hostname to the current ingress
reached its default certificate without a matching Canvas identity. These
checks do not establish service readiness or mTLS behavior. Missing DNS can
also prevent the configured certificate issuance. This DNS operation changes
neither VMs nor ingress configuration; those earlier readiness observations
are not resolved by successful DNS queries.

The active [service deployment change](../../add-openhands-deployment/design.md)
retains service acceptance and its TLS, issuance, and authentication limitations.
DNS verification also does not resolve the separately recorded
[PVE host/API reachability failure](../../../../../../projects/alwaldend.com/openspec/changes/archive/2026-09-13-adopt-dns-records/design.md#earlier-blocked-attempt).
