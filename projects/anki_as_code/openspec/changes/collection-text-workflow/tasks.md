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

- [x] Use a standalone Cargo project with unmodified upstream manifests/build scripts; the shared build-anki-native-backend change owns dependency acquisition and native compilation.
- [x] Use direct Cargo compilation with pinned Git dependencies and standard cc_import; remove custom orchestration and the upstream patch while retaining existing native behavior and clean streams.

Extraction PR #138 was closed at the user's request on 2026-10-12. Validation
and Anki-format packages remain in PR #128; any future extraction will use a
new pull request.

Review covered production packages and their call sites. Reuse the shared name
normalizer in SQLite collation, compute cloze matches once per note, and retain
the final input hash check in archive publication while removing the earlier
duplicate full-file hash. The existing transaction, filtered-card references,
unknown protobuf bytes, and staged publication guards remain necessary.

## Deferred field caches

- [x] Remove custom HTML normalization and cache calculation; preserve existing caches and initialize new-note caches empty/zero.
- [x] Document required Anki Check Database after changed fields or sort-field settings; verify the collection workflow and publish.

- [x] Remove custom field-cache HTML conversion and document Anki-owned rebuilds; delivery receipts own validation and publication.

## Readable filtered decks

- [x] Export/edit filtered settings with upstream protobuf schema, preserving unknown data and resident cards.
- [x] Verify readable TOML archive workflow, conflict rejection, convergence, and publish.

- [x] Expose the pinned upstream filtered schema as readable TOML; select workflow, codec, and API consumer checks.

## Reuse regular-deck schema

- [x] Replace custom DeckConfig with upstream Deck.Normal and preserve unknown configuration.
- [x] Verify readable optional regular settings through archive editing and publish.

- [x] Reuse upstream regular-deck configuration and audit remaining messages for equivalent contracts.

- [x] Reuse upstream note-type and template configs as nested fields; preserve appearance-only management and verify rejection of unrelated edits.

## Data structure review and compact export

- [x] Rename appearance wrappers to NoteTypeSettings and TemplateSettings, align runtime naming, and remove redundant note-type/settings and ordered-field representations.
- [x] Extend the existing export/edit/rebuild workflows before implementation to verify omitted defaults, retained optional zero values, empty message presence, named empty fields, and convergence.
- [x] Review serialized, shared runtime, and adapter-private structures and select exact-candidate gates for both PRs; delivery receipts own validation and publication.

Reviewed all serialized messages and shared models, plus adapter-private,
command, and fixture structures. Removed duplicated settings snapshots, ordinal
maps, and per-note field-name arrays. Named note fields and positional Bodies
remain source and derived persistence projection; deck Kind retains original
wire bytes for lossless preservation. Collection manifests and plans remain CLI
contracts, while deck/note-type/template configuration reuses upstream protos.
The later exported-symbol audit also identifies an unused note source-path field.

## Validation review follow-up

- [x] Replace field-specific unknown-data copying with recursive preservation and reject edits that discard unknown nested configuration; exercise failure cases before implementation.
- [x] Store note-type settings in an explicit field and split desired-state, deck, note, card, and settings validation into focused helpers.
- [x] Remove cloze-content inference; validate explicit cards against desired note types and preserve existing card identities and scheduling during structural edits.
- [x] Add note-type creation, renaming, field/template structure and configuration editing, and removal of unreferenced types; select schema-marker, preservation, filtered-card, and convergence checks for both candidates; delivery receipts own exact-candidate validation and guarded replies.

## Exported-symbol audit follow-up

- [x] Keep archive, repository, and textstore implementation details private; migrate workflow fixtures to local text editing and the public identity operation.
- [x] Remove the unused collection identity wrapper and note source-path field.
- [x] Verify remaining handwritten exports have real consumers or required contracts, preserve the integration use of `PreserveUnknown`, and validate the exact candidate.

## Anki Rust backend and derived cards

- [x] Inspect the pinned backend, current persistence/validation boundaries, and existing toolchains; prepare the exact locked dependency inventory for approval.
- [x] Obtain approval for the locked Rust dependencies before adding or installing them. The user approved the exact proposal on 2026-10-12.
- [x] Write backend-compatible E2E failure cases before implementation for card generation without explicit identities, read-only derived plans, template migration/errors, cache rebuilding, unknown configuration, scheduling/history/preset preservation, filtered-card protection, and convergence.
- [x] Build the pinned rslib engine with existing toolchains and a bounded first-party C ABI/cgo bridge; verify actual native execution.
- [x] Replace reconciliation SQL and duplicated Anki semantic validation, retaining only desired-state and unmanaged-resource guards.
- [x] Make card entries deck assignments, remove card identity allocation from copied notes, and derive lifecycle changes through Anki.
- [x] Use the same private-copy backend execution for planning and applying; report derived changes and persist required managed normalization only after publication.
- [x] Update README, package docs, architecture, examples, and existing workflows to describe and verify the final backend behavior.
- [x] Select exact-candidate gates and verify representative native artifacts; repo-delivery receipts own final verification, publication through PR #128, and review state.

The first pre-implementation E2E cases compile and fail for the intended gaps:
`TestDerivedCardWorkflow` reports no generated card for a new cloze note without
explicit card identities (`derived-card-generation-missing`), and
`TestBackendTemplateValidationWorkflow` reports that plan accepts an unknown
field reference (`template-validation-missing`). Configured formatting and
strict OpenSpec validation pass. Native generation and invalid-template cases now pass, alongside the existing
preservation, structural migration, and convergence workflows. The dependency proposal and failure
logs are under `out/anki-rust-backend/`; the approved locked dependencies are wired through checksum-verified Bazel acquisition.
The user approved the locked Rust dependency proposal on 2026-10-12.
Implementation is complete; repo-delivery receipts own exact-candidate verification and publication.

Earlier native acceptance evidence: the full project run passed both test targets in
`out/anki-rust-backend/project-tests-2.log`; semantic lint passed in `lint-1.log`.
The locked projection contains 475 external versions, all within the 493-version
approved proposal. Module reconciliation and both strict OpenSpec checks passed.
Exact-candidate delivery gates and receipts remain authoritative for publication.
The final full project run also passes the paginated command-stream E2E in
`project-tests-3.log`. The attached collection exports, builds a private round-trip
archive, and produces one clean JSON plan with zero changes.

