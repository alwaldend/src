## Why

The shared Traefik pin is 3.7.1. The user requested extracting the already
verified 3.7.16 upgrade and toolchain-schema correction into their own PR.

## What Changes

- Pin the existing Linux AMD64 archive to 3.7.16 with its upstream SHA-256.
- Move execution-platform constraints from `rules` to the supported
  `toolchain` field.
- Preserve the current binary alias, platform, and Ansible packaging.

## Capabilities

No repository requirement changes. This is maintenance of an existing dependency.

## Impact

All consumers of the shared binary receive the new pin on their next build.
Routing, role hardening, and live deployment belong to
[PR #202](https://github.com/alwaldend/src/pull/202). This extraction performs
no live deployment and introduces no new dependency.
