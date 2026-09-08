# legacy-migration evidence

## Criterion

> Legacy goal migration creates a fresh valid identity, retains source path,
> digest, and raw bytes, maps fields explicitly, and preserves unmapped prose
> without modifying the source.

## Claim

The migration now retains workspace-relative source path, the imported-byte
content digest, and explicit mapping provenance (`mappingVersion: v1`,
`extractionMode` extracted/overridden). Idempotence is rebound from the
frozen imported attempt snapshot (plan + evidence bytes), so a records whose
imported bytes no longer derive the stored digest is refused. Fresh identity,
full-byte preservation, and non-destructive source handling were already
provided and are still covered by the existing suite.

## Implemented

- `projects/goal/api/v1alpha1/types.go`: `MigrationStatus` gains
  `sourcePath`, `mappingVersion`, `extractionMode`.
- `projects/goal/api/v1alpha1/validation.go`: populated migration status
  requires an unversioned source format, a non-empty normalized
  workspace-relative `sourcePath` that stays under the workspace root, a
  canonical `sha256:` source digest, mapping provenance
  (`mappingVersion: v1`, extraction mode `extracted` or `overridden`), and an
  RFC3339 `migratedAt`.
- `projects/goal/internal/fsstore/migrate.go`:
  - `SourcePath` is computed as the workspace-relative portable reference
    (`portableOwnerRoot`) to the unversioned source directory.
  - `MappingVersion` is fixed to `v1`; `ExtractionMode` records whether title
    or criteria were extracted from the README or supplied as overrides.
  - `matchExistingMigration` now also verifies `sourcePath`, mapping
    provenance, migration timestamp, and rebinds the legacy digest from the
    frozen `imported-unversioned` attempt: `plan.md` is the README and
    evidence/ files are the remaining root Markdown bytes. A valid record
    whose imported snapshot no longer derives the stored digest is refused
    with `existing migration imported snapshot does not match provenance`.
- `projects/goal/cmd/goal/command.go`: migrate long help documents retention
  of source path, digest, mapping version, and extraction-vs-override.

## Tests

- `store_test.go::TestUnversionedMigrationIsNonDestructiveAndIdempotent`
  asserts `SourcePath` (`out/task/legacy/legacy-goal`), `MappingVersion`
  `v1`, `ExtractionMode` `extracted`, and source digest on the migrated
  record while the legacy source remains byte-identical.
- `migrate_test.go::TestMigrationIdempotenceRebindsDigestFromImportedSnapshot`
  morphs the frozen imported plan into different bytes and updates the
  attempt manifest so the record stays structurally valid; a retry is then
  refused because the rebound digest no longer matches provenance.
- `migrate_test.go::TestMigrationRecordsExtractionMode` asserts
  `extractionMode: extracted` for an implicit source and `overridden` when
  title/criteria are supplied.
- `validation_test.go::TestMigrationStatusRequiresPathDigestAndMappingProvenance`
  rejects absolute, empty, parent-escaping source paths; missing or unknown
  `mappingVersion`; unknown extraction mode; and missing/malformed digest.

## Validation run

- `bazel_agent test //projects/goal/api/v1alpha1:all //projects/goal/internal/fsstore:all` — 2/2 pass.
- `bazel_agent build //projects/goal/...` — success.
- `goal validate` — passes on the committed goal records including the new
  frozen attempt.

## Remaining open (not claimed)

- Crash-staging residue cleanup, atomic no-replace directory rename
  primitive, and post-check target-race guarantees remain documented open
  items in the goal evidence; this attempt does not claim them.
