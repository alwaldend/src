## Purpose

Provide transport startup, RPC access and shutdown without owning application
collection storage.

## ADDED Requirements

### Requirement: Expose the transport client lifecycle

Client SHALL expose only Connection and Stop. StartableClient SHALL embed Client
and add Start for callers that own startup. Factory adapters SHALL return their
constructed StartableClient directly without another lifecycle wrapper.
Concrete clients SHALL own readiness and cleanup of their transport resources.
Collection open, close, archive publication and source reservations SHALL remain
outside the generic client contract.

#### Scenario: Reuse a transport with explicit native collection calls

- **WHEN** a caller starts a client, opens and closes collections using generated
  native RPCs, and then stops the client
- **THEN** transport lifetime follows Start, Connection and Stop without managed
  collection phases or an additional client wrapper

### Requirement: Observe optional native completion

Concrete owned transports MAY expose final native completion through Done when
Stop can return before owned work finishes. That capability SHALL remain optional
and SHALL NOT expand Client or StartableClient or require a public AsyncClient interface.
Closing a raw network connection SHALL NOT claim to
terminate its caller-owned server or collection.

#### Scenario: Stop waiting for native cleanup

- **WHEN** native cleanup continues after a Stop caller's context expires
- **THEN** a caller needing final completion can observe the concrete Done method
  while ordinary transport consumers retain the three-method StartableClient contract

### Requirement: Destroy owned native resources outside asynchronous serving

The native server SHALL retain synchronous ownership of its initialized backend
until asynchronous serving and its runtime have finished. Backend resources that
own runtimes SHALL be destroyed outside an asynchronous execution context.

#### Scenario: Successful media operation followed by normal stdio shutdown

- **WHEN** a stdio client successfully adds media, closes its collection and stops
- **THEN** the worker exits successfully without a runtime-destruction panic
