## Why

Ansible packages currently declare syntax checks by hand. Some packages have
no check, and existing checks can drift from their deployment arguments.

## What Changes

- Generate a Bazel syntax test for every playbook command declared by
  `al_ansible_package`.
- Replace the standalone syntax targets in the migrated packages.

## Capabilities

### New Capabilities

- `packaged-playbook-syntax` provides syntax checks for packaged playbook commands.

### Modified Capabilities

## Impact

The shared Ansible package macro and its infrastructure consumers are affected.
Existing deployment commands retain their arguments and configuration.
