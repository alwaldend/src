## Context

MLflow discovers project backends through Python entry points. The repository
already has a Go client for Vast.ai and uses Gazelle's Go dependency extension.

## Decisions

- A Python adapter will implement MLflow's backend and submitted-run interfaces.
  It will delegate Vast.ai API calls and rental lifecycle management to a Go
  helper that can reuse the repository's Vast.ai client and Stateless.
- Inject the existing `afero.Fs` interface for all helper-owned file
  operations, following the
  [Anki as Code filesystem pattern](../../../../anki_as_code/pkg/ankiascode/app.go).
  Supply `afero.NewOsFs()` at the command boundary and use
  `afero.NewMemMapFs()` for tests that need isolated storage.
- Keep MLflow in a project-local uv lock to isolate its dependency graph from
  consumers of the shared Python environment. Use MLflow's own console script
  to verify the environment. Plugin registration needs an importable adapter.
- Resolve Stateless through the existing root Go module so the future helper
  can share the same workspace as the Vast.ai client.
