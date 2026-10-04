---
title: MLOps 2026
description: MLOps course
statuses:
  - active
languages:
  - py
  - bzl
---

![Recommender runtime and training](assets/architecture.svg)

An MLOps course project with a factory-based FastAPI application named
`recommender`, served by Uvicorn. App and router instances are created inside
factories. `RecommenderApplication` owns the app and receives its service in the constructor;
`RecommendationRouter` owns the versioned routes and service reference.
The injected `RecommendationService` implements `__aenter__` and `__aexit__`
for model loading and deterministic shutdown. Its idempotent `close()` also
runs from `__del__` as a fallback. Application, router and service state is
private; app and router access use explicit read-only properties. Imports
do not start services or create shared application instances. Response contracts
are exported from `internal/models`.
Implementation lives under `internal/recommender/`, with other future apps owning
sibling packages under `internal/` and separate `cmd/` entry points.

## Run recommender

```sh
bazel_agent bazel run //projects/mlops_2026/cmd/recommender -- --host 127.0.0.1 --port 8000
```

The defaults are `127.0.0.1:8000`. GET `/health` returns `{"status": "ok"}`.
FastAPI serves the generated schema at `/openapi.json`, Swagger UI at `/docs`,
and ReDoc at `/redoc`; the documentation UIs use FastAPI's default external
assets. GET `/api/v1/recommend?track_name=Time&n=5` returns `requested_track`
and song recommendations containing `track_name`, `artists`, and `album_name`.
`track_name` is required; `n` defaults to 5 and must be between 1 and 20.
An unknown track returns an empty list. Invalid queries return HTTP 422.

The pinned upstream `data.csv` and `dataset.csv` are included in the executable
runfiles and resolved independently of the working directory. Startup checks
that both exist. Each app loads its packaged model once during startup using
the upstream `load_model` function, then releases it on shutdown. Health checks
do not invoke recommendation computation. Source and integrity are
owned by `third_party/com_github_yandex_practicum_mlops_freetrack/`.

## Build the course model

```sh
bazel_agent bazel build //projects/mlops_2026:model
```

The sandboxed target runs the upstream `train_and_save_model.py` unchanged,
stages its CSV inputs in an action directory, and produces
`bazel-bin/projects/mlops_2026/model.pkl`. It exposes a declared build artifact
consumed by the recommender executable. The upstream script is
also buildable as
`//third_party/com_github_yandex_practicum_mlops_freetrack:train_and_save_model`;
running that script directly requires a working directory with both CSVs and
writes `model.pkl` there. Prefer the model build target for isolated execution.

The full training data joins to 6,788 unique tracks. Its dense similarity matrix
uses about 0.34 GiB; the serialized artifact is about 359 MiB. The model target and isolated model
verification are marked manual, but building or testing the recommender now
requires this artifact and can trigger training when it is not cached. Shared Python dependencies are declared and locked under `tools/py/`.

## Verify

```sh
bazel_agent bazel test //projects/mlops_2026/test/e2e:recommender_test
bazel_agent bazel test //projects/mlops_2026/test/e2e:model_test
```

These checks exercise real HTTP, factory isolation, startup conflict handling,
shutdown, service context exit, repeated and deletion cleanup, course runfiles
and checksums, and a recommendation from the built
model. Repeatable JSON artifacts are saved in each test's undeclared outputs
under `bazel-testlogs/projects/mlops_2026/test/e2e/`.
There is no deployment or release target.

## Links

- [Yandex Practicum MLOps course](https://practicum.yandex.ru/mlops/)
- [Upstream course source](https://github.com/Yandex-Practicum/mlops-freetrack)
