## Why

The runner scaffold must address Forgejo directly after the local Git URL
becomes a browser redirect.

## What Changes

- Point the runner instance URL at `https://forgejo.alwaldend.com`.
- Preserve the existing token reference, disk and inactive runner role.

## Capabilities

### Modified Capabilities

- `infra-forgejo-runner`: Canonical instance URL in packaged configuration.

## Impact

One Ansible variable and its specification. No runner deployment or registration.
