# collection-store Specification

## Purpose

Coordinate source reservations and archive publication around managed native collection connections.

## Requirements

### Requirement: Callback-scoped collection access

collection.Store.WithCollection SHALL accept CollectionOptions selecting Path, Address and WriteBack and inject a prepared grpc.ClientConnInterface into its callback. Callers SHALL complete operations sequentially before callback return. Store SHALL reserve canonical source paths and file identities before startup and reject duplicate aliases. It SHALL expose no Collection service container, reference IDs or separate lookup API. The guarantee SHALL apply within one Store; external writers remain caller-owned.

#### Scenario: Several operations in one callback

- **WHEN** several operations use the connection injected by WithCollection
- **THEN** they share one prepared native collection, and retained connections reject calls after cleanup

#### Scenario: Reject duplicate opens

- **WHEN** a reserved source is reopened through its original path, a symlink or a hard link
- **THEN** the second open fails without starting another client

### Requirement: Caller-owned service assembly

Callers SHALL construct resource/state stores and exporters on the injected connection and MAY reuse local instances for multiple operations. Store SHALL retain no service or namespace caches. Services SHALL read fresh native data and SHALL NOT own another native lifetime. App SHALL delegate managed path access to Store without issuing native collection open/close calls; borrowed access SHALL follow collection-source-api.

#### Scenario: Use services across operations

- **WHEN** a caller writes state, reads it again and exports configs in one scope
- **THEN** each service uses the same prepared connection and subsequent reads observe the writes

### Requirement: Scope publication and cleanup

Store SHALL delegate package write-back to its native collection client during
cleanup according to CollectionOptions.WriteBack. Enabled write-back SHALL
attempt publication even after callback errors or cancellation; disabled
write-back SHALL leave the source archive unpublished. Callbacks SHALL return
only an error. Direct database writes SHALL remain committed after later failure.
Store SHALL expose no destination-selection or output-reservation API. Callers
exporting other destinations SHALL do so through Anki and MAY disable source
write-back. Collection cleanup SHALL precede client cleanup.

Offline databases are opened in place without a temporary database copy or replacement publication.

Managed package replacement SHALL use Anki's native atomic export. Collection scopes SHALL NOT duplicate that export with application-owned archive hashing, staging, synchronization or rename; separate staged builds SHALL follow archive-build.

#### Scenario: Explicitly export another destination

- **WHEN** a caller disables source write-back, performs state operations and directly exports a distinct archive through Anki
- **THEN** the source archive remains unchanged, cleanup succeeds, and the output reopens through Anki with the changes

#### Scenario: Failure after native changes

- **WHEN** a package callback fails or is cancelled after writing state with write-back enabled
- **THEN** cleanup attempts to publish committed native state while preserving operation and cleanup errors
- **AND** the source reservation remains held until native cleanup finishes

#### Scenario: Native write followed by caller failure

- **WHEN** a direct database operation commits a native API write and then fails
- **THEN** cleanup releases the database and the committed change remains in it

### Requirement: Managed client lifetime

Store SHALL delegate address interpretation to its configured ClientFactory,
compose ankiclient.NewCollection over the selected client, and own startup and
cleanup under the collection-lifecycle contract. Without a configured factory,
an explicit registered address SHALL be required; Store SHALL NOT discover a
worker through bundled files or PATH. Reservations SHALL remain held until
asynchronous native cleanup completes. Client construction and factory defaults
SHALL follow client-connection-delegation.

#### Scenario: Inject a client factory

- **WHEN** a caller supplies a client factory to Store
- **THEN** the collection and its services use that factory's selected client
- **AND** reservations are released only after owned cleanup finishes

### Requirement: Typed application state workflows

App SHALL expose typed ColStateDump, ColStateImport, ColStateUpdate and ColStateReplace workflows with per-operation CollectionOptions as defined by collection-source-api. It SHALL load import configurations, construct collection-bound services and delegate managed publication through collection.Store. It SHALL NOT expose WithCollection or raw connection callbacks. CLI flag parsing SHALL remain outside App. CLI code SHALL own argument parsing, state serialization and output while preserving existing command syntax and formats.

#### Scenario: Manage state through App

- **WHEN** a Go caller imports, updates, replaces or dumps state on a database or collection package
- **THEN** App manages the complete scope and native publication internally while preserving namespace isolation and input objects

#### Scenario: Reject an invalid state workflow

- **WHEN** a workflow is cancelled or a batch fails validation
- **THEN** it returns an error, preserves validated state-operation semantics and releases the source reservation after cleanup; archive publication follows WriteBack
