## Context

The [host DNS root](../../../../tf/README.md) retains enabled record ownership after
verified adoption. The [proposal](proposal.md) identifies the ownership scope.

## Deployment boundary

Use `//users/simeonwarren/host_bot/tf:tf.import`, `:tf.plan`, and `:tf.apply` with
the canonical host DNS declaration. This DNS-only root selects `dns_al`, the
`src_users_simeonwarren_host_bot` AppRole, and `tf=main` injection for its own
Vault backend and both DNS providers. Scoped DNS access is provisioned through
the owning Vault Terraform flow. Preserve Ansible's existing `al` configuration and
`user_simeonwarren` authority.

The shared [cutover procedure](../../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns ordering, writer coordination, provider matching and import rules,
credential handling, recovery prerequisites, and rollback. Apply it to this
host's records in both views through this exact root. Require a reviewed plan
with zero record additions, changes, replacements, or deletions before applying;
then verify both DNS views and unrelated record preservation.

## Adoption evidence

The scoped Vault bootstrap created nine AppRole resources at
`2026-09-13T02:32:30.078495Z`. Adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled `tf/dns.tf` SHA256
`256cd961cd69fcc1def910396ee1ad0813fe04a7e1b1cff6b69d6bdcbd243179`.
At `2026-09-13T02:34:09.165104+00:00`, seven existing records were imported:
four RouterOS and three Cloudflare bindings. The full saved plan and apply
changed no managed resources. All 76 RouterOS and 80 Cloudflare records remained
equal as complete JSON objects keyed by provider ID.

At `2026-09-13T02:34:15Z`, five representative DNS queries matched the preserved
records without truncation: two internal A lookups through `192.168.1.1` and
three public A/CNAME lookups through `1.1.1.1`. This is representative resolution
coverage; provider inventory comparison covers all seven managed records.

The [shared operational evidence](../../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns authenticated recovery results and the writer audit's coverage limits.
It does not establish universal disablement of historical checkouts.

Baseline ownership-linter results cover unchanged canonical declarations.
Successful stage builds, source shape/hash checks, and the full no-op plan
support activation at the recorded adoption snapshot. Final source formatting
and aggregate quality gates for the current source remain with repository
delivery. Raw provider inventories, state, plans, and credentials
remain in private ignored scratch.