Session ergonomics: the first prototype used separate native and Go SQLite
connections and failed to observe a transplanted note-type identity. The adapter
now uses Anki's unpaginated database service for snapshots and preservation,
removing the production Go SQLite dependency. The native generation and existing
preservation workflows verify that boundary. Cargo now owns source acquisition
and dependency resolution; the custom inventory generator has been removed.

The attached collection exports and yields an empty native plan on a private copy.
Its paginated snapshot exposed an upstream debug print on stdout. A new public
paginated command workflow failed before the fix (`native-pagination-stdout`);
the previous upstream patch removed only that diagnostic and unused allocation.
Direct Cargo now consumes unpatched upstream source; the unpaginated database
command avoids that diagnostic without a native result cache.

## Consume the merged Anki gRPC backend

- [x] Rebase PR #128 onto merged server/client support, preserve collection features, and remove superseded Cargo/C ABI files and CI prerequisites.
- [x] Route collection operations through generated Anki clients and the shared managed Go backend; retain the database-service facade only for snapshots and preservation.
- [x] Propagate startup context, close/checkpoint the collection before stopping the worker, and verify existing generation, migration, preservation, large-result and convergence E2E artifacts.
- [x] Update current architecture/build contracts and documentation, review the complete rebased PR, and select exact-candidate delivery. Delivery receipts own validation and publication state.

Migration failure cases: missing worker runtime data, invalid generated service
routes, startup cancellation, lost native error details, truncated or imprecise
SQL responses, uncheckpointed publication, leaked worker lifetime, and changed
scheduling/history/configuration. Existing collection/CLI E2E and the shared
backend lifecycle E2E provide behavior coverage; delivery receipts own results.

Migration evidence: existing collection, built-CLI, and shared backend lifecycle
E2E targets passed after the gRPC migration. The CLI has no C ABI, native Cargo
manifest/lock, or CI bootstrap changes. Anki remains owned by third-party Bazel
targets; the shared Go backend owns server process lifetime. SQL remains on
Anki’s existing connection for exact snapshots, opaque/unmanaged preservation,
and declared identities. No new dependency was added.

Review corrected the external-note identity E2E fixture: configure its path
through the schema and give the selected copy a lower ID than the undeclared
neighbor. This verifies referenced-root reservation independently of quoting
and the selected note’s own ID. No production behavior changed in this fix.

## API-only reconciliation and fixtures

The user supersedes previous raw-storage and unknown-wire preservation
requirements. Earlier completed tasks are historical evidence. Current specs
and design use Anki APIs for both production and E2E fixtures.

- [x] Define API-fixture workflows before implementation, including managed edits, generated IDs, history/preset retention, filtered-card guards, structural changes, large fields, stream output, discovery, and convergence.
- [x] Replace snapshots and reconciliation with generated APIs; remove SQL, raw preservation, and serialized deck-kind transport. Record Anki-generated IDs and keep plan addition references stable.
- [x] Replace SQLite E2E fixtures with supported API creation/inspection and produce repeatable text/archive artifacts.
- [x] Update docs, diagrams, build dependencies, and OpenSpec contracts; validate, review the complete PR for further API-based simplification, and select guarded delivery.

Failure cases: missing runtime server; omitted filtered decks or cardless notes
in searches; field/template ordinal drift; ID remapping breaking references;
API normalization preventing convergence; partial failure changing input files;
closing after caller cancellation; large gRPC values; invalid templates;
filtered-card conflicts; unexpected changes to supported scheduling/history or
presets; duplicate identity and recursive path ownership errors.

Review evidence: the complete PR diff was reviewed with review-agent. Removed
field/tag/presentation validation owned by Anki and the full snapshot used only
for card assignment; CardsOfNote/GetCard provides that view directly. ID
references are remapped in one pass to avoid cascading collisions, and plan
references use indexed lookups. API workflows verify cloze generation, entity
deletions, template deck references, content normalization, reviewed scheduling,
review history, presets, media, external layouts, and convergence. The
normalization workflow failed before redundant Go checks were removed. The
final delivery receipt owns exact-candidate check and publication results.

Full-PR review also found a filtered-card guard that compared template names
but missed upstream identity swaps. The API E2E failed when such a swap was
accepted. Reuse the reconciliation identity matcher in that ownership guard;
the rejected plan/apply workflow also verifies archive bytes remain unchanged.

Review replaces repeated desired-note scans during deletion discovery with an
identity index, avoiding quadratic work for large collections. Existing
creation/deletion E2E verifies this reconciliation behavior.

## GUID-based read-only reconciliation

- [x] Write native E2E cases before implementation for missing GUID rejection, GUID generation that preserves existing identities, new notes and note types in one apply, unchanged source files, ordinal assignments, and repeated convergence.
- [x] Consume the fixed-hierarchy Config schemas and internal/config from #151; remove textstore and source normalization. Match notes by GUID and cards by ordinal with runtime-only numeric identities.
- [x] Resolve/create note types before notes, keep collection publication guarded, and compare actual API snapshots for stable plans without rewriting TOML.
- [x] Update current README/site/OpenSpec boundaries, review the full candidate, validate native and CLI E2E plus required quality/lint gates, and publish #128 independently against master.

Structural-convergence evidence: Anki's first field migration populated browser
front templates while the unchanged source declared them absent; a second apply
removed them. Reapplying declared configuration against the migrated template
ordinals fixes that transition in the same reconciliation, through Anki APIs.
The source files stay unchanged and native E2E converges after one apply/build.

## Unreleased version 1

- [x] Update version-1 config dump and GUID fixture expectations before implementation; replace legacy-layout cases with current-schema/version checks and verify export through existing native workflows.
- [x] Keep format version 1, remove the earlier-version header pass and protobuf compatibility reservations, and align docs/specs without migration guidance.
- [x] Review the correction, validate exact-candidate behavior and required quality/lint/build gates, and publish the updated PR.

