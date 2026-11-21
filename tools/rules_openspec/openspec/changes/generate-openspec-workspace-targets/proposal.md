## Why

The central OpenSpec validation map can omit workspaces and requires separate archive registration. Generate the source and test in each workspace's own Bazel package so ordinary package testing covers its specifications and changes.

## What Changes

- Add an OpenSpec Gazelle language that generates an owner-local source target and validation target.
- Migrate existing workspaces and remove the central validation package and inventory projection.
- Document generation and validation in the repository OpenSpec skill.

## Capabilities

### Modified Capabilities

- `project-rules-openspec` gains Gazelle generation and owner-local validation coverage.

## Impact

Affects the root Gazelle command, owner-local OpenSpec packages, and release manifests for modules that include an OpenSpec workspace. The central test labels are replaced by per-package `:validate_test` labels.
