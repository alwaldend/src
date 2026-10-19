## Context and scope

The pure ankigrpc package owns connection-string routing and reusable client construction. The executable owns worker entrypoints, supported subprocess schemes and defaults. Anki owns storage semantics; Go owns explicit client and collection lifetimes. This feature extends the Client naming and optional executable composition merged through PR #186.

## Factory contract and concrete results

ClientFactory exposes NewClient(address) returning the lifecycle interface for runtime-selected clients. ClientFactoryFunc[T] adapts functions returning concrete client types to that contract, preserves error causes and rejects nil functions across instantiations. Fixed subprocess/network constructors return concrete pointers. Generic parameters do not replace runtime scheme dispatch or grpc.ClientConnInterface's generated-client contract.

Only TCP and Unix are built in. Executables register subprocess and supported qualified worker schemes through RegisterClientFactory using NewSubprocessClient. The library constructor requires an explicit executable URL; empty library addresses fail. The injected executable factory handles empty operation addresses using its explicit worker URL. Optional native imports advertise no worker entrypoint. NewFactory snapshots package registrations and retains its own explicit default and private registrations. Registry locks cover lookup/registration only, and duplicates or invalid registrations fail without replacing an owner. Client exposes only Start, Connection and Stop; resource ownership stays inside each implementation. Neither clients nor registrations classify application access policies.

## Address and scope behavior

Identify the routing scheme without normalizing or serializing transport data. Constructors receive the selected original address unchanged and parse it once. NewGRPC(address) validates its address before startup and retains a decoded transport and endpoint. Preserve bare, raw and escaped scoped IPv6 zones, including reserved characters inside brackets, while rejecting suffix delimiters outside brackets. Decode explicit Unix URL paths once, including relative socket names containing spaces. Subprocess URLs preserve exact argv, including zero arguments, without Go platform flags in Rust.

ankigrpc.WithClient(ctx, factory, address, callback) owns raw addressed client startup, connection borrowing and bounded independent cleanup. The factory remains a construction dependency; object implementations need no lifecycle methods. Application collection Store owns managed path preparation, archive publication and source reservations, outside the raw scope helper. Public libraries remain pure Go.

## Direct native execution

Share Rust dispatch through explicit owned command handles and result buffers. Go borrows bytes only during the synchronous C call. The native Server owns its handle; Client owns its connection and cleanup. Importing native registers only inprocess and inprocess+cgo, as the requested exception to explicit initialization. Equal labels never share state.

Register native work before dispatch and initialize outside the lifecycle mutex. Cancellation stops waiting without freeing active buffers. Stop rejects new calls and drains native work; its deadline bounds waiting while the canonical ankiclient.AsyncClient Done capability signals final cleanup, including shutdown overlapping startup. The shared lifecycle contract records single-attempt startup, terminal repeatable Stop, retained connection rejection, non-interruptible native initialization and subprocess reaping. Collection reservations remain held until Done. Native panics retain the accepted abort policy. Preserve the merged Client API, pure libraries and optional executable build composition; no Backend compatibility aliases are added.
