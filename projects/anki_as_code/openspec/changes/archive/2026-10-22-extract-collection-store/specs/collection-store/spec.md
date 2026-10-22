## Purpose

Coordinate source reservations and archive publication around managed native collection connections.

## ADDED Requirements

### Requirement: Callback-scoped collection access

collection.Store.WithCollection SHALL accept a source path and inject a prepared grpc.ClientConnInterface into its callback. Callers SHALL complete operations sequentially before callback return. Store SHALL reserve canonical source paths and file identities before startup and reject duplicate aliases. It SHALL expose no Collection service container, reference IDs or separate lookup API. The guarantee SHALL apply within one Store; external writers remain caller-owned.

#### Scenario: Several operations in one callback

- **WHEN** several operations use the connection injected by WithCollection
- **THEN** they share one prepared native collection, and retained connections reject calls after cleanup

#### Scenario: Reject duplicate opens

- **WHEN** a reserved source is reopened through its original path, a symlink or a hard link
- **THEN** the second open fails without starting another worker

### Requirement: Caller-owned service assembly

Callers SHALL construct resource/state stores and exporters on the injected connection and MAY reuse local instances for multiple operations. Store SHALL retain no service or namespace caches. Services SHALL read fresh native data and SHALL NOT own another native lifetime. App SHALL delegate collection access to Store without constructing workers or issuing native collection open/close calls.

#### Scenario: Use services across operations

- **WHEN** a caller writes state, reads it again and exports configs in one scope
- **THEN** each service uses the same prepared connection and subsequent reads observe the writes

### Requirement: Scope publication and cleanup

Store SHALL publish package scopes back to their source path after a successful callback requests publication, using Anki's native export. Callbacks that return false, callback errors and cancelled scopes SHALL skip package export. Store SHALL NOT expose destination selection or output reservations. Callers SHALL export other destinations directly through Anki as their final operation and return false to skip source publication. Cleanup SHALL stop the collection before its worker on all paths. Direct database writes SHALL remain committed after later failure.

#### Scenario: Explicitly export another destination

- **WHEN** a caller performs state operations, directly exports a distinct archive through Anki and returns false
- **THEN** the source archive remains unchanged, cleanup succeeds, and the output reopens through Anki with the changes

#### Scenario: Failure after native changes

- **WHEN** a package callback fails or its scope is cancelled after writing state
- **THEN** the source archive remains unchanged, native owners are cleaned up and its source reservation is released

### Requirement: Manager-owned native lifetime

Store SHALL accept an optional backend.Manager dependency selected before use and use backend.ManagerCollection for native preparation and cleanup. Its zero value SHALL retain bundled-server/PATH defaults. Store SHALL borrow the prepared connection without constructing or starting/stopping native backends. Reservations SHALL remain held until manager cleanup completes. Callback errors, standard cancellation errors, custom cancellation causes and panics SHALL remain observable; cancellation SHALL prevent archive publication.

#### Scenario: Inject a transport manager

- **WHEN** a caller supplies a transport manager to Store
- **THEN** the collection and its services use that manager's native backend, cleanup leaves transport ownership with the manager and reservations are released after scope exit

### Requirement: Typed application state workflows

App SHALL expose typed ColStateDump, ColStateImport, ColStateUpdate and ColStateReplace workflows with explicit collection paths and namespace keys. It SHALL load import configurations, construct collection-bound services and determine publication through collection.Store. It SHALL NOT expose WithCollection, raw connection callbacks or CLI flag options. CLI code SHALL own argument parsing, state serialization and output while preserving existing command syntax and formats.

#### Scenario: Manage state through App

- **WHEN** a Go caller imports, updates, replaces or dumps state on a database or collection package
- **THEN** App manages the complete scope and native publication internally while preserving namespace isolation and input objects

#### Scenario: Reject an invalid state workflow

- **WHEN** a workflow is cancelled or a batch fails validation
- **THEN** it returns an error without publishing a partial package or retaining the source reservation
