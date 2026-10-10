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
