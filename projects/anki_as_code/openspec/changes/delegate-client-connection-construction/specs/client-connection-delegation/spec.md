## Purpose

Keep application workflows independent of optional engines while delegating connection-string interpretation to the Anki client library.

## ADDED Requirements

### Requirement: Application delegates client construction

The application SHALL pass the requested address to a library-owned client constructor. Its reusable library SHALL remain pure Go independently of the executable's compiled engines. Existing collection publication and external-server lifetime guarantees SHALL remain unchanged.

#### Scenario: One application uses different endpoints

- **WHEN** successive operations select configured subprocess, TCP and Unix addresses
- **THEN** each operation uses its requested endpoint without retaining the previous selection

#### Scenario: Native-enabled executable with pure application library

- **WHEN** native support is enabled for the command
- **THEN** consumers of only the application library can still build without cgo

### Requirement: Client scope and default ownership

The client library SHALL own startup and bounded cleanup for raw addressed API access. The collection store SHALL retain managed collection scopes, publication and reservations. The application SHALL supply its injected factory to those scope owners and SHALL NOT duplicate client startup or default worker settings. Injected factories SHALL receive each operation's address unchanged, including an empty address for configured defaults. Explicit addresses SHALL override the factory's default without changing later operations.

#### Scenario: Configured default after explicit API access

- **WHEN** an application uses explicit TCP, Unix or subprocess addresses for API and collection calls and then performs a path-only operation
- **THEN** the same injected factory receives each requested address, selects its default for the empty input, and collection persistence and client cleanup remain intact
