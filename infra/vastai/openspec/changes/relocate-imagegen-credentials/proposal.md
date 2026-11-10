## Why

The alwaldend team root owns key issuance. Image generation is moving from a
personal rental into infra/image_gen with its own Vault identity and state.

## What Changes

Destroy the legacy dedicated key and current personal delivery entry, then
declare a fresh component key delivered into src_infra_image_gen/vastai.
The existing team issuer and unrelated team resources remain unchanged. The
Vault owner supplies the new exact-path delivery grant; the component consumes
its generated key with its own AppRole.

## Impact

The user authorized legacy destruction, which is complete. Replacement key
creation is not part of source delivery. No state handoff, import, key-value
reuse or GPU rental is required. Historical handoff receipts remain context.
