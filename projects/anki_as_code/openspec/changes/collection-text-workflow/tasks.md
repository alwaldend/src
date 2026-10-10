## 1. CLI

- [x] 1.1 Write E2E failure cases before implementation for round trips, content/deck edits, duplicate identities, additions/deletions, convergence, wrong bases, and malformed text.
- [x] 1.2 Implement export/plan/apply/build/generate-id and document usage; verify E2E tests and CLI help through Bazel.

## 2. Representative workflow

- [x] 2.1 Verify the CLI against the supplied collection; user-owned export publication belongs to its separate OpenSpec change and PR.
- [x] 2.2 Verify the rebuilt archive and text round trip; collection acceptance evidence is maintained in users/simeonwarren/anki/openspec/changes/export-organized-collection.

## 3. Delivery

- [x] 3.1 Validate OpenSpec, review formatted source and representative output, and select aggregate delivery scope and validation gates. The repo-delivery receipt owns subsequent exact-candidate checks, commit/push, and PR publication state.

## Verification evidence

The SQLite v1.60.1 E2E suite passes, including exact field round trips,
compact single-line TOML, multiline escaping, filtered cards, additions and
deletions, restored identities without stale tombstones, and repeated apply.
Independent collection comparisons and cardinalities belong to the collection
owner's export change. Raw logs, archive output, and delivery receipts remain
in ignored out/anki-as-code. The original attachment is unchanged.
An actual Anki desktop import has not been performed.

## Session ergonomics

- GO-SOURCE-LAYOUT (live user feedback): compressed Go was rejected; source now
  uses the repository Go formatter and focused loading/validation/reconciliation
  helpers split by responsibility.
- SQLITE-UPGRADE-FRESHNESS (live registry evidence): v1.26.0 was selected without
  checking freshness; replaced with user-approved v1.60.1. The separately owned
  dependency skill now prefers verified current stable releases.
- RESULT-OUTPUT-BOUNDING (observed tool output): printing the delivery status
  array produced a truncated result of roughly 256,000 tokens. Subsequent
  inspection selects individual diagnostic fields and counts.

## Review corrections

- [x] 4.1 Reproduce missing empty notes directories and filtered home destinations before fixes.
- [x] 4.2 Treat absent directories as empty, create them for first notes, and require normal card homes; verify E2E convergence and rejection.
- [x] 4.3 Verify plan/build from a clean Git materialization and select the correction in the aggregate delivery receipt; that receipt owns subsequent validation and publication state.

Clean Git materialization reports no planned differences and builds the archive
successfully without any untracked empty notes directories.

- [x] 4.4 Reproduce stale allocation and per-card sibling positions, and archive mode loss before fixes. Preserve runtime nextPos, allocate once per note, retain archive modes, and verify E2E convergence.

- [x] 4.5 Reproduce and correct runtime cursor storage: write its JSON as an Anki config BLOB and verify its SQLite storage class in the allocation E2E case.

- [x] 4.6 Reproduce incomplete media caches, wrong lesson preset inheritance, and unbounded direct SQLite inputs; preserve all Anki media filenames, inherit the nearest normal parent, reserve declared deck IDs during allocation, enforce the database limit, and verify E2E archive outputs.

- [x] 4.7 Split user collection data and large-inventory delivery support into an independent master-based export PR; keep project and dependency changes in PR #128.

## 5. Documentation and landing

- [x] 5.1 Correct README metadata, supported inputs, command workflow, and documentation ownership.
- [x] 5.2 Register the project landing and package README documentation through the existing main Hugo site.
- [x] 5.3 Build the site and inspect the project index, landing, documentation, and taxonomy links; select the new exact-candidate delivery gates.

## 6. Cobra CLI review correction

- [x] Replace manual flag parsing with existing pinned Cobra, using command-specific flags.
- [x] Verify command help, flag validation, and E2E collection behavior before delivery.

CLI command verification passed for root/subcommand help, scoped and required
flags, unexpected arguments, and unknown commands. The Cobra CLI reports an
empty plan against the rebuilt collection and exported desired state.

## 7. Identity workflow review corrections

- [x] Write nested/overlapping path and invalid-batch E2E cases before implementation.
- [x] Replace new-note with generate-id on existing files/directories and remove the one-off lesson command and production organizer.
- [x] Verify copied notes rebuild and converge while unselected scheduling is preserved; update README, landing, and maintained specifications.

