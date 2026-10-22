## Context and scope

The pure ankigrpc package owns connection-string routing and reusable client construction. The executable owns worker entrypoints, supported subprocess schemes and defaults. Anki owns storage semantics; Go owns explicit client and collection lifetimes. Both PRs target master independently.

## Factory contract and concrete results

ClientFactory exposes NewClient(address) returning the lifecycle interface for runtime-selected clients. ClientFactoryFunc[T] adapts functions returning concrete client types to that contract, preserves error causes and rejects nil functions across instantiations. Fixed subprocess/network constructors return concrete pointers. Generic parameters do not replace runtime scheme dispatch or grpc.ClientConnInterface's generated-client contract.

Only TCP and Unix are built in. Executables register subprocess and supported qualified worker schemes through RegisterClientFactory using NewSubprocessClient. The library constructor requires an explicit executable URL; empty library addresses fail. The injected executable factory handles empty operation addresses using its explicit worker URL. Optional native imports advertise no worker entrypoint. Registry locks cover lookup/registration only, and duplicates or invalid registrations fail without replacing an owner.

## Address and scope behavior

Normalize addresses before one registry lookup. Preserve bare, raw and escaped scoped IPv6 zones, including reserved characters inside brackets, while rejecting suffix delimiters outside brackets. Share URL decoding with network construction but keep shorthand inference in routing. Subprocess URLs preserve their exact argv, including zero arguments, without Go platform flags in Rust.

ankigrpc.WithClient(ctx, factory, address, callback) owns raw addressed client startup, connection borrowing and bounded independent cleanup. The factory remains a construction dependency; object implementations need no lifecycle methods. Application collection Store owns managed path preparation, archive publication and source reservations, outside the raw scope helper. Public libraries remain pure Go.

## Native hosting and naming

ClientSTD, ClientGRPC and ClientCollection retain the existing Start/Connection/Stop contract. Remove Backend compatibility aliases. The native package owns the C ABI CLI bridge and hosting; build-selected command files choose imports, worker schemes and runServer. Keep in-process command handles and direct execution in independent PR #184. No new dependencies or WASM runtime are introduced.
