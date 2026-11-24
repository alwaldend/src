---
title: MLflow Vast.ai backend
description: Planned MLflow Projects execution backend for Vast.ai rentals
statuses:
  - maintenance
languages:
  - py
  - go
---

The Python plugin connects MLflow Projects to a Vast.ai job server using gRPC.
The Go job server assigns jobs to existing Vast.ai rentals. A worker on each
rental runs jobs and records their outcomes.

## Go CLI

```sh
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- worker serve
```

The CLI machine parses arguments and starts the application child with a typed
command. The application delegates worker commands to its worker child.
`worker serve` starts the worker's listener machine, and the listener uses
the SQLite-backed job machine. `server serve` starts the job coordinator.

Export the composed application graph as DOT or build the themed SVG on Linux.

```sh
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- state-machine graph
bazel build //projects/mlflow_vastai_backend/cmd/mlflow-vastai:state_machine_svg
```

The SVG is at `bazel-bin/projects/mlflow_vastai_backend/cmd/mlflow-vastai/state-machine.svg`.
Bazel builds Graphviz from pinned source. The SVG embeds the shared diagram
font and uses the common theme's light and dark palettes. The
[state machine page](docs/state-machines.md) embeds the generated graph.

## Job server

```sh
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- server serve \
  --state-dir /path/to/server-state --workers training,12345 \
  --tracking-uri https://mlflow.example.com --ssh-user root
```

Inject `VAST_API_KEY` for read-only rental discovery. Set `MLFLOW_TRACKING_TOKEN`
when the configured tracking endpoint needs bearer authentication. Use a
tracking endpoint reachable from both the server and the rentals. The server
accepts workloads using that configured HTTP or HTTPS tracking endpoint.

`--workers` lists allowed instance IDs and exact labels in preference order.
If the plugin specifies a worker, its selector must match an entry. A label
must identify a single existing rental. Install `mlflow-vastai` on the
rental's SSH account PATH and register its SSH host key with the coordinator's
account. The SSH client uses the rental worker's default state directory.

The gRPC listener binds to `127.0.0.1:50051` by default and accepts only
loopback addresses. Use the local endpoint from a trusted environment, or
front it with an authenticated transport for remote access.

The server saves the selected rental before submitting a workload. It retries
ambiguous submissions on that rental with the original run ID, observes jobs
after clients disconnect, and publishes confirmed terminal results to MLflow.

## Rental worker

Run `worker serve` as the rental's SSH user with `mlflow` on `PATH` and the
project's training runtime installed. The worker uses that MLflow CLI and its
tracking credentials to download the project snapshot. It runs the prepared
entry point with `MLFLOW_TRACKING_URI`, `MLFLOW_EXPERIMENT_ID`, and
`MLFLOW_RUN_ID` set.

Use `worker --state-dir /path/to/state serve` to set the state directory.
The default is `$XDG_STATE_HOME/mlflow-vastai` or
`$HOME/.local/state/mlflow-vastai`. An existing directory must be private to
its owner. The listener keeps `worker.sock`, `jobs.sqlite3`, and per-job
directories under that path. Workload output goes to `stdout.log` and
`stderr.log`. Artifact download and workload errors are logged by the listener.

`worker submit`, `worker status`, and `worker cancel` use the same state
directory. They read one binary protobuf from stdin and return a binary `Job`
on stdout. Submission takes `SubmitJobRequest`. Status and cancellation take
`JobReference`.

The listener reserves the MLflow run ID in SQLite and immediately starts
preparing the workload. Status is `RUNNING` from acceptance through execution.
Cancellation stops the local workload, and the listener records the confirmed
outcome. Graceful shutdown stops active workloads and records their results.
After a crash, a new listener reports unfinished jobs as unknown and preserves
the original reservations.

## Run a project

Give the plugin a JSON backend configuration with `server` set to the gRPC
address. Set `worker` to request a configured worker selector. Use TLS through
an authenticated gateway or connect directly to a trusted local listener as
shown here.

```json
{
  "server": "127.0.0.1:50051",
  "worker": "training",
  "insecure": true
}
```

```sh
export MLFLOW_TRACKING_URI=https://mlflow.example.com
bazel run //projects/mlflow_vastai_backend:mlflow -- run /path/to/project \
  --backend vastai --backend-config /path/to/backend.json --env-manager local
```

Use an HTTP or HTTPS tracking endpoint reachable from the submitting machine,
coordinator, and rental. The rental also needs access to the run's artifact
store.

The adapter uploads a compressed project snapshot to the MLflow run at
`vastai-source/project.tar.gz` before submitting the job. The request carries
the artifact path, rendered entry point, and MLflow tracking context. The worker
needs read access to the run artifacts and a preconfigured training runtime.
The snapshot uses the current contents of Git-tracked files under the project
directory. Set `project_files` to a JSON list of additional project-relative
regular files. For a project outside Git, supply the full file list this way.

If the client reports an unconfirmed submission, use the reported MLflow run
ID to locate the existing job before submitting again.

This version does not stage MLproject `path` parameters or provision Docker,
conda, or virtualenv environments.

## Verify

```sh
bazel test //projects/mlflow_vastai_backend/test/e2e:plugin_test
bazel test //projects/mlflow_vastai_backend/test/e2e:worker_test
bazel test //projects/mlflow_vastai_backend/test/e2e:server_test
bazel build //projects/mlflow_vastai_backend/api/jobpb:jobpb
```

The plugin test uses a local gRPC fixture. The worker test uses the real CLI,
Unix gRPC socket, SQLite state, and an MLflow artifact store. Refresh the
project's Python lock with
`bazel run //projects/mlflow_vastai_backend:uv_lock.update`.

## References

- [MLflow Projects](https://mlflow.org/docs/latest/ml/projects/)
- [MLflow plugin interface](https://mlflow.org/docs/latest/plugins/)
