## Why

The user requested the latest stable Traefik for the local inference host.
The shared dependency is 3.7.16; upstream marks 3.7.17 as latest stable.

## What Changes

- Update the existing Linux AMD64 binary lock to 3.7.17 using the publisher's
  checksum manifest and matching release asset digest.
- Preserve the binary alias, platform constraints, and shared role packaging.
- Exercise the new binary against the prepared T3 mTLS route.

## Capabilities

No new capability contract: maintenance of an existing dependency.

## Impact

Every consumer receives the new binary on its next build. No running service
is changed. Deployment remains outside the authorized scope.
