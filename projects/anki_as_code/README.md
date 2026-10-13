---
title: Anki as code
description: Inspect and write Anki configuration files
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code represents notes, decks, and note types as protobuf-backed TOML.
This change provides the pure configuration API and `config dump`. Collection
export, state import, planning, apply, and rebuilding belong to
[PR #128](https://github.com/alwaldend/src/pull/128); collection commands remain
explicit placeholders here. Configuration commands require neither Anki nor Rust.

## Configuration commands

```sh
bazel_agent bazel run //projects/anki_as_code -- config dump \
  --text "$PWD/path/to/collection.anki.toml"
bazel_agent bazel run //projects/anki_as_code -- config dump \
  --text "$PWD/path/to/collection.anki.toml" --format toml
```

`--text` accepts a root file or its containing directory. JSON is the default;
TOML is the other supported format. Dump emits the root `config`, deck groups
with derived `title` and marker `path`, and note/note-type entries with `path`
and `config`. Paths are collection-relative and use `/`. Notes in a deck inherit
its marker path for omitted card deck overrides; root notes have no default deck.
Explicit empty overrides stay empty. Default protobuf values are omitted;
optional zero values and empty named fields remain present.

Dump leaves files unchanged. It checks syntax, protobuf types, version, and
filesystem layout, without resolving references, validating Anki semantics,
normalizing content, allocating identities, or opening a collection. Unfinished
notes and definitions remain inspectable. Success uses stdout; errors use stderr.
Missing or unsupported subcommands fail.

## Fixed hierarchy

```text
collection.anki.toml
note_types/
  Basic.notetype.anki.toml
notes/
  unattached.note.anki.toml
decks/
  Languages/
    deck.anki.toml
    notes/
      introduction.note.anki.toml
    decks/
      French/
        deck.anki.toml
        notes/
          greeting.note.anki.toml
```

The root contains `version`, currently 1. Optional moves describe file, directory,
or nested definition address changes for reconciliation:

```toml
version = "1"
[[moved]]
from = 'decks/Old'
to = 'decks/Languages'
```

`decks/Languages/decks/French/` declares `Languages::French`. Every deck directory
requires `deck.anki.toml`. Titles come from the hierarchy and are absent from
markers. A normal deck references an existing scheduling preset by readable name:

```toml
preset = 'Default'
[normal]
description = 'Language lessons'
```

Notes reference note types by configuration path and contain exact named fields:

```toml
note_type = 'note_types/Basic.notetype.anki.toml'
[fields]
Front = 'Bonjour'
Back = 'Hello'
[[cards]]
ordinal = "0"
```

An omitted card deck inherits its declaring deck's marker path. Explicit card
and template overrides use paths such as `decks/Languages/deck.anki.toml`.
Note types declare names, upstream configuration, ordered fields and templates.
`sort_field` selects a field by name; template `target_deck` references a deck path.
Array order determines field/template positions.

Paths identify resources. Editable configs contain no top-level Anki IDs or note
GUIDs. A note at a new path declares a new resource; identity-preserving moves
belong in the root's `moved` list. PR #128 initializes native identity bindings
inside collection configuration on export and supports explicit state import.
The identity generation command has been removed.

The loader lists immediate files in `notes/` and `note_types/`, descending only
through child `decks/`. Other files and deeper note directories are ignored.
Missing resource directories are empty sets. Empty markers are parseable;
business logic determines whether they are usable. Symlinks, non-directory
resource roots, and missing deck markers fail.

The format is unreleased and stays version 1 while the schema evolves. Only the
current schema is supported. Normal/filtered settings and note-type presentation
reuse Anki messages. Native-only fields in reused messages are omitted by collection
export and rejected by reconciliation; pure dump/export retains supplied protobuf
values without semantic checks.

## Serialization and package APIs

TOML decodes to generic values, then JSON, then protobuf using ProtoJSON.
Unknown keys and incompatible values fail. Unquoted TOML dates/times fail rather
than silently becoming strings; quote dates used as note text. Output converts
ProtoJSON directly to TOML, preserving optional presence and exact text values.
64-bit values use quoted decimal strings, including `version = "1"`.

`internal/config` exposes:

- `Load(input string) (*collectionpb.CollectionConfigSet, error)` parses the hierarchy.
- `Dump(input, format string) ([]byte, error)` loads and encodes JSON or TOML.
- `Export(configs *collectionpb.CollectionConfigSet, output string) error` writes
  the root marker and supplied path/config pairs as TOML files.
- `IsConfigFile(name string) bool` recognizes declaration filenames for caller-owned cleanup.
- `Version` supplies the supported schema version to collection exporters.

Export writes directly, preserves input messages, overwrites declared files, and
leaves unrelated/stale files alone. Callers own valid, unique paths, layout,
validation, staging, cleanup, and publication. Export does not load or normalize
its output; failures may leave partial files. Consumers use this facade;
the child packages remain implementation APIs with restricted Bazel visibility.

`config/layout` owns the fixed hierarchy, discovery, file I/O, and direct file
export. `config/codec` converts TOML/JSON through protobuf without filesystem
access. Dependencies flow from config to layout/codec, and from layout to codec.
Collection relationship validation and native operations remain outside them.

## Managed boundaries

Configs express note content/tags, generated-card deck assignments, deck settings,
and note-type fields/templates/presentation. The config package parses and writes
those declarations. Anki owns persistence, native rules, normalization, structural
migrations, and card generation.

Scheduling presets, review history, card scheduling, media payloads, collection/UI
settings, and resident filtered cards are excluded from editable desired state.
They represent study progress, external resources, or Anki's operational state.
Collection-backed reconciliation state is application metadata, not a user setting.

## Build and architecture

```sh
bazel_agent bazel build //projects/anki_as_code
bazel_agent bazel test //projects/anki_as_code/cmd/anki_as_code:config_workflow_test
```

Use absolute paths with `bazel run`, which may execute outside the checkout.

| Package                  | Responsibility                                               |
| ------------------------ | ------------------------------------------------------------ |
| `api/collection`         | Generated resource, configuration, plan, and state contracts |
| `internal/config`        | Public configuration facade and dump operation               |
| `internal/config/layout` | Fixed hierarchy, discovery, file I/O, and direct export      |
| `internal/config/codec`  | Filesystem-independent protobuf/TOML/JSON conversion         |
| `internal/model`         | Shared runtime snapshot models, without business logic       |
| `internal/app`           | Configuration adapter and collection placeholders            |
| `cmd/anki_as_code`       | Cobra flags, cancellation, output, and errors                |

The loader produces protobuf config sets directly. Resource messages use `Config`,
file wrappers retain source paths, and upstream messages keep their upstream names.

## Anki backend API

The reusable Go backend and Rust server belong to
[projects/anki_grpc](../anki_grpc/README.md#go-backend). The Go package owns the
server process and exposes a managed gRPC connection. Callers use Anki's generated
service clients for typed methods. Anki owns collection validation, note updates,
card generation, and persistence; this CLI owns TOML and reconciliation.

## Integration approaches considered

Direct storage editing has the simplest integration: an application edits
SQLite and protobuf itself, taking responsibility for Anki's data invariants and
business logic. The API-based approaches are Python with Anki's official SDK,
Rust with Anki's Rust API, and another language with a translation layer around
the Rust API. The two Rust API approaches require building Anki from source.
A translation layer can be a C ABI or a gRPC sidecar; it does not remove that
build requirement.

This project uses a translation layer: the Go CLI calls the Rust server in
`projects/anki_grpc`, which forwards operations to Anki's Rust API. Go owns the
server process and its failures.

![Anki backend integration approaches](/docs/projects/anki_as_code/assets/backend-approaches.svg)

The maintained [Mermaid source](backend-approaches.mmd) is rendered by the
`:backend_approaches` Bazel target and published with the project documentation.

## Verification

The built-CLI E2E runs with an empty executable search path and emits repeatable
JSON, TOML, and exported-directory artifacts. It checks nested ownership, fixed
traversal, path references, removed identities/command, exact Unicode/multiline
fields, optional zero values, empty fields, source preservation, schema/filesystem
errors, and direct object export without validation or message changes.

This PR targets master independently. Merge #151 first, then rebase #128 to remove
the temporary prerequisite overlap. Delivery receipts own publication/review status.
