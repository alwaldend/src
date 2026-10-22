## Why

The existing Gazelle 0.47.0 does not support ephemeral Go module caching.
The user selected Gazelle 0.53.0 and requested verification of repository
cache reuse across worktrees before adopting that host setting.

## What Changes

- Upgrade the existing Gazelle BCR dependency to the approved 0.53.0 release.
- Reconcile the module lock through Bazel's owning workflow.
- Verify representative Go consumers and shared completed-repository reuse
  from fresh output bases; record cache-miss limitations.

## Capabilities

No specification delta: this updates existing build tooling.
Host configuration is owned by the host-bot change
`enable-ephemeral-go-cache`.

## Impact

The shared Gazelle pin and generated module lock affect root-workspace Go
repository fetching and Gazelle consumers. This is a regular PR against
`master`, independent of the lesson and linter PRs.
