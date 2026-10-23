## Context

The client factory already chooses transport by address. Collection storage and
source reservations previously added further Client implementations around that
transport. The user selected one client with collection methods and opening options.
The linked Anki as Code change consumes the same client from ColOpen.

## Goals / Non-Goals

Keep one lifecycle owner for transport and optional storage while retaining
native persistence, safe cleanup and reuse between collections. Application
source reservations and session-context policy remain outside the library.
No native protocol, archive format or dependency changes are needed.

## Decisions

Factory adapters return an ankiclient.Client owning one Transport. Transport is
the existing Start/Connection/Stop capability of process, network and in-process
adapters. Client adds OpenCollection, CloseCollection and Done. Collection state
is private data, not another Client. An already-owned Client returned by a custom
factory passes through unchanged.

Start prepares the transport once. OpenCollection records a storage phase before
preparing it so concurrent cleanup can see failed or interrupted opening. Its
CollectionOptions value owns Path, RestoreDirectory, WriteBack and CloseTimeout.
Connection returns gated collection access while a phase exists and raw access
between phases. A successful CloseCollection removes the phase, invalidates its
retained connections and keeps the transport available for another collection.
A failed close retains its diagnostic and storage state for terminal cleanup.
Keep private state and lifecycle methods limited to those needed for readiness,
active native calls, publication and one bounded cleanup attempt; remove redundant
client-era guards and indirection without changing those guarantees.

Stop initiates cleanup once. It closes the collection under CloseTimeout (default
ten minutes), then stops the transport under a separate five-second budget.
A caller's context bounds waiting, not owned cleanup. Joining an existing close
attempt must honor that joiner's context so an unbounded concurrent close cannot
prevent terminal Stop from proceeding to transport shutdown. Done observes final native
transport completion even after a cleanup error or cancelled wait. Explicit
CloseCollection uses the supplied context directly. Start's context controls
startup only.

The consuming Store starts the factory-created client, opens its collection and
returns it unchanged. A plain source-reservation record observes Done and calls
Stop on session cancellation. Reservation deletion compares the record identity;
new opens reap completed records so immediate source reuse does not depend on
observer scheduling. MCP retains that client for each explicit handle.

## Risks / Trade-offs

Cancellation does not establish native completion. Keep restored files and source
exclusion until completion can be established, and preserve cleanup errors.
Collection operations must remain exclusive with lifecycle changes while ordinary
RPCs retain concurrency. Existing E2E fixtures cover consuming native calls,
publication failures and uncertain work; the preimplementation unified-client
case verifies transport reuse, old-connection rejection and final completion.

## Migration Plan

Migrate all repository callers of separate collection constructors to methods on
their existing client. Keep raw transport fixtures on Transport where they test
transport-specific behavior. Update application opening options and remove its
reservation-client wrapper. Validate real native, WASM and MCP workflows before
publication. No storage migration or live deployment is required.
