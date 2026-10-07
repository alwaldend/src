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

Anki as code exports collections as protobuf-backed TOML, plans changes, and
reconciles them through Anki's native API. It accepts modern `.colpkg` archives
containing `collection.anki21b` and offline collection databases such as
`collection.anki2`. Keep Anki and other writers closed during export, state operations,
apply, and build. Plan, config dump, and state dump use stdout; status and errors use stderr.

Configuration paths identify objects, similar to Terraform. Native IDs live in
Anki and collection-backed state; note GUIDs remain Anki metadata. New files
create new objects; explicit moves retain existing identities and review progress. Apply/build leave
all configuration files unchanged. New notes and their note types can be created
in the same apply, without generating or writing identities into source.

Collection operations live under `col` (alias `collection`): `col export`,
`col plan`, `col apply`, `col build`, and `col state`. The namespace flag
`--state-key` belongs to this group. Text-only operations use `config` (alias
`conf`), such as `conf dump`. Collection source flags remain specific to each
operation: `--collection` for existing collections and `--base` for builds.

## Commands

```sh
bazel_agent bazel run //projects/anki_as_code -- col export \
  --collection /absolute/path/collection.colpkg --output /absolute/path/text
bazel_agent bazel run //projects/anki_as_code -- config dump \
  --text /absolute/path/text
# Bind existing objects with state import before editing the exported configs.
# Edit the TOML files before planning and applying.
bazel_agent bazel run //projects/anki_as_code -- col plan \
  --collection /absolute/path/collection.colpkg --text /absolute/path/text
bazel_agent bazel run //projects/anki_as_code -- col apply \
  --collection /absolute/path/collection.colpkg --text /absolute/path/text
bazel_agent bazel run //projects/anki_as_code -- col build \
  --base /absolute/path/collection.colpkg --text /absolute/path/text \
  --output /absolute/path/rebuilt.colpkg
```

Export the initial `.colpkg` from Anki. Use `col state import` to bind exported paths
to existing objects before editing configs. Export does not initialize state.
After build, import
`rebuilt.colpkg` into Anki; the collection-backed bindings travel with it.
Apply updates the supplied archive, while build writes a separate one.

Use absolute paths with `bazel run`, which can start outside the checkout.
Current-user `~` and `~/...` paths are supported. `--text` accepts a collection
root file or its directory. `config dump` defaults to JSON; `--format toml` is
also available. The CLI locates the repository-built `anki_grpc` server through
runfiles, with a PATH fallback for installed deployments. No Python is required.

## Configuration-set export

`config export [FILE] --output DIRECTORY` writes a serialized
`CollectionConfigSet` as the fixed TOML hierarchy. Input defaults to JSON;
`--format toml` accepts a TOML-encoded config set. Omit FILE or use `-` to read
stdin. The `conf` alias works for both dump and export. These commands operate
without opening Anki.

```sh
anki_as_code config dump --text ./configs > configs.json
anki_as_code config export configs.json --output ./copy
anki_as_code config dump --text ./configs | anki_as_code config export --output ./copy
```

The caller supplies complete configs with valid, unique collection-relative
paths. Schema decoding happens before writes; collection semantics and path
validation remain the caller's responsibility. Export overwrites declared files
and retains unrelated entries. Filesystem failures can leave partial output;
callers can export to a temporary directory when they need staged publication.
Status goes to stderr and stdout stays empty.

`App.ConfigExport(ctx, output, configs)` exposes the same writer for Go callers
through `config.Store.Export`, without changing the supplied protobuf object.

## State operations

State commands accept an existing offline `.anki2` database or `.colpkg` archive
through `--collection`. Commands call App's typed state workflows. App borrows
one prepared connection from `collection.Store`, which reserves the source and
constructs the native worker and collection backends.
Collection cleanup finishes before transport cleanup and source reservation
release. State stores borrow the connection and own no native lifetime or files.

Databases are opened directly and committed writes persist in place, including
when a later operation fails. Packages are restored into a temporary directory
and changed results are exported through Anki's atomic package API including media.
Dump and unchanged writes retain the source archive bytes. Failures before export
discard restored changes. Normal cleanup removes restored files; uncertain native
completion follows the backend library's retention/error contract.

Missing files and non-regular inputs fail before startup. Anki validates package
contents during restore. Cleanup uses a bounded context independent of command
cancellation, and operation and cleanup errors are preserved. Keep Anki and other
writers closed during these operations. Collection export and reconciliation use the same native workflow.

