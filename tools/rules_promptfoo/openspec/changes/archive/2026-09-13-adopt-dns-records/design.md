## Context

See the [proposal](proposal.md). The [project DNS root](../../../../tf/README.md)
owns the canonical Terraform source and retains enabled record ownership after
verified adoption.

## Deployment boundary

Use the documented
[nested-project adapter](../../../../../../infra/dns/projects/README.md):
`//infra/dns/projects:rules_promptfoo_tf.<operation>` operates the owning
`projects/rules_promptfoo/tf` stage through its existing `tf=main` flow.
The project's [AL configuration](../../../../al.lua) owns its AppRole and backend;
the [DNS declarations](../../../../dnsconfig.json) remain canonical. This preserves
the reusable module boundary and avoids a second operational definition.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, access
prerequisites, import and plan review, record preservation, and recovery.
Use that procedure for this owner and retain a separate adoption receipt here.
Writer-control acceptance uses the shared
[writer audit](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md#writer-audit),
including its coverage limits for external hosts, CI, and manual historical
checkouts.

## Risks / Trade-offs

Applying before complete import could change existing records; require the
shared procedure's no-change adoption plan before applying. Keep ownership
enabled after import because restoring the disabled default would propose
record deletion.

## Adoption evidence

At `2026-09-13T01:59:30.581481Z`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/variables.tf` SHA256
`867691558b30fb00ff12a1e552094ad50f14070953842968df410be0568f9e2c`.
The scoped Vault stage created nine AppRole resources. The project imported its
one existing Cloudflare record; the reviewed plan and apply changed no managed
resources. All 80 public-zone records remained equal as JSON objects keyed by
provider ID, with no additions, removals, or changes. Public CNAME resolution
through `1.1.1.1` returned `alwaldend.github.io`; the provider record
retains TTL 600 and proxying disabled.

The baseline ownership-linter result covers unchanged DNS declarations.
Current adapter builds, source shape/hash checks, and the adoption plan cover
the activation. Final aggregate source validation remains with repository
delivery. Acceptance is tracked in [tasks](tasks.md). Raw plans, inventories,
and credentials remain in private ignored task scratch.
