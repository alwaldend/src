## Purpose

Reconcile offline Anki collections through path-addressed configurations and persistent collection state.

## ADDED Requirements

### Requirement: Delegate prepared collection access

App collection methods SHALL accept source path strings and delegate native
access to collection.Store.WithCollection. Scoped Go callers SHALL use the
prepared grpc.ClientConnInterface supplied by that store. Callers SHALL construct
and reuse local services within the callback. The collection-store capability
SHALL remain authoritative for source reservations and archive publication;
backend managers SHALL own native preparation and cleanup. App SHALL expose no collection reference selector or
scope lookup API.

#### Scenario: Reuse native services

- **WHEN** a caller performs several operations in one collection.Store callback
- **THEN** the operations share its injected connection and locally constructed services

### Requirement: Stable noncolliding note-type filenames

Export SHALL reserve tracked note-type filenames before naming untracked types, regardless of native ID order. Tracked types SHALL retain their filenames after native renames. Untracked types SHALL receive deterministic unused filenames with numeric suffixes, considering case-insensitive collisions, occupied suffixes and fallback names. Exported notes SHALL reference the allocated paths without modifying the source collection.

#### Scenario: A new type reuses a tracked type's former name

- **WHEN** a tracked type is renamed in Anki and a new untracked type uses its former name
- **THEN** export retains the tracked filename, gives the new type an unused filename, and repeats the same output after its bindings are explicitly imported

### Requirement: Path-addressed configuration

The collection SHALL declare version 1 in collection.anki.toml and use fixed note_types/, notes/, decks/, and nested child decks/ directories. Each deck directory SHALL require deck.anki.toml and derive its title from directory names. Parsed resource wrappers SHALL retain collection-relative source paths. Unknown fields, unsupported versions, unsuitable filesystem entries, and malformed syntax/types SHALL fail; native relationship validation SHALL remain outside internal/config. Historical aliases and removed layouts SHALL NOT be supported.

#### Scenario: Filesystem ownership

- **WHEN** notes, deck markers, and note types are loaded from the fixed hierarchy
- **THEN** their parsed wrappers retain unique relative paths without modifying sources or starting Anki

### Requirement: Identity-free resource configurations

Deck/note-type IDs and note GUIDs SHALL NOT be editable fields. Note-type/deck references SHALL use paths, preset references SHALL be readable names, and sort fields SHALL use names. Template order SHALL follow the array. Upstream native identity/provenance and derived metadata SHALL be omitted from export and remain outside the editable semantics. Card ordinals and semantic field-role tags SHALL retain their meaningful behavior. config generate-guid SHALL be removed.

#### Scenario: New configuration

- **WHEN** a user writes a new note or deck file without an ID or GUID
- **THEN** it is inspectable as a path-addressed declaration without an identity-generation step

### Requirement: Pure configuration operations

Load and Dump SHALL leave sources unchanged, support JSON default and TOML output, preserve exact values and optional presence, and perform only syntax/schema/filesystem checks. Export SHALL directly write supplied path/config objects without normalization, identity allocation, validation, cleanup, or directory publication. ProtoJSON SHALL be converted directly into TOML using existing codecs. TOML datetimes SHALL require quoted string content.

#### Scenario: Direct export

- **WHEN** a caller exports a supplied configuration set
- **THEN** each resource is written at its declared relative path, the supplied object is unchanged, and unrelated files are retained

### Requirement: Explicit moves

The root config SHALL support moved entries with from/to paths. Config parsing SHALL preserve these entries without applying state transformations.

#### Scenario: Inspect a move

- **WHEN** config dump loads a directory-prefix or individual-file move
- **THEN** it prints the declared move without changing any source or requiring Anki

### Requirement: Collection-backed state

The CLI SHALL store schema-backed path/native bindings through Anki configuration APIs in the collection. Export SHALL write configs only without changing the collection or generating state files. It SHALL read stored bindings only for naming and preset aliases. State export SHALL remain a separate operation; state dump SHALL return only stored bindings. State dump and replace SHALL default to TOML and support --format json. State import SHALL bind declared paths to existing native objects in one batch without modifying their content. State replace SHALL explicitly install a complete state object. apply/build SHALL persist native object changes and state together without rewriting source configurations. State-only changes SHALL be publishable. Reconciliation SHALL use collection-backed state after explicit installation; no SQL or GUID-generation command SHALL be required.

