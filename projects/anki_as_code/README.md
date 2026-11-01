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

The [CLI reference](/docs/projects/anki_as_code/docs/cli/) has a short index and
one child page per root subcommand, each containing its complete nested help,
including the generated Anki API tree. Bazel generates these pages from the
Cobra definitions with `tools/generate_cli_docs` and packages the hierarchy
with the documentation; generated Markdown is not checked in. Build it with
`al bazel build //projects/anki_as_code/docs:cli_reference`.

## Optional server engines

### WebAssembly CLI

Build the Go CLI itself for WASI Preview 1 with both bundled servers disabled:

```sh
al bazel build -c opt --//projects/anki_as_code:wasm=true //projects/anki_as_code
```

This produces a `wasip1/wasm` module. The public alias uses `select()` to choose
a `go_cross_binary` target that cross-compiles the existing server-free command.
Other targets in the same invocation keep their normal platform and server-engine settings. This
selection excludes both bundled servers regardless of the `anki_grpc` engine flags.
Find the module path with
`al bazel cquery -c opt --//projects/anki_as_code:wasm=true //projects/anki_as_code --output=files`.
Run it with a WASI Preview 1 runtime and grant access to the configuration
directory. For example, with the repository-pinned wazero CLI:

```sh
al bazel run @com_github_tetratelabs_wazero//cmd/wazero -- run \
  -mount "$PWD/configs:/configs:ro" /absolute/path/to/wasi \
  config dump --text /configs
```

Help and file-only `config export`/`config dump` commands work without Anki.
No native server, embedded Rust WASM server, or self-spawned default worker is
included. Standard Go WASI cannot spawn subprocesses or dial TCP/Unix sockets,
so collection and API operations require additional host integration. This is
a WASI command, without browser JavaScript bindings.

### Bundled server selection

The independent booleans `//projects/anki_grpc:cgo` (default true) and
`//projects/anki_grpc:wasm` (default false) compose available engines in the
public command. Build a pure-Go host executable with the bundled WASM server:

```sh
al bazel build -c opt //projects/anki_as_code --//projects/anki_grpc:cgo=false --//projects/anki_grpc:wasm=true
```

The default worker uses cgo when available and otherwise runs the complete Rust
CLI/server through WASM in the same executable. Both generic subprocess and
in-process strategies prefer cgo. With neither engine, the command supports explicit external
clients and reports that local engines are unavailable. The executable owns
its scoped client factory, worker arguments and supported subprocess schemes.
`subprocess+wasm://` and `subprocess+cgo://` select their named compiled engine
for an explicitly encoded application worker, including when both are enabled.
Their factories supply `--server-engine=wasm` or `--server-engine=cgo` before
the worker arguments when the encoded executable resolves to this application,
including filesystem aliases. An independently supplied executable retains its
exact arguments; its caller owns engine selection. For example, a standalone
`anki_grpc` worker still receives just `serve`. The generic subprocess scheme
also preserves its argument vector.
For direct server hosting, use `anki_as_code --server-engine=cgo server run serve`
or select `wasm` the same way; omit the selector for the preferred compiled engine.

Use `--address inprocess+wasm://collection` to execute the same Rust server in
the application process over in-memory gRPC streams. Generic `inprocess://`
prefers cgo when available and otherwise uses WASM; qualified schemes select
their engine explicitly. Each client owns fresh state and closes its runtime
after the operation. The default remains the supervised subprocess.

The WASM adapter owns building and embedding its guest and exposes `Run(arguments)`;
application code handles no module bytes. Only the guest targets WASI. Guest
and host follow normal Bazel compilation modes; `-c opt` optimizes both. Offline
collection, package and local media operations are supported. Network-dependent
Anki services return explicit unsupported errors. Keep exclusive access to
collection and media paths: the WASI SQLite VFS does not coordinate process locks.
`ANKI_WASM_CACHE_DIR` optionally selects a caller-protected compiled-code cache.

### Native server

The native Linux x86_64 application build links the Rust server into its executable. By default,
collection operations launch the same executable with `server serve` and communicate over stdio.
The installed CLI needs no companion worker, PATH lookup, or executable extraction.

The CLI also accepts `--address inprocess://collection` or
`--address inprocess+cgo://collection` for collection workflows and API calls.
Each operation gets fresh native state without launching a worker. A managed API
path can accompany any client address. A network address selects the server where
the collection is opened and closed; address-only access leaves collection lifetime
with the caller. Native work can outlive cancellation, so source reservations remain
held until native cleanup finishes. Configured subprocess execution remains the default.

Reusable application libraries stay pure Go. Callers that want direct execution
import the optional `anki_grpc/pkg/ankigrpc/native` subpackage for registration.
App and Store use the shared `ankigrpc.NewClient` constructor; they keep no native
engine map or concrete adapter constructor. Subprocess connection strings encode their exact
executable and arguments with `ankigrpc.SubprocessAddress`; local in-process labels
remain opaque and never identify shared state.

