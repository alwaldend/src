# Client connection platforms Specification

## Purpose

Allow callers to construct managed Anki clients from addresses while choosing optional native dependencies through explicit imports.

## Requirements

### Requirement: Optional native dependencies require opt-in imports

The generic client API SHALL remain usable without cgo even when native support is enabled elsewhere in the build. Native adapters SHALL have separate subpackages.

#### Scenario: Pure caller in a native-enabled build

- **WHEN** a consumer builds only the generic client or application library with native support enabled
- **THEN** its dependency graph excludes the native engine and it can build with pure Go enabled

#### Scenario: Executable composes its worker

- **WHEN** an executable imports native hosting and explicitly registers its packaged worker constructor
- **THEN** the shared NewClient constructor accepts that worker scheme while generic library imports remain independent of native hosting

### Requirement: Addresses select managed clients

The client library SHALL accept a connection string as its constructor input, preserve TCP and Unix connections and select only explicitly registered worker schemes. Explicit unavailable engines SHALL fail without fallback. The resulting lifecycle owner SHALL expose the existing gRPC connection abstraction and SHALL close only resources it owns.

#### Scenario: Explicit subprocess command

- **WHEN** the address includes a worker executable and repeated argument parameters
- **THEN** startup invokes that exact executable and argument vector without shell evaluation

#### Scenario: External server survives client shutdown

- **WHEN** a client for a running TCP or Unix server is stopped
- **THEN** only that client connection closes and the server remains available

#### Scenario: Explicit empty worker arguments

- **WHEN** a caller supplies an explicitly empty worker argument vector
- **THEN** startup preserves zero arguments instead of selecting the address default command

#### Scenario: Unavailable native engine

- **WHEN** a generic caller requests an explicitly native-qualified address
- **THEN** construction reports that the engine requires an opt-in package

### Requirement: Native integration has one package owner

The optional ankigrpc/native package SHALL own the native C ABI bridge and server hosting. The executable SHALL register its subprocess schemes because it owns worker entrypoints. Executables SHALL own optional imports and worker entrypoint registration while the generic library owns connection construction. The constructor registry SHALL resolve qualified scheme availability before launch. Independent workers SHALL receive exactly the encoded argument vector. Executables MAY adapt their own packaged worker arguments for explicit engine selection as specified by the wasm-worker capability.

#### Scenario: Qualified constructor preserves worker arguments

- **WHEN** a caller requests subprocess+cgo for an independent executable accepting native server arguments
- **THEN** Go selects the registered factory and launches the configured worker with exactly its encoded arguments, rejecting unavailable schemes without fallback

#### Scenario: Scoped IPv6 endpoint

- **WHEN** an external TCP address includes an IPv6 zone identifier in bare, raw TCP URL or escaped TCP URL form
- **THEN** Address construction preserves that zone, including question marks and hashes inside brackets, and the client can connect to the selected server while URL suffixes outside the brackets remain invalid

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

### Requirement: Raw client callback scope

The client library SHALL provide a raw addressed-client scope using an injected
factory. It SHALL start the selected client, lend its prepared connection and
attempt cleanup after startup failure, callback error, cancellation or panic
unwinding. Cleanup SHALL use a bounded context independent of caller cancellation.
Operation and cleanup error identities SHALL remain discoverable; cancellation
checks SHALL preserve standard and custom cancellation causes. This scope SHALL
NOT prepare collections, reserve source paths or own archive publication.

#### Scenario: Failed scoped operation

- **WHEN** a scoped callback returns an error and stopping its client also fails
- **THEN** both error causes remain discoverable after bounded cleanup

#### Scenario: Cancellation or panic

- **WHEN** the callback cancels its context or panics
- **THEN** independent bounded cleanup still runs and the cancellation cause or
  original panic remains observable
