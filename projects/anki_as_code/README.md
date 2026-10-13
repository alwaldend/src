---
title: Anki as code
description: CLI scaffold for managing Anki collections as TOML
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code is a Go CLI project for managing Anki collections through editable
TOML. This scaffold provides command discovery, flags, cancellation, and output
routing. Collection operations are placeholders: valid operation invocations
exit with an explicit "not implemented" error and do not read or write collections.

## Commands

| Command       | Arguments                                       | Intended operation              |
| ------------- | ----------------------------------------------- | ------------------------------- |
| `export`      | `--input`, `--output`                           | Export a collection to text     |
| `plan`        | `--input`, `--text`                             | Inspect desired changes         |
| `apply`       | `--input`, `--text`                             | Reconcile an offline collection |
| `build`       | `--base`, `--text`, `--output`                  | Build a collection archive      |
| `generate-id` | One or more note paths; optional `--collection` | Assign note and card identities |

All listed flags are required except `--collection`. Help and future plan JSON
use stdout; status messages and errors use stderr. Each subcommand supports
`--help`.

## Build and run

From the repository root:

```sh
bazel_agent bazel build //projects/anki_as_code:anki_as_code
bazel_agent bazel run //projects/anki_as_code:anki_as_code -- --help
```

## Project layout

- `cmd/anki_as_code/`: Cobra commands, flags, and output handling.
- `internal/app/`: App struct and placeholder collection operations.
- `site/content/`: Landing content for the shared website.
- `openspec/`: Project specifications and maintained changes.

The collection implementation is developed separately in
[PR #128](https://github.com/alwaldend/src/pull/128). Its end-to-end CLI test
requires real note-identity generation and belongs with that implementation.

## Anki backend API

The reusable Go backend and Rust server belong to
[projects/anki_grpc](../anki_grpc/README.md#go-backend). The Go package owns the
server process and exposes a managed gRPC connection. Callers use Anki's generated
service clients for typed methods. Anki owns collection validation, note updates,
card generation, and persistence; this CLI owns TOML and reconciliation.

## Integration approaches considered

Direct storage editing has the simplest integration: an application edits
SQLite and protobuf itself, taking responsibility for Anki's data invariants and
business logic. The API-based approaches are Python with Anki's official SDK,
Rust with Anki's Rust API, and another language with a translation layer around
the Rust API. The two Rust API approaches require building Anki from source.
A translation layer can be a C ABI or a gRPC sidecar; it does not remove that
build requirement.

This project uses a translation layer: the Go CLI calls the Rust server in
`projects/anki_grpc`, which forwards operations to Anki's Rust API. Go owns the
server process and its failures.

![Anki backend integration approaches](/docs/projects/anki_as_code/assets/backend-approaches.svg)

The maintained [Mermaid source](backend-approaches.mmd) is rendered by the
`:backend_approaches` Bazel target and published with the project documentation.

## Verification

The [backend E2E](../anki_grpc/README.md#verification) verifies native persistence,
typed service calls, and process lifecycle. CLI reconciliation remains in PR #128.

```sh
bazel_agent bazel test //projects/anki_grpc/pkg/backend:backend_test
bazel_agent bazel build //projects/anki_grpc:clippy
bazel_agent bazel build //projects/anki_as_code:backend_approaches
```
