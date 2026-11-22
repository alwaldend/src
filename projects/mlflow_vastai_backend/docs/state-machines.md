---
title: State machines
description: State transitions and entry actions in the MLflow Vast.ai backend
---

## CLI dispatch

![CLI state machine from argument parsing to the job server, worker service, or graph export](assets/state-machine.svg)

The CLI starts in `ready` and enters `parsing`. Cobra parses the arguments and
selects a transition. The destination state invokes the named entry action
shown in the diagram.

- `server serve` calls the job server entry point.
- `worker serve` calls the rental-local worker entry point.
- `state-machine graph` prints the configured machine as DOT.

Both service entry points currently return stub errors. Bazel regenerates the
diagram from the CLI when it builds this page.
