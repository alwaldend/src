## 1. Identity declaration

- [x] 1.1 Declare the Vast.ai AppRole and operator bootstrap membership through
      the shared module; verify the packaged root and Terraform formatting
      without accessing live Vault.

## 2. Integration validation

- [x] 2.1 Run Vault owner OpenSpec validation and affected semantic lint;
      inspect that no broad group membership or literal secret was added.

## Validation evidence

On 2026-11-04, Terraform formatting and the Vault owner OpenSpec validator
passed alongside repository quality checks (30 selected tests). The plan
launcher built with the AppRole source in its runfiles; affected semantic
lint passed. The declaration uses the shared module and existing operator
entity, disables Yandex policy attachment, and adds no broad group membership.
No live plan, apply, credential generation, or secret write was performed.
