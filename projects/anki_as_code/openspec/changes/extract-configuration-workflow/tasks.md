## 1. Configuration extraction

- [x] 1.1 Establish an independent master-based worktree, retain protobuf contracts, and record extraction scope.
- [x] 1.2 Consolidate configuration-file operations in internal/config, remove state-dependent modules from the extraction, and update app/build references.
- [x] 1.3 Keep syntax/schema/filesystem checks in configuration parsing and defer collection semantics; build grouped output from declaration ownership rather than titles.

## 2. Configuration commands

- [x] 2.1 Add config dump with JSON default and selectable TOML through CollectionConfigSet/DeckConfigSet protobuf messages.
- [x] 2.2 Wire generate-id to configuration manipulation and stage the complete replacement batch before publication.
- [x] 2.3 Verify semantic-invalid configurations remain inspectable, colliding deck titles retain separate notes, unfinished notes support identity generation, and file/schema errors still fail. Write behavior cases before implementation and emit repeatable artifacts.

## 3. Documentation and delivery

- [x] 3.1 Update current package/validation boundaries in README and landing documentation and validate OpenSpec strictly.
- [x] 3.2 Review the extraction, select required quality/lint and built-CLI/site checks, and prepare publication and applicable review replies for #151. Delivery receipts own exact-candidate validation, publication, and review facts.

The latest user decision replaces configurable paths with a fixed hierarchy; titles are derived from the hierarchy. Earlier identity-staging failure evidence remains relevant. Planning and target-validation failure evidence belongs to the reconciliation scope, which is now excluded from this extraction.

Boundary evidence: the prior CLI rejected duplicate references, colliding titles, invalid titles, and unfinished notes before the revised parser. The revised workflow accepts those declarations without changing sources and retains their original deck groups. Separate presence cases exposed and corrected explicit empty override replacement; schema/file rejection and identity staging checks continue to pass.

- [x] 3.3 Use the current schema with an explicit version; the unreleased version-1 correction below supersedes earlier-version header handling and migration guidance.
- [x] 3.4 Document protobuf numeric input mapping without duplicating schema validation. Prior export failure evidence remains in delivery scratch; snapshot export is now deferred by user decision.
- [x] 3.5 Remove snapshot export and its helpers/tests, unused normalization APIs, and the duplicate desired-state representation; use explicit protobuf fields in remaining shared models.
- [x] 3.6 Move generate-id under config, keep path expansion private, and review standard-library walking and numeric output conversion. Verify CLI behavior before implementation.
- [x] 3.7 Update documentation, review the full #151 diff, and select exact-candidate gates and reasoned thread replies. Delivery receipts own publication and resolution state.
- [x] 3.8 Reproduce dangling note-root and intermediate-link acceptance before correcting missing-root discovery; verify real missing roots remain empty and preserve source files. Delivery receipts own exact-candidate check results.
- [x] 3.9 Reproduce datetime coercion and final-int64 allocation failures; reject unsupported TOML datetimes during JSON conversion and fix inclusive identity allocation without wrapping arithmetic.
- [x] 3.10 Verify collection-root card deck absence in JSON/TOML before fixing default inheritance to require a nonempty declaring deck.

## 4. Fixed hierarchical layout

- [x] 4.1 Write built-CLI behavior cases for nested decks, isolated note ownership, fixed flat resource lists, missing markers, removed path/title fields, and identity reservation across the hierarchy before implementation.
- [x] 4.2 Remove configurable path/title fields, derive titles from directory names, and load only fixed decks/, notes/, and note_types/ locations with direct directory reads and recursion through child decks/. Preserve default inheritance and root-note presence.
- [x] 4.3 Update format compatibility, README/landing/package docs, verify repeatable output, review the diff, and select required exact-candidate gates for #151. Delivery receipts own publication and review state.

## 5. Uniform schemas and simplification

- [x] 5.1 Extend the existing CLI workflow to inspect nested note-type configuration in JSON/TOML before changing the codec; reuse identity/staging behavior checks.
- [x] 5.2 Rename editable protobuf messages consistently and update their actual Go consumers and docs; consolidate TOML encoding, remove the staging callback, and simplify selected-note/GUID handling.
- [x] 5.3 Review the full #151 diff, fix in-scope defects, verify representative outputs, and select exact-candidate gates. Delivery receipts own publication and remote review facts.