Collection-source commands SHALL use `--collection` for Anki input; build SHALL retain its explicit `--base` archive option. State commands SHALL accept both `.anki2` and `.colpkg` through the regular collection workflow. Archive restore and export SHALL use Anki APIs; callers SHALL NOT unpack files themselves. Dump and unchanged writes SHALL retain archive bytes, and changed state SHALL use native atomic export. The state store SHALL only borrow a prepared connection.

#### Scenario: Persistent bindings

- **WHEN** exported bindings are explicitly installed and the collection is closed, reopened, or rebuilt into a colpkg
- **THEN** native identities remain bound to paths and subsequent planning converges

### Requirement: Readable stable note filenames

Collection export SHALL derive new note filenames from the configured sort field's visible text, preserve readable Unicode, sanitize unsafe characters, bound filename length, and use a fallback for empty/media-only content. It SHALL disambiguate duplicates and reserve bound filenames before allocating new ones. Re-export SHALL reuse existing bound filenames despite content or sort-field edits whenever current deck destinations allow it. Naming SHALL NOT modify stored note fields or use content as identity. Initial text conversion SHALL use Anki's HtmlToTextLine RPC, with only filesystem sanitization performed locally. Bound filenames SHALL bypass text conversion.

#### Scenario: Re-export edited and duplicate notes

- **WHEN** a bound note's content changes and new notes have matching sort-field text
- **THEN** the bound note retains its filename, each new note receives a distinct readable filename, and planning converges

### Requirement: State-driven lifecycle and moves

Reconciliation SHALL resolve explicit moves idempotently before object operations, create unbound declarations, and delete only tracked objects absent from configuration. It SHALL reject conflicting bindings, unsafe moves, and destructive changes affecting untracked/filtered resources before publication. Field and template moves SHALL rename members within the same note type, without transferring native identity between note types. Native APIs SHALL retain scheduling, history, presets, media, and metadata outside the managed surface.

#### Scenario: Rename a deck directory

- **WHEN** a moved entry accompanies a directory rename
- **THEN** the deck, descendant bindings, notes, and retained cards preserve native identities and source files remain unchanged

#### Scenario: Import an existing object

- **WHEN** import binds a config path to an existing native ID
- **THEN** its content/review state is unchanged and later reconciliation updates that same object

### Requirement: Native archive and direct text export

Archive operations SHALL restore through BackendCollection and call Anki's atomic package export directly. The CLI SHALL NOT add archive hashing, staging, sync or rename. Failures before package export SHALL leave source archives unchanged. Database inputs SHALL use direct native writes. Text export SHALL write only config files directly without reloading generated files. Filesystem failures may leave partial output. It SHALL perform no resource or state writes; source archives SHALL remain byte-identical.

#### Scenario: Conflicting move

- **WHEN** a move targets another bound resource
- **THEN** plan/apply fail without changing the collection or configuration sources

### Requirement: Native collection package operations

The collection workflow SHALL call Anki's package import/export services directly,
including media. Repository operations SHALL borrow an already-open connection
without managing process or collection lifetime. No handwritten package codec or
separate archive resource abstraction SHALL be required.

#### Scenario: Native packaging and atomic replacement

- **WHEN** native packages are reconciled or rebuilt
- **THEN** supported media, identities, history, and state survive native reimport,
  while callback failure or cancellation before export leaves the source package untouched

### Requirement: Resolve import resources before interpreting identifiers

CLI import parsing and the Go import API SHALL carry identifiers as strings.
The CLI SHALL check binding syntax and duplicate paths without assuming integer
identifiers. State import SHALL resolve the target resource before validating
its identifier against the selected native API type.

#### Scenario: Invalid identifier for an undeclared resource

- **WHEN** an undeclared resource path is supplied with a nonnumeric identifier
- **THEN** resource resolution fails before numeric conversion

### Requirement: Nested state objects and selectable collection keys

The persisted protobuf contract SHALL be named State. Its Resources field SHALL contain separate
maps of config paths to DeckState, NoteState and NoteTypeState, and aliases to
PresetState objects. NoteTypeState SHALL own typed FieldState and TemplateState
maps. A generalized resource binding or kind discriminator SHALL NOT be used.
Paths SHALL be uniquely owned across entity maps; native IDs SHALL be unique
within each entity map except preset aliases. Synthetic member fragments and
preset addresses SHALL NOT appear as config path keys.

