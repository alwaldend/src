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

No live plan, apply, registration, Vault write, host change, credential rotation,
or deletion was performed. The existing writer AppRole's namespace is reused;
worker deployment and its reader policy are separate ownership concerns.

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
