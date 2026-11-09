## Why

The existing alwaldend root already owns team API keys. The personal rental must only consume its dedicated credential.

## What Changes

Move the dedicated token and user Vault write into infra/vastai/alwaldend, reusing the existing issuer and user identity. Remove the separate credential root. Preserve the existing key through an explicitly authorized state handoff rather than generating a replacement.

## Impact

No new AppRole, dependency, key rotation, or GPU rental. The exact-path delivery policy is owned and deployed by infra/vault. The authorized state handoff preserved the existing key and the team plan reports no changes.
