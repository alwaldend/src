## MODIFIED Requirements

### Requirement: Preserve live records during transfer

Operational cutover SHALL audit active and scheduled central writers, stop
identified competing writers, and coordinate one owner at a time. Evidence
SHALL state the audit's coverage and unavailable observations. Each existing
record SHALL be imported using its actual provider ID into exactly one owner
state, followed by a no-change adoption plan before writes. Fresh complete
inventories MAY establish missing declarations for an exact additions-only
plan that preserves every existing record. Endpoint-owning batches SHALL
verify an independent authenticated recovery path first.

The shared apex and mail records owned by `infra/dns` SHALL remain disabled
by default in its `tf` root before the authorized adoption revision and SHALL
retain enabled ownership after import. Reconciliation against the adopted state
and unchanged declarations SHALL propose no record additions, changes,
replacements, or deletions and SHALL preserve unrelated provider records.

#### Scenario: An owner is adopted

- **WHEN** its authorized cutover runs
- **THEN** exact-ID imports preserve existing records
- **AND** the operator coordinates exclusive active writers and verifies
  unrelated records remain unchanged

#### Scenario: Deploy verified missing declarations

- **WHEN** complete provider inventories establish that declared names are absent
  and conflict-free
- **THEN** the reviewed DNS plan creates only those missing declarations
- **AND** every existing provider identity and attribute remains unchanged
- **AND** a follow-up DNS plan contains no changes

#### Scenario: Inspect the prepared shared-record defaults

- **WHEN** the shared `tf` root uses its checked-in defaults before its
  authorized adoption revision
- **THEN** shared-record ownership is disabled and the root retains the
  canonical shared declaration and module inputs

#### Scenario: Inspect adopted shared-record defaults

- **WHEN** the adopted shared `tf` root uses its checked-in source defaults
- **THEN** ownership is enabled for the shared apex and mail records in
  their declared views

#### Scenario: Reconcile existing shared records

- **WHEN** the shared `tf` root plans against adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
- **AND** unrelated provider records remain unchanged
