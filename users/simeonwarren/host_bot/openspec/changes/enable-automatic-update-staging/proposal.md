## Why

The authorized fleet maintenance must keep Host Bot automatic updates enabled in its owning desired state.

## What Changes

Declare automatic updates enabled and reboots disabled; reuse canonical backend selection and automatic-update tasks. Preserve the two explicitly approved appliance SSH keys in the existing known_hosts file without replacing other entries.

## Capabilities

### New Capabilities

- `automatic-update-staging`: Unattended staging without reboot.

## Impact

Host Bot Ansible configuration and canonical updater role. No new dependency.
