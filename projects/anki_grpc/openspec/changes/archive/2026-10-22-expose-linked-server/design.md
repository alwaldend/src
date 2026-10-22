## Context

The Rust command owns a small opaque gRPC codec, command dispatch and transport selection. Its native calls already use Tokio blocking tasks. Existing pinned rules_rust supports static libraries with CcInfo; rules_go accepts CcInfo through cgo dependencies. See proposal.md for motivation.

## Goals / Non-Goals

Share one Rust implementation between the standalone executable and a linked Go entry point. Keep Rust responsible for its runtime, transport and native work. No protobuf binding changes or new dependencies.

## Decisions

- Move the server, parser and descriptor command generator to pkg/ankiserver. Keep the standalone Rust main as a thin argument/exit-code adapter.
- Export one argc/argv C ABI function, copy incoming strings into Rust and return an exit status. The Go wrapper owns argument allocations for the call. Rust panics continue to abort the hosting server process.
- Name the Go host package ankiserver and the connection/process/collection owner ankiclient, updating import paths and consumers together. Preserve its existing exported lifecycle types. Require callers to declare their own server executable runtime inputs; the client library supplies no standalone worker.
- Link a standard Rust static library through cgo using CcInfo. The native backend client library remains independent of this linked server library.
- Keep native command dispatch in Rust: bare invocation and root help list serve, serve selects the existing transports, and serve --help documents their flags. Go forwards the complete native argument vector through its generic server run command. No new dependencies are needed.
- Proceed with child-process hosting: it preserves forced termination and reaping for native blocking work. In-process socket-pair hosting saves another Go runtime but introduces new shutdown and uninterruptible-work ownership concerns. The current binary/extraction approach is simpler to build but cannot meet the no-extraction requirement.

## Risks / Trade-offs

- Mixed-language link flags and native system libraries: validate the pinned rules and optimized installed artifact, not only compilation.
- Foreign argument lifetimes: copy to Rust-owned strings and expose no retained Go pointers or protobuf objects.
- Process-aborting native failures: host the linked server in the existing isolated worker process for offline operations.