The CLI and Go APIs SHALL select a state namespace. The default key SHALL be
anki_as_code/state/v1. Operations SHALL read and write only the selected key and
preserve other namespaces through native package round trips.

#### Scenario: Run independent state workflows

- **WHEN** two stores select different keys in the same collection
- **THEN** dumping, importing, and replacing one state leaves the other unchanged

#### Scenario: Preserve nested identities through a move

- **WHEN** a bound note type or one of its members is renamed through reconciliation
- **THEN** the nested native identity and ordinal remain bound without synthetic state paths

### Requirement: Merge an external state object

The Go Store.Update method and state update command SHALL merge a supplied
version-1 State into the selected collection namespace. Omitted map keys SHALL
remain; supplied keys SHALL replace complete typed entries. Invalid input or
merged identity/path conflicts SHALL fail before writing. Neither input objects,
native resources nor other namespaces SHALL be modified. Empty or identical
updates SHALL perform no state write and retain archive bytes. CLI input SHALL support file/stdin,
TOML by default and optional JSON through the shared codec.

#### Scenario: Merge exported bindings with existing state

- **WHEN** update supplies valid new bindings or replacements for existing keys
- **THEN** those entries are installed together while omitted bindings remain

#### Scenario: Reject a conflicting update

- **WHEN** input is valid alone but conflicts with retained native/path ownership
- **THEN** the complete update fails without publishing any changes

### Requirement: Explicit prepared collection lifecycle

collection.Store SHALL delegate preparation and cleanup to backend.ManagerCollection
over an injected transport manager, with the bundled server as its default. The shared
backends SHALL own readiness, prepared connections and temporary restore files.
The CLI SHALL NOT duplicate lifecycle flags, native open/close, restore ownership
or journal checks. Feature operations SHALL borrow connections and SHALL NOT own
native lifetime. Read-only and unchanged operations SHALL skip archive export.

#### Scenario: Borrow a prepared collection

- **WHEN** the shared backend prepares a supported collection successfully
- **THEN** callers obtain its prepared gRPC connection until export or close

#### Scenario: Failed preparation releases resources

- **WHEN** native preparation fails or is cancelled
- **THEN** operations receive no connection, both backends are stopped, and cleanup follows the backend library contract without removing supplied databases

### Requirement: Direct database access

BackendCollection SHALL open a supplied offline database in place without making
a temporary copy. Direct native writes SHALL persist without replacement
publication. Backend shutdown SHALL preserve the supplied database. Package inputs
SHALL use BackendCollection's default restore ownership and Anki's atomic export.
Direct writes SHALL remain committed if a later operation fails.

#### Scenario: Native write followed by caller failure

- **WHEN** a direct database operation commits an API write and then fails
- **THEN** backend shutdown releases the database and the committed change remains in it

### Requirement: Read-only declared planning

Plan SHALL compare an API snapshot with resolved declarations in memory. It
SHALL NOT apply resource updates, write state or create a disposable database.
It SHALL report requested resource and known card-assignment changes, indicate
deferred card generation, and leave normalization and final template/generation
validation to Anki during apply. It SHALL NOT predict exact generated-card
counts.

#### Scenario: Plan changes to an existing database

- **WHEN** edited and new note declarations are planned repeatedly
- **THEN** each plan describes the same pending changes while database contents
  and bindings remain unchanged and no database copies are created

### Requirement: Retained cards survive deck removal

Plan and preflight SHALL reject deleting a deck referenced by retained normal cards
without a declared surviving destination. Omitted card assignments SHALL NOT
request deletion of those cards or their note. Rejection SHALL precede native
writes for both database and archive inputs.

#### Scenario: Retained root note without a destination

- **WHEN** a tracked note moves to root notes with no card assignments and its old deck is removed
- **THEN** plan and apply reject the request without changing the collection; assigning a surviving normal deck makes the move applicable

### Requirement: Converged database apply is a no-op

Apply SHALL compare editable note-type values independently of runtime hints and
unmanaged metadata, while retaining explicit member identity/order semantics.
A converged database apply SHALL leave source bytes and state unchanged.

#### Scenario: Apply unchanged exported configs

