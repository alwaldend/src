## Why

The ComfyUI role rejects absolute paths under hidden directories such as
`.cache`. Bazel's runner layout supplies those paths to two acceptance tests.

## What Changes

- Allow dots in ordinary absolute path components.
- Continue to reject relative paths and dot or dot-dot components.

## Capabilities

### Modified Capabilities

- `project-ansible-collection` defines the accepted paths for ComfyUI.

## Impact

The ComfyUI role uses the same path validation for application and state
directories.
