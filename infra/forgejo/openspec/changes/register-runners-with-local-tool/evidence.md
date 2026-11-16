## Offline acceptance

`//infra/forgejo/cmd/runner_registration:acceptance_test` passes with the actual
Go executable, production Terraform provisioner declaration, built-in Terraform
resources, and private TLS Forgejo/Vault fixtures. Its sanitized receipts are
`runner-registration.json` and `terraform-bootstrap.json` in Bazel undeclared
test outputs.

The fixture verifies one creation and one Vault write, Vault-cache reuse without
a Forgejo request, no mutation during plan, no rerun after changing the executable,
explicit replacement reusing credentials, removal retaining external objects,
same-name collisions across pages, lost-response and failed-save retries,
permission failures, invalid inputs, trusted TLS, redirect refusal, and credential
exclusion from command output and fixture state. These observations concern
isolated dummy services, not deployed service health.

## Packaging and checks

The service Terraform launcher's declared inputs and runfiles manifest include
the executable at the path used by its explicit interpreter. The packaged binary's
help command exits successfully. The component documentation and canonical
infrastructure skill build. Terraform formatting, owner-local specification
validation, offline skill evaluation configuration and affected semantic lint
pass. Final root quality and candidate review are recorded in the task-local
delivery receipts and review report.

The initial implementation phase performed no live plan, apply, registration,
Vault write, host change, credential rotation, or deletion. The existing writer AppRole's namespace is reused;
worker deployment and its reader policy are separate ownership concerns.

## Authorized deployment on 2026-11-16

The user explicitly requested application of the existing implementation.
The active Vault OAuth2 source was verified as ID 2 through the existing
maintenance SSH target. The full service apply stopped during planning because
an old checkout mirror and runner account would be removed despite their
prevent_destroy guards. No deletion was approved.

Targeting the runner while depending on managed repository resources still
included the protected mirror, including with a literal instance reference.
The final source derives the existing repository identity from the shared
catalog and leaves service-resource reconciliation independent.

The owning interactive apply with target
`terraform_data.runners["secure1"]` displayed one addition, zero changes and
zero deletions. After review, it was approved and Terraform reported completion
with those counts. Forgejo runner ID 2 was created for `alwaldend/src`, named
`secure1.dev.alwaldend.com`. The tool reported its initial Vault write, and a
separate owning Vault read verified matching identity metadata, UUID/token
presence and KV version 1 without displaying either credential. The subsequent
owning runner-only plan reported no changes. Runner-host deployment was not run.

Both completed Terraform operations were followed by the existing Forgejo login
plugin's cleanup error, `Forgejo logout status: 303`, causing the wrapper to exit
with status 1. This does not reverse the completed registration or saved state;
the independent verification establishes their persisted result. Full service
reconciliation and the logout response handling remain separate issues.

## Session ergonomics

- `terraform-bootstrap-invalidation`: user corrections established creation-only
  semantics. The maintained infrastructure reference and evaluation cases now
  explicitly omit timestamp and tool-content replacement triggers. Fixture tier:
  changed executable plus repeat apply performs no registration.
- `terraform-target-discovery`: one guessed validate-test label did not exist;
  owner target discovery established the available formatting check. No live
  operation was substituted for missing offline validation.
- `runfiles-manifest-inspection`: the host emits manifest-only build outputs.
  Inspect the declared manifest entry and executable rather than assuming a
  materialized runfiles directory. The full provisioner path is also exercised
  in the isolated Terraform fixture.

PR #312 merged during implementation. The follow-up branch starts at current
master `57d30f1ef751a1722ad37d85698f65618669c216`; earlier delivered work is
preserved by that merge. Publication identity belongs to the delivery receipt.