- **WHEN** exported bindings are installed and the unchanged configs are applied repeatedly
- **THEN** plan remains empty and database bytes remain unchanged

### Requirement: Native creation batches retain declaration identities

Reconciliation SHALL use Anki's existing AddNotes operation for new-note batches
and bind each returned native ID to its corresponding declaration. Notes SHALL
resolve existing native IDs through a single precomputed index. Command lifecycle
SHALL remain owned by the shared collection workflow.

#### Scenario: Multiple notes in different decks

- **WHEN** a batch creates distinct notes across declared decks
- **THEN** each path binds its own native note, content and generated-card assignments survive re-export, and subsequent planning converges

### Requirement: Capability package ownership

App SHALL route command workflows and compose the shared native backends. Export
SHALL own prepared-connection export projection and direct filesystem output,
have no direct codec dependency, and trust serialization without reloading
generated declarations. The shared backends SHALL own process/collection lifetime.
Exporter SHALL receive existing resource.Store and state.Store through New, without a reconciliation dependency; its Export method
SHALL take only context and output path as operation inputs.
Config SHALL own config files and layout;
the shared codec SHALL own filesystem-independent protobuf/TOML/JSON conversion.
internal/plan SHALL own reference resolution, ownership guards, comparison and
operation preparation with private runtime models. internal/reconcile SHALL
execute the supplied protobuf plan, bind Anki-returned identities and persist the
resulting state. Export SHALL write only declared config files, leaving
other entries untouched. Callers SHALL own cleanup, synchronization and publication. State SHALL retain
its lifecycle-free prepared-connection Store. Reconciliation Apply SHALL return
only a changed flag, derived from native OpChanges and persisted binding changes,
and SHALL NOT fetch another full note/card snapshot to verify identities.
Direct database commands SHALL NOT compute unused source fingerprints.

#### Scenario: Use a prepared collection for reconciliation

- **WHEN** an app workflow borrows a prepared collection connection
- **THEN** internal/plan prepares operations from supplied configs, snapshot and state; internal/reconcile executes the supplied plan without owning files/processes; export owns text output and Anki owns archive publication

### Requirement: Empty deck before changing kind

Plan and apply SHALL reject conversion of an occupied normal or filtered deck to
the other kind before native writes. An empty deck SHALL remain convertible,
without altering cards, scheduling or history already moved to a normal deck.

#### Scenario: Convert a normal deck containing retained cards

- **WHEN** retained notes move to root notes without card assignments and their occupied normal deck becomes filtered
- **THEN** plan and apply reject the conversion without changing database or archive bytes

### Requirement: Preserve unrelated export permissions

Collection export SHALL leave all unselected files and directories untouched,
including their bytes and permission bits. It SHALL NOT copy, stage, clean, replace
or roll back the output tree.

#### Scenario: Refresh an export containing private unrelated entries

- **WHEN** an existing export contains an unrelated 0700 directory and 0600 file
- **THEN** direct writes leave their contents and modes unchanged

### Requirement: Persistence belongs to the state store

state.Store SHALL own persisted dump, import, update and replacement. plan.Build SHALL accept caller-supplied configs, native snapshot and stored state; moves and proposed bindings SHALL remain private to reconciliation without a store preview API or persisted writes. Apply SHALL record actual native identities through validated store replacement; provisional identities SHALL NOT pass persisted-state validation.

#### Scenario: Plan without persisting bindings

- **WHEN** a caller supplies configs, captured native resources and stored bindings to plan.Build
- **THEN** private move resolution leaves those inputs and persisted state unchanged, and a later apply records actual identities through state.Store.Replace

### Requirement: Staged copy-and-apply builds

Build SHALL copy the base collection package to a private temporary directory on the destination filesystem, invoke the normal Apply workflow on the copied archive, and atomically replace the destination only after successful apply and cleanup of the collection. It SHALL remove staging on success or failure, preserve the base archive and any existing destination on failure before publication, reject source aliases and require caller-owned exclusive destination access, and leave other open scopes publishing only their own source archives.

#### Scenario: Build independently of an open source

- **WHEN** a caller builds from a source currently open in a scope and later changes that source scope
- **THEN** Build returns with its completed independent output, later changes affect only the source, and the output retains the built state

#### Scenario: Failed apply leaves no partial destination

- **WHEN** staged apply fails with an absent or existing destination
- **THEN** the destination remains absent or byte-identical, and the temporary archive directory is removed