The identity E2E emits generated-identities.txt. The export owner retains the
completed one-off lesson grouping; it is not a CLI capability.

## 8. Archive review corrections

- [x] Write large SQLite/archive and field E2E cases before removing fixed size cutoffs.
- [x] Stream database extraction/copying without a total-size cutoff and allow large TOML fields; retain schema validation.
- [x] Name the archive database entry and metadata constants; trace error/success cleanup ownership for the archive-close question.
- [x] Verify large collection artifacts and select current README/landing and exact-candidate delivery checks.

Large-collection E2E passes for 513 MiB SQLite and streamed Zstandard archive
inputs and a 17 MiB field. It emits large-database.txt and large-field.txt, and
all rebuilt outputs converge. These replace the previous fixed-size rejection
contract at the user's request.

## 9. Config-driven discovery and settings references

- [x] Write rearranged-layout and card-appearance E2E failure cases before implementation.
- [x] Implement explicit root/deck/note markers, configurable recursive discovery, and identity reservation.
- [x] Export, plan, validate, and reconcile card appearance with lossless strings and protobuf preservation.
- [x] Migrate the separate user export, update documentation/landing, and verify convergence and representative settings.
- [x] Select both independent owner delivery receipts and exact-candidate validation gates; those receipts own subsequent validation and publication state.

Config-driven E2E artifacts verify arbitrary relative paths, external identity
reservation, recursive notes, CSS/template edits, unknown protobuf fields,
unchanged card scheduling, and convergence. Independent supplied-collection
verification confirms exact unchanged notetypes/templates/deck_config tables.

## 10. Preserve unmanaged review history

- [x] 10.1 Exclude history from text while preserving all review rows during apply/build.
- [x] 10.2 Verify round-trip and edited archives preserve exact review rows; rebuild the supplied collection and verify its original history.

## 11. Readable literal TOML strings

- [x] 11.1 Verify HTML quotes, backslashes, delimiter conflicts, controls, and newline boundaries through archive E2E cases before implementation.
- [x] 11.2 Prefer lossless literal strings in note fields and card appearance, with basic-string fallbacks; update the README.
- [x] 11.3 Regenerate the separate export without changing field values or collection behavior and verify both candidates.

## 12. Consistent deck quoting

- [x] 12.1 Extend the deck export/rebuild E2E to cover quoting and exact description values before implementation.
- [x] 12.2 Apply the shared quoting policy to all deck string values, including generated-identity rewrites.
- [x] 12.3 Regenerate deck files and verify parsed settings, archive state, and convergence.

## 13. Separate note type documents

- [x] 13.1 Extend appearance E2E cases for per-note-type files, omitted collection settings, relocated/nested discovery, and missing/duplicate rejection before implementation.
- [x] 13.2 Export/load a note_types directory and leave collection settings unmanaged, preserving the computed new-card allocation cursor.
- [x] 13.3 Migrate the user export, update documentation, and verify appearance and archive state remain unchanged.

## 14. Group path references

- [x] 14.1 Verify grouped decks, notes, and note_types tables through relocated-root and archive E2E cases.
- [x] 14.2 Replace flat path properties with typed path tables and migrate the exported collection and documentation.
- [x] 14.3 Validate zero reconciliation drift, unchanged configuration/scheduling/media, and all 126,668 review records. Delivery receipts own subsequent candidate validation and publication.

## 15. Reconciliation review fixes

- [x] 15.1 Reproduce late archive changes, template-only sync markers, deck title collisions, and occupied filtered-deck conversions through behavioral workflows before implementation.
- [x] 15.2 Check the archive immediately before replacement, mark parent note types for template changes, stage deck renames without unique-index collisions, and reject occupied filtered-deck conversion.
- [x] 15.3 Verify preserved input/state, convergence, and repeatable artifacts; select the updated project delivery receipt and validation gates. The receipt owns subsequent validation and publication.

The behavioral run reproduced all four review defects before implementation.
The full collection suite then passed: late archive writes are retained and
replacement is rejected; template-only edits mark the parent for sync without
changing CSS or scheduling; title swaps and reuse of deleted titles converge
with Anki's unique name index; occupied filtered-deck conversion is rejected
without changing the input. Reports remain in out/anki-as-code/review-fixes-\*.

