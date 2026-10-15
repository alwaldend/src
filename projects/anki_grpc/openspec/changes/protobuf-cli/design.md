## Context

The server target is visible throughout the repository for any native consumer.
The server maps original protobuf service paths to Anki's native dispatcher.
The upstream proto library provides the complete schema; Cobra, protobuf, gRPC,
and a managed Go backend are already pinned. See proposal.md for scope.

## Decisions

- A project-local Bazel rule passes the upstream descriptor set to a Go generator.
  It emits formatted Go service/method registrations and embeds the schema in a
  build output. No generated files or copied schemas enter Git.
- Fully qualified service names and exact method names avoid ambiguous renaming.
  All declared unary methods are exposed; unsupported streaming fails generation.
  Implicit backend aliases need no duplicate commands: original service routes
  already expose those operations.
- A shared runner parses ProtoJSON into dynamic protobuf messages and invokes
  the selected full method path. Requests use --request with JSON, @file, or -
  for stdin; responses are ProtoJSON. No custom per-field conversion is needed.
- Default transport is a managed stdio worker. --collection opens a local
  collection for a single RPC. Known collection-consuming operations are closed
  by Anki; other one-shot operations close before stopping the backend.
- --target connects to an existing plaintext TCP/Unix server. The CLI owns only
  that connection; collection and server lifetime remain with the caller.
  --target and local collection/worker options cannot be mixed.
- Timeouts cover startup and the call. Cleanup gets a fresh bounded context and
  uses the backend's process ownership. Cancellation does not promise rollback.

### Failure cases before implementation

- A schema service/method is missing from help or its RPC path/type is wrong.
- Help starts a worker; malformed or unknown JSON is sent to Anki.
- Inline, file, and stdin inputs encode different requests; response JSON is invalid.
- Native updates do not persist across one-shot invocations.
- Native package import/export, media, or lifecycle methods cannot be called.
- Native failures, transport errors, invalid flags, or deadlines return success.
- Managed worker files are removed or process exits before cleanup; a network
  invocation stops the shared server or closes a caller-owned collection.
- New dependency or generated source blobs enter the repository.
- Long isolated-delivery scratch paths exceed the Unix socket address limit;
  use a relative fixture socket and retain startup diagnostics.

## Risks / Trade-offs

These are low-level Anki operations and may be destructive as their API specifies.
Network transports are plaintext, matching the server. Large byte fields use
ProtoJSON base64. Local collection/media paths refer to the server filesystem.
The pinned API is unary; generation fails loudly if that changes.
