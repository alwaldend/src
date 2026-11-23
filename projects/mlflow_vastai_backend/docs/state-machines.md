---
title: State machines
description: State transitions and entry actions in the MLflow Vast.ai backend
---

## CLI dispatch

[![CLI, application, and worker state machines](assets/state-machine.svg)](assets/state-machine.svg)

The CLI enters `cli-parsing` and uses Cobra to parse the arguments.
`cli-dispatch` calls the selected application method. The application calls
the worker through its named methods. Each component adds its own states and
child calls to the graph. Dashed arrows show those calls.

- `server serve` calls the job server entry point.
- `worker serve` calls the rental-local worker entry point.
- `state-machine graph` prints the three machine definitions and their calls as one DOT graph.

Both service entry points currently return stub errors. Bazel regenerates the
diagram from the CLI when it builds this page.