DELIVERY-LAUNCHER-EXPIRY (live): the retained task launcher referenced removed
Bazel runfiles. Regenerating it restored the owning delivery entry point. Its
cold build took about 153 seconds; the full collection test took about eight
seconds. Source and shared tooling were not changed to work around this.

## 16. Resolve card destinations by deck identity

- [x] 16.1 Reproduce swapped titles with explicit unchanged card titles, including filtered-card homes, before implementation.
- [x] 16.2 Compare existing and desired home-deck IDs in planning and reconciliation; preserve filtered membership and scheduling.
- [x] 16.3 Verify convergence and repeatable artifacts; select exact-candidate validation and publication through the delivery receipt.

The full collection suite passes after reproducing both missed and unnecessary
card moves. The deck-title artifacts record planned moves and resolved home IDs
for normal and filtered cards, preservation of filtered membership/original due,
and repeated-apply convergence. The delivery receipt owns aggregate checks and
publication; raw evidence remains in out/anki-as-code/card-deck-identity-\*.

## 17. Address outstanding human review threads

- [x] 17.1 Verify rebuilding against a changed compatible archive while preserving its unmanaged state; export configs without a hash and accept legacy hash metadata without binding.
- [x] 17.2 Remove manifest hash binding and manual TMPDIR reads; document why protobuf wire patching preserves Anki settings.
- [x] 17.3 Select exact-candidate validation/publication and guarded thread replies/resolution. The delivery and reply receipts own subsequent publication and review state.

The unbound-text E2E artifact verifies omitted export hashes, compatibility with
legacy metadata, preserved current scheduling/history/configuration, and empty
plans. Read-only operations and builds no longer hash archives; only apply
computes transient hashes for concurrent-write detection. The existing protobuf
wire dependency is retained because it preserves unknown Anki settings.

## 18. Rebase after the shared prerequisite merges

- [x] 18.1 Synchronize with master after PR #130 merges; preserve its shared dependency and skill records when resolving the two task-record conflicts.
- [x] 18.2 Verify that the Anki implementation matches the previous published candidate and that shared dependency and skill changes are absent from the project diff.
- [x] 18.3 Select exact-candidate validation and scoped delivery for the reduced project PR. Delivery receipts own check results, publication, and review observations.

The previous project head was 33cc18d264fc36a5b0f9b075c7afa49352d8f524;
the fetched master was 2e42480c2e78dc2c17b0e7e95ba1d15abbb11461.
The only conflicts were the two shared skill task records, resolved to their
merged delivery ownership. The project source was unchanged by the rebase.
Exact refs and bounded evidence remain under ignored out/anki-rebase.

## 19. Use Anki's upstream protobuf schemas

- [x] 19.1 Specify E2E preservation cases before implementation: note type/template unknown fields, known unmanaged fields, deck extra settings, empty/filtered configurations, edits, and convergence.
- [x] 19.2 Pin unmodified Anki schemas and generate Go types with the existing Bazel protobuf toolchain; replace manual wire handling throughout collection code.
- [x] 19.3 Select exact-candidate preservation, generated-schema consumer, formatting, and documentation checks through repo-delivery; its receipts own subsequent validation and publication state.

The preimplementation E2E reproduced a false deck change for reordered fields
and explicit defaults. Generated messages now preserve unchanged blobs and
retain optional zero limits, nested unmanaged settings, and unknown fields
during edits. Bounded evidence and delivery receipts are under ignored
`out/anki-schema`; the repeatable artifact is `schema-preservation.txt`.

## 20. Extract the upstream schema prerequisite

- [x] 20.1 Select a separate master-based dependency PR with only the unchanged Anki package, generated root include, and shared change record; bounded diff and byte comparison verify the split scope.
- [x] 20.2 Select exact-candidate checks and guarded PR projection linking the prerequisite; delivery receipts own subsequent check and publication state.

Following the user's earlier split preference, both PRs target master. The
project retains matching dependency files until the prerequisite merges,
so required CLI checks continue to pass. Rebase after that merge removes
the overlap. The independent dependency change record is
`infra/src/openspec/changes/extract-anki-upstream-schemas`; receipts remain
under ignored `out/anki-schema-pr-link` and the dependency worktree's
`out/anki-schema-split`.

## 21. Synchronize schema review fixes

