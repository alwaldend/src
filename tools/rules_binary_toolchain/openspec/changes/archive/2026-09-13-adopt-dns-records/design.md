## Deployment boundary

The project's [Terraform root](../../../../tf/README.md) owns its DNS state and
canonical declarations. Use its documented nested adapter
`//infra/dns/projects:rules_binary_toolchain_tf.<operation>` for the `tf` stage,
including import, plan, and apply. The adapter preserves the reusable module
boundary and selects the normal `tf=main` credential injection.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
Require exact provider IDs and a plan with zero record additions, changes,
replacements, or deletions before applying. Keep ownership enabled after
adoption; disabling it would propose deletion of the adopted records.

## Adoption evidence

At `2026-09-13T01:53:27.187375+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/variables.tf` SHA256
`867691558b30fb00ff12a1e552094ad50f14070953842968df410be0568f9e2c`.
The scoped Vault stage created 9 AppRole resources. The project imported its
1 existing Cloudflare record; the reviewed plan and apply changed no managed
resources. All 80 public-zone records remained equal as JSON objects keyed by
provider ID, with no additions, removals, or changes. Public CNAME resolution
through `1.1.1.1` returned `alwaldend.github.io`; the provider record
retains TTL 600 and proxying disabled.

The baseline ownership-linter result covers unchanged DNS declarations.
Current wrapper builds, source shape/hash checks, and the adoption plan cover
the activation. Final aggregate source validation remains with repository
delivery. Raw plans, inventories, and credentials remain in private ignored
task scratch.