`al bazel build -c opt //projects/anki_as_code` produces an optimized
executable with final symbol and debug information stripping. It uses the host
Linux runtime libraries, including OpenSSL 3.

`server` shows the available server operations. `server run` forwards arguments
to the linked Rust CLI, which owns native subcommands and their help.
For `server serve`, CLI prepends `serve` before calling App.ServerRun.
Cobra forwards arguments without interpreting native flags:

```sh
anki_as_code server run --help
anki_as_code server serve --help
anki_as_code server serve --transport tcp --listen 127.0.0.1:50051
anki_as_code server serve --transport unix --listen ./anki.sock
```

Stdio is the default and reserves stdout for gRPC. TCP and Unix listeners use
plaintext gRPC and require an explicit address. The child process keeps native
failures and shutdown separate from the application workflow.

Collection commands accept optional `--address ADDRESS` to use an already running
Anki gRPC server. Omit it or use an empty value to keep the embedded default.
The flag belongs to `col` (alias `collection`) and is inherited by export, plan,
apply, build and every state command:

```sh
anki_as_code col plan --address 127.0.0.1:50051 --collection /shared/collection.colpkg --text ./text
anki_as_code col apply --address tcp://127.0.0.1:50051 --collection /shared/collection.colpkg --text ./text
anki_as_code col state dump --address unix:///shared/anki.sock --collection /shared/collection.anki2
```

External addresses accept `host:port`, `tcp://host:port`, `unix://PATH` or an
absolute Unix socket path. Local addresses accept `subprocess://` and, when the
native engine is compiled, `subprocess+cgo://`; WASM-enabled commands also accept
`subprocess+wasm://`. Qualified application-worker schemes select that engine
for same-executable workers before forwarding the encoded server arguments;
independent executables retain their arguments. External connections use plaintext gRPC. An explicit address never
falls back to the embedded worker; interrupt the command to cancel connection
startup if the server is unavailable.

The CLI opens and closes the collection for each command and closes its client
connection afterward. It leaves the external server process running. Reserve
the server exclusively for the command, with no collection already open.
The CLI and server must access the same absolute filesystem paths, including
collection/media files, CLI `TMPDIR` for package restores, and build output
directories for staged archives. Files are not transferred over gRPC.

Build-selected command source files import `ankigrpc/native` for native hosting
and client scheme registration; `main.go` owns signals and worker configuration.
The reusable application library imports only pure Go packages, including when
`//projects/anki_grpc:cgo` is enabled. Set that independent boolean to false to
exclude native linkage; local hosting remains available when WASM is enabled. The executable name stays `anki_as_code`.

## MCP server

`anki_as_code mcp serve` exposes workflow tools and an individual tool for every
declared native Anki RPC over stdio. The pinned catalog contains 281 tools:
269 declared RPCs, ten workflows and two collection-session controls. Forwarded
aliases are omitted; methods declared on Backend services remain available.
An MCP client launches the executable; stdout contains protocol messages and
stderr contains command diagnostics. Discovery needs no running Anki instance,
collection, or optional server engine.

![MCP discovery, collection handles and cleanup](assets/mcp.svg)

A client configuration can launch an installed executable like this:

```json
{
  "mcpServers": {
    "anki_as_code": {
      "command": "/absolute/path/to/anki_as_code",
      "args": ["mcp", "serve"]
    }
  }
}
```

Workflow tools are `config_dump`, `config_export`, `col_export`, `col_plan`,
`col_apply`, `col_build`, `col_state_dump`, `col_state_import`,
`col_state_update`, and `col_state_replace`. Inputs use server-side paths:
`collection`, `text`, `output`, and `base` have the same meanings as their CLI
flags. `state_key` selects the binding namespace. Config export takes a structured
`configs` object; state update/replace take a structured `state` object; state
import takes a `bindings` object mapping config paths to native ID strings.

Use `col_open` with `collection` and an optional `address` to open a
collection once. It returns an opaque `handle`, which native tools and the
collection/state workflows accept instead of `collection` and `address`:

```json
{ "handle": "collection-1", "request": { "did": "1" } }
```

Pass the handle to `col_close` when finished. Its worker and connection
remain open between calls, and enabled archives are published on close.
Disconnect closes all remaining handles. `write_back` defaults to true on
`col_open`; set it to false there to disable archive publication.
Database writes remain immediate. Handle calls cannot include `collection`,
`address`, or `write_back`; `col_build` continues to require an independent base
archive path. Configuration-only tools do not use a handle.

Native tool names are `api_<service>_<method>`, replacing CLI hyphens with
underscores and retaining `backend_` for methods declared on those services.
Canonical names remain unchanged when forwarded aliases are omitted. For example,
`api_notes_get_note` accepts:

```json
{
  "collection": "/absolute/path/to/collection.anki2",
  "request": { "nid": "1234567890123" }
}
```