```sh
anki_as_code col state dump --collection /path/to/collection.anki2 > state.anki.toml
anki_as_code col state import --collection /path/to/collection.anki2 --text /path/to/configs \
  decks/Default/deck.anki.toml=1 notes/example.note.anki.toml=123
anki_as_code col state update --collection /path/to/collection.anki2 state.anki.toml
anki_as_code col state replace --collection /path/to/collection.anki2 state.anki.toml
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
Explicit `col state --help` succeeds.

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

`pkg/ankiascode` provides the concrete `ankiascode.App`. Its zero value is
ready to use. Operations follow the CLI namespaces: `ConfigDump`, `ConfigExport`, `ColExport`,
`ColPlan`, `ColApply`, `ColBuild`, and `ColStateDump`, `ColStateImport`,
`ColStateUpdate`, `ColStateReplace`. Collection inputs are ordinary file paths;
there are no collection reference objects, lookup IDs or selector messages.

`internal/cli` owns Cobra commands, flags, binding syntax, formats and streams.
Its root constructor receives one implementation of its private `application`
interface. The interface uses strings, protobuf messages and standard-library
types. This package does not import App or the native backend. `cmd/anki_as_code`
constructs `ankiascode.App`, passes it to `cli.New` and handles process signals.
Isolated command tests supply one fake without constructing App or starting Anki.
Commands pass namespaces explicitly, and App retains no command configuration.

```go
application := &ankiascode.App{}
// ConfigDump(ctx, path, format) works without a native backend.
value, err := application.ColStateDump(ctx, absoluteCollectionPath, key)
err = application.ColStateImport(ctx, absoluteCollectionPath, absoluteConfigPath, key, bindings)
err = application.ColStateUpdate(ctx, absoluteCollectionPath, key, external)
err = application.ColStateReplace(ctx, absoluteCollectionPath, key, replacement)
```

App's path-based operations return after cleanup. Internal callers needing
several operations in one scope can use collection.Store directly.

`internal/collection.Store` reserves existing sources and publishes changed archives.
Its zero value is ready to use. `WithCollection` injects a prepared
`grpc.ClientConnInterface` into its callback. Construct the services you need on
that connection and reuse local instances for several operations.

Store locates the bundled server or falls back to PATH and constructs a fresh
stdin/stdout worker and collection backend for each scope. It closes native
storage before stopping the worker. Each cleanup has a five-second bound
independent of callback cancellation; source reservations outlive cleanup.

```go
store := &collection.Store{}
err := store.WithCollection(ctx, path, func(connection grpc.ClientConnInterface) (bool, error) {
    resources := resource.New(connection)
    states := state.New(connection, key)
    exporter := export.New(resources, states)
    if err := exporter.Export(ctx, directory); err != nil {
        return false, fmt.Errorf("export configs: %w", err)
    }
    changed, err := states.Replace(ctx, replacement)
    if err != nil {
        return false, fmt.Errorf("replace bindings: %w", err)
    }
    return changed, nil
})
```

Callbacks finish all work sequentially before returning. Retained connections
reject calls after cleanup. Resource/state stores and exporters borrow the
connection; each operation reads fresh data. A Store or App must not be copied
after use. Duplicate canonical paths and file aliases are rejected within one
Store; callers coordinate access across different stores and external writers.

A true result after successful callback completion publishes changes back to the
source archive through Anki. False, errors and cancellation skip publication.
Database writes persist immediately, including before later errors; the scope is
not a transaction. To export another archive destination, use Anki's
`ExportCollectionPackage` with the injected connection as the final operation and
return false. Native export consumes the open collection; backend cleanup handles it.

## State Go API

![State store operations and Anki storage](assets/state-store.svg)

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

## Collection lifecycle

App delegates native access to `collection.Store.WithCollection`. The store
reserves the source and constructs `BackendSTD` and `BackendCollection`, using
the bundled server or PATH. BackendCollection restores packages into temporary
storage or opens databases directly. Callbacks receive the prepared connection
and construct their own services. Store closes the collection before stopping
the worker, including after callback errors, cancellation and panic unwinding.

BackendCollection owns readiness, native import/open/close and restore cleanup.
Its failure contract can retain restore files when native completion is uncertain,
reporting their path. The CLI and App add no restore owner or lifecycle state machine.

Store reserves each source before startup and retains that reservation through
cleanup. A successful callback returning true publishes only its own source
archive through Anki's `ExportCollectionPackage`; false, errors and cancellation
skip publication. Anki writes and atomically replaces that archive. Database
operations persist in place. Build separately copies its base to a temporary
directory beside the destination, applies to that copy and atomically replaces
the destination after success. Failed builds remove staging and leave the
destination untouched. Callers own exclusive access to the base and destination.
State, export and reconciliation borrow prepared services without owning lifetime.

## Export and explicit moves

Export writes configuration TOML without changing the source collection or
creating a state file. Existing stored bindings retain note and note-type
filenames on re-export. New note types receive unused filenames with numeric
suffixes when their names collide with tracked paths, including case-only
differences and fallback names. Without bindings, native edits can change
content-derived filenames.

Export selects a note's directory from its first regular card. Card generation or
changed destinations can therefore change an exported path, including when Anki
chooses a deck for a new root note. When adopting these new paths, declare
`[[moved]]` entries from their stored paths before applying. Export does not
rewrite bindings. Applying the same root-note configs again preserves native
identity and converges; Anki owns their initial deck selection.

Use `col state import` to bind config paths to existing native objects before the
first reconciliation. New configs without bindings describe new objects.
State dump exports stored bindings separately; it does not discover bindings
for an untracked config tree:

```sh
bazel_agent bazel run //projects/anki_as_code -- col state dump \
  --collection /absolute/path/collection.colpkg > /absolute/path/state.anki.toml