## 6. Stable note GUIDs

- [x] 6.1 Write CLI behavior cases for file/collection/deck selections, overlapping inputs, missing-only generation, byte-preserved existing GUIDs, and removed numeric note/card fields before implementation.
- [x] 6.2 Remove numeric note and card identities and note-type IDs from note configs; rename the command to config generate-guid and traverse only fixed config directories, preserving already assigned GUIDs.
- [x] 6.3 Update schemas, documentation, and OpenSpec; review and validate the exact candidate before publishing #151.

## Unreleased version 1

- [x] Update version-1 config dump and GUID fixture expectations before implementation; replace legacy-layout cases with current-schema/version checks and verify export through existing native workflows.
- [x] Keep format version 1, remove the earlier-version header pass and protobuf compatibility reservations, and align docs/specs without migration guidance.
- [x] Review the correction, validate exact-candidate behavior and required quality/lint/build gates, and publish the updated PR.

- [x] Replace handwritten valid fixture setup with protobuf messages and the existing TOML encoder, preserving independent parser checks and deliberately invalid syntax inputs. CLI behavior passes; delivery receipts own publication and review resolution.

## Public configuration loader

- [x] Export Load with its existing protobuf return type, route dump and GUID generation through it, and document the package API and validation boundary.
- [x] Verify the existing configuration CLI workflow and select required exact-candidate quality, lint, and build checks. Delivery receipts own validation, publication, and review state.

## Configuration object export

- [x] Add object-to-files-to-CLI E2E cases before implementation for nested layout, exact values, caller-message preservation, unchecked input, untouched other files, and filesystem errors; emit repeatable exported artifacts.
- [x] Expose Export as a direct TOML writer using existing codecs, without loader validation, normalization, GUID allocation, stale-file cleanup, or directory publication.
- [x] Document the writer contract, verify existing and new configuration workflows, and select required exact-candidate quality/lint/build gates. Delivery receipts own publication and review state.

## Direct JSON-to-TOML conversion

- [x] Update existing E2E expectations before implementation to retain ProtoJSON numeric strings in TOML while checking exact values, field content, presence, and missing-only GUID generation.
- [x] Remove descriptor traversal and scalar/integer conversion; encode ProtoJSON values directly through the existing TOML library and remove unused imports/dependencies.
- [x] Align current docs/specs, verify existing workflows and artifacts, and select exact-candidate delivery gates. Delivery receipts own validation and publication.

## Collection version field

- [x] Update existing configuration E2E fixtures/expectations before renaming the root protobuf field and Go consumers to version.
- [x] Align examples and deck-creation identity documentation, verify config/native workflows and representative output, and select exact-candidate delivery gates.

## Collection-backed path identities

The path-based version-1 contract supersedes historical identity/generation work above.
Native state behavior is maintained in #128's collection-text-workflow change.

- [x] Write configuration E2E failure cases before implementation, preserving codec and filesystem coverage.
- [x] Replace identity-bearing configs with path wrappers/readable references and remove GUID generation and its helpers.
- [x] Keep parsing/writing pure; define shared state contracts while leaving state persistence, import, and reconciliation to #128.
- [x] Align docs with the current APIs and inspect repeatable JSON/TOML/export-directory artifacts from passing E2E.
- [x] Complete source/output review and select exact-candidate quality/lint/build gates for guarded publication of #151 against master. Delivery receipts own final validation/publication status.

## Configuration package boundaries

- [x] Split filesystem-independent codecs and fixed filesystem layout into child packages, retaining the configuration facade and existing behavior in one Go module.
- [x] Update package/API docs and Bazel ownership, verify existing config CLI E2E and representative outputs, and select exact-candidate delivery gates. Delivery receipts own publication status.

Package-boundary evidence: existing config CLI E2E passes, preserving exact
values/presence, source bytes, hierarchy/defaults, syntax/schema errors, datetime
policy, and unchecked direct object export. The matching #128 consumer passes
native/CLI E2E and the 7,218-note collection retains a zero-change plan. The
codec contains no filesystem access; Bazel visibility restricts both child
packages to config implementation consumers. Evidence is under
`out/anki-config-modules/`; delivery receipts own final gate/publication status.