Native calls require a `handle`, `collection` or `address`; the latter two can select a managed source
and its backend. `address` alone uses that backend's current collection.
Requests and responses use protobuf field names and ProtoJSON values, including
string-encoded 64-bit integers and base64 bytes. The compiled descriptors supply
each tool's nested request and response schema. Responses include structured
JSON and a text JSON fallback; operations without a return value produce `{}`.
The SDK's stdio reader limits individual request frames to 16 MiB; clients can
apply separate limits to responses.

Keep Anki and other collection writers closed. Path-based access retains
per-call collection lifetime and publication. Explicit handles reserve their
source until close, serialize operations on that collection, and reuse the
prepared connection through the existing App methods. The shared App rejects
other managed access to the same source until native cleanup finishes.
Existing backend capability and path-based explicit open/close restrictions
apply. Native `write_back` defaults to true and controls archive publication;
false does not undo database writes. Workflow archive reads can repackage their
source during cleanup. Errors and cancellation do not imply rollback, and the
server does not automatically retry mutations. Cancellation and disconnect allow
existing bounded native cleanup before shutdown.

## Upstream API commands

`api <service> <method>` constructs its command tree from the compiled Go
protobuf descriptors for all 25 pinned Anki schemas. Bazel automatically generates
the registration imports into `internal/rpc` from the pinned descriptor set using
the project-local `tools/generate_anki_schemas` generator;
no separate generation step or runtime schema files are needed. Service/method names use
kebab-case, omit the `Service` suffix and retain `backend-` prefixes. Implicit
backend methods are included. Help lists request fields and message types;
no proto files or server discovery are needed.

```sh
anki_as_code api tags all-tags --collection ./collection.colpkg
anki_as_code api notes get-note --collection ./collection.anki2 \
  --input request.toml --input-format toml --output-format json
anki_as_code api notes get-note --address unix:///path/to/anki.sock \
  --input request.json
```

Exactly one nonempty `--collection` or `--address` is required. An omitted or
empty address retains embedded access with the supplied collection. Offline access starts
the same executable with `server serve` and owns the collection lifetime. `--address` accepts a TCP
host:port, `tcp://host:port`, `unix:///absolute/path` or an absolute Unix socket
path. It uses the server's current collection and closes only the client
connection. Explicit open/close/import RPCs require connected access.
Unix URI paths can be relative, such as `unix://anki.sock`.
App resolves methods, validates payloads and owns remote client connections directly; `collection.Store` manages only
offline collection scopes.

Input defaults to stdin (`--input -`); JSON is the default input and output
encoding. An empty request body is accepted for empty request messages.
`--timeout` bounds input reading, startup and the request, defaulting to 30 seconds. Responses
use stdout after cleanup succeeds; request, publication and cleanup errors use
stderr and a nonzero exit status. Native error details are preserved.

Offline archive write-back is enabled by default. `--write-back` enables it
and `--no-write-back` disables it; the last occurrence wins. Explicit boolean
values such as `--write-back=false` are also accepted. Disabling write-back
skips automatic `.colpkg` publication; native RPCs keep their normal effects.
`.anki2` mutations always commit directly. Requests are not
transactions: earlier committed changes can persist even when a later RPC fails.

## Onboarding an existing collection

Export a full `.colpkg` from Anki, then close Anki and other collection writers.
Use a fresh text directory so stale declarations from earlier exports are not
included. Export configs and their native bindings together:

```sh
al bazel run //projects/anki_as_code -- col export \
  --collection /absolute/path/collection.colpkg --output /absolute/path/text \
  --state-output /absolute/path/state.anki.toml
al bazel run //projects/anki_as_code -- col state update \
  --collection /absolute/path/collection.colpkg /absolute/path/state.anki.toml
al bazel run //projects/anki_as_code -- config dump \
  --text /absolute/path/text
# Edit the TOML files before planning and applying.
al bazel run //projects/anki_as_code -- col plan \
  --collection /absolute/path/collection.colpkg --text /absolute/path/text
al bazel run //projects/anki_as_code -- col apply \
  --collection /absolute/path/collection.colpkg --text /absolute/path/text
# Alternatively, build a separate archive instead of applying in place.
al bazel run //projects/anki_as_code -- col build \
  --base /absolute/path/collection.colpkg --text /absolute/path/text \
  --output /absolute/path/rebuilt.colpkg
```

Export reads collection resources and stored state through Anki's API; managed
source archives may be repackaged. Callers choose the output locations and own
their existing contents. Declared config files are overwritten in place, while
unrelated entries remain untouched. The optional state file contains version-1
TOML bindings for every exported deck, note and note type, including
field/template identities and referenced
scheduling preset aliases. It uses the exact generated filenames, including
collision suffixes and names retained from the selected state namespace.
`col state update` installs those bindings into the same collection before you
edit configs. This preserves existing object identities and review progress
when you later apply changes. Use the same `--state-key` for export, state
update, plan, apply and build if selecting a custom namespace.

