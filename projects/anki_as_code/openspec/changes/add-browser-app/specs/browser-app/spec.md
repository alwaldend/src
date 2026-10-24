## Purpose

Allow browser consumers to run configuration and collection workflows locally
through a reusable asynchronous application API and explicitly save completed
work across reloads.

## ADDED Requirements

### Requirement: Browser distribution

The project SHALL provide a browser WebAssembly distribution with a JavaScript module, worker and matching runtime support. Loading SHALL resolve only when the application is ready and SHALL reject on startup failure.

#### Scenario: Load the packaged application

- **WHEN** a consumer serves the distribution over HTTP and creates an application
- **THEN** it receives a usable application and workspace without invoking a CLI or external server

#### Scenario: Invalid WebAssembly resource

- **WHEN** the requested WebAssembly resource cannot be fetched or instantiated
- **THEN** creation rejects with an error and releases its worker

### Requirement: Direct configuration operations

The browser application SHALL provide promise-based configuration export and dump methods with the existing schema and application semantics. Export SHALL accept a protobuf JSON object; dump SHALL return JSON or TOML text. Invalid input and application errors SHALL reject the affected operation without preventing subsequent valid operations.

#### Scenario: Configuration round trip

- **WHEN** a consumer exports a collection containing Unicode and multiline fields and then dumps its output tree
- **THEN** the loaded configuration preserves those field values and established collection path semantics

#### Scenario: Recover after an invalid request

- **WHEN** a malformed configuration or missing input path causes a request to fail
- **THEN** that request rejects and a subsequent valid request can succeed in the same application

### Requirement: Persistent session workspace

Each application SHALL own an isolated virtual filesystem that persists across its operations. Consumers SHALL be able to write and read binary files, list files recursively and remove paths. Operations SHALL execute in request order. File transfer SHALL copy bytes without detaching the caller's input. Workspace paths SHALL be absolute virtual paths; storage SHALL not consult the host filesystem.

#### Scenario: Edit and reread

- **WHEN** a consumer imports TOML bytes, dumps them, edits a file and dumps again
- **THEN** the second dump reflects the edit in the same workspace

#### Scenario: Independent workspaces

- **WHEN** two applications use the same virtual path
- **THEN** their files and configuration operations remain independent

#### Scenario: Reject incompatible file writes

- **WHEN** a consumer writes bytes to a directory or a path with a regular-file ancestor
- **THEN** the write rejects before changing any existing files or directories
- **AND** a subsequent save and restore retains their paths and contents
- **AND** creating new files and replacing regular files remain supported

#### Scenario: Rename a completed file safely

- **WHEN** a consumer renames a regular file to an absent path or another regular file with an existing directory parent
- **THEN** the file is moved atomically
- **AND** directory sources or destinations are rejected without changing any files, including after save and restore

### Requirement: Explicit application lifetime

The returned application SHALL provide disposal that terminates its worker and
rejects pending and future requests. Unexpected worker or runtime failure SHALL
reject outstanding work. Active workspace data SHALL last only for the
application's lifetime; completed local snapshots SHALL survive disposal.
Normal collection operation completion SHALL join owned backend cleanup before
the next operation begins.

#### Scenario: Dispose with work pending

- **WHEN** a consumer disposes an application immediately after submitting a request
- **THEN** the pending request and subsequent requests reject and repeated disposal is safe

### Requirement: Reuse collection workflows in the browser

The browser API SHALL invoke the existing App's collection workflows through
the managed browser client and Rust's JavaScript command API without CLI parsing
or duplicate business logic. Go and Rust SHALL access one Afero workspace.
Asynchronous requests SHALL execute serially through their complete lifecycle.
Backend loading and creation SHALL support asynchronous completion; Rust commands
SHALL execute synchronously in the worker and SHALL NOT promise preemption by
cancellation. Collection packages and media SHALL remain binary data owned by
Anki's existing package semantics.

#### Scenario: Plan and apply a collection edit

- **WHEN** a consumer imports a collection with binary media, exports configs, edits them and invokes planning followed by application
- **THEN** the existing App performs those operations and publishes a completed package retaining the collection's media and unaffected scheduling and review history

### Requirement: Lazy backend activation

The distribution SHALL load the callable Rust WASI reactor only when a collection
operation requires it. Configuration and collection workflows SHALL operate
without JSPI. Configuration-only workflows SHALL remain usable without a
successfully fetched backend artifact. Backend startup failures SHALL reject the
requesting operation with an observable error.

#### Scenario: Configuration-only session

- **WHEN** a consumer exports, edits and downloads configuration without opening a collection
- **THEN** the application completes those operations without fetching the Rust server

#### Scenario: Collection operation without JSPI

- **WHEN** a consumer invokes a collection workflow in a browser without JSPI
- **THEN** the application loads the callable backend and completes the operation through its JavaScript command API

#### Scenario: Unavailable backend

- **WHEN** a collection operation cannot load or start its backend
- **THEN** that operation rejects and configuration editing and downloads remain available

### Requirement: Atomic local snapshots

The browser API SHALL provide explicit save and restore operations using
IndexedDB. A saved snapshot SHALL contain editable configuration files and,
when a collection is selected, its completed package from one completed
operation boundary, committed atomically with revision metadata.
Configuration-only snapshots SHALL be supported without backend activation.
Runtime files under reserved `/tmp` SHALL be excluded from saved/restored
workspace data. Save SHALL report success only after the
transaction commits. It SHALL NOT snapshot an open SQLite database or silently
overwrite work saved by a conflicting writer.

#### Scenario: Save and reload

- **WHEN** a consumer saves an applied collection and configs, disposes the app, then creates an app and restores the saved workspace
- **THEN** it receives the same completed collection, binary media and exact editable config bytes

#### Scenario: Save configuration-only work

- **WHEN** a consumer saves and restores configuration files without opening a collection
- **THEN** it receives the exact editable bytes without activating the backend or including runtime temporary files

#### Scenario: Failed save

- **WHEN** quota exhaustion or a transaction abort prevents a save
- **THEN** the previous saved snapshot remains complete, current edits remain available and the save rejects

#### Scenario: Conflicting writer

- **WHEN** two app instances supply the same expected revision from a prior save or restore and one saves before the other
- **THEN** the second save cannot silently replace the newer snapshot and reports a conflict without discarding either instance's current work
- **AND** revision zero permits creation only, including from a fresh replacement worker

#### Scenario: Failed restore

- **WHEN** a saved snapshot cannot be read or restored completely
- **THEN** restore rejects and the existing active workspace remains intact

#### Scenario: Conflicting saved paths

- **WHEN** a saved snapshot declares a file at a directory path or beneath a regular-file ancestor, in either entry order
- **THEN** restore rejects without replacing current files or changing stored snapshots
- **AND** the retained workspace can still be saved and restored without losing bytes
- **AND** compatible child entries before their explicit directory entries remain supported
