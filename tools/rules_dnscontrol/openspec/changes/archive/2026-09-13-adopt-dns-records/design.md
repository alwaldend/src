## Deployment boundary

The project's [Terraform root](../../../../tf/README.md) owns its DNS state and
canonical declarations. Use its documented nested adapter
`//infra/dns/projects:rules_dnscontrol_tf.<operation>` for the `tf` stage,
including import, plan, and apply. The adapter preserves the reusable module
boundary and selects the normal `tf=main` credential injection.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
Require exact provider IDs and a plan with zero record additions, changes,
replacements, or deletions before applying. Keep ownership enabled after
adoption; disabling it would propose deletion of the adopted records.

## Adoption evidence

At `2026-09-13T01:54:26.694482+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/variables.tf` SHA256
`867691558b30fb00ff12a1e552094ad50f14070953842968df410be0568f9e2c`.
The scoped Vault stage created nine AppRole resources. The project imported its
one existing Cloudflare record, `f33b4321380e83cda82426b130d06f69`; the reviewed
plan and apply changed no managed resources. All 80 public-zone records remained
equal as complete JSON objects keyed by provider ID, with no additions, removals,
or changes.

At `2026-09-13T01:58:14Z`, public CNAME resolution for
`rules-dnscontrol.alwaldend.com` through `1.1.1.1` returned `alwaldend.github.io`.
The preserved provider record retains TTL 600 and proxying disabled.

The baseline source and ownership-linter checks cover unchanged canonical
declarations. Successful adapter builds, source shape/hash checks, and the
adoption plan cover activation. Final aggregate source validation remains with
repository delivery. Raw plans, inventories, and credentials remain in private
ignored task scratch.
