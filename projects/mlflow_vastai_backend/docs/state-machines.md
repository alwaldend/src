---
title: Application state machines
description: Command dispatch, listener lifecycle, and job transitions in the MLflow Vast.ai backend
---

## CLI dispatch

[![CLI, coordinator, and worker state machines](assets/state-machine.svg)](assets/state-machine.svg)

See the [plugin architecture](../architecture/) for the service boundaries
and protocols.

The CLI machine enters `cli-parsing` and uses Cobra to parse the arguments.
`cli-dispatch` calls the selected application method. The application calls
the worker through its named methods. Each component adds its own states and
child calls to the graph. Dashed arrows show those calls.

- `server serve` runs the coordinator and its job machine.
- `worker serve` runs the worker's listener child machine. Its states open the listener
  in `worker-starting`, serve RPCs in `worker-serving`, and release resources
  in `worker-stopping`.
- `state-machine graph` prints the machine definitions and their calls as one DOT graph.

`worker submit`, `worker status`, and `worker cancel` enter `forward-rpc`
and send protobuf requests over the worker's Unix socket. Bazel regenerates
the diagram from the CLI when it builds this page.

## Coordinator jobs

The coordinator reserves an assignment before contacting the worker. Its
job machine advances from reserved to active or to a recorded result. The
server reconciles unresolved assignments and publishes confirmed terminal
results to MLflow while their tracking runs are active.

## Worker jobs

The listener delegates submission, status, and cancellation to the job runner.
The runner reserves the run ID in SQLite and starts the workload machine as
soon as that transaction commits. The workload machine stages the project and
runs the prepared command. The job runner signals active workloads and records
their outcomes. It retries failed writes using the observed result.

The listener closes the job runner before closing SQLite. The runner stops
active workloads and gives pending writes a bounded time to finish. On
startup, the job runner marks interrupted jobs as unknown through the same
result transition and preserves their run IDs.
