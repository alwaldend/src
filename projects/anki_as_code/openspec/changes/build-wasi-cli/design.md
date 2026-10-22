## Context

The public Go command composes server adapters through two independent Bazel
booleans. Its startup currently resolves `os.Executable` before parsing commands.
Go's WASI implementation cannot resolve that path or spawn subprocesses.

## Goals / Non-Goals

Provide project-local selection through the existing public target. Browser JavaScript
integration and WASI networking extensions are outside this change.

## Decisions

- **Verdict: revise.** The public alias uses `select()` on the project-local
  `//projects/anki_as_code:wasm` boolean, following `anki_grpc`'s existing alias
  selection pattern. It selects `go_cross_binary` wrapping the existing
  `no_engine` variant with platform `@rules_go//go/toolchain:wasip1_wasm`.
  The built-in rules_go transition cross-compiles that binary and its dependencies; sibling targets
  retain the invocation's platform and server-engine settings.
- The selected variant has no server tags or adapter dependencies, independently
  of the two `anki_grpc` engine flags. It reuses the existing command-source
  declarations. The pinned rules_go macro rejects `goos`/`goarch` alongside the
  existing link-option `select()` (`go-platform-select-restriction`), so the
  supported cross-binary wrapper applies the transition without replacing those
  link options. Selecting the complete variant keeps the server-free contract explicit.
- A global rc profile would work for a single target but changes the platform
  and engine flags of every target in a mixed invocation. Remove that profile
  instead of maintaining two ways to select the same project behavior.
- Generate schema imports through a project-owned rule whose descriptor input
  and generator both use `cfg = "exec"`. The previous `run_binary` input used
  target configuration, incorrectly requiring a WASI Rust toolchain for the
  upstream descriptor build script (`wasi-schema-generator-platform` failure).
- Resolve the executable only when a compiled engine needs a same-executable
  worker. Without engines, leave the factory default address empty so the
  existing client constructor diagnoses missing addresses.
- Exercise the produced module with the repository's existing pinned wazero CLI,
  passing the exact public target built through the project-local flag. Retain
  command output and the exported TOML tree under ignored task scratch.
- Apply a small post-Gazelle BUILD patch to the pinned gRPC module. Gazelle's
  platform list omits WASI, so regeneration loses the `grpclog` dependency in
  `internal/syscall` (`grpc-wasi-build-metadata` failure). The patch adds the
  existing dependency to the default selection branch; versions stay unchanged.

## Risks / Trade-offs

- WASI Preview 1 cannot dial TCP/Unix sockets or create subprocesses with standard
  Go. Document that collection/API operations need additional host integration;
  this build supports file-only config workflows and command discovery.
- Conditional startup must preserve same-executable native/WASM worker behavior.
  Run the existing native CLI workflows alongside the WASI command check.

## Migration Plan

Opt in with `--//projects/anki_as_code:wasm=true`; the default public target
retains its configured host engines. Remove the project flag, alias selection,
and WASI variant to roll back the additional build mode.
