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
The shared protobuf contract defines both the job server
and the rental worker service. The Go services are pending.

## Go CLI

```sh
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- server serve
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- worker serve
```

The CLI machine passes the parsed command to the application machine. The
application delegates worker commands to its worker child. Both service entry
points currently return a nonzero exit status with a stub error.

Export the composed machines as DOT or build the themed SVG on Linux.

```sh
bazel run //projects/mlflow_vastai_backend/cmd/mlflow-vastai -- state-machine graph
bazel build //projects/mlflow_vastai_backend/cmd/mlflow-vastai:state_machine_svg
```

The SVG is at `bazel-bin/projects/mlflow_vastai_backend/cmd/mlflow-vastai/state-machine.svg`.
Bazel builds Graphviz from pinned source. The SVG embeds the shared diagram
font and uses the common theme's light and dark palettes. The
[state machines page](docs/state-machines.md) embeds the generated graph.

## Run a project

Give the plugin a JSON backend configuration with `server` set to the gRPC
address. Set `worker` to request a configured worker selector. Connections use
TLS by default. Set `insecure` to `true` for a trusted local test service.

```json
{
  "server": "jobs.example.com:443",
  "worker": "training"
}
```

```sh
export MLFLOW_TRACKING_URI=https://mlflow.example.com
bazel run //projects/mlflow_vastai_backend:mlflow -- run /path/to/project \
  --backend vastai --backend-config /path/to/backend.json --env-manager local
```

Use a tracking endpoint and artifact store reachable from both the submitting
machine and the rental. MLflow's default local SQLite store is only suitable
when both sides share the same store and artifact paths.

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
bazel build //projects/mlflow_vastai_backend/api/jobpb:jobpb
```

The end-to-end test runs an MLflow project against a local gRPC fixture and
checks tracking results, submission replay, and cancellation. Refresh the
project's Python lock with
`bazel run //projects/mlflow_vastai_backend:uv_lock.update`.

## References

- [MLflow Projects](https://mlflow.org/docs/latest/ml/projects/)
- [MLflow plugin interface](https://mlflow.org/docs/latest/plugins/)
