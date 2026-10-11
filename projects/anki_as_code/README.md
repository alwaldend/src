---
title: Anki as code
description: CLI scaffold for managing Anki collections as TOML
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code is a Go CLI project for managing Anki collections through editable
TOML. This scaffold provides command discovery, flags, cancellation, and output
routing. Collection operations are placeholders: valid operation invocations
exit with an explicit "not implemented" error and do not read or write collections.

## Commands

| Command       | Arguments                                       | Intended operation              |
| ------------- | ----------------------------------------------- | ------------------------------- |
| `export`      | `--input`, `--output`                           | Export a collection to text     |
| `plan`        | `--input`, `--text`                             | Inspect desired changes         |
| `apply`       | `--input`, `--text`                             | Reconcile an offline collection |
| `build`       | `--base`, `--text`, `--output`                  | Build a collection archive      |
| `generate-id` | One or more note paths; optional `--collection` | Assign note and card identities |

All listed flags are required except `--collection`. Help and future plan JSON
use stdout; status messages and errors use stderr. Each subcommand supports
`--help`.

## Build and run

From the repository root:

```sh
bazel_agent bazel build //projects/anki_as_code:anki_as_code
bazel_agent bazel run //projects/anki_as_code:anki_as_code -- --help
```

## Project layout

- `cmd/anki_as_code/`: Cobra commands, flags, and output handling.
- `internal/app/`: App struct and placeholder collection operations.
- `site/content/`: Landing content for the shared website.
- `openspec/`: Project specifications and maintained changes.

The collection implementation is developed separately in
[PR #128](https://github.com/alwaldend/src/pull/128). Its end-to-end CLI test
requires real note-identity generation and belongs with that implementation.

## Validation and format packages

`internal/validation` checks desired notes, cards, deck identities and hierarchy,
and complete desired note-type definitions against a collection snapshot. It
resolves new note-type identities and validates ordered fields and template
structure. It rejects identity conflicts outside the managed contract, including
nonpositive or duplicate note IDs and conflicting normal/filtered declarations.
`internal/ankiformat` encodes upstream deck configuration, retains unmanaged
protobuf members and folds names for comparisons.
Deck conversion takes and returns the editable `collection.Deck` protobuf
declaration and the encoded upstream payload; it does not mutate a persistence
model. Callers own identity, marker paths, and assigning the converted payload.
Both use the merged collection API and runtime models; neither depends on
SQLite, archive handling, or filesystem storage.

These packages are available independently of the collection implementation
in [PR #128](https://github.com/alwaldend/src/pull/128). Commands remain
placeholders until that implementation merges. Build the packages with:

```sh
bazel_agent bazel build //projects/anki_as_code/internal/validation //projects/anki_as_code/internal/ankiformat
```

Run the independent validation and deck-codec workflow with:

```sh
bazel_agent bazel test //projects/anki_as_code/internal/validation:workflow_test
```

The workflow emits `validation-workflow.json` in Bazel test outputs.

## Rebuild derived note caches in Anki

After importing an archive with changed note fields or sort-field settings, run
**Tools → Check Database** in Anki before relying on browser sorting or duplicate
detection. The CLI preserves existing notes' sort-field (`sfld`) and first-field
checksum (`csum`) caches and initializes new notes with empty/zero caches. It
does not interpret HTML or rebuild these derived values.

Stale caches do not change card rendering, review history, or scheduling; cards
use the updated field content. Check Database rebuilds the caches using Anki's
own rules, but also performs broader database repairs and may generate missing
cards. This step is required for the collection implementation in PR #128;
placeholder commands do not modify collections.

## Filtered deck configuration

Filtered deck markers expose Anki's upstream settings directly:

```toml
id = 2
title = 'Filtered'
notes_path = 'notes'
[filtered]
reschedule = true
preview_delay = 10
preview_again_secs = 60
preview_hard_secs = 120
preview_good_secs = 0
delays = []
[[filtered.search_terms]]
search = 'deck:grammar is:due'
limit = 100
order = 'RANDOM'
```

Declare exactly one of `[normal]` or `[filtered]`. Filtered configuration is
editable; resident card membership, original home decks, and temporary card
scheduling remain owned by Anki. Changing a search does not rebuild the filtered
card set automatically: rebuild it in Anki to apply the search to membership.
The old `filtered_base64` declaration is unsupported; re-export existing text.

## Editable schema ownership

Regular and filtered deck configuration reuse Anki's pinned `Deck.Normal` and
`Deck.Filtered` protobuf messages directly. All known settings are readable
TOML, including optional daily limits and retention overrides for regular decks.
The preset referenced by `config_id` remains unmanaged. Unknown protobuf fields
are retained from the base archive; `extra_fields_base64` is no longer accepted.
Re-export old text containing that field before editing it.

The collection config is CLI-owned: format version and filesystem paths have
no upstream Anki equivalent. Notes provide named fields and managed card
declarations instead of Anki's positional field list and sync metadata. Card
declarations exclude scheduling state. Note-type declarations reuse upstream
configuration and ordered field definitions. Template ordinals are explicit;
validation does not parse HTML or infer cloze cards. Current snapshots index the
generated NoteTypeSettings directly, without a runtime note-type wrapper.

The implementation in #128 supports creating/renaming/removing note types and
editing fields, templates, and configuration. All affected note fields and
explicit card declarations must agree with the desired schema. Field list order
defines storage order; template ordinals must be contiguous from zero. Scheduling
and review history are preserved, and structural changes mark Anki's schema for
synchronization. Resident filtered cards require unchanged template position and
note-type kind. Archive bytes and synchronization metadata stay repository-private.

Recursive protobuf preservation retains unknown data in corresponding nested
messages. Removing a nested message or list item carrying unknown data fails.
The independent workflow emits deck-preservation-workflow.json as well as its
identity/codec artifact. Commands remain placeholders in this PR.
