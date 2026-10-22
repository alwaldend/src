## Context

The store reserves paths and decides whether to export its source archive. Managers currently bracket a fresh worker and a borrowed collection backend. Only tests inject a manager.

## Goals / Non-Goals

Keep scope behavior and errors while removing callback manager abstractions. Preserve backend storage and transport contracts. Pooling, new transports and changes to restore retention are outside this refactor.

## Decisions

The store constructs BackendSTD and then BackendCollection. A private application helper brackets each Start/Stop pair and uses a five-second cleanup context derived with context.WithoutCancel. Nested scopes close the collection before reaping its worker. Source reservations outlive both scopes.

Private BackendSTDOptions permit selecting a fixture executable and arguments; an empty binary keeps bundled-server/PATH discovery. No public factory, manager, alternate backend protocol or test-only exported API is added.

Write failure coverage before changing production code. Replace manager spies with a PID-recording fixture that executes the real native worker. Verify success, callback error, custom cancellation, panic, crash plus callback error, failed startup and source reuse. Existing native tests continue checking persistence, publication and expired connections.

## Risks / Trade-offs

- Cleanup order or cancellation regressions: retain native E2E with process reaping and temporary-restore observations.
- Removed manager injection: the only production selection is unchanged; future transport requirements can introduce a concrete constructor when needed.
- Native uncertainty: preserve existing retained-restore behavior and error identities rather than deleting uncertain data.
