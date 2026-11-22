## Why

Four component-owned Python lockfiles are registered in the shared Python
module fragment. Keep each hub declaration beside the lockfile its owner maintains.

## What Changes

- Move the Ansible, YubiKey Manager, Vast.ai, and Autoscroll pip declarations
  into the Bazel module fragments of their respective components.
- Regenerate the root module's include list.
- Set a shared Python version for the moved pip hubs and the repository-wide hub.
- Document component-local pip declarations in the `repo-python` skill.

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

Existing repository labels, lockfiles, and Python versions continue to apply.
