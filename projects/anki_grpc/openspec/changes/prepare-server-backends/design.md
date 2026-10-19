## Context

pkg/ankiserver currently combines native Rust transport and a Go cgo wrapper. The WASM implementation remains in PR #181.

## Goals / Non-Goals

Keep the native adapter API and supervised process boundary. Extract prerequisites only.

## Decisions

Move FFI and Go wrapper together into pkg/ankiservercgo. Retain Rust transport in pkg/ankiserver and expose its std-only CLI parser as server_cli so both backends can call the same Rust CLI. Give command implementations explicit native targets and preserve public aliases. Declare backend with cgo as its sole supported value until the WASM implementation lands. Reuse the existing native/WASI SQLite fixture and built-in rules_rust platform selection to verify the toolchain without a server.

## Risks / Trade-offs

A package import changes; update its actual application consumer atomically. No duplicate compatibility wrapper. Validate installed native CLI and SQLite cross-runtime reopen.

## Validation

Native installed CLI, transport, client lifecycle, application workflow and SQLite E2E passed. Rust Clippy, generated CLI reference and the rendered package diagram passed; the diagram was inspected. The SQLite fixture retains a database and verification report in Bazel test outputs. Required quality and semantic lint run through the exact-candidate delivery receipt. Both PRs target master; #181 retains overlapping setup until this prerequisite merges.

The user requested removal of wasi_module. Explicit native and WASI probe targets use the same Rust source with platform=None or @rules_rust//rust/platform:wasip1. No custom transition or output rename rule remains.

The SQLite native/WASI cross-reopen E2E passed after this removal, retaining its database and verification artifact.
