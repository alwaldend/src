## Context

VaultSshFetcher assigned the same .pub path to the generated public key and signed certificate.

## Decisions

Use the conventional -cert.pub certificate path and keep private/public/certificate files distinct. Keep existing environment field references and ephemeral cleanup. No persistent key, policy or secret change is introduced.

## Evidence

Authenticated publication and strict-host-CA Ansible deployment succeeded on both authorized download hosts after this correction; injector cleanup checks pass.
