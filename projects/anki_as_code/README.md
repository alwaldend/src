---
title: Anki as code
description: Export and reconcile Anki collections as TOML
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code is a Go CLI that turns an existing collection into editable TOML.
Export can refresh an existing directory: it replaces managed declarations with
the selected collection, removes stale declarations, and preserves unrelated
files. Export generates the default layout; reconciliation follows configured
paths. Paths accept current-user `~` and `~/...`, including paths in configs.
An explicit root config file may have any name; directory inputs select
`collection.anki.toml`.
Status messages and errors use stderr; plan JSON and help use stdout.
Change note fields, tags, cards, decks, or card appearance, inspect the
plan, and reconcile the collection with that desired state.

It supports modern `.colpkg` archives containing a schema-18
`collection.anki21b` database and offline schema-18 SQLite collections such as
`collection.anki2`. Existing scheduling, note type structure, scheduling presets, and archive
media are preserved. Card CSS and template formats are editable. Review history stays unchanged in the collection and is excluded from text.

## Managed and unmanaged entities

The text describes the complete desired state for these managed entities:

- **Notes:** add/remove notes and edit named fields and tags. Fields remain
  shared by every card generated from that note.
- **Regular-deck cards:** add/remove explicit card declarations and change their decks.
  The CLI validates identities and template ordinals against the base collection.
- **Decks:** add/remove/rename decks through explicit titles and edit normal or
  filtered deck configuration. Directory names do not define Anki titles.
- **Card appearance:** edit note type CSS and existing templates' front/back
  HTML, browser formats, and browser font settings.

These entities and properties are not exposed as editable desired state:

- **Cards resident in filtered decks:** omitted from text and preserved unchanged.
  Filtered membership, home assignments, and temporary scheduling belong to Anki.
  Empty the filtered deck in Anki before managing those cards or removing a
  referenced note/deck.

- **Collection settings:** omitted from text and preserved in the database, as
  requested. The CLI only advances the computed `nextPos` cursor when adding cards.

- **Scheduling presets:** retained byte-for-byte, as requested. Appearance
  edits must not change learning steps, daily limits, or scheduling algorithms.
- **Existing card scheduling and state:** due dates, intervals, ease, review
  counts, suspension/queue state, flags, and opaque card data are preserved to
  retain study progress. Newly declared cards start with new-card scheduling.
- **Note type and template structure:** IDs/names, field definitions/order,
  note type kind, and template identities/ordinals stay fixed. Changing them
  requires schema migration and card generation beyond this CLI's editing scope.
- **Other note type and field settings:** LaTeX configuration, editor fonts,
  sticky fields, and unrelated protobuf fields remain unchanged so editing
  appearance does not reset other Anki behavior.
- **Media payloads:** copied from the base archive, preserving files used by
  note references. File references can change in fields; media file management
  is outside the text reconciliation surface.
- **Anki UI and synchronization bookkeeping:** collapsed tag/deck state and
  unrelated metadata are retained. Change markers and deletion tombstones are
  computed when needed, rather than accepted as user-authored desired state.
- **Review history:** preserved unchanged in the collection and omitted from text,
  because it records actual reviews rather than desired configuration.

Template edits do not automatically create or remove cards. Card declarations
in note files remain explicit so a rebuild has a reviewable resource plan.

## Build and run

Run from the repository root:

```sh
bazel_agent bazel build //projects/anki_as_code
bazel_agent bazel run //projects/anki_as_code -- --help
```

Use absolute file paths with `bazel run`, which can start the executable outside
the repository directory. Close Anki before using a SQLite collection; files
with a WAL or rollback journal are rejected.

## Export, edit, and rebuild

Keep the original archive as the base. Export creates a new destination or
refreshes an existing one while preserving unrelated files. This workflow uses
the same archive throughout:

```sh
mkdir -p "$PWD/out/anki"
cp /absolute/path/collection.colpkg "$PWD/out/anki/base.colpkg"

bazel_agent bazel run //projects/anki_as_code -- export \
  --input "$PWD/out/anki/base.colpkg" --output "$PWD/out/anki/text"

# Edit the exported TOML files, then inspect the proposed changes.
bazel_agent bazel run //projects/anki_as_code -- plan \
  --input "$PWD/out/anki/base.colpkg" --text "$PWD/out/anki/text/collection.anki.toml"

bazel_agent bazel run //projects/anki_as_code -- build \
  --base "$PWD/out/anki/base.colpkg" --text "$PWD/out/anki/text/collection.anki.toml" \
  --output "$PWD/out/anki/edited.colpkg"
```

`plan` prints JSON additions, updates, and deletions. Review history remains
unchanged. Plan does not change either input. `build` creates a new archive from
the selected compatible base and desired text; it leaves the base unchanged.
All other ZIP entries are copied verbatim.

To reconcile an existing archive or offline SQLite file in place, use `apply`:

```sh
cp "$PWD/out/anki/base.colpkg" "$PWD/out/anki/working.colpkg"
bazel_agent bazel run //projects/anki_as_code -- apply \
  --input "$PWD/out/anki/working.colpkg" --text "$PWD/out/anki/text/collection.anki.toml"
```

