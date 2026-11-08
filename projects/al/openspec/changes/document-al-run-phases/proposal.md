## Why

The al run documentation needs one diagram showing which work happens during
Bazel Build and which happens during Run.

## What Changes

- Add one Mermaid source with Build and Run phases and a nested Runfiles
  subgraph to the lifecycle page.
- Show merged configuration dependencies, command runfiles, and AL reading
  configuration and running plugins and commands. Show plugins using external services.
- Publish themed light/dark renders through the repository AL integration.

## Capabilities

None. Documentation only; skip_specs is enabled.

## Impact

AL documentation and its exported diagram source. Repository render packaging
is recorded in infra/src/openspec/changes/publish-al-run-diagram.
