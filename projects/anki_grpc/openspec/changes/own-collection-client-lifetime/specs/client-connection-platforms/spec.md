## MODIFIED Requirements

### Requirement: Executable-owned worker registration

The library SHALL supply unregistered subprocess construction requiring an
explicit executable. Only an executable or fixture owning its worker entrypoint
SHALL register subprocess schemes and choose default worker arguments.
ClientFactory SHALL be a func(string) (ankiclient.StartableClient, error).
ToClientFactory[T] SHALL adapt concrete client constructors to this function
type with contextual error wrapping and SHALL return nil for a nil constructor.
Registration SHALL reject nil constructor functions.

#### Scenario: Explicit owner registration

- **WHEN** a caller imports generic and native libraries without registering a subprocess scheme
- **THEN** a subprocess address and empty library address fail rather than discovering or assuming a worker entrypoint

#### Scenario: Registered worker arguments

- **WHEN** the executable registers its worker constructor and supplies an explicit executable URL
- **THEN** generated Anki calls succeed through that worker and zero argument values remain zero

#### Scenario: Nil typed constructor

- **WHEN** a caller adapts a nil typed constructor
- **THEN** ToClientFactory returns nil
- **AND** attempting to register that result panics

### Requirement: Configured factory isolation

Concrete Factory values SHALL snapshot package registrations and accept an
explicit default address. Their registrations SHALL remain independent of other
factories and package defaults. Explicit operation addresses SHALL override the
default for only that construction. Routing SHALL preserve the original selected
address for one constructor-owned parse. Factory.NewClient SHALL be injectable as
a ClientFactory method value. Returned transports SHALL implement StartableClient;
resource ownership SHALL remain encapsulated by their implementations.

#### Scenario: Private worker defaults

- **WHEN** two factories register the same private worker scheme and alternate default and explicit native calls
- **THEN** each invokes its own constructor, performs real Anki calls and retains its configured default without changing package registrations

#### Scenario: Relative Unix socket path

- **WHEN** a server advertises a relative Unix socket path containing spaces in raw or escaped form
- **THEN** a configured client connects to the advertised socket and persists collection state

### Requirement: Raw client callback scope

The client library SHALL provide a raw addressed-client scope using an injected
constructor function. WithClient SHALL select the package NewClient constructor
when that function is nil. The scope SHALL start the selected client, lend its
prepared connection and attempt cleanup after startup failure, callback error,
cancellation or panic unwinding. Cleanup SHALL use a bounded context independent
of caller cancellation. Operation and cleanup error identities SHALL remain
discoverable; cancellation checks SHALL preserve standard and custom cancellation
causes. This scope SHALL NOT prepare collections, reserve source paths or own
archive publication.

#### Scenario: Failed scoped operation

- **WHEN** a scoped callback returns an error and stopping its client also fails
- **THEN** both error causes remain discoverable after bounded cleanup

#### Scenario: Cancellation or panic

- **WHEN** the callback cancels its context or panics
- **THEN** independent bounded cleanup still runs and the cancellation cause or original panic remains observable

#### Scenario: Default constructor

- **WHEN** a caller invokes WithClient without an injected constructor
- **THEN** the package NewClient constructor resolves the supplied address