```

Export writes declared files directly through `config.Store.Export` and does not reload
generated files. Other entries remain untouched, including old configs and state
files. The destination must not contain the source collection. Filesystem errors
may leave partial output.

Use a fresh destination for a clean snapshot. For directory cleanup or
synchronization, export to a temporary directory and use rsync or another
external tool; those policies belong to the caller.

To rename or move a resource without replacing it, move its file/directory and
add a root `[[moved]]` pair. Directory moves include descendant notes/decks; type
moves include their members. Fields/templates use addresses such as
`note_types/Basic.notetype.anki.toml#fields/Front`. Update readable references in
the same edit. Moves require independent, nonoverlapping endpoints, reject
occupied destinations, and are idempotent when retained. The old source must no
longer be declared. Remove a completed move before reusing its old path.

## Managed and unmanaged entities

- **Notes:** add/remove tracked notes; edit exact named fields and tags shared by
  their generated cards. Only tracked notes missing from configs are removed.
- **Regular-deck assignments:** move generated cards using note defaults and
  template/cloze ordinals. Card entries assign decks and do not declare existence.
  Before deleting a deck containing retained cards, declare a surviving destination
  for them. An omitted assignment keeps the card and its existing deck.
- **Decks:** create, rename/move, remove, and edit normal or filtered configuration.
  Titles derive from the fixed hierarchy. A deck must be empty in Anki before
  conversion between normal and filtered kinds. Native APIs must support the requested
  removal; required native objects cannot silently become untracked.
- **Note types:** create, rename, remove, and edit ordered fields, templates,
  CSS, browser fonts, and upstream presentation settings. Explicit member moves
  preserve field/template identities. Edit affected named note fields together.
  Existing notes cannot switch to another native note type through this CLI.

These are excluded from editable desired state:

- **Resident filtered cards:** membership, home decks, and temporary scheduling
  belong to Anki. Their assignments are omitted. Empty the filtered deck before
  removing referenced resources or changing a resident card's template structure.
- **Scheduling presets:** referenced by readable names; never created, deleted,
  or edited. The CLI calls no preset-writing API.
- **Existing card scheduling/state and review history:** preserved in the database
  to retain study progress. Anki initializes cards it generates for new notes/templates.
- **Media payloads:** preserved through Anki’s native package APIs. Field references are
  editable; media creation/deletion is outside the workflow.
- **Collection/UI/sync settings and native metadata:** owned by Anki and omitted
  from text. API operations update their normal bookkeeping as needed.
- **Untracked objects:** preserved during apply. Import them before edits that
  would destroy their references, such as deleting a deck with untracked cards,
  renaming a parent with untracked children, or changing shared field/card structure.

## Publication and planning

Plan reads an API snapshot, resolves config references and moves, checks owned
relationships, and compares requested resources in memory. It performs no native
mutations, writes no bindings and creates no database copy. It reports paths,
resource changes and known card deck assignments. `card_generation_deferred`
indicates that Anki determines generated-card effects during apply; the planner
does not claim exact creation/removal counts for those effects.

