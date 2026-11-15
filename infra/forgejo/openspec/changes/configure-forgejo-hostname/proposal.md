## Why

Forgejo must use `forgejo.alwaldend.com` as its canonical hostname. Local
`git.alwaldend.com` should redirect there while public requests retain the
download server's GitHub redirect and Go discovery behavior.

## What Changes

- Set Forgejo's HTTP, SSH and root URL hostname to `forgejo.alwaldend.com`.
- Redirect local Git browser requests to the canonical URL, preserving paths
  and queries. Retain TLS coverage for both names and the inventory hostname.
- Preserve DNS declarations and coordinate the Vault and client URL updates.
- Replace the historical architecture diagram with the current XCP-ng deployment,
  canonical hostname, split Git routing, and Vault authentication.

## Capabilities

### Modified Capabilities

- `infra-forgejo`: Canonical Forgejo hostname and local Git redirection.

## Impact

Forgejo's Ansible variables, proxy templates and owner documentation.
Offline source validation only; no deployment or service migration.
