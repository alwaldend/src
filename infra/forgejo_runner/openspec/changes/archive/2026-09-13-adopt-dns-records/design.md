## Context

See [the proposal](proposal.md) for the adoption scope. The checked-in
[`tf_setup` package](../../../../tf_setup/README.md) retains enabled DNS resources and
packages the owner's canonical declaration with the shared DNS module.

## Deployment boundary

Use `//infra/forgejo_runner/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply`.
These targets select `dns=1` alone and target `module.dns` in the same owning
setup root and backend, using the existing `src_infra_forgejo_runner` AppRole.
Imports are declared through the root's optional record-ID map and reviewed in
the saved plan. Apply only that reviewed plan, then require a no-change follow-up
DNS plan. The ordinary setup targets retain their `tf=setup` flow.

The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns prerequisites, ordering, writer coordination, credential handling, and
recovery. Keep its procedure canonical; this change records the owner's outcome.
The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns the writer audit and authenticated recovery observations with their stated
coverage limits. Adoption remains coordinated one owner at a time.

## Decisions and risks

Import the matching existing provider identities into this owner's state so the
transfer preserves the live records. Keep `dns_enabled=true` in the adopted
revision and update the setup documentation to match that supported state;
restoring a disabled default after import would propose record deletion.

The adoption plan must propose no DNS additions, changes, replacements, or
deletions. Targeting limits the operation to DNS and its dependencies; it does
not establish PVE or runner service health or reconcile the full setup root.

## Adoption evidence

At `2026-09-13T02:52:28.571475+00:00`, the adoption receipt recorded enabled
`tf_setup/dns.tf` SHA256
`a960168c202dac64309abdd85f27ea609b23500d9c7bc2db0fdabbf6f7bbf424`.
Authenticated operations through the existing AppRole imported RouterOS row
`*53` at `module.dns.module.dc1[0].routeros_ip_dns_record.records["runner1/A/dc1"]`.
The inspected saved plan contained one import with a no-op resource action;
apply reported one import and zero additions, changes, or deletions. The
follow-up scoped plan contained one no-op record and no pending import.
All 76 RouterOS records remained equal as complete JSON objects keyed by row ID.
The owner manages only the `dc1` view; Cloudflare inventory was not part of this
adoption.

At `2026-09-13T02:53:20.456385+00:00`, the saved DNS probe through
`192.168.1.1` resolved `runner1.forgejo-runner.alwaldend.com` to `192.168.10.100`,
matching the canonical declaration and imported plan. Its receipt and canonical
source digests matched the inspected inputs.

The enabled source shape/hash and all three scoped wrapper builds passed their
checks. The unchanged canonical declaration reuses the passed baseline source
and ownership-linter evidence. These are adoption snapshot checks; final
aggregate source validation remains with repository delivery. Raw state, plans,
inventories, and credentials remain in task-private ignored scratch.
