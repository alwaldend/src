## Why

The Flux Operator UI needs a shared ingress route with the existing
client-certificate and upstream TLS verification settings.

## What Changes

- Add the Operator hostname and Flux-owned private upstream to the Traefik route map.

## Capabilities

### Modified Capabilities

- infra-ingress adds the Flux Operator route.

## Impact

The shared ingress Ansible package consumes the Flux DNS declarations.
