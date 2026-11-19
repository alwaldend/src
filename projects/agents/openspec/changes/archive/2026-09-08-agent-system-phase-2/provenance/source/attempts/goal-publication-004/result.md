# goal-publication-004: legacy-migration

## Outcome

This attempt closes with `refine`. It completes the `legacy-migration`
criterion: migration now retains a workspace-relative source path, the
imported-byte content digest, and explicit mapping provenance
(`mappingVersion: v1` with extracted or overridden title/criteria), and
idempotence is rebound from the frozen imported snapshot so copied provenance
fields can no longer misclassify a retry as the completed import.

## Implemented

- `MigrationStatus` gains `sourcePath`, `mappingVersion`, and
  `extractionMode` (`projects/goal/api/v1alpha1/types.go`).
- Validation requires populated migration status to carry a normalized
  workspace-relative source path, canonical `sha256:` digest, mapping
  provenance version `v1`, extraction mode `extracted` or `overridden`, and
  RFC3339 `migratedAt` (`projects/goal/api/v1alpha1/validation.go`).
- `Migrate` records the portable source path, mapping version, and
  extraction mode; `matchExistingMigration` verifies them and rebinds the
  source digest from the frozen `imported-unversioned` attempt's plan and
  evidence bytes (`projects/goal/internal/fsstore/migrate.go`).
- Migrate CLI long help documents the retained provenance
  (`projects/goal/cmd/goal/command.go`).

## Verification

- `bazel_agent test //projects/goal/api/v1alpha1:all //projects/goal/internal/fsstore:all` — 2/2 pass.
- `bazel_agent build //projects/goal/...` — success.
- New tests: digest rebinding rejection, extraction-mode recording, source
  path/mapping provenance assertions, and negative validation cases.

## Not in scope

- Reward: `bounded-catalogs` (remaining slices), `system-index`,
  `context-capsule`, `runtime-isolation`, `resource-baseline`.
- Migration crash-residue cleanup, atomic no-replace rename primitive, and
  post-check target-race guarantees remain documented open items in the goal
  evidence.