- [x] Replace handwritten valid fixture setup with protobuf messages and the existing TOML encoder, preserving independent parser checks and deliberately invalid syntax inputs. CLI behavior passes; delivery receipts own publication and review resolution.

## Configuration writer boundary

- [x] Extend the shared configuration E2E before changing the writer to verify object-to-files-to-CLI export, exact values, unchanged messages, unchecked inputs, and direct filesystem behavior.
- [x] Share the direct writer from #151; move snapshot default-assignment projection, export checks, staging, cleanup, and publication into collection while preserving collection export guarantees.
- [x] Update docs and package wiring, verify native/CLI workflows and representative outputs, and select exact-candidate delivery gates. Delivery receipts own publication and review state.

## Direct JSON-to-TOML conversion

- [x] Update the shared E2E expectations before implementation for quoted ProtoJSON integer output and unchanged values, presence, and GUID behavior.
- [x] Use the direct JSON-to-TOML codec from #151, remove scalar/descriptor handling, and align docs/specs.
- [x] Verify native and configuration workflows, generate real note/deck/collection preview artifacts, and select exact-candidate delivery gates. Delivery receipts own validation and publication.

## Collection version field

- [x] Update existing configuration E2E fixtures/expectations before renaming the root protobuf field and Go consumers to version.
- [x] Align examples and deck-creation identity documentation, verify config/native workflows and representative output, and select exact-candidate delivery gates.

## Collection-backed path identities

This version-1 path/state contract supersedes historical ID/GUID workflows above.

- [x] Write path/state and configuration E2E failure cases before implementation, preserving native and codec coverage.
- [x] Replace editable native identities with path wrappers/readable references and remove identity generation.
- [x] Persist state via Anki's config API; initialize on export, import existing objects, resolve explicit moves, and retain untracked objects.
- [x] Verify Anki's actual full-package export/import/re-export, guarded failures, native structural migrations, import convergence, and unchanged source configs through E2E.
- [x] Inspect a 7,218-note real export and before/after configs, complete source/output review, and select exact-candidate quality/lint/build/E2E gates. Delivery receipts own final validation/publication status.

## Readable stable note filenames

- [x] Extend native export E2E before implementation for sort-field selection, HTML/entities/Unicode, duplicates, unsafe/empty/long content, and retained names after edits and new native notes.
- [x] Derive safe initial filenames from sort-field content, reserve existing bound names, and document state-backed stability without modifying note values.
- [x] Verify representative exports and convergence, and select exact-candidate E2E/quality/lint/build gates for the update to #128. Delivery receipts own validation and publication status.

Filename evidence: native CLI E2E failed before implementation because readable
names were absent, then passed for sort-field selection, HTML/entities/Unicode,
image/audio-only fallback, unsafe/long labels, imported-name reservations,
content edits, and native additions after export. The 7,218-note attached
collection exports with a zero-change plan; repeated export preserves every
config path/value and source-archive bytes. Evidence is under
`out/anki-note-filenames/`; the exact-candidate delivery receipt owns gates.

## Configuration package boundaries

- [x] Consume #151's configuration facade/layout/codec package split without changing native collection behavior or introducing another Go module.
- [x] Verify native and configuration E2E, document the package boundaries, and select exact-candidate delivery gates. Delivery receipts own publication status.

Package-boundary evidence: native/CLI/config E2E pass, and the 7,218-note
collection retains a zero-change plan. The configuration subtree matches #151;
codecs perform no filesystem access, layout owns discovery/file I/O, and all
collection callers retain the configuration facade. Evidence is under
`out/anki-config-modules/`; delivery receipts own final gate/publication status.

## Configuration prerequisite merge

- [x] Rebase #128 onto master after #151 merges, preserve configuration APIs and implemented collection operations, and retain existing integration documentation.
- [x] Inspect the reduced aggregate diff and select exact-candidate native/config E2E, quality, lint, and build gates. Evidence belongs to `out/anki-config-merge-rebase/`; delivery receipts own validation and publication status.

## Configuration diagrams

- [x] Document the configuration structure, hierarchy, and references with a schema-level Mermaid diagram, and show pure configuration operations in a separate diagram.
- [x] Render site SVGs and dark WebP previews, inspect their accuracy and readability, and include the diagrams in exact-candidate delivery gates.
- [x] Simplify the dense diagrams: separate hierarchy from references, use short labels, and show Load/Dump/Export as input/output edges. Inspect all dark previews before delivery.
- [x] Compare alternative diagram types, keep a class diagram for config messages, and give each Go API operation its own process chart with start/end events and gateways. Inspect dark WebPs and update exact-candidate gates.
- [x] Verify mixed diagram types are unsupported by the pinned Mermaid renderer, combine the three API processes in one chart, and keep the class model separate. Inspect dark renders and select updated delivery gates.
- [x] Add file-location annotations to config messages, covering root/deck notes and recursive deck directories, and put the operation chart above the model in the README.
- [x] Replace routine validation branches with discovery/serialization steps in Load and substantive JSON/TOML output paths in Dump; keep the model to types, fields, layout, and references without field comments.
- [x] Inspect dark renders, retain layout and path references, and select exact-candidate quality, lint, E2E, and site-build gates for publication. Delivery receipts under `out/anki-config-field-docs/` own validation and publication evidence.
- [x] Restore the class model with file-location notes and no field comments, inspect its dark render, and select exact-candidate delivery gates. Delivery receipts under `out/anki-config-class-model/` own validation and publication evidence.
- [x] Show deck loading as a loop that discovers each deck's child directories; render and inspect a dark preview without publication. Evidence belongs to `out/anki-config-deck-loop/`.
- [x] Correct Load and Dump process order, show decoding and defaults inside deck loading, keep depth-first child recursion explicit, and inspect a clearer dark render without publication. Evidence belongs to `out/anki-config-operation-flow/`.
- [x] Replace inline serialization arrows with individual process nodes for ProtoJSON encoding, JSON decoding, and TOML encoding; render and inspect without publication. Evidence belongs to `out/anki-config-encoding-nodes/`.

## Direct Anki package operations

- [x] Record native workflow boundaries and failure cases; update E2E fixtures and
      publication scenarios before production code.
