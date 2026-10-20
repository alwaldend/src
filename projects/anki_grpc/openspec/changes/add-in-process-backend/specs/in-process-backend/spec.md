## Purpose

Allow generated Anki clients to execute native operations inside their Go process with explicit ownership and the established native error contract.

## ADDED Requirements

### Requirement: Direct native execution

The backend SHALL implement the existing start, connection and stop contract and execute upstream unary commands through a C ABI without starting a process or transport listener. Generated clients SHALL retain native persistence, concurrent access, original service routes and serialized native errors. Unsupported streaming and call options SHALL fail explicitly.

#### Scenario: Native collection workflow

- **WHEN** generated clients create a note, close and reopen its database through the direct backend
- **THEN** the note and generated cards persist and native failures retain their status and error trailers

### Requirement: Native lifetime

The backend SHALL permit one start attempt and reject calls once stopping begins. Cancellation SHALL stop waiting without freeing resources still used by a native call. Stop SHALL release native state after active calls complete; its deadline SHALL bound waiting without pretending unfinished work was terminated. Repeated Stop calls SHALL be safe.

#### Scenario: Cancellation during native work

- **WHEN** a native media request remains active after its caller cancels and Stop times out
- **THEN** retained clients reject calls, native resources remain valid, and later Stop completes after the native work finishes

#### Scenario: Shutdown overlaps native initialization

- **WHEN** synchronous native startup is still running and Stop reaches its deadline
- **THEN** Stop returns the deadline error, Done waits for startup cleanup, and startup does not expose a connection after shutdown

### Requirement: In-process address selection

An address with the `inprocess://` prefix SHALL select a fresh direct backend. Subprocess addresses SHALL encode worker configuration and preserve their exact arg values without injected flags. Zero arg values SHALL mean zero arguments. An empty library address SHALL fail because executable defaults belong to the calling executable. In-process suffixes SHALL be opaque labels and SHALL NOT create shared state or network resources. Other addresses SHALL retain existing network behavior.

#### Scenario: Independent named instances

- **WHEN** callers start two backends selected with the same in-process address
- **THEN** each has independent collection state and stopping either leaves the other usable

### Requirement: Explicit optional native imports

The generic client and application libraries SHALL compile without cgo regardless of native support in the executable. Direct execution SHALL require opting into the native subpackage through an explicit import that registers only its direct client constructors. Explicit unavailable engines SHALL fail without fallback.

#### Scenario: Generic caller without native support

- **WHEN** a caller imports only the generic constructor and requests in-process execution
- **THEN** construction fails clearly and the caller can build with pure Go enabled

#### Scenario: Registered client strategy

- **WHEN** an adapter registers a client factory under a custom scheme
- **THEN** managed collection scopes use its client without transport classification or backend ownership methods

#### Scenario: Scoped IPv6 endpoint

- **WHEN** an external TCP address includes an IPv6 zone identifier in bare, raw TCP URL or escaped TCP URL form
- **THEN** construction preserves that zone, including question marks and hashes inside brackets, and the client can connect to the selected server while URL suffixes outside the brackets remain invalid

### Requirement: Executable-owned worker registration

The library SHALL supply unregistered subprocess construction requiring an explicit executable. Only an executable or fixture owning its worker entrypoint SHALL register subprocess schemes and choose default worker arguments. Generic adapter ClientFactoryFunc[T] SHALL accept concrete client constructors while implementing the shared runtime factory contract.

#### Scenario: Explicit owner registration

- **WHEN** a caller imports generic and native libraries without registering a subprocess scheme
- **THEN** a subprocess address and empty library address fail rather than discovering or assuming a worker entrypoint

#### Scenario: Registered worker arguments

- **WHEN** the executable registers its worker constructor and supplies an explicit executable URL
- **THEN** generated Anki calls succeed through that worker and zero argument values remain zero

### Requirement: Configured factory isolation

Configured factories SHALL snapshot package registrations and accept an explicit default address. Their registrations SHALL remain independent of other factories and package defaults. Explicit operation addresses SHALL override the default for only that construction. Routing SHALL preserve the original selected address for one constructor-owned parse. Client lifecycle SHALL remain Start, Connection and Stop; resource ownership SHALL be encapsulated by implementations.

#### Scenario: Private worker defaults

- **WHEN** two factories register the same private worker scheme and alternate default and explicit native calls
- **THEN** each invokes its own constructor, performs real Anki calls and retains its configured default without changing package registrations

#### Scenario: Relative Unix socket path

- **WHEN** a server advertises a relative Unix socket path containing spaces in raw or escaped form
- **THEN** a configured client connects to the advertised socket and persists collection state
