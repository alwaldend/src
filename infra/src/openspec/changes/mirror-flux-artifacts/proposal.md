## Why

Provide local copies of the Flux artifacts and publish those same versions
to the repository's Forgejo package registries.

## What Changes

- Add explicit Bazel pull and push targets for the Flux and Flux Operator
  chart archives and images.
- Pin upstream artifacts and retain their exact archive and image contents.

## Capabilities

### New Capabilities

- `third-party-flux-mirror` provides reproducible Flux artifact transfer.

## Impact

The targets use the pinned Helm and ORAS tools and the third-party Forgejo
authentication configuration.
