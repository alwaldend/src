## Purpose

Expose native Anki collections as editable TOML through a command without changing collection resources or stored bindings.

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

### Requirement: Isolated command application

The CLI SHALL depend on its private application interface. The concrete App SHALL own native lifetime; its ColStateImport method SHALL load import configs and use the same collection registry. Commands SHALL own flags, binding syntax, encoding and streams. The root SHALL accept only the application interface; collection operations SHALL receive namespace strings explicitly without mutable shared options. The complete command tree SHALL be executable against a fake application without starting Anki.

#### Scenario: Dispatch a state import without native access

- **WHEN** an isolated command receives a collection path, config path and several PATH=ID bindings
- **THEN** it forwards the context, paths and opaque ID strings to one application operation without reading configs or constructing a native store

### Requirement: Export offline collections

The export command SHALL accept an existing offline database or collection package through --collection and a destination through --output. It SHALL write the fixed TOML config hierarchy, use the selected state namespace for names, and release native resources after success or failure. Success status SHALL use stderr without polluting stdout.

#### Scenario: Export either input type

- **WHEN** export receives a valid database or package and output directory
- **THEN** it writes loadable collection, deck, note-type and note TOML configs and releases native resources

### Requirement: Preserve source and state

Export SHALL NOT modify source resources or stored bindings, initialize state, or generate a state file. Package source bytes SHALL remain unchanged. Caller-owned state files and unrelated output entries SHALL remain untouched.

#### Scenario: Repeat export with caller-owned files

- **WHEN** an existing destination includes unrelated files and state files
- **THEN** export overwrites its declared configs while retaining those other files and collection state

### Requirement: Direct output and input failures

The command SHALL reject output physically containing the source collection before writes, including directory aliases on either path. Missing input or native read failures SHALL leave existing declarations unchanged. It SHALL preserve ordinary directory creation and existing permission modes. Filesystem write failures SHALL report errors and may leave partial config output.

#### Scenario: Source lies inside destination

- **WHEN** the destination contains the supplied source collection
- **THEN** export fails without changing either source or destination

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
