## Context

The team SDK can resolve an owned team by numeric ID from personal-account
key metadata. Existing creation stores a generated management credential and
a private creation-confirmation flag; failed creation must never adopt an
existing same-name team.

## Decision review

**Proceed** with read-only identity import, keeping imported ownership separate
from confirmed creation. It satisfies adoption without live resource changes.
Generating a replacement management credential would expand live mutations;
fabricating credential values would corrupt lifecycle state. Reporting missing
import support alone does not satisfy the requested team reconciliation.

## Decisions

Import resolves a canonical numeric ID through the existing SDK and records
the name, account ID, null credentials, and a separate private imported flag.
Refresh handles that flag independently from failed creation. Verified external
absence removes imported state because no generated key needs cleanup.

Imported-team deletion fails before mutation and retains state. The caller
also uses prevent_destroy. Detach ownership with Terraform's explicit state
removal workflow when necessary. Imported teams cannot supply credentials to
project-key creation; their key outputs remain null. No existing creation
behavior or SDK API needs to change.

## Risks / Trade-offs

- Imported teams have identity tracking but no destructive lifecycle without
  management credentials. Document that limitation explicitly.
- Read failures must not be mistaken for absence. Test malformed metadata and
  unavailable endpoints, plus missing ownership and invalid IDs.
- Failed creation and import must remain distinct. Re-run the existing
  ownership fixture and add import scenarios before implementation.

## Validation

Use the real Terraform CLI fixture. Assert unchanged plans, retained state
after refused destruction, external absence, and GET-only API traces. Preserve
a JSON operation report as a repeatable test artifact.