- [x] 21.1 Synchronize #131's tagged archive, upstream Go import paths, and wrapper-only visibility; verify the prerequisite package remains identical across both branches.
- [x] 21.2 Update the CLI imports and Gazelle mappings, then select exact-candidate E2E, lint, build, and publication through delivery receipts.

## 22. Rebase after the schema prerequisite merges

- [x] 22.1 Rebase onto master after #131 merges and verify schema/module changes disappear from the project diff while CLI source remains unchanged.
- [x] 22.2 Update the aggregate PR description and select exact-candidate checks and publication; delivery receipts own validation, remote state, and review observations.

The previous project head was a6f2c7435ad644d3579a623a88b13b6eec74a996.
Master at 3d75eec5647be3dcaf55e6496c203bd8bc7076b8 contains merged #131.
Rebase was conflict-free; all project source and integration files matched
the previous candidate. The merged schema PR is unchanged. Bounded evidence
and delivery receipts remain under ignored `out/anki-rebase-current`.

## 23. Address API and structure review

- [x] 23.1 Specify repeated export and home-relative path behavior with a complete E2E workflow before implementation.
- [x] 23.2 Replace init registration with explicit database setup; centralize database reads and reconciliation in a repository type and reuse generated note-type structures.
- [x] 23.3 Co-locate deck methods with their type, remove the root filename restriction, retain standard WalkDir discovery, and use the standard random identifier generator.
- [x] 23.4 Implement repeatable exports and home expansion, simplify TOML string escaping through existing encoding support, and verify preservation/convergence artifacts.
- [x] 23.5 Publish separate Go skill guidance for explicit initialization and type/method placement; select exact-candidate delivery and reasoned review replies for both PRs.

## 24. Address CLI review

- [x] 24.1 Specify CLI stream and flag behavior with a process-level workflow before implementation.
- [x] 24.2 Send status messages to stderr and reuse flag definitions.
- [x] 24.3 Select exact-candidate validation, publication, and review replies through delivery receipts.

## 25. Address serialization and package boundary review

- [x] 25.1 Retain end-to-end coverage of field values, archive preservation, appearance edits, and convergence; remove obsolete literal/multiline formatting assertions before changing serialization.
- [x] 25.2 Dump complete TOML objects through the existing encoder and remove manual serializers; verify lossless E2E artifacts.
- [x] 25.3 Extract model, text storage, archive resources, and SQLite repository packages with an orchestration facade; keep appearance SQL in repository methods and verify the dependency graph and complete workflows.
- [x] 25.4 Select exact-candidate validation, guarded publication, and reasoned resolution of the three outstanding review findings through delivery receipts.

## 26. Rebase after the CLI scaffold merges

- [x] 26.1 Replay the task-owned implementation on merged #133, preserve pending review artifacts, retain the merged App-backed command boundary, and connect App methods to existing collection operations; verify the shared project scaffold disappears from the aggregate diff.
- [x] 26.2 Select exact-candidate workflow/lint/build checks and lease-protected publication through delivery receipts; receipts own subsequent validation and remote state.

## 27. Rebase after delivery tooling merges

- [x] 27.1 Replay the task-owned implementation onto merged #134 and verify the replay preserves the previous feature patch before continuing package-boundary corrections.
- [x] 27.2 Select the existing E2E preservation workflows, all extracted package lint targets, repository quality, OpenSpec, CLI build, and website build for exact-candidate delivery. Delivery receipts and the PR own subsequent validation, publication, and review-resolution outcomes.

## 28. Document architecture

- [x] 28.1 Add a high-level architecture section and authoritative Mermaid diagram of current internal package dependencies.
- [x] 28.2 Package the generated diagram with README documentation, inspect the rendered preview and built page, and select exact-candidate delivery checks.

## 29. Centralize Gazelle mappings

- [x] 29.1 Define upstream Anki import resolutions once in the project BUILD file, remove repeated package declarations, and verify inherited generation leaves target dependencies unchanged.
- [x] 29.2 Select exact-candidate checks and guarded publication/reply operations for both project-level mapping comments; delivery receipts own their subsequent outcomes.

## 30. Keep models declarative

- [x] 30.1 Move all functions and methods out of model into validation, planning, and Anki format packages without changing serialized types or reconciliation behavior.
- [x] 30.2 Update dependencies, README architecture and its diagram; inspect the rendered preview and select existing preservation/convergence workflows and representative-output verification.
- [x] 30.3 Select exact-candidate validation and guarded publication of the refactor; delivery receipts own subsequent outcomes.

