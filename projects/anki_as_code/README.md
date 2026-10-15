---
title: Anki as code
description: Manage collection-backed state and inspect Anki configurations
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code represents notes, decks, and note types as protobuf-backed TOML.
Configuration APIs work without Anki. State commands use a managed Anki backend
and modify only a selected collection metadata key. Collection
export, planning, apply and archive handling remain in
[PR #128](https://github.com/alwaldend/src/pull/128).

## State operations

State commands require `--input` pointing to an existing `.anki2` database. The
CLI initializes the application once after parsing command inputs. `App.Start`
starts the packaged Anki server and opens the database through `BackendCollection`.
`WithState` borrows that prepared connection; it does not start or close resources.
The command boundary always calls `App.Stop`, which closes the collection before
the server, including on operation failure. CLI cleanup has a bounded context
independent of command cancellation; operation and cleanup errors are preserved.

Missing files, directories and other extensions are rejected before backend
startup. State commands do not create collections or handle archives. Database
changes persist in place; archive export and reconciliation remain in #128.

```sh
anki_as_code state dump --input /path/to/collection.anki2 > state.anki.toml
anki_as_code state import --input /path/to/collection.anki2 --text /path/to/configs \
  decks/Default/deck.anki.toml=1 notes/example.note.anki.toml=123
anki_as_code state update --input /path/to/collection.anki2 state.anki.toml
anki_as_code state replace --input /path/to/collection.anki2 state.anki.toml
```

Dump returns stored bindings without writes or discovery. An absent namespace
returns empty version-1 state. TOML is the default; `--format json` selects JSON.
Update and replace read FILE or stdin when omitted/`-`. Status and errors use
stderr; dump uses stdout.

Import accepts multiple PATH=ID arguments, resolves declared config paths before
interpreting native identifiers, checks the requested existing objects and saves
once after the complete batch succeeds. IDs remain strings at the CLI/Go API
boundary; Anki currently requires positive int64 IDs. Note-type imports include
native field/template bindings. Existing scheduling presets use `presets/<alias>`
addresses. Import does not edit config files or native objects.

Update merges external version-1 state: omitted map keys remain, supplied keys
replace their complete entry, including note-type member maps. Replace installs
the whole object and removes omitted bindings. It can repair invalid saved state;
invalid replacement input and native read failures still fail before writing.
Validation rejects conflicting path/native ownership before writes.
Empty updates and identical operations
perform no state write. Native notes, cards, scheduling, history, presets and
media remain owned by Anki.

Missing or unknown state subcommands fail with nonzero status and stderr output.
Explicit `state --help` succeeds.

`--state-key KEY` selects a namespace; the default is `anki_as_code/state/v1`.
Only that key is written. Callers coordinate concurrent writes to the same key.
State uses separate typed maps under `State.Resources`:

```toml
version = "1"
[resources.decks."decks/Default/deck.anki.toml"]
id = "1"
[resources.notes."notes/example.note.anki.toml"]
id = "123"
[resources.note_types."note_types/Basic.notetype.anki.toml"]
id = "456"
[resources.note_types."note_types/Basic.notetype.anki.toml".fields.Front]
id = "789"
[resources.presets.Default]
id = "1"
```

The CLI owns flags and output. `App.Start` and `App.Stop` compose the lifecycle
implemented by `projects/anki_grpc/pkg/backend`; only that backend creates and
closes the gRPC connection. Collection `Input` belongs to `app.Options` and is read
during startup. `app.StateOptions` contains only `Key` and is passed by value to
`WithState`. The CLI binds application options directly; formats and streams
remain command concerns.

Go callers configure `App.Options`, call `Start` once, reuse `WithState` for any
number of state operations, then call `Stop` with a suitable cleanup context.
Call `Stop` even if startup fails. Lifecycle calls must be serialized. State
access before successful startup or after shutdown fails without invoking its
callback. Repeated startup is rejected without replacing an initialized backend;
state operation errors leave cleanup to the caller.

## State Go API

`internal/state` borrows a prepared connection and owns no process, collection
lifecycle or files. The selected namespace defaults when the key is empty:

```go
store := state.New(connection, "")
value, err := store.Dump(ctx)
changed, err := store.Import(ctx, configs, map[string]string{
    "notes/example.note.anki.toml": "123",
})
changed, err = store.Update(ctx, external)
changed, err = store.Replace(ctx, replacement)
```

Dump, Update and Replace exchange `*collectionpb.State`. Import accepts a parsed
`*collectionpb.CollectionConfigSet` and path/ID strings. Multiple calls reuse the
same prepared connection. The changed flag reports whether state was written.
`State.Resources` groups typed deck, note, note-type and preset maps; note types
own field/template maps.

`config/codec` shares file serialization with config operations:
`Encode(message proto.Message, format codec.Format)` and
`Decode(data []byte, format codec.Format, message proto.Message)`. Use `codec.TOML`
or `codec.JSON`; an empty format defaults to TOML. The state store has no codec
or filesystem dependency.

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
belong in the root's `moved` list. PR #128 exports native identity bindings
as a separate state object; explicit state operations install it in the collection.
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

![Load, Dump, and Export configuration processes](assets/config-operations.svg)

These operations work with configuration files and protobuf objects without
opening Anki. Load checks syntax, schema, version, and filesystem layout; callers
own collection semantics. Dump calls Load and serializes the complete parsed set.
Export can write a caller-supplied object directly without loading it first.
Read, encoding, and filesystem failures end the operation with an error. Export
may leave partial output; it does not validate relationships or clean up files.

![Configuration model](assets/config-references.svg)

The class diagram shows root and resource configuration messages. References
use paths; card and template deck references target normal decks. Each note stores
shared field content and tags, while note types define ordered fields, templates,
and presentation. The root marker contains the version and explicit moves.
Attached notes identify each file's location in the fixed hierarchy. The
operation chart shows discovery and conversion; the JSON/TOML branch changes
how Dump encodes the complete configuration set.

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
bazel_agent bazel test //projects/anki_as_code/internal/app:app_test \
  //projects/anki_as_code/internal/state:state_test \
  //projects/anki_as_code/cmd/anki_as_code:cli_workflow_test \
  //projects/anki_as_code/cmd/anki_as_code:config_workflow_test
```

Use absolute paths with `bazel run`, which may execute outside the checkout.

| Package                  | Responsibility                                                                    |
| ------------------------ | --------------------------------------------------------------------------------- |
| `api/collection`         | Generated resource, configuration, plan, and state contracts                      |
| `internal/config`        | Public configuration facade and dump operation                                    |
| `internal/config/layout` | Fixed hierarchy, discovery, file I/O, and direct export                           |
| `internal/config/codec`  | Filesystem-independent protobuf/TOML/JSON conversion                              |
| `internal/state`         | Stored-state dump, batch import, update and replacement                           |
| `internal/model`         | Shared runtime snapshot models, without business logic                            |
| `internal/app`           | Configuration adapter, explicit application lifecycle and collection placeholders |
| `cmd/anki_as_code`       | Cobra flags, cancellation, output, and errors                                     |

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

This project uses a translation layer: application wiring obtains the managed
connection from `projects/anki_grpc`, whose Rust server forwards operations to
Anki's Rust API. The backend library implements server and collection lifecycle;
the application calls its Start, Connection and Stop methods.

![Anki backend integration approaches](/docs/projects/anki_as_code/assets/backend-approaches.svg)

The maintained [Mermaid source](backend-approaches.mmd) is rendered by the
`:backend_approaches` Bazel target and published with the project documentation.

## Verification

The built-CLI E2E runs with an empty executable search path and emits repeatable
JSON, TOML, and exported-directory artifacts. It checks nested ownership, fixed
traversal, path references, removed identities/command, exact Unicode/multiline
fields, optional zero values, empty fields, source preservation, schema/filesystem
errors, and direct object export without validation or message changes.

State E2E uses a managed native worker for Go API checks and lends its closed
database to managed CLI invocations with an empty executable search path. Native
API checks after each CLI call verify cleanup and persistence. It covers batched
operations, rejected updates, namespace isolation, both formats, invalid-state
recovery and native-object preservation, and emits TOML, JSON and a closed native
database. General backend
lifecycle and archive tests belong to their owning components.

The application E2E starts once, exercises nested and sequential state access,
checks rejected repeated startup and access before/after lifecycle boundaries,
and verifies callback/cancellation error handling. It reopens the collection in a
new application to verify persistence, emitting JSON and a closed native database.
