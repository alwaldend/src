## Context

See proposal.md for motivation. The application already supports borrowed open
connections, but the CLI selects paths and owns their collection lifecycle.
The shared gRPC client provides separate stdio, network and collection owners.
Build delegates its staged archive to the same application apply operation.

## Goals / Non-Goals

Use one concrete collection options value and inject one application into the CLI.
Use existing transport and collection owners rather than adding native RPC or
file-transfer behavior to CLI handlers. TLS, authentication and remote-only
filesystems are outside this change.

## Decisions

- Proceed with managed path operations on the external server. Reusing the
  borrowed-open-connection API would require callers to pre-open a collection
  and would not cover build staging; transport selection belongs below App.
- Define one `ankiascodeapi.CollectionOptions` with `Path`, `Address`, `Connection`
  and `StateKey` for each collection operation. Remove `CollectionSource` and its
  constructors. Managed calls require a path and optionally an external address;
  borrowed calls require a connection and reject simultaneous path/address fields.
  Empty addresses and namespaces retain embedded/default behavior. Build requires
  a path and copies options with its staged path before apply.
- `APIRequest.Collection` uses the same options value instead of a source and a
  separate Backend field. Native API access retains exactly one managed path,
  external address, or borrowed connection. Namespace settings apply to Col
  workflows; raw RPC payloads continue to control native semantics.
- Construct App once in the executable and inject it into CLI. App retains no
  operation address, so a single instance can alternate external servers,
  embedded operations and borrowed connections. Borrowed connections retain
  their selected transport and lifetime. One store retains
  source reservations across all operations.
- Use `api <service> <method> --address` for the already-merged native RPC command,
  replacing its `--backend` flag spelling. Preserve its existing `APIRequest`
  validation and ownership: `--collection` uses an embedded managed collection;
  `--address` connects to the server's current collection and closes only the
  client. Explicit native collection lifecycle RPCs remain caller-controlled.
- Preserve master's archive publication during managed cleanup, including
  unchanged operations and failures. Test semantic convergence and unchanged
  build bases rather than archive byte identity after a managed operation.
- Reuse the existing network client's readiness and cancellation behavior.
  A failed external connection never falls back to a worker.
- The executable's native-hosting adapter ignores SIGPIPE. A direct probe
  against the linked listener failed after one disconnected client with the
  default policy (signal 13); inheriting SIGPIPE ignore permitted 30 sequential
  clients. Rust's library entry point does not install its executable policy.
  Keep signal handling with the executable rather than in App or the Go client.

## Risks / Trade-offs

- Server paths refer to its filesystem, while preparation and build staging use
  the CLI filesystem. Require shared absolute paths, including CLI TMPDIR.
- The server has one native collection shared across clients. Require exclusive
  use during each command; closing the client does not stop the server process.
- Plaintext transport has the same boundary as the existing server and client.
  Document it without adding an unrelated security configuration surface.
- A network endpoint may remain unavailable. Startup respects command
  cancellation; verify interruption through an unreachable endpoint.

## Migration Plan

Existing commands retain the embedded default. Opt in per command with
`--address`; omit it to return to embedded operation. No persistent host changes
or data migration are required.

## Collection options decision

Proceed with the requested consolidation. The source interface already required
runtime validation of its path/connection alternatives, while the separate
options value split path and address selection and silently ignored addresses
on borrowed connections. A concrete value reduces parameters and removes that
ambiguity. The main risk is invalid combinations; App validates them before
opening resources, including typed nil connections. Keeping separate source and
transport settings would preserve the overlap the user explicitly asked to fix.
Native API remote collection lifetime and CLI access modes remain unchanged.
