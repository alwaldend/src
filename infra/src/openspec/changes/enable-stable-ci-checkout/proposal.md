## Why

Forgejo creates temporary checkout paths, so Bazel starts new servers across
jobs. The dedicated runner can reuse a working directory across serial jobs.

## What Changes

- Add opt-in stable checkout preparation to the local repository CI action.
- Enable it in the existing workflow and allow other workflows to call it.

## Capabilities

### New Capabilities

- `repository-ci`: Reusable repository CI with an optional stable checkout.

## Impact

Changes the local CI action and Forgejo workflow. The first checkout still
provides the event revision, Git history, tags, and LFS objects.
