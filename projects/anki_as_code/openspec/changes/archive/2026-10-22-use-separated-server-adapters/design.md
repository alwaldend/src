## Context

Trunk composes optional native engines and scoped client factories in the
executable. App owns orchestration and starts the same executable as its
worker. The linked add-wasm-server change owns the Rust command and host.

## Goals / Non-Goals

Integrate WASM without changing reconciliation. Preserve the default supervised
worker, native in-process clients and external-only builds, and offer an explicit
in-process WASM client using the existing Rust server.

## Decisions

Keep one command package and public target. Independent Bazel booleans
`//projects/anki_grpc:cgo` (default true) and `//projects/anki_grpc:wasm`
(default false) select optional adapter imports. Command composition chooses
native for the supervised worker when present, then WASM, and otherwise
reports unavailable local hosting. Cgo imports retain native in-process
registrations when both engines are enabled. The executable's scoped factory
registers every compiled engine's qualified subprocess scheme. Each constructor
prepends the host's `--server-engine` selector for same-executable workers,
recognizing filesystem aliases; independent executable arguments remain unchanged.
The shared subprocess parser owns URL validation and exposes options for this
composition. The shared Cobra root declares the selector and consumes it before
opaque Rust CLI arguments, so generated help includes the same flag. Conflicting
selectors in qualified same-executable addresses are rejected. Generic subprocess
arguments remain unchanged.
It registers `inprocess+wasm` when WASM is enabled and generic `inprocess`
when cgo is unavailable; otherwise native registration owns the generic scheme.

The WASM-only target sets `pure = "on"`. The adapter owns guest building and
embedding and exposes `Run(arguments)`; commands handle no guest bytes. Optional
imports live under `anki_grpc/pkg/ankigrpc/native` and `anki_grpc/pkg/ankigrpc/wasm`.
Keep the host OS/CPU platform; only the Rust guest targets WASI. Fixed native,
WASM-only, combined and external-only fixtures verify engine combinations
independently of public flags and keep native package reopening native.
Go build constraints select optional source files; Bazel variants supply their
matching tags, dependencies and cgo mode without repeating source selection.

Command composition owns process signals. App and collection libraries keep
injected client construction and their existing lifecycle boundaries. The
full Rust CLI owns argument parsing and gRPC; Go supplies WASI byte streams.

### Strategy documentation

One maintained App/Server diagram contains three groups by dependency rather than
enumerating complete build configurations. External TCP/Unix clients require
no optional local engine. WASM enables its supervised subprocess and explicit
in-process server, with generic in-process selection preferring cgo;
cgo enables native in-process calls and the native subprocess, including in
combined builds. Diagram titles show qualified schemes only, keeping the generic
aliases and their preference rules in the owning README. Cgo remains enabled
by default. README prose records the bundled
worker's engine requirement and the precedence when both engines are enabled.
Subprocess strategies appear only in the engine groups.
The Mermaid source belongs in `docs/`; site SVG and dark WebP targets render it
with the repository's existing theme. A local Mermaid canvas-padding setting
leaves clear space around the outer group titles and borders.

## Risks / Trade-offs

An optional import can accidentally link native code into a pure-Go build.
Inspect dependency graphs and executable linkage. Exercise qualified schemes,
both-engine native clients, no-engine external clients, existing CLI and
collection E2E, and native reopening of the large WASM-produced package.