After editing, inspect `col plan`, then use `col apply` to update the supplied
collection or `col build` to write a separate archive. After apply, import the
updated archive into Anki; after build, import `rebuilt.colpkg`. Collection-backed
bindings travel with the archive.

Omit `--state-output` for config-only export. Use `col state import` for selective
manual bindings. Generated state belongs to the collection it was exported
from; do not install it into an unrelated collection. State update merges
bindings and retains omitted entries; `col state replace` intentionally removes
omitted bindings. For an already managed collection, review exported state
before installing it, especially if paths have changed.

`--state-output` names a TOML file, regardless of its extension. Current-user
`~` paths are supported. The file may live outside the text directory or inside
it under an unrelated name such as `state.anki.toml`. Export rejects state paths
that conflict with config declaration filenames or directories needed by exported
configs. State output is atomically replaced after config export succeeds.
Filesystem failures may leave partial configs; a failed config export retains
the previous state output. The config tree and state file are not published as
one filesystem transaction.

Use absolute paths with `bazel run`, which can start outside the checkout.
Current-user `~` and `~/...` paths are supported. `--text` accepts a collection
root file or its directory. `config dump` defaults to JSON; `--format toml` is
also available. The CLI starts its linked server through its own executable path.
No Python is required.

Collection and output paths resolve symbolic links before subsequent `..`
components. Configuration loading rejects symbolic links, including before `..`,
instead of silently selecting a lexically normalized configuration directory.

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
and written back during cleanup through Anki's atomic package API, including media.
Publication runs independently of request outcomes and may repackage unchanged
content. Normal cleanup removes restored files; publication failures or uncertain
native completion follow the backend library's retention/error contract.

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

`api/ankiascodeapi` provides `ankiascodeapi.App`, the shared interface for
configuration, collection, state and server workflows. `pkg/ankiascode` provides
the concrete `ankiascode.App`. Its zero value supports configuration operations
and borrowed connections. Path access requires a configured factory or an explicit address;
hosting the server requires an injected `ServerRunner`.
Supply `ankiascode.Options{FS: storage}` to use an `afero.Fs` for every
application-owned filesystem operation, including config discovery/read/write,
collection source inspection, export/state output and build staging/publication.
For example:

```go
storage := afero.NewMemMapFs()
application := ankiascode.New(ankiascode.Options{FS: storage})
err := application.ConfigExport(ctx, "/configs", configs)
data, err := application.ConfigDump(ctx, "/configs", "json")
```

Configuration operations run entirely in that storage and need no Anki server.
Custom filesystem paths retain their own namespace, including relative paths
and literal `~`; they do not use the host home or working directory. Nil selects
`afero.NewOsFs()` and preserves existing OS path expansion, including for a
zero-value App. Callers retain ownership of the filesystem and coordinate access.
`BasePathFs` aliases resolve within its virtual namespace, including directly
nested bases. Equivalent absolute and relative virtual names use one canonical
relative spelling. Parent traversal and symlink targets outside that base are
rejected during reads and writes, including configuration export.
`internal/filesystem.FS` keeps the Afero backend private and owns shared path
handling, file/directory I/O and temporary names. App constructs this wrapper once
and supplies it to every internal filesystem owner. `afero.Fs` in the public
options is an interface that callers can implement with their own storage.
Server and backend-library filesystem implementations remain separate: injecting
storage does not transfer files to a backend or make its native paths accessible.
Managed collection operations still require the backend to access the selected
collection and staging paths. CLI input files and standard streams belong to the
CLI adapter, outside App.

`ankiascode.New(ankiascode.Options{NewClient: factory.NewClient})` injects an
`ankigrpc.ClientFactory` constructor function. It receives each operation's address
unchanged. A concrete Factory's method value preserves its configured default
and private registrations. Use `ankigrpc.ToClientFactory` for a constructor
returning a concrete client, or pass a function returning StartableClient directly.
Without a configured factory, explicit addresses use `ankigrpc.NewClient`.
The executable registers its subprocess schemes and imports direct native adapters; App delegates raw
client startup and bounded cleanup to `ankigrpc.WithClient`; the collection store
owns managed collection scopes.
Collection workflows accept one `ankiascodeapi.CollectionOptions` value:
`Path` selects a managed database/archive and optional `Address` selects the client
used to prepare it. Empty per-operation addresses let the injected factory select its default.
Alternatively, `Connection` borrows an already-open collection; it cannot be
combined with a path or address. `StateKey` selects the namespace, with an empty
value selecting the default. App validates these combinations before opening
resources and retains no endpoint between calls. Path selects managed collection
lifetime, Address selects the client, and Connection borrows an already-open collection.
Managed API calls use the same source semantics as other collection operations.

