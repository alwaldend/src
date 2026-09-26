## Why

Unrestricted Ansible arguments can override the generated VM inventory or
connection and redirect a deployment play to the controller or another host.
The runner must reject this before executing any lifecycle playbook.

## What Changes

- Allow only the diagnostic and tag-skipping argument forms used by existing
  consumers; reject other flags and values during preflight.
- Cover connection, inventory, user, authentication, and extra-variable
  overrides in the real runner prerequisite cases.
- Exercise supported consumer arguments through the public smoke target.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `molecule-qemu-runner`: reject caller arguments that can replace the private
  inventory, connection, or authentication configuration.

## Impact

The runner preflight, README, scenario test declaration, acceptance cases, and
stable isolation requirement change. Existing host-role arguments remain valid.
No deployment or persistent host configuration is performed.
