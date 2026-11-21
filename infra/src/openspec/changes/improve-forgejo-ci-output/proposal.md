## Why

The Forgejo build and test step streams every Bazel phase into one log, making
it difficult to identify a failed workspace or find its diagnostics.

## What Changes

- Group each workspace's build and test logs in Forgejo.
- Fold routine Bazel diagnostics inside each phase and the CI launcher build.
- Show the outcome and duration of each phase in the job summary.
- Expose the pinned WASI SDK sysroot to LLVM as a single directory artifact.

## Capabilities

### Modified Capabilities

- `repository` gains native Forgejo presentation for repository CI results.

## Impact

The shared CI command and local action own the log presentation. The WASI SDK
dependency declaration uses the same LLVM repository rule as the Linux sysroots.
