## Context

The user approved github.com/tetratelabs/wazero v1.12.0 and the WASI SDK 34.0 sysroot archive. Existing LLVM lacks WASI C headers for SQLite.

## Goals / Non-Goals

Make build prerequisites independently reviewable. No Anki networking port, runtime host, timezone bridge or WASM server.

## Decisions

Pin only the WASI SDK sysroot archive with verified SRI and reuse repository LLVM. Declare wazero through existing Go module wiring. Extend the existing Anki crate universe to wasm32-wasip1 for SQLite without changing dependency versions. Regenerate MODULE.bazel.lock using Bazel mod deps.

## Risks / Trade-offs

The sysroot adds approximately 119 MB plus checksum/version maintenance; wazero adds runtime maintenance and execution overhead for later consumers. The standalone SQLite fixture proves headers, libraries and pure-Go execution now.

## Validation

Native installed CLI, transport, client lifecycle, application workflow and SQLite E2E passed. Rust Clippy, generated CLI reference and the rendered package diagram passed; the diagram was inspected. The SQLite fixture retains a database and verification report in Bazel test outputs. Required quality and semantic lint run through the exact-candidate delivery receipt. Both PRs target master; #181 retains overlapping setup until this prerequisite merges.
