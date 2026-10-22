## Purpose

Allow application and command-line workflows to select owned native execution through the existing backend address option.

## ADDED Requirements

### Requirement: Managed direct execution

Collection workflows SHALL accept `inprocess://<name>` and configured `subprocess://` addresses alongside their managed database or archive path. They SHALL preserve collection preparation, namespace selection, archive publication and cleanup. Direct execution SHALL NOT require a worker executable. Each operation SHALL own a fresh backend and SHALL NOT retain an earlier address. Empty operation addresses SHALL delegate to the executable-configured factory default.

#### Scenario: Direct collection workflow

- **WHEN** export, build, plan, apply and state commands select an in-process address
- **THEN** they operate on their managed sources with the existing persistence and publication policies

#### Scenario: Optional native composition

- **WHEN** the executable enables in-process collection access
- **THEN** the anki_grpc library constructs the compiled adapter from the connection string, through the shared constructor and an opt-in native package import, without application-owned adapter construction or engine registries

### Requirement: Direct native API selection

Native API calls SHALL accept any client address alone or alongside a managed collection path. Path SHALL select managed collection lifetime and Address SHALL select its transport. Managed paths SHALL retain native collection-lifetime restrictions. Borrowed connections SHALL reject every additional source. CLI help SHALL describe the direct selectors.

#### Scenario: Managed native API

- **WHEN** a native API call supplies a database and an in-process address
- **THEN** it prepares that database, executes the call and closes its owned resources

#### Scenario: Backend-only native API

- **WHEN** a native formatting operation supplies only an in-process address
- **THEN** it succeeds without opening a collection or starting a worker

#### Scenario: Managed network source

- **WHEN** a native API call supplies a database and a caller-reserved TCP or Unix server address with shared paths
- **THEN** it opens that collection on the selected server, executes the call and closes the collection while leaving the server available

### Requirement: Client scope and default ownership

The client library SHALL own startup and bounded cleanup for raw addressed API access. The collection store SHALL retain managed collection scopes, publication and reservations. The application SHALL supply its injected factory to those scope owners and SHALL NOT duplicate client startup or default worker settings. Injected factories SHALL receive each operation's address unchanged, including an empty address for configured defaults. Explicit addresses SHALL override the factory's default without changing later operations.

#### Scenario: Configured default after explicit API access

- **WHEN** an application uses explicit TCP, Unix or subprocess addresses for API and collection calls and then performs a path-only operation
- **THEN** the same injected factory receives each requested address, selects its default for the empty input, and collection persistence and client cleanup remain intact

#### Scenario: Injected factory selection

- **WHEN** an injected factory constructs a client for an unregistered or network-looking address alongside a managed API path
- **THEN** the factory is called once with unchanged input and its client executes the managed call without global classification

### Requirement: Scope callbacks report errors only

Collection and application scope callbacks SHALL return only an error. Archive publication SHALL depend on WriteBack rather than a callback change flag or its success. Domain operations SHALL retain meaningful change results for their own consumers.

#### Scenario: Committed write followed by callback failure

- **WHEN** a collection callback writes native state and then returns an error with WriteBack enabled
- **THEN** cleanup preserves the error and publishes committed state without requiring a callback change result