Apply uses a private database copy, checks integrity and convergence, checks the
input hash immediately before replacement, and atomically replaces the target
while preserving its permission bits. Archive compression finishes before the
final hash check. Keep Anki and other writers closed: the check detects observed
changes but cannot atomically exclude an uncooperative writer. Text is independent
of archive checksums for export, plan, apply, and build. Legacy `base_sha256`
metadata is accepted and ignored. Repeated reconciliation produces an empty plan.

The text describes the complete desired collection. Removing a note file
removes its note and cards; removing card or deck declarations removes
those resources. Deletions record Anki synchronization tombstones. A filtered
deck cannot be deleted while cards still occupy it.
An occupied filtered deck also cannot be converted to a normal deck; empty it
in Anki first. Deck title swaps and reuse of removed titles retain declared IDs.

## Config-driven discovery

`collection.anki.toml` is the root config. Plan, apply, and build accept its path
or its containing directory. Paths resolve relative to the config declaring
them; absolute paths are also supported. The export starts with this layout,
which you can rearrange freely:

```text
collection.anki.toml
note_types/
  Chinese.notetype.anki.toml
decks/
  parent/
    deck.anki.toml
    child/
      deck.anki.toml
      notes/
        greeting.note.anki.toml
```

The root config declares the locations to read using flat path keys:

```toml
format_version = 2
decks_path = "decks"
note_types_path = "note_types"
notes_path = "notes"
```

`notes_path` in the root is optional and contains cardless notes. The CLI
recursively searches `decks_path` for files named exactly `deck.anki.toml`.
Each marker declares its Anki title and its own notes location:

```toml
id = 0
title = 'parent::child'
notes_path = '../../cards/child'
[normal]
config_id = 1
extend_new = 0
extend_review = 0
description = ''
markdown_description = false
```

The exporter initially maps `parent::child` to `decks/parent/child` and escapes
unsafe filesystem characters. The loader uses `title` exactly as declared.
Moving a marker or changing its directory name does not rename the Anki deck.
A new deck can use `id = 0`; successful apply/build records its allocated ID in
the marker. All Anki parent decks must be declared, wherever their markers live.

Notes are discovered recursively under each declared notes path and must be
named `<name>.note.anki.toml`. Other files are ignored. Missing notes directories
represent empty decks; missing configured deck/note-type roots are errors.
Symlinks and a note referenced by multiple decks are rejected.

Each note lists shared fields and cards. A card's optional `deck` override uses
an exact Anki title. Without an override, it belongs to the marker declaring
that notes path. Moving notes between declared roots changes that default;
renaming a file or rearranging referenced directories preserves its identity.

Card destinations must be normal decks. Cards currently in filtered decks are
not exported or reconciled. Their notes remain editable, but deleting a note
or deck they reference or declaring a conflicting card is rejected.
Existing note IDs, GUIDs, note types, and card template ordinals remain fixed.
Editing an existing cloze note must retain its cloze ordinal set.

## Note fields and card appearance

Keep the metadata and every named field retained from an exported or copied note,
including empty fields. Fields retain the original Anki HTML and media
references. Notes, decks, root configs, and appearance files are encoded as
complete objects by the TOML library, using its standard quoting and layout.
String values round-trip exactly, including quotes, backslashes, controls, and
leading or trailing newlines. A note with Front and Back fields can contain:

```toml
[fields]
Front = "<div class=\"example\">Hello</div>"
Back = "First line\nSecond line"
```

You can use any valid TOML string syntax when editing these values. Markdown
conversion and external field files are unsupported.

Normal deck settings expose the configuration profile ID, extra daily limits,
and description. Deck titles, paths, descriptions, and opaque string settings
use the same TOML encoder as note fields. Recording a newly allocated deck ID
re-encodes the complete marker. Filtered settings and unknown deck protobuf fields are retained
as base64. `note_types_path` references a directory containing one
`*.notetype.anki.toml` file per existing note type. Each file has root `id`, `name`,
and `css`, plus `[[templates]]` entries with ordinal, name, front/back HTML,
browser formats, and browser font settings. Files are discovered recursively;
their names and nesting do not define identities. Missing or duplicate note type
declarations fail before mutation. CSS and template content use the same
lossless TOML encoding as note fields. Collection settings are neither
exported nor reconciled.

Database configuration uses Go messages generated from Anki's unmodified
26.09.3 schemas, fetched from an integrity-pinned upstream revision during the
Bazel build. The schemas and generated Go files are not copied into the project.
Appearance edits assign only the relevant message members. Note type identities,
field structure, and template ordinals stay fixed; unknown configuration and
unedited database bytes are retained. Template edits do not generate or remove
cards: the note files declare them explicitly. Scheduling presets are not
exported or modified. Runtime `nextPos` allocation state remains in the database
and advances when new positions are allocated. Media stays in the base archive.

## Generate note identities

