## Why

The shared dependency-graph repair must load Autoscroll's native uv lock without
exposing macOS-only PyObjC packages as Linux source builds. uv can preserve the
platform distinction in package resolution markers without a custom importer.

## What Changes

- Partition uv resolution into macOS and non-macOS environments, whose union
  preserves every previously supported operating system.
- Regenerate the native lock through the existing Bazel uv update target.
  Preserve existing package versions and distribution hashes; include the
  user-approved tomli 2.5.0 transitive dependency for Python 3.10.
- Document why the partition is needed and validate unpatched native imports.

## Capabilities

No new application behavior. This component record implements the shared
`python-lock-platform-reachability` contract in
`infra/src/openspec/changes/repair-site-dependency-loading`.

## Impact

Autoscroll's uv configuration, generated lock, and dependency documentation.
No new requirements-file workflow, host changes, or deployment.
