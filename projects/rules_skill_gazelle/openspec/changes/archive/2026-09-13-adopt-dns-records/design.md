## Context

See the [proposal](proposal.md). The [project DNS root](../../../../tf/README.md)
owns the canonical Terraform source and retains enabled record ownership after
verified adoption.

## Deployment boundary

Use the documented
[nested-project adapter](../../../../../../infra/dns/projects/README.md):
`//infra/dns/projects:rules_skill_gazelle_tf.<operation>` operates the owning
`projects/rules_skill_gazelle/tf` stage through its existing `tf=main` flow.
The project's [AL configuration](../../../../al.lua) owns its AppRole and backend;
the [DNS declarations](../../../../dnsconfig.json) remain canonical. This preserves
the reusable module boundary and avoids a second operational definition.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, access
prerequisites, import and plan review, record preservation, and recovery.
Use that procedure for this owner and retain a separate adoption receipt here.

Writer coordination relies on the bounded
[shared operational evidence](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md#writer-audit)
and the root coordinator's serialized adoption sequence. The documented audit
coverage limitations continue to apply.

## Risks / Trade-offs

Applying before complete import could change existing records; require the
shared procedure's no-change adoption plan before applying. Keep ownership
enabled after import because restoring the disabled default would propose
record deletion.

## Adoption evidence

At `2026-09-13T02:01:31.625396+00:00`, adoption used confirmed source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/variables.tf` SHA256
`867691558b30fb00ff12a1e552094ad50f14070953842968df410be0568f9e2c`.
The scoped Vault stage created nine AppRole resources. The project imported its
one existing Cloudflare record, `578739344d0b40d3e451b3f99d2e0077`, through
`//infra/dns/projects:rules_skill_gazelle_tf.import`. The reviewed plan and apply
changed no managed resources. All 80 public-zone records remained equal as
complete JSON objects keyed by provider ID, with no additions, removals, or changes.

At `2026-09-13T02:03:09Z`, public CNAME resolution for
`rules-skill-gazelle.alwaldend.com` through `1.1.1.1` returned
`alwaldend.github.io`, matching the preserved provider record with TTL 600 and
proxying disabled.

The passed baseline source checks and ownership-linter result cover unchanged
DNS declarations. Enabled wrapper builds, source shape/hash checks, and the
adoption plan cover activation. Final aggregate source validation remains with
repository delivery. Raw state, plans, inventories, and credentials remain in
private ignored task scratch.
