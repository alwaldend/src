## Why

The MLOps course needs a runnable recommendation service and reproducible model
training, with an inspectable API and independent application instances.

## What Changes

- Name the application and command recommender; use internal/recommender for
  app-specific implementation and internal/models for shared response contracts.
- Construct fresh app, routers, and model service inside factories. Load and
  release the model during lifespan, with no module-level application objects.
- Serve GET /health and GET /api/v1/recommend with required track_name and
  bounded n (default 5, range 1–20); preserve default OpenAPI and docs endpoints.
- Package pinned upstream CSVs and model.pkl in executable runfiles.
- Build the upstream trainer and declare its model artifact through Bazel.
- Add upstream requirements to shared Python dependencies and generated metadata.
- Add the README Mermaid runtime/training diagram and truthful run instructions.
- Require and enforce namespaced imports through the Python skill and the existing
  pinned ast-grep tool in repository quality; migrate first-party imports.

## Capabilities

### New Capabilities

- health-api: HTTP service, recommendation contract, OpenAPI, and independent
  application lifecycle.
- course-data: Reproducible upstream CSV packaging and model training artifact.

### Modified Capabilities

- course-placeholder: Truthful course documentation, supported commands, and
  absence of a release deployment requirement.

## Impact

Project source, models, E2E checks, diagram, and docs live in projects/mlops_2026.
Shared dependency and Python convention changes have a companion record in
infra/src/openspec/changes/archive/2026-10-04-add-mlops-training-dependencies. The approved source
https://github.com/Yandex-Practicum/mlops-freetrack.git is pinned under third_party.
Its loader receives an optional path compatibility patch, preserving the default.
Add scikit-learn, pandas, numpy, and joblib with the upstream version minima,
and preserve FastAPI and uvicorn[standard]. No deployment, authentication,
online model training, or unrelated application behavior changes are included.
