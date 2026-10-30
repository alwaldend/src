## Why

Use Bazel 9 repository caching and the existing protobuf prebuilt compiler support to reduce repeated dependency setup and compiler builds.

## What Changes

- Enable experimental remote repository contents caching for configured cache backends.
- Prefer prebuilt protoc in the root workspace.
- Document the existing default local repository cache and scope limitations.

## Capabilities

### New Capabilities

None. Build tooling configuration only; skip_specs is enabled.

### Modified Capabilities

None.

## Impact

Shared tools/bazelrc configuration and documentation. No dependency version changes, new cache service, or host configuration changes.
