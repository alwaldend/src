## Why

Forgejo Actions run 38 failed the repository formatting check, two ComfyUI
fixtures, and the repository Git history scan.
The runner exhausted memory and the host killed a Bazel JVM as the job ended.

## What Changes

- Format the affected Bazel packages with Buildifier.
- Validate the ComfyUI path and Git scan fixes in their owning components.
- Shut down each workspace's Bazel server after its build and test phases on
  the Forgejo runner.

## Capabilities

### Modified Capabilities

- `repository` defines Bazel server lifetime on the serial Forgejo runner.

## Impact

The repair restores coverage of existing repository CI checks.
