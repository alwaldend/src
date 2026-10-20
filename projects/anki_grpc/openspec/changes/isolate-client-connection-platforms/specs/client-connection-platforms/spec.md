## Purpose

Allow callers to construct managed Anki clients from addresses while choosing optional native dependencies through explicit imports.

## ADDED Requirements

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

The optional ankigrpc/native package SHALL own the native C ABI bridge and server hosting. The executable SHALL register its subprocess schemes because it owns worker entrypoints. Command build composition SHALL select optional imports without constructing clients or interpreting addresses. The Go constructor registry SHALL resolve qualified scheme availability before launch. The worker SHALL receive exactly the encoded argument vector without injected platform flags.

#### Scenario: Qualified constructor preserves worker arguments

- **WHEN** a caller requests subprocess+cgo for an executable accepting native server arguments
- **THEN** Go selects the registered factory and launches the configured worker with exactly its encoded arguments, rejecting unavailable schemes without fallback

#### Scenario: Scoped IPv6 endpoint

- **WHEN** an external TCP address includes an IPv6 zone identifier in bare, raw TCP URL or escaped TCP URL form
- **THEN** URL normalization preserves that zone, including question marks and hashes inside brackets, and the client can connect to the selected server while URL suffixes outside the brackets remain invalid

### Requirement: Executable-owned worker registration

The library SHALL supply unregistered subprocess construction requiring an explicit executable. Only an executable or fixture owning its worker entrypoint SHALL register subprocess schemes and choose default worker arguments. Generic adapter ClientFactoryFunc[T] SHALL accept concrete client constructors while implementing the shared runtime factory contract.

#### Scenario: Explicit owner registration

- **WHEN** a caller imports generic and native libraries without registering a subprocess scheme
- **THEN** a subprocess address and empty library address fail rather than discovering or assuming a worker entrypoint

#### Scenario: Registered worker arguments

- **WHEN** the executable registers its worker constructor and supplies an explicit executable URL
- **THEN** generated Anki calls succeed through that worker and zero argument values remain zero
