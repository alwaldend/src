## Why

The Forgejo build and test step streams every Bazel phase into one log, making
it difficult to identify a failed workspace or find its diagnostics.

## What Changes

- Group each workspace's build and test logs in Forgejo.
- Show the outcome and duration of each phase in the job summary.

## Capabilities

### Modified Capabilities

- `repository` gains native Forgejo presentation for repository CI results.

## Impact

The shared CI command supplies the presentation to the existing action.
