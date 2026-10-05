## Why

The user requested the upstream MLOps training dependencies and reproducible
third-party source so the course recommender and trainer are buildable.

## What Changes

- Pin Yandex-Practicum/mlops-freetrack under third_party and include its module.
- Add its requirements to tools/py, regenerate the lock and Gazelle manifest.
- Preserve existing FastAPI and uvicorn[standard] dependencies with upstream minima.

## Capabilities

Add python-import-validation for the repository quality syntax check.
Application and model contracts belong
to projects/mlops_2026/openspec/changes/archive/2026-10-04-add-fastapi-health-app.

## Impact

MODULE.bazel, tools/py dependency manifest and lock, gazelle_python.yaml, and
third_party/com_github_yandex_practicum_mlops_freetrack. Validate dependency
consistency, semantic lint, and actual HTTP/model consumers.

The user requested enforcement of namespaced imports. Add the pinned ast-grep
quality check and syntax fixtures, migrate existing first-party Python imports,
and update repo-python and its offline evaluation. Configure Ruff first-party
module names so import order is consistent in both source and lint sandboxes.
