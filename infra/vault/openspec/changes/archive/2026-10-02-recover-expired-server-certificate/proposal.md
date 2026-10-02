## Why

An expired Vault server certificate prevents credential injection for the
Ansible workflow that renews it. The existing broad role tag also updates
packages, binaries, and configuration unrelated to certificate recovery.

## What Changes

- Expose a certificate-only Ansible selection that stops on issuance or write
  failures and suppresses private-key output.
- Document temporary, invocation-scoped TLS verification bypass and subsequent
  verification with normal certificate validation restored.
- Connect through the existing DNS-owned hostname while retaining the inventory
  identity, and include that connection hostname in the certificate names.
- Renew the active Vault host's certificate within the user's explicit request.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-vault`: support scoped recovery of an expired server certificate.

## Impact

The Vault Ansible role, inventory, certificate variables, and
`infra/vault/README.md`. No new dependencies, persistent insecure defaults, or
authentication-policy changes are needed.
