---
title: MLflow Vast.ai backend
description: Planned MLflow Projects execution backend for Vast.ai rentals
statuses:
  - maintenance
languages:
  - py
  - go
---

This project is a scaffold for executing MLflow Projects on Vast.ai GPU rentals.
A Python adapter will implement MLflow's project backend and submitted-run
interfaces. It will delegate rental operations to a Go helper using the
repository's [Vast.ai client](../../internal/vastai_sdk/README.md) and
[Stateless](https://github.com/qmuntal/stateless).

The scaffold contains documentation and dependencies. The adapter, Go helper,
and plugin registration are pending.

## Verify the dependencies

Run MLflow's upstream CLI and build the Go state-machine library:

```sh
bazel run //projects/mlflow_vastai_backend:mlflow -- run --help
bazel build @com_github_qmuntal_stateless//:stateless
```

The Python environment is locked in this project. Refresh it with
`bazel run //projects/mlflow_vastai_backend:uv_lock.update`.

## References

- [MLflow Projects](https://mlflow.org/docs/latest/ml/projects/)
- [MLflow plugin interface](https://mlflow.org/docs/latest/plugins/)