- [x] Remove internal/archive; use native package APIs from collection orchestration
      and make repository clients borrow the open connection.
- [x] Verify native packaging, media/history/state, large fields/media, no-op behavior,
      cancellation, journal guards, and external-writer protection through E2E.
- [x] Update the README and architecture, review the candidate, and select native/CLI,
      quality, lint, and site gates. Preserve unpublished config-diagram previews.

Evidence belongs to out/anki-direct-api/; delivery receipts own publication state.

## Extract remaining shared bindings and config documentation

- [x] Keep third-party service bindings and configuration diagrams in independent master-based PRs; record their ownership and prerequisite links while keeping #128 buildable until merge.
- [x] Select exact-candidate quality, lint, native/CLI/config E2E, and site gates for the updated #128 projection, then restore only the three previously unpublished diagram-preview hunks. Delivery receipts own publication state.

## Rebase after shared binding and documentation merges

- [x] Rebase #128 onto the merged #154/#155 base, resolve the documentation-only conflicts while retaining config docs and collection architecture, verify native source bytes are preserved and extracted files leave the aggregate diff, and select exact-candidate E2E, quality, lint, and site gates. Keep the three existing preview hunks unpublished; delivery receipts under out/anki-rebase-bindings-docs own publication state.

## Export and state-import extraction

- [x] Publish export and state import independently as #156 against master; retain #128's working overlap until merge and record the focused repository boundary.
- [x] Link #156 as #128's prerequisite, inspect the scoped documentation/message change, and select exact-candidate E2E, quality, lint, CLI, and site gates. Keep the original three diagram-preview edits unpublished; delivery receipts under out/anki-export-import-split own validation/publication state.

## Shared export review corrections

- [x] Reproduce #156's permission and preset-alias findings through native E2E before updating the overlapping production code.
- [x] Keep #156's fixes and focused native coverage in #128, verify alias-based reconciliation convergence, and select new exact-candidate delivery gates. Retain the three existing diagram previews locally; receipts under out/anki-export-import-split own publication state.

## Explicit state package and read-only export

- [x] Write state CLI/Go API and read-only export/state-artifact E2E cases before production changes; preserve native convergence coverage with explicit fixture state installation.
- [x] Integrate #156's state package and shared worker orchestration; make export emit TOML plus state JSON without changing source bytes, and verify native E2E convergence after explicit replacement.
- [x] Review the updated boundaries and APIs, document current commands, and select exact-candidate quality/lint/E2E/site gates. Preserve the original preview patch for restoration; receipts under out/anki-batch-import own final publication/restoration evidence.

## State file encoding

- [x] Default state dump/replace to TOML with optional JSON; export state.anki.toml by default or state.anki.json with --state-format json. Verify both encodings, unchanged input bytes, explicit installation, and refreshed exports through native E2E; select exact-candidate delivery gates.

## Import identifier interpretation

- [x] Pass import identifiers as strings through CLI and Go APIs; resolve the declared resource before validating its identifier against the native API type. Verify syntax-only CLI parsing and resource-first errors in native E2E.

## Shared format codec

- [x] Move format-dispatch Encode/Decode out of state into codec, decode into any supplied protobuf message, route config/state callers through the shared API, hide format helpers, and update boundaries/docs. Verify existing native and pure-config E2E, select exact-candidate delivery checks, and prepare replies to reviewed codec threads. Delivery receipts own final publication and resolution evidence.

## Nested state and selectable namespaces

- [x] Write E2E cases first for nested member bindings, independent state namespaces, selected-key CLI/API behavior, native package round trips, and namespaced reconciliation.
- [x] Rename the protobuf contract to State, store members beneath resource objects and presets beneath aliases, and propagate a configurable versioned state key through all consumers without a compatibility layer.
- [x] Verify moves, native study-data preservation, default TOML/optional JSON, and namespace isolation; document APIs and select exact-candidate delivery gates and review replies.

Delivery receipts under out/anki-state-schema own final validation, publication, and review-thread resolution evidence.

## Literal note-type filenames

- [x] Verify that a real filename containing # after another type's filename is preserved during that type's move; keep fragment interpretation limited to declared field/template move addresses.

## Typed state maps

- [x] Adapt native E2E before implementation for separate deck, note, note-type and preset state objects, nested typed members, and rejected cross-map ownership.
- [x] Replace generalized resource bindings and discriminators with typed protobuf maps and direct consumers; retain namespace, format, move and atomic import behavior.
- [x] Verify native round trips and preservation, inspect representative TOML, document the schema, and select exact-candidate delivery gates for both PRs.

## Grouped resources and state update

- [x] Write native CLI/API E2E before implementation for State.Resources and atomic external-state update, preserving omitted entries, replacing supplied map entries, no-op bytes, namespace isolation, rejected conflicts and native data.
- [x] Group typed maps in Resources and propagate direct access; add Store.Update and state update with existing file/stdin and TOML/JSON handling. Preserve filesystem-path import and version 1 without old-shape compatibility.
- [x] Verify native/config/CLI E2E, inspect emitted state, document merge semantics and select exact-candidate quality/lint/build checks for both PRs.

## Explicit collection lifecycle

- [x] Write native E2E first for connection gating before preparation, after failure/export/close, repeated close, failed preparation cleanup and unchanged inputs, with repeatable native output artifacts.
- [x] Introduce internal/collection/repository.Collection to own preparation, worker/native open-close/export and private cleanup; keep staging/hash checks/publication in collection.Run and let feature stores borrow its prepared connection.
- [x] Verify native/config/CLI behavior, document ownership and API contracts, and select exact-candidate delivery checks for both PRs without publishing the existing diagram previews.

## Direct database access

- [x] Write native E2E first for opening supplied databases without temporary copies and persistent API writes, including before a later caller error.
- [x] Open databases directly and bypass replacement staging for direct operations; retain native restore/guarded output for package inputs and document cleanup/commit boundaries.
- [x] Replace mutating plan execution with snapshot/config comparison, expose deferred native card-generation effects, verify repeated read-only plans, and retain native apply validation and actual-result checks.
- [x] Select exact-candidate validation for both PRs and inspect native artifacts; delivery receipts own validation and publication status.

