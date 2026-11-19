## Context

`infra/forgejo/dnsconfig.json` remains the declaration source. `tf_setup` owns DNS state
through the shared `dns_records` module and the existing Vault HTTP backend flow.
The [DNS migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/design.md) owns shared semantics and cutover.

## Goals / Non-Goals

Package equivalent owner records, explicit provider dependencies, and an offline
format check. Live adoption, credential provisioning, and cutover remain separate
operator work under the shared change.

## Decisions

- Call the shared module from `tf_setup` and default `dns_enabled` to false.
- Use `src_infra_dc1_forgejo1` for DNS authentication and inject only secret references.
- Scope RouterOS DNS credentials to the `dns` provider alias so unrelated provider
  operations keep their current credentials.
- Retain the existing AL stage labels and backend paths. DNS credential injection
  selects the existing Terraform stage label, with real Vault credentials required
  before operational use of this revision.

## Risks / Trade-offs

New roles, policies, and credential fields require separately authorized bootstrap
before operational Terraform calls. `dns_enabled=false` prevents DNS resource
creation; it does not prevent provider configuration or live authentication.
Once records are imported and management is enabled, keep it enabled; reverting
to false would propose deletion.

## Migration Plan

Follow the shared DNS adoption procedure. Provision prerequisites, import existing
record identities, review the owner plan, and enable management.

## Validation

Source inspection on 2026-09-12 verifies the canonical DNS declaration SHA-256
against the shared migration inventory, the shared module source resolves to its
repository owner, and the declaration is packaged in runtime data.
`git diff --check` passes. Package format, package build, provider-lock updates,
and repository format checks remain pending.

## Open Questions

None for offline implementation; live adoption is outside this change.

## Source validation

The coordinator passed the 53-target Terraform/skill package test batch and
built all 45 owning Terraform wrappers plus Vault and the infrastructure skill.
The runtime linter discovered 46 canonical files and printed 97 declaration
rows without ownership conflicts. Provider mocks cover both views; no live
provider, backend or Vault operation was run. This change completes source
preparation; operational adoption remains in the DNS coordination change.