The CLI registers its worker schemes and supplies its own executable and
`[]string{"server", "serve"}` through the injected factory for its default. Operations follow the CLI namespaces: `ConfigDump`, `ConfigExport`,
`ColExport`, `ColPlan`, `ColApply`, `ColBuild`, `ColStateDump`, `ColStateImport`,
`ColStateUpdate`, `ColStateReplace`, `APICall` and `ServerRun`.
Creating an options value performs no I/O. `ColBuild` requires a base archive in
`options.Path` and forwards a copy of those options with its staged path to apply.

`internal/cli` owns the complete Cobra tree, flags, binding syntax, formats and
streams. Configuration, collection, API and server commands call
`ankiascodeapi.App`; MCP commands use a private interface owned by CLI.
The App interface uses shared options values, strings, protobuf
messages and standard-library types. Commands wrap collection path flags
without opening resources. App owns server delegation and composes
collection/configuration workflows; the lower-level packages retain their domain
and lifecycle responsibilities. CLI passes parsed operation options into the
injected App; the executable owns concrete construction and native adapters.

`cmd/anki_as_code` constructs App and the CLI, owns its scoped client factory,
and injects the configured server runner through `Options.ServerRunner`. The
default worker prefers cgo and otherwise uses WASM; qualified schemes select
their compiled engine explicitly.
The executable also owns process signal handling. App.ServerRun checks startup
cancellation, builds the argument vector and invokes that runner. Serving owns
the worker process until return or termination; its context controls entry only. App does not link Rust/cgo or change
process-global signals. Missing runner configuration fails before native hosting.
Isolated command tests supply one fake without constructing App or starting Anki.
Commands pass namespaces explicitly, and App retains no command configuration.

```go
factory := ankigrpc.NewFactory(ankigrpc.SubprocessAddress(serverPath, "serve"))
factory.RegisterClientFactory("subprocess", ankigrpc.ToClientFactory(ankigrpc.NewSubprocessClient))
var application ankiascodeapi.App = ankiascode.New(ankiascode.Options{NewClient: factory.NewClient})
options := ankiascodeapi.CollectionOptions{Path: absoluteCollectionPath, StateKey: key}
// ConfigDump(ctx, path, format) works without a native backend.
value, err := application.ColStateDump(ctx, options)
err = application.ColStateImport(ctx, options, absoluteConfigPath, bindings)
err = application.ColStateUpdate(ctx, options, external)
err = application.ColStateReplace(ctx, options, replacement)
```

App's path-based operations return after cleanup. Borrowed connections remain
owned by their caller, including after failure or cancellation. Operations read
fresh collection data and neither close the borrowed connection nor publish a
package. Callers coordinate access and decide when to save and close; direct
database writes persist immediately and multiple calls are not a transaction.
One App can alternate between several open collections:

```go
first := ankiascodeapi.CollectionOptions{Connection: firstConnection, StateKey: key}
second := ankiascodeapi.CollectionOptions{Connection: secondConnection, StateKey: key}
err := application.ColStateUpdate(ctx, first, external)
report, err := application.ColPlan(ctx, first, absoluteConfigPath)
other, err := application.ColStateDump(ctx, second)
```

The protobuf contract lives in `api/ankiascodepb`, with Go package
`ankiascodepb` and protobuf namespace `anki_as_code`. Message names and fields
retain the existing config and state formats; Go imports and protobuf fully
qualified type names use the renamed API.

`ColExport(ctx, options, output, stateOutput)` accepts an empty
`stateOutput` for config-only export or a TOML file path for snapshot bindings.
Callers own both output locations and their existing contents for managed paths
and borrowed connections. Config/state output conflict checks still apply.

`internal/collection.Store` reserves existing sources and owns enabled archive publication.
Configure its default client or pass a per-operation address before collection access. `WithCollection` injects a prepared
`grpc.ClientConnInterface` into its callback. Construct the services you need on
that connection and reuse local instances for several operations.

`collection.New(collection.Options{NewClient: factory.NewClient})` configures client
construction and default selection. No executable or server arguments are
discovered by the application or store. Without a configured factory, path-only
access fails before client startup; explicit addresses remain usable. Each scope
constructs a fresh transport client and a single-use application collection scope. `ankigrpc.WithClient` uses the same
factory and bounded transport cleanup for raw addressed API access without
entering the store or reserving a collection. Managed scopes close native
storage before stopping the worker. Collection publication and close have a ten-minute default budget, configurable
with Options.CollectionCloseTimeout; transport cleanup has a separate five-second budget
independent of callback cancellation; source reservations outlive cleanup.

```go
store := collection.New(collection.Options{NewClient: factory.NewClient})
err := store.WithCollection(ctx, ankiascodeapi.ColOpenOptions{Path: path, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
    resources := resource.New(connection)
    states := state.New(connection, key)
    exporter := export.New(resources, states, filesystem.FS{})
    if err := exporter.Export(ctx, directory, ""); err != nil {
        return fmt.Errorf("export configs: %w", err)
    }
    if _, err := states.Replace(ctx, replacement); err != nil {
        return fmt.Errorf("replace bindings: %w", err)
    }
    return nil
})
```

