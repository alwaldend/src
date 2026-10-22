# in-process-selection Specification

## Purpose

Allow application and command-line workflows to select owned native execution through the existing backend address option.

## Requirements

### Requirement: Managed direct execution

Collection workflows SHALL accept `inprocess://<name>`, `inprocess+cgo://<name>`
and configured `subprocess://` addresses alongside their managed database or archive
path when the selected engine is available. They SHALL preserve collection
preparation, namespace selection, archive publication and cleanup. In-process
execution SHALL NOT require a worker executable. Each operation SHALL own a fresh
client; in-process clients SHALL own fresh native state. Operations SHALL NOT retain
an earlier address. Empty operation addresses SHALL delegate to the
executable-configured factory default.

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

### Requirement: Scope callbacks report errors only

Collection and application scope callbacks SHALL return only an error. Archive publication SHALL depend on WriteBack rather than a callback change flag or its success. Domain operations SHALL retain meaningful change results for their own consumers.

#### Scenario: Committed write followed by callback failure

- **WHEN** a collection callback writes native state and then returns an error with WriteBack enabled
- **THEN** cleanup preserves the error and publishes committed state without requiring a callback change result
