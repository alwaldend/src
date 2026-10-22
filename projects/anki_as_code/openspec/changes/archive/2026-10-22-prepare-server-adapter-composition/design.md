## Context

The application already supervises its own executable and injects a Rust CLI runner.

## Goals / Non-Goals

Keep command composition and argument forwarding shared while isolating the native adapter binding. No WASM command or implementation.

## Decisions

Bind runServer to ankiservercgo.Run in server_native.go. Main owns signals and same-executable worker configuration. Keep native as an explicit implementation target behind the public alias. Document SIGINT reset and SIGPIPE ignore at the signal commands.

## Risks / Trade-offs

Preserve Rust help, transport validation, signals and collection workflows through existing installed CLI E2E. A later PR adds the second binding and build selector.

## Validation

Native installed CLI, transport, client lifecycle, application workflow and SQLite E2E passed. Rust Clippy, generated CLI reference and the rendered package diagram passed; the diagram was inspected. The SQLite fixture retains a database and verification report in Bazel test outputs. Required quality and semantic lint run through the exact-candidate delivery receipt. Both PRs target master; #181 retains overlapping setup until this prerequisite merges.
