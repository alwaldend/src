## Why

The shared AL Forgejo login URL must follow the service's canonical hostname
so API authentication does not depend on a browser redirect.

## What Changes

- Update the existing `tools/al` Forgejo login URL to `https://forgejo.alwaldend.com`.
- Preserve the OAuth source, executable and credential injection flow.

## Capabilities

No new contract. This updates an existing consumer's endpoint reference;
`skip_specs: true` keeps service requirements with `infra/forgejo`.

## Impact

Shared AL configuration source and its offline packaging checks. No login or
live authentication operation.
