## Boundary

internal/state owns the application metadata key, typed bindings, ownership
validation and native identity lookups. Store.New borrows an already-open gRPC
connection. It owns no worker, native collection lifecycle, archive or files.
App.Start composes the existing BackendSTD and BackendCollection. WithState
obtains BackendCollection.Connection() from that prepared application and lends
it to the store without starting or stopping resources. The CLI passes flags,
values and output callbacks; it does not dial, configure or close gRPC connections.
Collection lifecycle remains implemented in projects/anki_grpc, with application
wiring responsible for calling Start and Stop in the correct order.
app.Options owns collection Input; app.StateOptions owns only Key and is passed
by value to App.WithState.
The CLI binds flags to that application type; format and stream handling stay
in the command layer.

## Typed state and operations

State contains version 1 and Resources with independent typed maps of config
paths to deck, note and note-type state, plus preset aliases. Note types own
field/template maps with native IDs and ordinals. Native IDs are unique per
entity map, paths are uniquely owned across maps, and preset aliases may share
IDs. There is no old-shape compatibility layer.

Dump reads the selected collection config key and returns empty version-1 state
when absent. Import accepts parsed configs and opaque path/ID strings, resolves
resources before checking their native ID representation, verifies only the
requested objects, binds note-type members and writes once after the complete
batch succeeds. It changes no native objects or config files.

Update retains omitted map keys and replaces complete supplied entries using
protobuf map merge semantics. Replace installs the complete object. Validate
input and merged ownership before writing. Empty/identical operations write
nothing; input messages and unrelated namespaces remain unchanged. An empty
key selects anki_as_code/state/v1. Callers coordinate writes to a selected key.

Replace validates its input first, then reads the selected raw metadata value.
An undecodable or invalid discarded object does not prevent repair or clearing.
Native read failures remain errors; they are never treated as invalid metadata.
Decode the previous value only to retain no-op detection for valid identical
state. Dump, Import and Update continue to require valid stored state.

## CLI and encoding

State commands require --input for an existing regular .anki2 database. Reject
missing sources and other extensions before starting a backend so state commands
cannot accidentally create a collection or discard changes in a temporary archive
restore. The application resolves the packaged server, starts its managed backend,
starts BackendCollection once through App.Start, then invokes state operations
through WithState. Repeated and nested calls reuse the prepared backend. Start
rejects another startup attempt after resources are initialized; WithState rejects
use before startup or after shutdown. App.Stop closes the collection before stopping the worker, preserves operation and
cleanup errors, and has a bounded context independent of operation cancellation.
The CLI calls Start once after parsing command inputs and defers Stop around
command execution, including failures. There is no separate transport or
client-lifecycle implementation in the CLI.

Dump emits TOML by default or JSON with --format json. Update/replace read a file
or stdin when omitted/-. Import accepts multiple PATH=ID arguments plus --text.
Status/errors use stderr; dump uses stdout. --state-key selects the namespace.
Missing or unknown state subcommands fail with stderr output and nonzero status;
explicit --help remains successful.

config/codec exposes Encode and Decode with a typed Format. Fixed TOML callers
use codec.TOML; codec.JSON selects ProtoJSON. CLI strings convert at the adapter
boundary and unsupported values fail in the codec. Empty format defaults to
TOML. No new serializer, reflection layer or dependency is introduced.

## Verification

Write managed-backend E2E cases before production changes; observe the old CLI
reject --input. The test fixture uses a managed worker for Go API checks and closes
its collection before giving the same database to each built CLI invocation.
Reopen through the native API after success or failure to verify persistence,
cleanup, native-object preservation and namespace isolation. Run commands with an
empty executable search path to verify packaged server resolution. Reject missing,
directory and archive inputs without creating a database. Verify Go API connection
reuse, both text formats, batch atomicity and invalid-state recovery. Emit TOML,
JSON and a closed native database.

General collection/archive lifecycle tests remain with their owners. Validate
pure config and command syntax E2E, quality, affected-target lint, CLI/site builds
and strict OpenSpec against the exact prepared candidate. Delivery receipts own
publication and review resolution evidence.

The application E2E prepares a native database and explicitly starts App once.
Nested and sequential WithState operations must share the collection, including
after a rejected second Start and callback/cancellation errors. Calls before
Start or after Stop must reject without invoking their callbacks. A second
application reopens the database to check persistence and emits JSON plus a
closed native database. General process lifecycle mechanics remain covered by
projects/anki_grpc.
