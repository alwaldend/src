## Source implementation

- [x] Inventory every canonical source and the owning Terraform stage.
- [x] Verify pinned provider types, TTLs, multiplicity and import identifiers.
- [x] Implement shared normalization and provider modules with offline tests.
- [x] Prepare all 45 roots and 33 missing AppRole modules with least-privilege reads.
- [x] Discover DNS declarations at runtime, print their table, and reject
      multiple files owning the same canonical name across types and views.
- [x] Remove exporter, central ownership registry and central DNSControl
      entrypoints; label retained snapshots as historical.
- [x] Record baseline implementation formatting, package/output checks and strict source specifications in the validation evidence.

Final rollout candidate quality checks remain pending repository delivery.
Git delivery receipts and the pull request own exact-candidate validation and publication state.

## Authorized operational adoption

The user authorized deployment one project at a time with existing records and
services preserved. Track adoption in each owner's change; `projects/sri` is
the pilot.

- [x] Provision reviewed missing AppRoles through owning IaC and verify existing roles' scoped access and authenticated provider discovery.
- [x] Verify independent authenticated provider, Vault and backend recovery paths.
- [x] Audit central writers and record coverage, import the pilot owner's exact
      record IDs, and verify its no-change plan and provider preservation.
- [x] Complete all remaining owners with independent receipts, exact existing-record adoption or verified missing-record creation, and preservation of every preexisting provider record; retain rollback evidence before archiving this change.