Copy an exported note with the required note type and card declarations, edit
its fields, then assign fresh identities. The command accepts one or more note
files or directories and recursively finds note TOML in directories:

```sh
cp "$PWD/out/anki/text/decks/parent/child/notes/existing.note.anki.toml" \
  "$PWD/out/anki/text/decks/parent/child/notes/greetings.note.anki.toml"
bazel_agent bazel run //projects/anki_as_code -- generate-id \
  "$PWD/out/anki/text/decks/parent/child/notes/greetings.note.anki.toml"
```

Replace the example paths with paths in your export. For a batch, pass a
directory containing the copied notes. Overlapping paths are processed once;
only `*.note.anki.toml` files are selected when scanning a directory.
Malformed notes, missing paths, and symlinks fail explicitly. The command
parses the batch before writing, and reserves identities from the containing
export when its `collection.anki.toml` is present. For notes outside that
config directory, pass `--collection /path/to/collection.anki.toml` to reserve
identities from every referenced note root.

Each selected note receives a new numeric note ID, a random GUID, and new IDs
for its existing cards. Filenames, named fields, tags, note types, card ordinals,
and deck overrides remain unchanged. The command does not create note files or
render templates to discover cards. Edit the copied note's fields, inspect the
plan, then apply or build.

Run this on copies when adding notes. Regenerating identities on an existing
note makes reconciliation delete its previous note/cards and add new ones;
those new cards receive new scheduling. Repeating the command generates fresh
identities again. Renaming a file by itself keeps its identities.

## Architecture

The CLI uses focused Go packages within one repository module.
`api/collection/collection.proto` owns serialized field names and types.
TOML is decoded into a value map, marshaled as JSON, and parsed as a generated
message through ProtoJSON. Export uses ProtoJSON and converts integer fields to
TOML numbers without floating-point conversion. Plan output uses ProtoJSON with
original field names; 64-bit counts are JSON strings.
Generated messages are used directly for configurations, cards, appearance,
and plans. Notes and decks wrap those messages with runtime metadata;
`internal/model` does not duplicate serialized fields. Arrows show
internal package dependencies; external libraries and test fixtures are omitted.

![Anki CLI package dependencies](assets/architecture.svg)

The command tree calls `app`, which delegates to `collection`. Collection
coordinates text discovery, archive resources, database reconciliation,
and planning. Repository combines shared validation rules with guards for
unmanaged database resources. Model contains runtime wrappers and shared snapshots. Repository-only preservation bytes and scheduling cursors remain private to repository. Validation owns desired
state rules; planning compares current and desired state; ankiformat encodes Anki
configuration and normalizes names and field-cache text. The storage adapters
do not depend on each other. Archive handles files without knowing collection
entities or SQLite.

| Package               | Responsibility                                                         | Internal dependencies                                       |
| --------------------- | ---------------------------------------------------------------------- | ----------------------------------------------------------- |
| `api/collection`      | Protobuf contracts for editable TOML and plan JSON                     | None                                                        |
| `internal/model`      | Runtime wrappers and shared snapshots                                  | Collection schema                                           |
| `internal/validation` | Desired-state rules and validation                                     | Model, Anki format                                          |
| `internal/planning`   | Resource comparison and change plans                                   | Model, Anki format, collection schema                       |
| `internal/ankiformat` | Anki configuration codecs, name normalization, field-cache text        | Model, collection schema                                    |
| `internal/textstore`  | TOML encoding, discovery, identity generation, staged text exports     | Model, validation, Anki format, collection schema           |
| `internal/archive`    | ZIP/Zstandard resources, offline file copies, guarded file publication | None                                                        |
| `internal/repository` | SQLite setup, snapshots, transactions, all reconciliation SQL          | Model, validation, planning, Anki format, collection schema |
| `internal/collection` | Export, plan, apply, and build orchestration and resource lifetimes    | Planning, storage adapters, collection schema               |
| `internal/app`        | Operations consumed by the command tree                                | Collection, collection schema                               |

Bazel visibility restricts the storage adapters to collection orchestration.
The model cannot import storage adapters; text storage and archive handling
cannot import SQLite. Database handles remain private to the repository.
End-to-end fixtures open their own raw SQL handles to verify preserved rows.

## Validation

```sh
bazel_agent bazel test //projects/anki_as_code/internal/collection:collection_test
bazel_agent bazel build --config=lint //projects/anki_as_code/...
```

The E2E suite verifies field round trips, additions/deletions, filtered cards,
media and scheduling preservation, copied-note identity generation, arbitrary
configured layouts, CSS/template edits, and convergence. It emits verification artifacts in Bazel's undeclared test
outputs. There are no fixed CLI cutoffs for database or note-file sizes. Database
extraction and copying stream to scratch disk; TOML documents and collection
metadata use memory. Available disk and RAM govern practical capacity.
Unsupported schemas fail explicitly.

The project uses the repository's pinned TOML and Zstandard modules and the
approved SQLite driver v1.60.1. User collection data and its acceptance evidence
are maintained separately in the Anki export PR.
