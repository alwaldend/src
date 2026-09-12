## Context

See [proposal](proposal.md) and the
[shared module change](../../../../../tf_modules/openspec/changes/archive/2026-09-13-add-dns-records-module/proposal.md).

## Decisions

The `tf` service stage consumes the project JSON directly through the shared
global module. Its Bazel data includes the owner JSON and module sources so
source-checkout and runfiles execution use the same inputs. AL selects the
project AppRole, its own HTTP backend path, and Cloudflare injection with the
stage's `tf=main` label. The public-only module does not require RouterOS.

Record ownership defaults to disabled, while the optional zone input defaults
to null. These source defaults do not bypass provider setup or authentication.
Operational use requires the AppRole and Vault provider fields documented in
the [migration prerequisites](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prerequisites-and-order).

## Migration Plan

The [DNS migration owner](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/design.md)
owns central-writer freezing, source activation, imports, and live acceptance.
Provider import identifiers, deployed ownership, and live parity remain
unobserved; these artifacts prepare the source integration only.

## Validation

Earlier strict OpenSpec and source-reference checks are historical evidence.
This revised integration awaits the coordinator's combined configured
formatting, package builds, provider locks, and offline module validation.
No live Terraform or DNS operations are part of this work.

## Source validation

The coordinator passed the 53-target Terraform/skill package test batch and
built all 45 owning Terraform wrappers plus Vault and the infrastructure skill.
The runtime linter discovered 46 canonical files and printed 97 declaration
rows without ownership conflicts. Provider mocks cover both views; no live
provider, backend or Vault operation was run. This change completes source
preparation; operational adoption remains in the DNS coordination change.
