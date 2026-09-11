## Why

The goal tool has been deprecated since the 2026-09-08 migration and its
replacement, the pinned OpenSpec CLI, is in active use. Keeping the deprecated
store, CLI, skill sources, landing site, and DNS record costs maintenance and
contradicts the documented workflow, which already tells agents not to create
goal records.

## What Changes

- **BREAKING** Remove the `projects/goal` component: CLI, API package, local
  filesystem store, docs and diagrams, landing site, disabled skill sources,
  Bazel targets, and DNS record.
- Remove the `goal` entry from the project registry that drives landing-site,
  DNS-config, and Pages-repository generation, along with its committed zone
  record.
- Withdraw the `project-goal` capability and its specification, since the
  component no longer exists.
- Leave the removed bytes in git history; the change neither copies nor
  restates them.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `repository`: The legacy goal migration requirement changes from "disable and
  mark deprecated" to "removed", and the component-coverage requirement no
  longer expects a `projects/goal` workspace.

## Impact

`projects/goal` is deleted; `projects/projects.bzl`, `projects/README.md`,
`infra/src/openspec/validation/BUILD.bazel`, and `infra/dns/zones` drop their
goal references. The removed bytes stay retrievable from the removal commit's
parent in git history; other OpenSpec `changes/`, `specs/`, and archived
`provenance/` bytes are unaffected.