## Rebase after collection-backend merge

- [x] Verify the projected backend files match merged #157, retain only application integration after rebasing onto master, and select exact-candidate quality, lint, native/CLI/config E2E and site checks. Review the complete rebased PR diffs; delivery receipts own final observations.

## State store diagram preview

- [x] Render and inspect an unpublished dark WebP showing state.Store loading configs and calling anki_grpc for import, update and replace, with State inside Anki storage. Keep the preview local. Evidence: out/anki-state-store-diagram/state-store.webp, rendered by docs:state_store_preview; labels, containment and edges visually verified.
- [x] Clarify the preview: list Dump, Import, Update and Replace inside state.Store, label its server connection gRPC, show configs as a document used for import, and connect anki_grpc directly to State with read/write. Inspect the dark WebP and retain it without publication. Evidence: out/anki-state-store-diagram/state-store-revised.webp; labels, document shape and storage containment visually verified.
- [x] Give Dump, Import, Update and Replace separate nodes inside a state.Store group. Render and inspect the dark WebP without publication. Evidence: out/anki-state-store-diagram/state-store-methods.webp; method nodes, group titles, labels and edges visually verified.
- [x] Arrange the methods as a compact two-by-two block, connect Import directly to a plain Parsed configs node, and annotate State as a collection config key. Render and inspect the dark WebP without publication. Evidence: out/anki-state-store-diagram/state-store-compact.webp; method grouping, config edge, labels and storage containment visually verified.
- [x] Remove the State annotation, align both method rows, reduce group whitespace and straighten the main connection. Inspect the dark WebP and retain the unpublished preview. Evidence: out/anki-state-store-diagram/state-store-connection.webp; ELK renders aligned method rows and a horizontal main connection without overlapping labels.
- [x] Group the anki_grpc node and Anki storage container inside an anki_grpc connection graph. Inspect the dark WebP and retain it without publication. Evidence: out/anki-state-store-diagram/state-store-connection.webp; both group titles and nested storage containment visually verified.
- [x] Show each method's read/write access to the anki_grpc connection group, rename its inner node to read/write and connect it to Anki storage without an edge label, and add unlabeled Input configs and Input state dependencies. Render and inspect the unpublished dark WebP. Evidence: out/anki-state-store-diagram/state-store-access.webp; all four access edges, input dependencies, labels and group containment visually verified.
- [x] Move Input configs and Input state to the left with arrows into their consuming methods, restore the inner node name anki_grpc, and inspect the unpublished dark WebP. Evidence: out/anki-state-store-diagram/state-store-inputs-left.webp; input positions, all seven external arrows, labels and nested groups visually verified.
- [x] Label Replace read/write to reflect its read-before-write behavior. Compare method and edge arrangements through the layout previews below rather than choosing a final arrangement before user selection.
- [x] Replace the anki_grpc connection container and its inner node with one connection node pointing to a separate Anki storage group. The following direct-State adjustment owns the final edge target.
- [x] Point the simplified connection node directly to State inside Anki storage. All three layout candidates render this edge with no intervening backend node.
- [x] Put connection on a new line beneath anki_grpc in the simplified node label. All three layout candidates render the two-line label.
- [x] Render and inspect three unpublished layout candidates with the same operations, inputs and storage relationships: method grid, method column and vertical flow. Show all dark WebPs for user selection. Evidence: out/anki-state-store-diagram/layouts/; labels, input arrows, individual access edges and group containment visually verified. The vertical candidate flows bottom-to-top to avoid routing an arrow through the storage title. Final layout selection remains with the user; no preview is published.
- [x] Adopt the user's selected layout B as state-store.mmd, remove the alternative source files and render targets, and inspect its dark WebP. Keep the selected diagram unpublished. Evidence: out/anki-state-store-diagram/state-store-selected.webp; the canonical target renders layout B with the selected labels and direct State connection.

## Rebase after state-management merge

- [x] Rebase #128 onto merged #156, retain the merged state store and prepared-connection API, and restore the existing database/archive CLI E2E coverage.
- [x] Keep one shared state namespace option and reuse the existing note-type binding projection across state import and reconciliation.
- [x] Include all previously unpublished config-process and selected state-store diagram sources at the user's explicit request; publish maintained SVGs with the documentation and retain dark WebP render targets.
- [x] Verify the rebased native/config/CLI workflows, inspect changed diagram renders, and select exact-candidate quality, semantic lint and site checks. Delivery receipts under out/anki-rebase-state-156 own final validation and publication state.

## Restore archive state integration

- [x] Restore the original native CLI E2E coverage for both .anki2 and .colpkg before correcting the rebase integration; keep no-op bytes, failure guards, explicit exported-state installation and namespace coverage.
- [x] Route state commands through the existing collection workflow, returning each operation's changed flag for native archive publication. Keep Store lifecycle-free; callers sharing a prepared connection reuse a store and close their backend owner.
- [x] Restore native database/archive workflow coverage, document the complete archive round trip and select exact-candidate native/config/CLI, quality, lint and site checks. Delivery receipts own final verification and publication.

## Explicit collection input naming

- [x] Update existing CLI/native workflow calls before implementation to use --collection for Anki .anki2/.colpkg inputs.
- [x] Rename the shared flag and application collection option, update current examples/help and keep build's explicit --base archive option.
- [x] Select native/config/CLI E2E and delivery gates for the renamed command surface; exact-candidate receipts own the final results.

## Review corrections and simplification

- [x] Add native E2E cases before implementation for retained-note deck deletion, raw-database no-op apply, and multiple created-note identities/deck assignments.
- [x] Reject removal of a deck still holding retained cards without a declared destination during plan/preflight, before native writes.
- [x] Share editable note-type equality with planning while retaining member-identity/order checks, so converged apply performs no note-type updates.
- [x] Remove the unused App Start/Stop/WithState lifecycle and its test-only API; keep the collection workflow as the command lifecycle owner.
- [x] Index current notes by native ID once during resolution and reuse the resolved previous note for type checks.
- [x] Create new notes with Anki's AddNotes batch API and map its returned IDs to the correct declarations.
- [x] Correct the landing-page bootstrap flow and document current lifecycle and no-op guarantees.
- [x] Verify the native reproductions and E2E workflows, inspect output artifacts, and select exact-candidate quality/lint/site and workflow checks. Receipts under out/anki-fix-review-128 own final validation/publication state.

