## Scope and dependency boundary

App coordinates operations and validates sources. It receives one ClientFactory dependency from the executable and supplies the same factory to raw client scopes and collection Store. Factory default selection and worker configuration are not App or Store settings. The paired backend change owns generic adapters, constructor registration and connection-string semantics.

## Raw clients and managed collections

App prepares RPCs and delegates address-only access to ankigrpc.WithClient(ctx, factory, address, callback). That library helper constructs, starts and lends the client connection, then performs bounded cleanup independently of request cancellation. It reserves or opens no collection. ClientFactory remains a construction contract; injected objects need no lifecycle methods.

Store owns WithCollection, source alias reservations, collection preparation and archive publication. It constructs clients through the injected factory, using an explicit address or passing the empty operation address unchanged for factory default selection. Without a configured factory, path-only access fails before client startup; explicit registered endpoints remain usable. Collection publication/close and transport shutdown retain separate budgets, and operation and cleanup errors remain discoverable. Store and application scope callbacks return only an error: WriteBack controls publication independently of any domain change result. Domain methods and the fixture export predicate retain meaningful changed flags.

## Executable composition

The executable owns subprocess schemes, its entrypoint and server arguments. It registers worker constructors explicitly through the common RegisterClientFactory API and generic ClientFactoryFunc[T] adapter. Build-selected command sources choose available worker schemes and native hosting. The configured factory selects its same-executable worker URL only for empty operation input. Neither generic nor native library imports advertise a subprocess entrypoint. Fixture owners construct their declared default workers directly or explicitly register packaged worker schemes.

The executable also injects ServerRunner and owns process-global signals. App and Store import no optional native implementation. Configuration, state, resource, export, planning, reconciliation and rebuilding keep their existing domain responsibilities. This feature extends the Client API and optional build composition merged through PR #186.

## Direct native sources

Path selects managed collection preparation and cleanup; Address selects the client regardless of transport. Network servers must share the managed paths and be reserved by their caller for that scope. A managed API path opens and closes its collection on the selected server without stopping that server. Address-only API access leaves collection lifetime with the caller. Borrowed Connection remains exclusive of Path and Address. No backend ownership method, flag or address classifier is needed. Direct clients still own fresh native state.

The configured Factory snapshots opt-in registrations, owns private worker registrations and selects its default without per-caller wrapper functions.

Store uses ankiclient.AsyncClient and retains source reservations until asynchronous native cleanup signals Done, including when bounded Stop returns while work is active. Existing synchronous subprocess reaping remains unchanged. Merged Client naming, pure imports and optional executable build composition are preserved; native execution is owned by add-in-process-backend.
