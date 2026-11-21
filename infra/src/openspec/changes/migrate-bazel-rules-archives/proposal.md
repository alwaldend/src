## Why

Nested Bazel workspaces add independent builds and tests to repository CI.

## What Changes

- Move the shared license package to `internal/license`.
- Move all remaining nested Bazel modules into the root development build.
- Preserve standalone release manifests and existing consumer labels.

## Capabilities

### Modified Capabilities

- `repository` gains additional root-built rulesets with standalone archives.

## Impact

The remaining rulesets retain their repository namespaces under root Bazel
resolution. CI builds and tests those repositories in the root invocation.
