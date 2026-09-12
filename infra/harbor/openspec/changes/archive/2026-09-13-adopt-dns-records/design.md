## Context

The [`tf_setup` root](../../../../tf_setup/README.md) retains enabled ownership of
Harbor's canonical DNS declaration after verified adoption. The
[proposal](proposal.md) identifies this owner's scope.

## Deployment boundary

Use `//infra/harbor/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply` with the
existing `src_infra_harbor` AppRole and setup backend. These wrappers select
`dns=1` and target `module.dns`; apply accepts the reviewed saved plan.
Successful authenticated adoption verifies the existing role's DNS access.
Ordinary setup wrappers retain their service authentication behavior.

The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns matching, import maps, saved-plan controls, recovery and rollback. The
[shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns writer-audit and recovery coverage; it does not establish universal
disablement of historical checkouts.

## Adoption evidence

At `2026-09-13T02:53:43.110373+00:00`, scoped adoption imported two existing
RouterOS records with no resource additions, changes or deletions. Every managed
action in the saved import plan was `no-op`; the subsequent scoped plan also
contained no changes. All 76 RouterOS records remained equal as complete objects
keyed by provider ID. This owner uses only the internal DNS view.

The receipt records source baseline `9ec1bf2d11863518d7a02c57693ca9c016403f65`
and enabled `tf_setup/dns.tf` SHA256
`a960168c202dac64309abdd85f27ea609b23500d9c7bc2db0fdabbf6f7bbf424`.
At `2026-09-13T02:54:40.359299+00:00`, both declared A queries matched through
`192.168.1.1`; the DNS report is tied to the adoption receipt and canonical
declaration by SHA256.

Source shape/hash checks, scoped wrapper builds and the actual no-change plans
support this historical adoption snapshot. Baseline ownership-linter evidence
covers the unchanged canonical declaration. Current aggregate source formatting
and quality gates remain with repository delivery. DNS targeting does not
establish PVE VM or service health or full-root reconciliation. Keep
`dns_enabled=true` after adoption; disabling it would propose deleting managed
records. Raw plans, state, inventories and credentials remain in private
ignored scratch.