## 31. Document model purposes

- [x] 31.1 Add Go doc comments explaining the purpose of every model type and relevant representation or preservation constraints.
- [x] 31.2 Select existing exact-candidate delivery checks; delivery receipts own validation and publication outcomes.

## 32. Clarify model contracts

- [x] 32.1 Audit model consumers and retain only serialized or inter-package types; move repository-only preservation bytes and scheduling cursors into private repository snapshot data.
- [x] 32.2 Apply review naming, co-location, and flat path feedback; adapt existing E2E workflows to the flat path keys.
- [x] 32.3 Select exact-candidate delivery checks and guarded review replies; delivery receipts own subsequent outcomes.

## 33. Exclude filtered cards from managed state

- [x] 33.1 Write E2E failure cases before implementation for filtered-card identity/ordinal conflicts and note deletion; adapt existing managed-card workflows and verify full filtered-card row preservation.
- [x] 33.2 Exclude filtered cards from export, planning, and reconciliation; protect their references through private repository validation shared by plan and apply.
- [x] 33.3 Update the managed/unmanaged contract, architecture, and exact-candidate delivery selection; receipts own subsequent validation and publication outcomes.

## 34. Document model fields

- [x] 34.1 Add purpose and representation comments to every field in internal/model, including the embedded upstream note type.
- [x] 34.2 Select configured formatting and exact-candidate delivery checks; receipts own subsequent validation and publication outcomes.

## 35. Define serialization through protobuf

- [x] 35.1 Define first-party TOML and plan contracts in a documented protobuf schema, generated with the existing toolchain; retain upstream Anki schemas for database configuration.
- [x] 35.2 Load TOML into generic values and parse generated protobuf messages; derive export names, types, and optional presence from descriptors, removing serialization tags from shared models.
- [x] 35.3 Verify existing lossless field, large-identity, appearance, layout, and preservation workflows; update architecture and select exact-candidate checks. Delivery receipts own subsequent publication outcomes.

## 36. Remove nested path compatibility

- [x] 36.1 Remove nested resource-path normalization; accept only protobuf-defined flat path keys and update the documented contract.
- [x] 36.2 Verify flat-path convergence and rejection of nested collection and deck paths; select exact-candidate validation through delivery receipts.

## 37. Use generated serialization models directly

- [x] 37.1 Replace duplicated configuration, appearance, card, and plan structs with generated protobuf messages; retain only runtime wrappers and snapshot metadata.
- [x] 37.2 Remove field-by-field projections and preserve cloning, optional values, and semantic comparisons.
- [x] 37.3 Update architecture and verify existing end-to-end behavior and real-collection convergence; delivery receipts own exact-candidate validation and publication.

## 38. Extract API and runtime models

- [x] 38.1 Publish exact API/model package copies in prerequisite PR #136 against master, with independent compilation and validation.
- [x] 38.2 Record the prerequisite relationship and rebase onto merged PR #136 to remove shared additions without breaking this implementation.

## 39. Document package boundaries

- [x] 39.1 Add doc.go to the eight implementation packages and document every top-level declaration and struct field, including private helpers and test fixtures. Include package docs in their Bazel targets.
- [x] 39.2 Verify package comments describe current responsibilities and boundaries; document merged packages and the Go skill separately against master.
- [x] 39.3 Select exact-candidate validation and publication of the updated implementation through delivery receipts.

## 40. Simplify serialization and collection internals

- [x] 40.1 Replace custom serialization validation with TOML/map/JSON/ProtoJSON decoding; keep numeric TOML export conversion local to textstore and use ProtoJSON for CLI plans.
- [x] 40.2 Review the complete implementation and simplify duplicated note-field ordering and slice comparisons while retaining preservation and publication guards.
- [x] 40.3 Update the documented architecture and serialization contract, verify end-to-end behavior and representative output, and select exact-candidate delivery checks. Delivery receipts own subsequent validation and publication outcomes.

Review covered production packages and their call sites. Reuse the shared name
normalizer in SQLite collation, compute cloze matches once per note, and retain
the final input hash check in archive publication while removing the earlier
duplicate full-file hash. The existing transaction, filtered-card references,
unknown protobuf bytes, and staged publication guards remain necessary.
