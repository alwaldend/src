# goal-publication-004: legacy-migration

## Goal

Complete the `legacy-migration` criterion (r1, required) for
`agent-system-phase-2`.

## Criterion statement

Legacy goal migration creates a fresh valid identity, retains source path,
digest, and raw bytes, maps fields explicitly, and preserves unmapped prose
without modifying the source.

## Gap analysis (from goal-publication-001 evidence)

Migration already derives a fresh identity, keeps a full-directory staging
copy, preserves README as `plan.md` and other Markdown files as evidence, and
never writes the source. Remaining defects:

1. **No portable source reference.** `MigrationStatus` has only
   `sourceFormat`, `sourceDigest`, `migratedAt`.
2. **Field mapping is not explicit.** No record of title/criteria
   extraction-vs-override or the mapping version.
3. **Idempotence is weakly bound.** `matchExistingMigration` compares stored
   metadata but does not recompute the legacy digest from the imported
   attempt's `plan.md` and evidence.

## Plan

### API

- Add `SourcePath` (json/yaml `sourcePath`) to `MigrationStatus`:
  normalized workspace-relative portable path to the unversioned source
  directory, never a host-absolute path.
- Add `MappingVersion` (json/yaml `mappingVersion`) and
  `ExtractionMode` (json/yaml `extractionMode`) to record the explicit field
  map provenance: parser version plus `extracted` / `overridden` per axis?
  The criterion asks "maps fields explicitly"; keep the migration record
  inspectable without over-engineering.

### Migration

- Compute a portable `SourcePath` in Workspace-relative form via
  `portableOwnerRoot`-style validation and store it in
  `MigrationStatus.SourcePath`.
- Set `MappingVersion` to a fixed `"v1"` and `ExtractionMode` from whether
  title/criteria were extracted (`extracted`) or supplied (`overridden`).
- In `matchExistingMigration`, recompute the legacy digest from the imported
  attempt's frozen artifacts (`plan.md` as README plus evidence basenames and
  bytes) and require it to equal the stored `SourceDigest` before accepting
  idempotent success.

### Validation

- `MigrationStatus.validate` requires `SourcePath` to be a normalized
  non-empty workspace-relative portable path when migration is populated, and
  requires `MappingVersion` and `ExtractionMode` to be non-empty (bounded
  set) when populated.
- Empty migration status remains valid (fresh goals without migration).

### Tests

- Store migration integration asserts `SourcePath`, `MappingVersion`,
  `ExtractionMode` values after migration.
- Validation unit tests cover populated and empty `MigrationStatus`,
  absolute-path rejection, missing mapping provenance rejection.
- Idempotence test asserts digest rebinding: a valid target whose imported
  attempt plan/evidence bytes differ from the stored source digest is
  rejected even if the metadata fields match.
- CLI docs (migrate command long help) mention source path and mapping
  provenance retention.

## Scope

- `projects/goal/api/v1alpha1/types.go`
- `projects/goal/api/v1alpha1/validation.go`
- `projects/goal/api/v1alpha1/validation_test.go`
- `projects/goal/internal/fsstore/migrate.go`
- `projects/goal/internal/fsstore/migrate_test.go`
- `projects/goal/internal/fsstore/store_test.go`
- `projects/goal/cmd/goal/command.go`

Not in scope: crash-residue cleanup promotion, atomic no-replace rename
primitive, and post-check target-race guarantees (documented as open in the
goal evidence; this attempt does not claim them).