## Capability layout and follow-up review corrections

- [x] Add native E2E failure cases before implementation for occupied normal-to-filtered deck conversion and unrelated export file/directory permissions.
- [x] Reject occupied deck-kind conversions before native writes and retain reviewed cards when an emptied deck changes kind.
- [x] Preserve permissions on unrelated copied export entries, including nested directories and private files.
- [x] Give app complete workflow/publication ownership, move the lifecycle owner into collection, and move the generic codec into internal/codec.
- [x] Consolidate resource access, planning, resolution, validation, export projection and private runtime models in reconcile; remove obsolete helper/model packages.
- [x] Return only a changed flag from native reconciliation, reuse Anki OpChanges, and verify native identity sets without a second full snapshot.
- [x] Restrict source hashing to archive replacement; direct database commands must not hash unused source bytes.
- [x] Update package/API documentation and the architecture diagram, and verify native/config/state/CLI workflows and representative outputs.
- [x] Select exact-candidate quality, semantic lint, E2E and documentation checks for delivery. Receipts under out/anki-package-layout own final validation/publication state.

## Shared state and codec extraction

- [x] Publish the standalone state/codec prerequisite against master and verify its scope contains only shared packages and required callers/docs.
- [x] Prepare the prerequisite reference in this change and the #128 description; keep required overlap until it merges so both master-based PRs remain buildable.
- [x] Select exact-candidate quality, semantic lint and workflow/documentation checks for the dependency update; receipts under out/anki-state-codec own final validation/publication results.

## Store-owned state operations

- [x] Consume Store.Preview for binding construction and explicit moves; remove raw State.Resources mutation and duplicated state helpers from reconciliation.
- [x] Preserve export snapshots, provisional planning identities, native member bindings, preset aliases and validated apply persistence through store operations; verify native E2E.
- [x] Update package/API documentation and OpenSpec ownership, then select exact-candidate quality, lint, workflow and docs checks; receipts under out/anki-store-boundary own final results.

## Resources protobuf preview input

- [x] Define native E2E cases before implementation for whole Resources replacement, nil retention, detached snapshots and invalid/nil entries; observe the old API fail.
- [x] Remove Bindings, accept Resources in Store.Preview, update native callers and keep state transformations/writes inside Store.
- [x] Update API docs, verify native workflows/representative artifacts/OpenSpec, and select exact-candidate delivery checks; receipts under out/resources-proto own final results.

## Rebase after state and codec merge

- [x] Rebase onto master after #158; resolve only duplicate/stale README blocks, retain the prior source tree, and remove already-merged state/codec changes from the PR diff.
- [x] Refresh the merged-prerequisite description and select exact-candidate quality, semantic lint, native E2E and documentation checks. Receipts under out/rebase-158-review-128 own final validation/publication state; the requested read-only review is reported separately.

## Stable note-type export filenames

- [x] Add native E2E cases before implementation for renamed tracked types, occupied suffixes, earlier untracked identities, case-insensitive names and fallback-name collisions; observe export failures.
- [x] Reserve tracked note-type filenames before allocating deterministic unused names, preserving note references and state bindings without changing the collection.
- [x] Verify stable re-export, emitted native/config artifacts and OpenSpec; select exact-candidate quality, semantic lint and workflow/documentation checks. Receipts under out/notetype-export-collisions own final validation/publication state.

## Export ownership and review corrections

- [x] Extend the existing native export-refresh workflow before implementation to cover stale nested decks, notes and note types; retain permission, failure and state-format coverage.
- [x] Move the prepared-connection export operation and filesystem publication out of app into internal/export; reuse codec.Format and remove generated-output reloading. Retain stale declaration removal and empty hierarchy pruning required by refresh.
- [x] Update package docs and the architecture diagram; verify native/CLI workflows, representative outputs and strict OpenSpec, then select exact-candidate delivery checks. Receipts under out/anki-export-boundary own final validation and publication state.
- [x] Prepare evidence-backed replies for the four export review threads explaining the implementation changes or why cleanup remains necessary; delivery receipts own final publication and thread-resolution state.

## Reuse native collection lifecycle and archive publication

- [x] Replace wrapper-specific lifecycle/publication cases with native workflow E2E coverage before implementation: package/database state persistence, borrowed-connection shutdown, normal restore cleanup, callback failures, cancellation and invalid packages.
- [x] Compose BackendSTD and BackendCollection directly in app; remove internal/collection, duplicate lifecycle flags, restore ownership, journal checks and custom archive hashing/staging/sync/rename. Call Anki's atomic ExportCollectionPackage directly and keep no-op export decisions.
- [x] Reconcile documentation/specification claims with native API ownership, rerender the architecture diagram, and select exact-candidate delivery checks. Receipts under out/anki-export-boundary own final validation and publication state.

## Constructor-injected exporter

- [x] Replace export.Run with Exporter, inject its reconciliation engine through New, and keep only context, destination and format on Export. Update callers, dependency declarations and documentation without changing export behavior.
- [x] Verify existing native export workflows and representative outputs, validate OpenSpec, and select exact-candidate delivery checks. Receipts under out/exporter-constructor own final validation and publication state.

## Independent prepared-connection export

- [x] Remove Engine.Export and export projection from reconcile; construct Exporter from a borrowed gRPC connection and state key. Read native protobufs directly through generated clients without a shared reader.
- [x] Extract the independent export library and native E2E into a master-based PR. Keep CLI integration and required prerequisite overlap in #128 until the extraction merges.
- [x] Verify library and collection workflows, representative TOML/state/native outputs, documentation and OpenSpec; select exact-candidate delivery gates for both PRs. Receipts under out/export-extraction own publication state.