Anki performs final template validation, normalization and card generation during
apply. Planning can therefore accept a template that apply rejects. Source values
that Anki normalizes can remain different from stored values in subsequent plans;
accepted values are never written back into TOML.

Database apply updates the supplied collection directly. Anki commits each native
operation, so later errors can leave earlier operations committed. Archive apply
uses Anki's atomic package export; complete no-ops retain original bytes by skipping
that call. Failures before export leave the source package unchanged. Build requires
a distinct archive output and stages its copy before Apply, publishing the completed
archive atomically after success. Keep other writers closed while using a collection;
the CLI does not detect concurrent archive replacement.

Collection export writes its declared config files in place. It does
not traverse, clean or replace the destination directory. Other entries remain
untouched, and filesystem errors may leave partial output. Callers can export to
a temporary directory and synchronize it externally when required.
Converged database apply leaves source bytes unchanged. Note-type comparison
ignores runtime ordinal hints and unmanaged metadata while retaining explicit
field/template identity changes. New notes use Anki's batch creation API.
Export neither writes source state nor produces a state file.
Failed native reads or projection leave existing output unchanged.

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
belong in the root's `moved` list. Explicit state operations install bindings
inside the collection; export writes configs only.
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
Attached notes identify each file's location in the fixed hierarchy. In Load,
each file is decoded as it is read, and deck titles and omitted card deck paths
are derived inside the deck loop. Dump selects JSON or TOML encoding for the
complete configuration set.

TOML decodes to generic values, then JSON, then protobuf using ProtoJSON.
Unknown keys and incompatible values fail. Unquoted TOML dates/times fail rather
than silently becoming strings; quote dates used as note text. Output converts
ProtoJSON directly to TOML, preserving optional presence and exact text values.
64-bit values use quoted decimal strings, including `version = "1"`.

`config.New(root)` returns a `*config.Store` bound to a configuration path without
accessing the filesystem. Use a directory for writing; reads also accept a root
config file. The store exposes:

- `Load() (*collectionpb.CollectionConfigSet, error)` parses the hierarchy.
- `Dump(format string) ([]byte, error)` loads and encodes JSON or TOML.
- `Export(configs *collectionpb.CollectionConfigSet) error` writes the root
  marker and supplied path/config pairs as TOML files beneath its directory.

```go
store := config.New(absoluteDirectory)
configs, err := store.Load()
data, err := store.Dump("json")
err = store.Export(configs)
```

Stores retain no parsed messages or open files. Each Load reads fresh values;
callers own filesystem coordination and the returned messages. `config.Version`
supplies the supported schema version to collection exporters.

Export writes directly, preserves input messages, overwrites declared files, and
leaves unrelated/stale files alone. Callers own valid, unique paths, layout,
validation, staging, cleanup, and publication. Export does not load or normalize
its output; failures may leave partial files. Consumers use the store;
the layout package remains an implementation API with restricted Bazel visibility.

`config/layout` owns the fixed hierarchy, discovery, file I/O, and direct file
export. `internal/codec` converts TOML/JSON through protobuf without filesystem
access. Config layout, config dump, state CLI commands, and collection export
reuse its format dispatch; TOML/JSON helpers are private. Dependencies flow from
config to layout/codec and from layout to codec. State stores have no serialization
dependency.
Collection relationship validation and native operations remain outside them.

## Architecture

The project keeps one Go module with focused packages. Configs and state are
protobuf contracts. Runtime resolution models stay private to `internal/plan`; packages expose only operations used by their consumers.

![Anki CLI package dependencies](assets/architecture.svg)

The dotted edge shows calls through the CLI-owned application interface. Solid
edges show component dependencies.

| Package                  | Responsibility                                                 |
| ------------------------ | -------------------------------------------------------------- |
| `cmd/anki_as_code`       | App construction, CLI invocation and process signals           |
| `internal/cli`           | Cobra commands, private application interface and streams      |
| `internal/collection`    | Source reservations and archive publication                    |
| `api/collection`         | Generated config, plan, and persisted-state contracts          |
| `pkg/ankiascode`         | Command routing and shared backend composition                 |
| `internal/export`        | Native-to-config projection and direct file output             |
| `internal/plan`          | Pure resolution, comparison, ownership and prepared operations |
| `internal/reconcile`     | Execute protobuf plans, bind native IDs and persist state      |
| `internal/rebuild`       | Apply to a staged archive copy and publish it atomically       |
| `internal/resource`      | Native generated clients and fresh resource snapshots          |
| `internal/config`        | Root-bound configuration Store for Load, Dump, and Export      |
| `internal/config/layout` | Fixed hierarchy, discovery, file I/O, and direct export        |
| `internal/codec`         | Filesystem-independent protobuf/TOML/JSON conversion           |
| `internal/state`         | Stored-state dump, import, update, and replacement             |

