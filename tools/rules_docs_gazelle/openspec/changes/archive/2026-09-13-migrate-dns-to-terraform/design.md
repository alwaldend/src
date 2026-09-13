## Context

See [proposal](proposal.md). The [project BUILD](../../../../BUILD.bazel) exports
[DNS declarations](../../../../dnsconfig.json) from a standalone Bzlmod module.

## Decisions

Keep canonical Terraform files in the project's `tf` directory and export
those files through Bazel. Root-workspace targets under `infra/dns/projects`
provide `rules_docs_gazelle_tf` and `rules_docs_gazelle_vault` operations. This avoids introducing a
parent-monorepo dependency into the reusable module while keeping one source
for project infrastructure. A central shared state would violate project
ownership; copying the Terraform source would create competing definitions.

This project declares only global records, so its root consumes the shared
`dns_records/global` module and requires only the Cloudflare provider.
Use the existing declaration as its input and an HTTP backend for project
state. The [AL configuration](../../../../al.lua) injects provider credentials
through the normal `tf=main` label. The project's AppRole and required Vault
credentials must be deployed before invoking its authenticated root commands.

The [runtime DNS linter](../../../../../../infra/dns/README.md) discovers each
owner's canonical declarations and prints its table directly. It does not
supply Terraform input or maintain a second declaration source.

## Risks / Trade-offs

Packaging crosses a Bazel repository boundary, so validate the produced
runfiles and canonical relative source paths. The initial `dns_enabled`
default is false. This suppresses DNS record creation; it does not bypass
AppRole authentication, backend access, or provider credential injection.
Prepared source is not evidence of live adoption.

Current source removes central DNSControl commands. Before imports or
reconciliation, freeze all older central deployment candidates as required by
the [DNS operating guide](../../../../../../infra/dns/README.md) and
[cutover runbook](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md).
The runbook owns adoption, no-change plan review, and rollback requirements.

## Implementation state

The canonical Terraform root, provider lockfile, AL configuration, source
exports, and root operational wrapper are present. Earlier source checks
verified the original declaration digest, module input, disabled default,
and project role and backend path. The revised AL integration still needs
its package checks and representative runtime output verified against the
candidate; record that evidence in [tasks](tasks.md). No live adoption is
established by this implementation.

## Source validation

The coordinator passed the 53-target Terraform/skill package test batch and
built all 45 owning Terraform wrappers plus Vault and the infrastructure skill.
The runtime linter discovered 46 canonical files and printed 97 declaration
rows without ownership conflicts. Provider mocks cover both views; no live
provider, backend or Vault operation was run. This change completes source
preparation; operational adoption remains in the DNS coordination change.
