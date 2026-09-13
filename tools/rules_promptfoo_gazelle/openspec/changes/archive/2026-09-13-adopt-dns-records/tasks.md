## 1. Access and preparation

- [x] 1.1 Verify or provision the dedicated AppRole and scoped DNS access through the owning Vault workflow, and verify authenticated backend and provider access.
- [x] 1.2 Review the [shared writer audit](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md#writer-audit) and its coverage limits, complete this owner's preparation under the shared cutover procedure, and record exact matches for its existing records.
- [x] 1.3 Activate the project's source default in its reviewed adoption revision, and verify its source checks and runtime ownership lint output.

## 2. Adoption and evidence

- [x] 2.1 Import the matched provider IDs through `//infra/dns/projects:rules_promptfoo_gazelle_tf.import`, and verify a plan with no record additions, changes, replacements, or deletions.
- [x] 2.2 Apply the reviewed plan through the owning adapter, and verify DNS resolution, unrelated record preservation, and no-change reconciliation.
- [x] 2.3 Record sanitized adoption evidence in [design](design.md), and verify the specification delta against the adopted source before synchronization and archive.