The independent prepared-connection exporter is merged in [PR #160](https://github.com/alwaldend/src/pull/160). This PR uses that library from master; its remaining export changes are CLI integration.

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

## Native filename conversion and export review

- [x] Extend the shared native export E2E before implementation and reproduce incorrect word splitting by the custom HTML parser.
- [x] Reuse Anki's HtmlToTextLine RPC through generated rendering clients; remove custom HTML/audio parsing, its direct dependency and redundant test helpers. Keep the implementation identical to #160.
- [x] Review the complete #160 diff for further correctness and simplification issues; retain its minimal constructor/export API and report no additional blockers.
- [x] Select native collection/export workflows, representative-output, strict OpenSpec and exact-candidate quality/lint/docs checks; out/native-export-names receipts own final results.

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

## Rebase after export merge

- [x] Rebase onto master after #160; resolve only README and Bazel integration conflicts, retain the merged exporter and all native Go source, and remove merged export/resource/schema additions from the feature diff.
- [x] Select native/CLI/config E2E, quality, semantic lint and integrated documentation gates. Receipts under out/rebase-review-160 own exact-candidate validation/publication evidence; report the requested read-only review separately.

## Rebase after configuration-store merge

- [x] Rebase onto master after #162, preserve the merged config.Store implementation and retain the detailed operation diagrams and collection architecture while resolving documentation conflicts.
- [x] Migrate remaining collection callers and E2E fixtures to root-bound config stores; verify existing native/config/state/CLI workflows and representative config output without changing the file format or native reconciliation behavior.
- [x] Select exact-candidate quality, semantic lint and integrated documentation gates. Receipts under out/rebase-162 own final validation and publication evidence.

## Rebase after collection-workflow merge

- [x] Rebase onto master after #163; retain WithCollection and direct state-store construction while preserving the remaining export and reconciliation implementation and E2E helpers.
- [x] Reconcile state/lifecycle documentation and merged prerequisite boundaries without restoring RunState or WithState.
- [x] Select exact-candidate native/CLI/config/state E2E, quality, affected-target lint and integrated documentation checks. Receipts under out/extract-export-cli own validation and publication evidence.

- [x] Remove duplicated state diagram target declarations from the automatic prerequisite merge; retain the merged diagram once and rerun quality checks.

## Export containment review correction

- [x] Reproduce source/destination alias bypasses in existing native export E2E before implementation; resolve the source and destination's nearest existing ancestor through symlinks and preserve new-directory exports. Fresh exact-candidate receipts under out/extract-export-cli own final verification and publication evidence.

## CLI application boundary

- [x] Carry the CLI-owned application interface and isolated command failure cases from #164 into the remaining feature branch before production edits. Keep concrete state stores, import config loading and native callbacks inside App methods; remove the redundant public lifecycle callback API and collection-options object.
- [x] Verify isolated commands, native workflows and representative outputs, then select fresh exact-candidate delivery gates under out/extract-export-cli.

## Test ownership

- [x] Move config parsing/layout/export behavior into internal/config tests and collection input failures into internal/app tests. Keep syntax, flags, dispatch, codecs and streams in CLI tests; update build targets and references.

## Operation arguments

- [x] First adapt isolated CLI cases to require an interface-only root, forwarded explicit namespace arguments, default keys, and repeated commands without retained configuration.
- [x] Remove StateOptions and its mutable App field. Pass the namespace as a string argument through collection methods, keeping root flag parsing local to CLI and state interpretation in App. Update native callers and documentation on both branches.

## Reusable collection scopes

- [x] Record native failure cases for shared collection operations, duplicate path aliases, expired references, overlapping access, cancellation and archive publication before implementation.
- [x] Add an App-owned collection registry and opaque CollectionRef selectors. WithCollection lends a reference; path selectors retain one-operation behavior. Reject duplicate opens and overlapping uses through the registry guard.
- [x] Reuse prepared connections, resource stores and namespaced state stores. Publish changed archives once after a successful scope and preserve direct database writes, cleanup and one-shot behavior.
- [x] Update CLI adapters, the independent export CLI prerequisite and current documentation; validate native artifacts and select delivery gates for both affected PRs. Delivery receipts own publication and review evidence.
- [x] Verify command isolation, native workflows and representative artifacts, then select exact-candidate delivery gates. Delivery receipts own publication and review evidence.

## Staged copy-and-apply builds

- [x] Define native E2E expectations before implementation: Build publishes its own copied archive independently of an existing source scope, later source mutations stay on that source, failed builds preserve existing/absent destinations, and temporary staging is removed. The original scope test reproduced output deferral before the implementation changed.
- [x] Implement Build as a streamed copy to a temporary archive in the destination directory, call the ordinary Apply workflow on that copy, then atomically rename it into place after success. Remove mutable output retargeting and retain source-alias protection.
- [x] Document source-only scope publication and staged build publication; select exact-candidate quality, semantic lint, native/CLI and documentation gates. Delivery receipts own publication and review evidence.

## Collection command group and aliases

- [x] Carry the prewritten isolated command expectations from #164 before changing the command tree: use col with alias collection, retain config with alias conf, and reject removed root spellings.
- [x] Group native collection and state commands under col, retain namespace flags there, and update native CLI calls, help and examples consistently with #164.
- [x] Verify the grouped command tree together with staged builds and select fresh exact-candidate delivery gates. Delivery receipts own publication and review evidence.

## Application handlers

- [x] First adapt isolated CLI and native App callers to Config and Collection handler accessors; keep workflow dispatch, errors, context, state and file behavior expectations intact.
- [x] Reduce App and the CLI root interface to Config() and Collection(). Define handler contracts and move config/collection implementation and command factories to their owning files, retaining zero-value construction and scoped native ownership.
- [x] Document handler ownership and verify grouped CLI, native workflows and staged builds where implemented. Delivery receipts own exact-head checks and publication.

## Separate state workflows

- [x] Move state orchestration into its own handler, sharing the existing collection registry, prepared connections, namespaces and collection lifecycle through Collection().State().

## Public application contract

- [x] Adapt existing isolated CLI and native scope cases before implementation to use pkg/ankiascodeapi contracts, including caller-defined paths, nil/empty references, foreign ownership and expired references.
- [x] Move App and handler interfaces to pkg/ankiascodeapi and the implementation/native E2E to pkg/ankiascode with only protobuf and standard-library dependencies. Make CollectionRef an interface, retain private managed-reference state and validation, and update callers and Bazel dependencies without aliases or adapters.
- [x] Update documentation and specifications, verify native artifacts and the API dependency boundary, and select exact-candidate delivery gates. Receipts under out/application-api own validation and publication evidence.

## Flat namespaced application API

- [x] Adapt existing isolated CLI and native workflow callers before implementation to one App interface with Config, Col and ColState method prefixes; preserve command behavior, scoped references, namespace forwarding and failure expectations.
- [x] Remove handler interfaces, accessors and the separate API package; keep concrete App and collection reference contracts in pkg/ankiascode and define the application interface privately in the CLI. Implement namespaced operations directly on App and retain shared collection lifecycle and service boundaries. Update current documentation and specifications on both application branches.
- [x] Verify isolated CLI dispatch, native artifacts and the shared public package; select strict OpenSpec and exact-candidate delivery gates. Receipts under out/flat-application-api own final validation and publication evidence.

## Direct collection services and isolated CLI package

- [x] Adapt existing CLI and native workflow cases before implementation to plain collection paths; reuse the collection.Store E2E cases for direct prepared services, lifetime, aliases, namespaces and source publication.
- [x] Remove collection references, lookup IDs and App lifecycle wrappers. Reuse collection.Store and its injected Collection directly, keeping ordinary App operations path-based and existing service business rules intact.
- [x] Move Cobra commands and their private application interface to internal/cli; keep cmd/anki_as_code as composition and verify the CLI package does not depend on pkg/ankiascode or native backend packages.
- [x] Update docs and OpenSpec, inspect representative native artifacts, and select exact-candidate checks. Delivery receipts under out/direct-collection-cli own final validation and publication evidence.

## Configuration-set export command

- [x] Define command failure cases and a subprocess JSON/file/stdin/TOML round-trip artifact before implementation.
- [x] Add ConfigExport and config export through the existing codec and config.Store.Export; keep native lifecycle and semantic validation outside this operation.
- [x] Document the command, verify focused behavior and representative output, and select exact-candidate delivery gates on both application branches. Delivery receipts under out/config-export-cli own final validation and publication evidence.

## Prepared connection callbacks

- [x] Adapt existing native cases before implementation to locally reused services on the callback connection; retain persistence, cleanup, source-alias and publication coverage.
- [x] Carry the simplified collection.Store boundary from #165, remove Collection/service caches and construct application services on the injected grpc.ClientConnInterface.
- [x] Update current documentation and select exact-candidate native workflow, artifact, quality, lint and build gates. Receipts under out/simplify-collection own check execution and publication.

## Synchronize after collection-store merge

- [x] Rebase onto merged #165, preserve its typed App state-workflow E2E in the moved package and retain the feature's CLI/application boundaries.
- [x] Resolve visibility/relocation conflicts, update current ownership documentation and select exact-candidate quality, lint, native/API/CLI and site-build gates. Review and delivery records under out/rebase-after-165 own their results.

## Synchronize after CLI integration merge

- [x] Rebase onto merged #164 and remove merged CLI/export integration from the remaining diff.
- [x] Preserve both upstream export E2E and reconciliation export/convergence cases, resolve duplicate test/helper declarations, and retain the implemented plan/apply/build methods over merged placeholders.
- [x] Inspect the reduced merge diff and select exact-candidate quality, lint, native workflows and integrated documentation gates. Records under out/rebase-after-164 own validation, review and publication results.

## Correct post-conversion card assignment

- [x] Extend the existing native kind-conversion E2E before implementation to convert and create assigned cards in the same apply, cover both database/package inputs and verify inherited/explicit destinations and convergence.
- [x] Classify generated cards using resulting deck kinds while preserving resident filtered-card exclusions.
- [x] Verify exact-candidate native artifacts, quality, lint, builds and review; publish the fix in #128.

## Synchronize after manager removal

- [x] Rebase onto merged #167 and use its directly constructed native owners without restoring manager APIs.
- [x] Preserve verified lifecycle completion records from the independent manager-removal worktree and select exact-candidate native, quality, lint and documentation gates. Records under out/rebase-after-167 own revalidation and publication status.

- [x] Extract and verify pure planning with caller-supplied configs, snapshot and state in independent master-based PR #168; integrate it without duplicate models or state reads.

## Simplify reconciliation helpers

- [x] Split resource resolution and comparison by resource type, use native note identities without GUID indexing, and remove post-apply identity enumeration.
- [x] Preserve untracked template deck references, align the state-operation requirement with pure planning, and verify native planning/application artifacts before delivery.

- [x] Align plan counter documentation with deferred native card generation and structural migrations; known deletion counts cover deleted notes only.

- [x] Reproduce cross-note-type field/template transfers before restricting member moves to same-type renames; simplify the typed map helper and verify native identity retention.

## Executable planning boundary

- [x] Write native E2E cases for detached protobuf plans, execution without resource discovery, dynamic identity binding, stale-state rejection, input immutability and convergence before implementation.
- [x] Move pure planning to internal/plan and return a protobuf execution plan with prepared resource operations and a change report. Keep runtime resolution models private.
- [x] Make reconciliation execute supplied plans without configuration resolution or diffing. Let App compose capture, build and execution in one collection scope.
- [x] Update package docs and OpenSpec; select native artifact, quality, lint, consumer and documentation gates. Exact-candidate receipts own execution and publication.

## Planner review corrections

- [x] Reproduce native member reordering, removed stable identities and unnecessary converged assignments before implementation in the extracted planner.
- [x] Use path references, expected state and proposed bindings; prepare operations and reports together and consume actual bindings without temporary-ID maps.
- [x] Validate the native saved-plan and full application workflows, update the normative package ownership requirement and select exact-candidate delivery gates.

- [x] Reproduce and reject replacement/removal of a template used by resident filtered cloze cards; preserve explicit identity moves, handle all cloze ordinals through template zero and select exact-candidate consumer gates.

- [x] Add and inspect a maintained planning diagram showing App borrowing a collection connection, config/resource/state capture and pure plan construction; publish themed SVGs and a dark WebP preview.