Callbacks return only an error and finish all work sequentially before returning.
WriteBack alone controls archive publication; domain change results do not affect
collection cleanup. Retained connections
reject calls after cleanup. Resource/state stores and exporters borrow the
connection; each operation reads fresh data. A Store or App must not be copied
after use. Duplicate canonical paths and file aliases are rejected within one
Store; callers coordinate access across different stores and external writers.

Store.Open and WithCollection use `ankiascodeapi.ColOpenOptions` for the source
path, address and WriteBack policy. Enabled write-back
publishes the source archive during cleanup independently of callback results, errors
and cancellation. Application workflows explicitly enable it; API flags select it.
Database writes persist immediately, including before later errors; the scope is
not a transaction. To export another archive destination, use Anki's
`ExportCollectionPackage` with the injected connection as the final operation and
disable WriteBack when source publication is unwanted. Native export consumes the
open collection; application scope cleanup handles it.

`App.APICall(ctx, ankiascodeapi.APIRequest{...})` exposes the same generic operation
without per-method wrappers. Supply a fully qualified gRPC Method, Input bytes and
JSON/TOML encodings. `Collection` is the shared `CollectionOptions` value: exactly
one of `Collection.Path`, `Collection.Address`, or `Collection.Connection` is
required for raw API calls. An address selects an app-owned connection to the
server's current collection. WriteBack defaults to false for Go callers; the CLI explicitly
selects true. Borrowed native connections retain caller-owned cleanup and
publication. Request validation precedes native startup.

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

Dump, Update and Replace exchange `*ankiascodepb.State`. Import accepts a parsed
`*ankiascodepb.CollectionConfigSet` and path/ID strings. Multiple calls reuse the
same prepared connection. The changed flag reports whether state was written.
`State.Resources` groups typed deck, note, note-type and preset maps; note types
own field/template maps.

`state.Export(configs, snapshot, identities)` generates version-1 bindings from
projected configs, their path-to-native-ID map and the same native snapshot.
State owns construction and validation of resource and member bindings,
including scheduling preset aliases. This pure operation performs no backend
calls and installs no state. Collection export owns config paths and file output.

## Collection lifecycle

For path sources, App delegates native access to `collection.Store.WithCollection`.
`App.ColOpen(ctx, ColOpenOptions{Path, Address, WriteBack})` delegates to
`Store.Open` and returns an already-open `ankiclient.Client` handle with
`Connection` and `Stop`. Its private scope owns one source and one transport.
`WithCollection` defers Stop around its callback, including error and panic paths.
`ankiclient.StartableClient` embeds Client and adds Start for transport owners.
Collection publication and cleanup remain application guarantees of ColOpen.

The application scope opens databases directly or restores packages through
Anki into temporary storage. Stop publishes enabled archives, closes storage and
stops the transport under independent cleanup budgets. Its caller can stop waiting
without cancelling owned cleanup. Session cancellation also initiates cleanup.
The Store keeps source reservations until private native-completion tracking
confirms release; no public restart, collection-phase or Done method is exposed.

Generated services borrow the prepared connection. Active calls remain concurrent;
cleanup drains registered callers, tracks consuming native close/export calls and
retains uncertain restores with a diagnostic. Owned transports can establish final
completion through concrete Done methods, observed by a private interface in the
collection package. Closing an external
network connection cannot establish completion of uncertain native work; recovery
requires external server cleanup before restoring application access.

Enabled archive publication runs independently of callback outcomes. Anki writes
and atomically replaces that archive. Database operations persist in place. Build
separately stages its base, applies to the copy and atomically replaces its output
after success. Failed builds preserve the destination. Callers coordinate access
across different applications and external writers.

## Export and explicit moves

Export reads the collection and writes configuration TOML to caller-owned paths.
`--state-output` optionally writes matching snapshot bindings for onboarding.
Existing stored bindings retain note and note-type filenames on re-export.
New note types receive unused filenames with numeric
suffixes when their names collide with tracked paths, including case-only
differences and fallback names. Without bindings, native edits can change
content-derived filenames.

Export selects a note's directory from its first regular card. Card generation or
changed destinations can therefore change an exported path, including when Anki
chooses a deck for a new root note. When adopting these new paths, declare
`[[moved]]` entries from their stored paths before applying. Export does not
rewrite bindings. Applying the same root-note configs again preserves native
identity and converges; Anki owns their initial deck selection.

