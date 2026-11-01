## Why

The recommender acceptance test overrides an obsolete canonical repository
name, so its missing-CSV cases start a healthy server and time out.

## What Changes

- Resolve negative-fixture manifest keys through Bazel's current runfiles mapping.
- Retain both missing-CSV startup checks and their existing process deadlines.

## Capabilities

No application requirements change. This fixture repair uses `skip_specs: true`.

## Impact

Only the existing recommender E2E fixture changes. No dependencies or application
behavior change. Acceptance requires the complete HTTP test and its JSON artifact.
