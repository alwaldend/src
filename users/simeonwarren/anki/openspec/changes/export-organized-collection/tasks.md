## 1. Export and organize

- [x] 1.1 Preserve the verified TOML export, organize titled lessons, and document rebuilding with the companion CLI.
- [x] 1.2 Verify every note/card and lesson membership, original settings/templates/scheduling, immutable ZIP assets, unchanged review history, and an empty reconciliation plan.

## 2. Delivery

- [x] 2.1 Validate owner OpenSpec, inspect representative output, and select independent export build/formatting gates. The repo-delivery receipt owns exact-candidate validation, commit/push, and separate PR publication targeting master.

## Acceptance evidence

The export contains 7,218 notes, 7,485 cards, 445 titled grammar lessons, and
486 decks. Independent SQLite/ZIP/TOML inspection verifies every original field,
note/card identity, original deck, setting, template, scheduling column, and all
257 non-database ZIP entries. The rebuilt archive has no review records and
passes integrity_check. Planning it against the text reports no changes.
The untouched base archive SHA-256 is recorded in collection.anki.toml. Raw evidence
and binary archives remain in ignored task output. Desktop import was not tested.

The CLI implementation and SQLite/dependency policy upgrade belong to PR #128.
This data PR targets master independently and has no build dependency on that
unmerged project. Rebuilding uses the separately validated CLI candidate.

The collection README is packaged through the existing user documentation tree;
the apex site build is included in the exact-candidate delivery checks.

The README reflects the companion CLI revision: generate-id replaces new-note;
lesson grouping remains a completed one-off adoption rather than a CLI command.

## Config-driven format revision

- [x] Migrate collection/deck markers and named note suffixes without changing lesson titles, IDs, fields, or scheduling presets.
- [x] Include unchanged card CSS/templates in referenced settings and verify the rebuilt archive, full-table preservation, and empty plan.
- [x] Update the owner README and select the independent export receipt and gates; the receipt owns exact-candidate validation and publication state.

## Literal string formatting

- [x] Regenerate note fields and appearance with safe single-quoted and triple-double-quoted multiline TOML strings, preserving every parsed value and all archive state.

## Deck string formatting

- [x] Regenerate deck string quoting without changing any parsed deck setting or collection behavior.

## Separate note type documents

- [ ] Store only per-note-type CSS/template documents under the referenced note_types directory; omit collection settings and preserve all archive configuration.

## 14. Group resource path references

- [x] 14.1 Migrate collection and deck references to tables containing path keys, preserving all declared values.
- [x] 14.2 Rebuild and verify zero drift, unchanged archive settings, and complete review history.
