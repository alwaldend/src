## 1. Extract packages

- [x] 1.1 Copy validation and its required format package without semantic changes.
- [x] 1.2 Document the available package boundaries while retaining placeholder commands.

## 2. Verify and deliver

- [x] 2.1 Compare copied source with PR #128 and select independent compilation, lint, and repository quality checks; delivery receipts own exact-candidate validation.
- [x] 2.2 Select publication against master and document the prerequisite relationship; delivery receipts own publication and follow-up metadata outcomes.

## 3. Review corrections

- [x] 3.1 Reproduce duplicate, zero, and negative note identities and conflicting normal/empty-filtered declarations through a public-package workflow, retaining valid-state and filtered round-trip checks.
- [x] 3.2 Reject invalid note identities and inspect filtered-field presence; preserve independent package boundaries.
- [x] 3.3 Select corrected-candidate checks and guarded publication/review replies; delivery receipts own validation, replies, and publication outcomes.

## 4. Deck codec boundaries

- [x] 4.1 Use config terminology and separate editable protobuf deck conversion from persistence-model mutation.
- [x] 4.2 Verify the existing workflow against the pure codec and publish/reply to the three review comments.

- [x] Remove custom field-cache HTML conversion and document Anki-owned rebuilds; delivery receipts own validation and publication.

- [x] Expose the pinned upstream filtered schema as readable TOML; select workflow, codec, and API consumer checks.

- [x] Reuse upstream regular-deck configuration and audit remaining messages for equivalent contracts.

- [x] Reuse upstream note-type and template configs as nested fields; preserve appearance-only management and verify rejection of unrelated edits.

## Data structure review

- [x] Rename NoteTypeAppearance and TemplateAppearance to settings terminology and remove duplicate runtime settings, ordinal, and field-name representations.
- [x] Review serialized/shared/private structure contracts and select the independent workflow and required gates; delivery receipts own validation and publication.

Reviewed all serialized messages and shared models, plus adapter-private,
command, and fixture structures in the implementation. Removed duplicated
settings snapshots, ordinal maps, and per-note field-name arrays. Named note
fields and positional Bodies remain source and derived persistence projection;
deck Kind retains original wire bytes for lossless preservation. Collection
manifests and plans remain CLI contracts while database configs reuse upstream
protos. The later exported-symbol audit also identifies an unused note source-path field.

## Validation review follow-up

- [x] Replace field-specific unknown-data copying with recursive preservation and reject edits that discard unknown nested configuration; exercise failure cases before implementation.
- [x] Store note-type settings in an explicit field and split desired-state, deck, note, card, and settings validation into focused helpers.
- [x] Remove cloze-content inference; validate explicit cards against desired note types and preserve existing card identities and scheduling during structural edits.
- [x] Add note-type creation, renaming, field/template structure and configuration editing, and removal of unreferenced types; select schema-marker, preservation, filtered-card, and convergence checks for both candidates; delivery receipts own exact-candidate validation and guarded replies.

## Exported-symbol audit follow-up

- [x] Keep note-body validation private and remove the unused note source-path field.
- [x] Verify remaining handwritten exports have real consumers or required contracts, preserve the integration use of `PreserveUnknown`, and validate the exact candidate.
