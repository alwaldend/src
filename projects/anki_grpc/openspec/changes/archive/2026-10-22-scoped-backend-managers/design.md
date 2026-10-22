## Context

BackendSTD owns one native process and transport. BackendCollection borrows a started Backend and owns native open/close and restore files. The user requests a separate additive PR for scoped managers, keeping family-first names and transport replacement independent of collection semantics.

## Goals / Non-Goals

**Goals:** Add Manager, ManagerSTD and ManagerCollection in pkg/backend; retain one owner for each scope's startup and shutdown; keep generated clients and error identities.

**Non-Goals:** Rename existing public backend types, add a TCP/Unix client implementation, change other PRs, pool idle processes, duplicate source reservations, or move publication/resource/state/config policy into anki_grpc.

## Decisions

withBackend(ctx, owned, cleanupTimeout, callback) owns a supplied unstarted backend for one scope. Managers supply the owned backend and cleanup timeout. The helper rejects a missing callback or already-cancelled context before startup, starts the backend, lends it to the callback and defers Stop once startup is attempted. Cleanup uses context.WithoutCancel and the explicit timeout to retain values independently of cancellation. Callback and cleanup errors remain discoverable and the original panic propagates after cleanup. This helper preserves existing uncertain-native-completion behavior; it does not add transport-termination finalization or change panic diagnostics.

Manager exposes WithBackend(context.Context, func(Backend) error) error. NewManagerSTD accepts existing BackendSTDOptions and copies its arguments. Each call creates a new BackendSTD and delegates to the shared private withBackend function using the backend's normalized shutdown timeout.

NewManagerCollection accepts Manager. Its WithCollection(ctx, path, callback) method constructs BackendCollection inside the delegated transport scope and calls the shared private withBackend function with a five-second cleanup timeout. This closes the collection before returning to the transport manager. Only this layer accepts collection paths; it need not implement the transport Manager interface. Invalid configuration and callback inputs are rejected before delegation.

Successful callbacks still report context cancellation observed before scope completion. errors.Join preserves callback and cleanup failures; propagated errors retain contextual wrapping. Panics propagate after deferred cleanup. The managers retain only immutable configuration and delegate native readiness, restoration, consuming operations, uncertain completion and reaping to the existing owners. Native uncertainty continues to retain restore files; a generic transport manager cannot assume that disconnecting an externally owned server stops its native work.

These managers provide scoped resource ownership. A factory describes only creation, and a store or pool would imply additional retention/retrieval responsibilities that are not needed here. Keep source reservations and archive publication in their existing anki_as_code owner. The new transport Manager is the deliberately reusable contract for alternative implementations.

## Risks / Trade-offs

- Callback borrowers could call lifecycle methods themselves; document that managers retain lifecycle ownership and borrowers finish all operations before returning.
- Cancellation does not roll back native database writes; preserve existing semantics and bounded cleanup.
- A custom transport manager must honor the scope contract and lend an exclusive native backend. Test composition by wrapping a real native manager and probing the transport after collection cleanup.
- The common scope helper removes duplicated lifecycle handling while leaving transport creation and collection composition with their managers.

## Verification

Before implementation, add native E2E cases for process reaping, argument copying, startup rejection, callback failure/cancellation/panic, combined crash and callback errors, delegated cleanup order, native persistence, archives and restore cleanup. Reuse the existing native executable/PID fixture and emit closed database/archive and JSON artifacts. Then run focused backend tests, existing anki_as_code consumers, repository quality, semantic lint and integrated documentation gates.

## Review correction: custom cancellation causes

Manager cancellation guards join ctx.Err with context.Cause so errors.Is can
identify both the standard cancellation sentinel and a caller-supplied cause.
Extend existing native pre-start and callback-cancellation cases with
WithCancelCause before changing the guards; retain the existing cleanup order
and bounded independent contexts.
