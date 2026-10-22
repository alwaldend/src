## Why

The collection CLI needs typed Anki service clients and process lifetime
management. Keep the reusable Go connection owner beside the Rust server in
projects/anki_grpc, with upstream protobuf services supplying caller methods.

## What Changes

- Add projects/anki_grpc/pkg/backend for gRPC HTTP/2 over stdin/stdout,
  process startup, cancellation, shutdown, reaping, and crash reporting.
- Consume the merged projects/anki_grpc runtime target from master.
- Generate Go gRPC clients from upstream protobuf services, without copied
  definitions or handwritten wrappers for each operation.
- Expose the managed gRPC connection to caller-constructed service clients.
- Expose Backend with context-aware Start and Stop and an error-returning
  Connection; configure BackendSTD through
  BackendSTDOptions with a binary path, copied arguments, and configurable
  shutdown timeout, starting the process only when Start is called.
- Limit the Start context to readiness and use the Stop context to bound graceful
  shutdown while preserving process reaping and context errors.
- Adapt original service paths in Rust using Anki's descriptor helpers, while
  retaining the existing backend-prefixed paths.
- Keep the considered-approaches Mermaid diagram in the Anki CLI documentation.
- Verify the Go API with a synthetic collection and repeatable database artifact.
- Keep the native build and locked dependency versions unchanged; update only
  service routing and schema generation needed by the typed clients.

## Capabilities

### New Capabilities

- `anki-cargo-api`: Typed Go service clients, transport, and lifecycle for Anki.

### Modified Capabilities

None.

## Impact

The Anki gRPC project's reusable Go API and service routing, upstream Go client
generation, project docs, diagram, and runtime inputs. Server ownership belongs
to projects/anki_grpc; native build ownership stays in the third-party package.
Existing approved dependency versions are preserved.
Collection reconciliation remains separately owned by PR #128.
