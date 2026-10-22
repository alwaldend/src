## 1. Export library

- [x] 1.1 Write native E2E cases before implementation for borrowed lifetime, exact configs/state, naming, refresh and failure publication; observe the missing API fail.
- [x] 1.2 Add a connection-injected Exporter reading native protobufs directly through generated clients with no reconcile/app/lifecycle dependency.
- [x] 1.3 Preserve stable filenames, native member bindings, filtered-card exclusion, preset aliases, TOML/JSON state and direct filesystem output.

## 2. Verification and delivery

- [x] 2.1 Document the prepared-connection API and extraction boundary; run native E2E and inspect emitted config/state/package artifacts.
- [x] 2.2 Validate OpenSpec, format and select exact-candidate quality, semantic lint and workflow/docs gates. Receipts under out/export-extraction own publication state.

## Direct filesystem export

- [x] Extend existing native E2E before implementation for direct overwrites, retained files/directories, both state formats and partial output after filesystem errors.
- [x] Replace staged refresh/publication with config.Export and one encoded state-file write. Remove staging, copying, stale cleanup, chmod traversal, backups and rename/rollback helpers.
- [x] Document ordinary filesystem-write semantics and caller-owned temporary export/synchronization; verify native workflows, outputs and strict OpenSpec, then select exact-candidate delivery gates. Receipts under out/simplify-export own final validation/publication state.

## Borrow the caller's state store

- [x] Update native export callers before implementation, then accept a prepared state.Store in New alongside the connection; App chooses the namespace and supplies that store.
- [x] Update API docs/specs and select native workflow, representative-output and exact-candidate delivery checks. Track review replies and final validation/publication results in out/simplify-export/store-injection receipts.

## Configuration-only export

- [x] Update native library and CLI workflows before implementation: no generated state file, retained caller state files, separate state setup/dump and stable names.
- [x] Remove export state projection, preview, serialization and format parameters/flags; write configs only through config.Export.
- [x] Update current API/docs/specs and bootstrap examples for separate explicit state operations.
- [x] Format and select exact-candidate native workflow, representative-output, quality, lint and docs checks; out/config-only-export receipts own final validation/publication evidence.

## Shared native resource prerequisite

- [x] Extract native client construction and snapshot reads into resource.Store; keep config projection, state and reconciliation policy in their owning packages. Borrow resource and state stores at operation constructors.
- [x] Update API/dependency docs, verify native integration and select exact-candidate gates. Keep identical resource-package prerequisite overlap until its separate master-based PR merges; delivery receipts own final results.

- [x] Include the selected native deck ID in shared preset-read errors while retaining the wrapped cause; select existing native workflow and exact-candidate gates for the correction.

## Rebase after resource-store merge

- [x] Rebase onto master after #161, retain the merged resource store unchanged and remove its prerequisite files from this PR diff. Resolve only duplicate documentation and the server-consumer visibility entry.
- [x] Refresh prerequisite references and select exact-candidate quality, semantic lint, native E2E and documentation checks. Receipts under out/rebase-161 own final validation and publication evidence.

## Native filename text and simplification review

- [x] Extend native export E2E before implementation for inline markup, line breaks, audio/media, Unicode and bounded safe filenames; observe the custom parser disagree with Anki.
- [x] Reuse Anki's HtmlToTextLine RPC through generated clients, remove custom HTML/audio parsing and retain filesystem naming only.
- [x] Review the full master-based export diff for correctness, simplification and removable APIs; fix in-scope blockers and record remaining findings.
- [x] Verify representative artifacts and strict OpenSpec; select exact-candidate quality, semantic lint, native workflow and docs checks. Receipts under out/native-export-names own final validation and publication evidence.

## Named validation character sets

- [x] Replace inline NUL-containing validation strings with documented private character-set constants in their owning packages; retain each rule's original character set and keep literal fixture data unchanged.
- [x] Inspect the behavior-preserving diff and select existing native workflow, quality, semantic lint and documentation gates. Receipts under out/validation-characters own final validation and publication evidence.

## Export data-flow diagram

- [x] Add the requested Exporter, resource/state stores, gRPC connection, Anki storage/State and Configs diagram with gRPC and Write edges; embed it in the export Go API documentation.
- [x] Render and inspect light/dark SVG and a dark WebP preview, then select exact-candidate quality, semantic lint and documentation checks. Receipts under out/export-diagram own final validation and publication evidence.
- [x] Add Resources beside State inside Anki storage and remove redundant read labels; retain Configs and inspect the updated render. Receipts under out/export-resources own validation and publication evidence.

## Review: projection stages

- [x] Split exportConfigs into private deck, note-type and note projection stages, keeping shared references local to each export and preserving native inputs, naming and output behavior.
- [x] Verify the existing native export and consumer workflows, inspect emitted artifacts, then select exact-candidate quality, lint and docs gates. Publication and the reasoned review reply are tracked under out/export-review-20261017.