Follow [onboarding](#onboarding-an-existing-collection) to bind the whole export
before the first reconciliation, or use `col state import` for selective
bindings. New configs without bindings describe new objects.
State dump exports stored bindings separately; it does not discover bindings
for an untracked config tree:

```sh
al bazel run //projects/anki_as_code -- col state dump \
  --collection /absolute/path/collection.colpkg > /absolute/path/state.anki.toml
```

Export writes declared files directly through `config.Store.Export` and does not reload
generated files. Other entries remain untouched, including old configs and state
files other than an explicit `--state-output` destination. Callers own output
locations and their existing contents. Filesystem errors may leave partial output.

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
  Planning rejects removal of the tracked native default deck (ID `1`), even
  when renamed; retain its declaration or use an explicit move.
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
uses Anki's atomic package export during cleanup. Earlier committed changes can
persist after later request failures; read operations may repackage unchanged content. Build requires
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
Export reads stored bindings and optionally writes a separate state file; it does
not install bindings through Anki's API.
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

`config.New(root, storage)` returns a `*config.Store` bound to a configuration path and
a `filesystem.FS` without
accessing the filesystem. Use a directory for writing; reads also accept a root
config file. The store exposes:

- `Load() (*ankiascodepb.CollectionConfigSet, error)` parses the hierarchy.
- `Dump(format string) ([]byte, error)` loads and encodes JSON or TOML.
- `Export(configs *ankiascodepb.CollectionConfigSet) error` writes the root
  marker and supplied path/config pairs as TOML files beneath its directory.

```go
store := config.New(absoluteDirectory, storage)
configs, err := store.Load()
data, err := store.Dump("json")
err = store.Export(configs)
```

Pass `filesystem.FS{}` for OS storage or `filesystem.New(storage)` for an Afero
backend. Discovery, regular-file and alias checks, directory
creation, reads and writes use the selected filesystem throughout.
Stores retain no parsed messages or open files. Each Load reads fresh values;
callers own filesystem coordination and the returned messages. `config.Version`
supplies the supported schema version to collection exporters.

Export writes directly, preserves input messages, overwrites declared files, and
leaves unrelated/stale files alone. Callers own output locations and existing
contents, valid and unique paths, layout, validation, staging, cleanup, and
publication. Export does not load or normalize its output; failures may leave
partial files. Consumers use the store;
the layout package remains an implementation API with restricted Bazel visibility.

`config/layout` owns the fixed hierarchy, discovery, file I/O, and direct file
export. `internal/codec` converts TOML/JSON through protobuf without filesystem
access. Config layout, config dump, state CLI commands, and collection export
reuse its format dispatch; TOML/JSON helpers are private. Dependencies flow from
config to layout/codec and from layout to codec. State stores have no serialization
dependency.
Collection relationship validation and native operations remain outside them.

## Architecture

![Anki as code runtime architecture](assets/architecture.svg)

The executable contains CLI, MCP and App. Ordinary CLI commands call App;
`mcp serve` starts MCP, whose tools call the same App while an agent communicates
over stdio. App reads and writes Configs and uses gRPC to access Server, which
owns Anki storage, including native resources and persisted application state.

Server abstracts the in-process, subprocess and external-server variants. Its
position here expresses the application boundary; the
[backend selection diagram](#backend-selection) shows where each variant runs
and how it connects.

### Package dependencies

The project keeps one Go module with focused packages. Configs and state are
protobuf contracts. Runtime resolution models stay private to `internal/plan`; packages expose only operations used by their consumers.

![Anki CLI package dependencies](assets/package-dependencies.svg)

Each package's `doc.go` describes its responsibilities and boundaries; this
diagram shows how the packages compose.

Dotted edges show calls through injected contracts. Solid
edges show component dependencies.

| Package                  | Responsibility                                                     |
| ------------------------ | ------------------------------------------------------------------ |
| `cmd/anki_as_code`       | App construction, CLI invocation and process signals               |
| `internal/cli`           | Cobra commands, application dispatch and streams                   |
| `internal/mcp`           | Stdio MCP tools, protocol envelopes and descriptor-derived schemas |
| `internal/collection`    | Source reservations and archive publication                        |
| `api/ankiascodepb`       | Generated config, plan, and persisted-state contracts              |
| `api/ankiascodeapi`      | Shared application interface and collection source contracts       |
| `pkg/ankiascode`         | Application workflows and service composition                      |
| `internal/export`        | Native-to-config projection and direct file output                 |
| `internal/plan`          | Pure resolution, comparison, ownership and prepared operations     |
| `internal/reconcile`     | Execute protobuf plans, bind native IDs and persist state          |
| `internal/rebuild`       | Apply to a staged archive copy and publish it atomically           |
| `internal/resource`      | Native generated clients and fresh resource snapshots              |
| `internal/config`        | Root-bound configuration Store for Load, Dump, and Export          |
| `internal/config/layout` | Fixed hierarchy, discovery, file I/O, and direct export            |
| `internal/filesystem`    | Injected storage, path expansion and filesystem-local aliases      |
| `internal/rpc`           | Compiled upstream descriptor catalog and protobuf RPC execution    |
| `internal/codec`         | Filesystem-independent protobuf/TOML/JSON conversion               |
| `internal/state`         | Stored-state dump, import, update, and replacement                 |

The command entry point constructs App and passes it to `internal/mcp.New`.
It passes App and the resulting concrete server into `internal/cli`. CLI owns
the private `mcpServer` interface with `Serve(context.Context) error` for its
MCP command group; protocol tools call App. This interface is local to CLI
and is not part of `ankiascodeapi`.
App has no MCP dependency or MCP hosting method.
App obtains a prepared connection from
the supplied source and constructs its services locally. Path sources use
collection.Store; borrowed sources retain caller-owned lifetime. `export.New(resources, states, filesystem.FS{})` constructs an
`Exporter` using caller-provided resource and state stores. Its
`Export(ctx, output, stateOutput)` method owns native-to-config projection,
filenames and direct config output, with optional snapshot bindings written as
TOML. State installation remains separate.
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
builder := rebuild.New(apply, filesystem.FS{})
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

Run `al bazel test //projects/anki_as_code/internal/rebuild:rebuild_test`.
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

Store exposes the generated deck, note-type, note, card and rendering clients
without forwarding wrappers. Search and preset clients remain private snapshot
dependencies. Snapshot returns upstream protobuf objects, including unmanaged
resources and resident filtered cards; Snapshot.Presets retains the public map of
preset identities to names.
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
exporter := export.New(resources, states, filesystem.FS{})
err := exporter.Export(ctx, absoluteOutputDirectory, absoluteStateFile)
```

Stored bindings retain filenames and preset aliases.
`Export(ctx, directory, stateFile)` writes configs and optional bindings from
the same snapshot as TOML. An empty state destination selects config-only
behavior. Export never installs state; state operations remain separate.
Callers own output locations and their existing contents.

New note filenames use Anki's `HtmlToTextLine` result for the configured sort
field. Go only sanitizes that text into a bounded Unicode filename. Bound names
are reused from state; filename conversion never changes stored note fields.

Export owns config projection and readable filenames and delegates file output
to `config.Store.Export`. It overwrites declared configs in place and leaves other
entries untouched, including stale configs and state files other than the
explicit state destination.
Filesystem errors may leave partial output. For cleanup or publication policies,
export to a temporary directory and synchronize it externally, for example with
rsync. The caller also owns exclusive offline access, backend startup,
collection preparation and shutdown.

Export and reconciliation share `resource.Store` for generated clients and fresh
native snapshots; their projection and policy remain separate.

## Backend selection

This diagram groups explicit strategies by dependency. External TCP/Unix clients
require no optional local engine. Every enabled engine supplies its subprocess
and in-process strategy, including in combined builds. Cgo defaults to enabled;
WASM defaults to disabled. Generic `subprocess://` and `inprocess://` both prefer
cgo and fall back to WASM when cgo is unavailable.

![Backend strategies grouped by engine dependency](/docs/projects/anki_as_code/assets/backend-strategies.svg)

Explicit TCP and Unix addresses connect to an external server without starting
or linking a local engine.

With `--//projects/anki_grpc:wasm=true`, `subprocess+wasm://` selects a worker
hosting the Rust server through wazero, including in combined builds.
`inprocess+wasm://` hosts the same
server inside the application process through wazero and in-memory gRPC streams.
Disable cgo as well for a pure-Go executable.

With `--//projects/anki_grpc:cgo=true`, native in-process clients and the
`subprocess+cgo://` worker remain available, including in combined builds.
In-process calls use the Go gRPC interface
and cross the C ABI directly, without a network transport or child process.
Generic `inprocess://` uses this native client when cgo is available and falls
back to WASM when only WASM is available. With neither engine it is unavailable.

## Anki backend API

The reusable Go backend and Rust server belong to
[projects/anki_grpc](../anki_grpc/README.md#go-client). The Go package owns the
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
al bazel build -c opt //projects/anki_as_code
al bazel test -c opt //projects/anki_as_code/...
```

E2E targets use Bazel's default shard count. Opt into parallel Go E2E execution
with a count appropriate for the host, for example:

```sh
al bazel test -c opt --test_sharding_strategy=forced=4 \
  //projects/anki_as_code/pkg/ankiascode:ankiascode_test \
  //projects/anki_as_code/pkg/ankiascode:wasm_workflow_test \
  //projects/anki_as_code/cmd/anki_as_code:cli_workflow_test \
  //projects/anki_as_code/cmd/anki_as_code:wasm_cli_workflow_test \
  //projects/anki_as_code/cmd/anki_as_code:composed_cli_workflow_test
```

Use `-c opt` for these suites: the default fastbuild mode compiles Rust and its
embedded WASM guest without optimization. Each shard shares a private WASM
compilation cache among its workers. Bazel's job and resource settings bound
the number of shards running concurrently.

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
Run `al bazel test //projects/anki_as_code/internal/export:export_test`.