The command entry point injects App into internal/cli. Commands invoke its private
application interface without importing the implementation. App borrows a prepared
connection through collection.Store and constructs its services locally. `export.New(resources, states)` constructs an
`Exporter` using caller-provided resource and state stores. Its
`Export(ctx, output)` method owns native-to-config projection,
filenames and direct config output. State operations remain separate.
Export has no dependency on reconciliation.

`plan.Build(configs, snapshot, stored)` constructs prepared protobuf operations
without I/O. `reconcile.New(resources, states)` constructs a `Reconciler` whose
`Apply(ctx, compiled)` method executes them and returns whether native data or
bindings changed. App shares config loading, capture and plan preparation between
plan and apply within one collection scope. Planning owns reference resolution,
comparison and ownership; execution owns native calls and returned identities.

Config and codec have no Anki runtime dependency. State owns stored path bindings,
not lifecycle or serialization. App uses the Go backends in
`projects/anki_grpc` and Anki's collection/package services; Resource owns generated resource clients and fresh native snapshots.
The executor borrows resource/state stores and owns no process or filesystem resources. There is no archive facade or handwritten ZIP/Zstandard
handling. Anki owns storage compatibility, normalization, structural migrations,
field caches, and card generation. No HTML cache rebuilding or SQLite driver is
implemented in Go. The application does not force Check Database after import;
normal Anki API updates maintain their caches.

## Planning Go API

![Planning orchestration and pure plan construction](assets/planning.svg)

`plan.Build(configs, snapshot, stored)` returns an executable protobuf
`CollectionPlan` from caller-supplied configs, a resource snapshot and stored
bindings. It performs no I/O and leaves all inputs unchanged. Callers capture
fresh inputs under exclusive offline access to the same collection and namespace.
Backend lifetime belongs to the caller; planning can run with the backend open or closed.

The plan separates `expected_state` from `proposed_bindings` and contains a change
report, prepared deck/type/note operations, name releases, tracked deletions and
card destinations. New objects have zero native IDs; prepared note-type and
template-deck references use configuration paths. Resource comparison produces
operations and report entries together. Converged plans have no native operations.
Card destinations are included only for known changes or relevant native generation.
Planning indexes native protobufs directly and resolves desired configs into the
prepared protobufs used for comparison and execution. Returned payloads are
independent of supplied inputs; path/ID indexes and ownership checks stay private
to `internal/plan`.
Member moves rename fields or templates within one note type; cross-type identity
transfers are rejected. Stable native member IDs are authoritative. Without IDs,
matching uses original stored names, including before explicit renames; stale
ordinals cannot claim another current member. Ambiguous bindings fail before a
plan is returned. Untracked resources and resident filtered cards stay unmanaged,
including deck targets referenced by untracked templates. All cloze card ordinals
share template zero; resident cards block template replacement or removal, while
explicit template moves that preserve identity remain allowed.

Anki owns native validation, normalization, structural migrations and card generation.
The report compares current card destinations and cards removed with deleted notes;
template-driven additions, removals and migrations remain deferred through
`card_generation_deferred`. Prepared operations do not simulate those native effects.
The CLI and App's `ColPlan` return the report; Go services pass the complete plan to
`reconcile.New(resources, states).Apply(ctx, compiled)`.

The executor checks expected stored state before writes, executes prepared operations,
resolves new native IDs and persists bindings. It never loads configs, captures a
resource snapshot, matches members or recomputes the plan. Native metadata and
newly generated cards are read only as needed to execute it. A private `applier`
reads the supplied plan unchanged, clones the proposed bindings, and copies only
note-type requests whose native identities or presentation need updating. Reports,
expected state and note bodies are not cloned. Callers retain exclusive access and
own freshness beyond stored-state checks.

```go
compiled, err := plan.Build(configs, snapshot, stored)
changed, err := reconcile.New(resources, states).Apply(ctx, compiled)
```

## Archive build Go API

`internal/rebuild` stages a separate collection archive using an injected apply operation:

