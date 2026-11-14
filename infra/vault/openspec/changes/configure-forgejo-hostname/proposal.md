## Why

Vault's Forgejo URL default still selects `git.alwaldend.com`. OIDC must
use the canonical Forgejo hostname while certificates cover the local alias.

## What Changes

- Default the Forgejo URL to `https://forgejo.alwaldend.com`.
- Derive the allowed callback from that URL without a duplicate literal.
- Retain certificate issuance permissions for both Forgejo and Git names.

## Capabilities

### New Capabilities

- `forgejo-hostname`: Canonical Forgejo OIDC callback and local TLS alias.

## Impact

Vault Terraform source only. No apply, credential rotation or live Vault writes.
