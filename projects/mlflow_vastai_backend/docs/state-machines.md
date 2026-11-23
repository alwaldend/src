---
title: State machines
description: State transitions and entry actions in the MLflow Vast.ai backend
---

## CLI dispatch

[![CLI, application, and worker state machines](assets/state-machine.svg)](assets/state-machine.svg)

The CLI enters `cli-parsing` and uses Cobra to parse the arguments.
`cli-dispatch` sends the typed command to the application machine.
`app-worker` invokes the worker child. Dashed arrows show calls between
the machines.

- `server serve` calls the job server entry point.
- `worker serve` calls the rental-local worker entry point.
- `state-machine graph` prints the three machine definitions and their calls as one DOT graph.

Both service entry points currently return stub errors. Bazel regenerates the
diagram from the CLI when it builds this page.