```go
builder := rebuild.New(apply)
err := builder.Build(ctx, baseArchive, outputArchive)
```

The non-nil callback receives a private `.colpkg` path and must finish all native
work and close its session before returning. Builder knows nothing about configs,
Anki clients or reconciliation. Callers own home-path expansion and exclusive
source/destination access; the destination parent must already exist.

Builder checks the regular source file and rejects source/destination aliases,
copies into a private directory beside the destination, then publishes the
completed archive with an atomic rename. Apply errors, cancellation observed
before publication and publication failures leave the original and existing
destination intact. Staging is removed on return, including panic unwinding;
cleanup errors are retained. Cancellation preserves both the context error and
its custom cause for `errors.Is`. Atomic publication does not guarantee power-loss
durability. The App callback delegates to `ColApply`.

Run `bazel_agent bazel test //projects/anki_as_code/internal/rebuild:rebuild_test`.
The native E2E changes real Anki state in the staged archive and emits original,
changed and unchanged `.colpkg` files plus a decoded JSON state object.

## Native resource Go API

`internal/resource` shares generated native clients and complete snapshots on a
caller-owned prepared connection:

```go
resources := resource.New(connection)
snapshot, err := resources.Snapshot(ctx)
note, err := resources.Notes.GetNote(ctx, &notes.NoteId{Nid: noteID})
```

Store exposes the generated deck, note-type, note, card, search, preset and rendering clients
without forwarding wrappers. Snapshot returns upstream protobuf objects,
including unmanaged resources and resident filtered cards, plus preset names.
Decks are ordered by title; note types, notes and cards by native ID. Each call
reads fresh values and returns caller-owned results. There is no cache or atomic
snapshot guarantee; callers provide exclusive offline access.

The resource store owns no backend lifecycle, state bindings, configs, projection,
planning or mutation order. Export and reconciliation retain those policies.
No SQL or duplicate storage logic is involved.

## Collection export Go API

![Export data flow](/docs/projects/anki_as_code/assets/export.svg)

`internal/export` borrows resource and state stores for the same prepared collection:

```go
resources := resource.New(connection)
states := state.New(connection, key)
exporter := export.New(resources, states)
err := exporter.Export(ctx, absoluteOutputDirectory)
```

Stored bindings retain filenames and preset aliases; export does not generate,
serialize, write or install state. State operations remain separate.
The destination must not contain the source collection.

New note filenames use Anki's `HtmlToTextLine` result for the configured sort
field. Go only sanitizes that text into a bounded Unicode filename. Bound names
are reused from state; filename conversion never changes stored note fields.

Export owns config projection and readable filenames and delegates file output
to `config.Store.Export`. It overwrites declared configs in place and leaves other
entries untouched, including stale configs and caller-owned state files.
Filesystem errors may leave partial output. For cleanup or publication policies,
export to a temporary directory and synchronize it externally, for example with
rsync. The caller also owns exclusive offline access, backend startup,
collection preparation and shutdown.

Export and reconciliation share `resource.Store` for generated clients and fresh
native snapshots; their projection and policy remain separate.

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

```sh
bazel_agent bazel build //projects/anki_as_code
bazel_agent bazel test //projects/anki_as_code/pkg/ankiascode:ankiascode_test \
  //projects/anki_as_code/internal/export:export_test \
  //projects/anki_as_code/internal/config:store_test \
  //projects/anki_as_code/internal/cli:command_test \
  //projects/anki_as_code/cmd/anki_as_code:cli_workflow_test
```

Synthetic native E2E verifies reviewed-card/history/preset/media retention, filtered
cards, field/template migrations, copied notes, config-only export and independent state operations,
import/moves, untracked-object safety, read-only configs, convergence, rejected
publication, large fields/archives, and Anki's actual full-package round trip.
CLI E2E verifies streams, flags, removed generation, and pure config behavior.
Repeatable TOML, JSON, and collection artifacts are emitted through Bazel's
undeclared test outputs. Exact-candidate delivery receipts own validation status.

The independent exporter E2E uses a caller-managed native backend and collection.
It checks exact Unicode/multiline fields, presentation, member bindings, stable
bound/discovered paths, native filename text conversion, filtered-card exclusion,
read-only state namespaces, untouched existing entries, cancellation and partial output
after filesystem errors. It emits a complete text tree and a native `.colpkg` artifact.
Run `bazel_agent bazel test //projects/anki_as_code/internal/export:export_test`.
