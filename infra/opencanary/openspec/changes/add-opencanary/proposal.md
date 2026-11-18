## Why

Deploy a broad set of OpenCanary decoys on a dedicated XCP-ng guest with a native Linux installation.

## What Changes

- Provision one Fedora VM through Xen Orchestra.
- Deploy OpenCanary using the repository's OpenCanary Ansible role and collect events locally.

## Capabilities

### New Capabilities

- `opencanary-vm`: A dedicated guest running OpenCanary under systemd.

### Modified Capabilities

None.

## Impact

The component uses the shared Xen cloud-init, Vault AppRole, and Ansible workflows. Live provisioning requires an assigned honeypot network and reviewed deployment inputs.
