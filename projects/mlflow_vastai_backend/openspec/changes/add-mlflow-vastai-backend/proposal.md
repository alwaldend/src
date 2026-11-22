## Why

Give the planned Vast.ai execution backend a project home and make its
dependencies available for development.

## What Changes

- Add a project README and register its landing page.
- Lock MLflow and expose its upstream CLI through Bazel.
- Pin Stateless in the repository Go module for future rental lifecycle work.

## Capabilities

This change adds no backend behavior.

## Impact

The new project owns its Python environment. The shared Go module exposes
Stateless through the existing Bazel Go dependency extension.