### Requirement: Collection command group and aliases

Native collection operations SHALL be exposed under col with alias collection: export, plan, apply, build, and state. The state-key flag SHALL belong to that group. Text-only config operations SHALL remain at the root with alias conf. Removed ungrouped collection command spellings SHALL fail before invoking application workflows.

#### Scenario: Dispatch canonical and aliased collection commands

- **WHEN** a caller invokes col export or collection export with the same arguments
- **THEN** both dispatch the same workflow with unchanged contexts, namespace and streams

#### Scenario: Inspect configs through either name

- **WHEN** a caller invokes config dump or conf dump
- **THEN** both run the text-only workflow without a collection or namespace flag

### Requirement: Flat namespaced application operations

The public pkg/ankiascode package SHALL provide concrete App with ConfigDump, ConfigExport,
ColExport, ColPlan, ColApply, ColBuild, ColStateDump, ColStateImport, ColStateUpdate
and ColStateReplace. Collection arguments SHALL be path strings. internal/cli
SHALL own a private application interface with those command operations, using
protobuf messages and standard-library types without importing App or native
backend implementations. cmd/anki_as_code SHALL construct App and inject it into
cli.New. The complete command tree SHALL remain executable against one fake.
Handler interfaces, accessor hierarchy, collection references and the separate
API package SHALL be removed. Service business rules SHALL remain unchanged.

#### Scenario: Isolated command dispatch

- **WHEN** the CLI constructor receives a fake implementing its private application interface
- **THEN** commands invoke the matching Config, Col or ColState operation with unchanged contexts, paths, errors and streams, without constructing App or starting Anki

### Requirement: Export a serialized configuration set

App.ConfigExport SHALL write a supplied CollectionConfigSet through config.Store.Export
without native collection access. config export [FILE] --output DIRECTORY SHALL
decode JSON by default and support TOML through --format toml; omitted FILE or -
SHALL read stdin. The conf alias SHALL work identically. Schema, format and read
errors SHALL fail before app invocation. The existing store SHALL own direct TOML
writes, preserve input objects and unrelated output files, and leave semantic
validation and valid representable paths to callers. Filesystem failures may leave
partial output. Success SHALL write status to stderr and leave stdout empty.

#### Scenario: Round-trip the config tree

- **WHEN** a caller dumps a config tree as JSON and exports that object from a file or stdin
- **THEN** loading the resulting TOML tree preserves the declared resources and exact content without opening Anki

#### Scenario: Reject unreadable or malformed input

- **WHEN** the input file cannot be read or the serialized object fails its protobuf schema
- **THEN** the command fails without creating the destination or emitting success output

### Requirement: Assign cards using the resulting deck kind

After converting an empty filtered deck to a normal deck, Apply SHALL treat cards generated there in the same operation as regular cards. It SHALL honor both inherited deck destinations and explicit ordinal overrides, persist their bindings and converge on the next plan. Existing reviewed cards SHALL retain scheduling and history.

#### Scenario: Convert and create notes together

- **WHEN** an empty filtered deck becomes normal and a new note declares regular-card assignments in that deck
- **THEN** plan reports the conversion and note creation, apply succeeds for databases and packages, and each generated card receives its declared destination

### Requirement: Native execution consumes a prepared plan

Planning SHALL prepare native operation payloads, member mappings, expected stored state and desired bindings as protobuf data. Apply SHALL execute that data without config resolution, resource discovery or diffing. It SHALL reject changed stored bindings before writes, preserve the input plan and bind Anki-returned identities before persistence. Callers SHALL own snapshot freshness and exclusive offline access.

#### Scenario: Execute a detached plan

- **WHEN** a captured plan is serialized and decoded before execution on the same prepared collection
- **THEN** execution performs the prepared work without native inventory reads, retains the original plan and records actual resource identities

### Requirement: Filtered cloze template identity

Planning SHALL protect the shared template used by resident filtered cloze cards at every cloze ordinal. Replacement or removal SHALL fail before execution; explicit moves that retain the native template mapping SHALL remain permitted.

#### Scenario: Preserve resident cloze cards

- **WHEN** a tracked cloze note has resident filtered cards at ordinals zero and one
- **THEN** both retain the shared template identity and planning rejects a replacement without managing those cards
