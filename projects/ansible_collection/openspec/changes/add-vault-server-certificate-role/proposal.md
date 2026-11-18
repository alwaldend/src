## Why

Keep Vault server-certificate signing and installation in a shared Ansible role.

## What Changes

Package a role that signs an existing server key, installs a verified certificate,
and notifies the caller's handler when the certificate changes. Use it in the
OpenCanary role.

## Capabilities

### New Capabilities

- `vault-server-certificate`: Install and renew a server certificate through Vault PKI.

### Modified Capabilities

None.

## Impact

OpenCanary uses the same Vault signing endpoint, guest key, and renewal threshold.
