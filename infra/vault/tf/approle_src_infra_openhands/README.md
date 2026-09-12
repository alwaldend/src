---
title: OpenHands AppRole
description: src_infra_openhands AppRole, SSH role, and server PKI role
tags:
  - terraform
  - vault
---

The `src_infra_openhands` AppRole used by `infra/openhands`. It grants the
component's Vault identity only the secrets, SSH host-key signing, and server
certificate issuance it needs:

- `approle`: the AppRole, its entity, its group, and the KV policy scoped to
  `alwaldend.com/vault1/approles/src_infra_openhands/*`.
- `ssh`: a `vault_ssh_server_role` allowing host-key signing for
  `openhands.alwaldend.com`.
- `pki_server`: a `vault_pki_server` role for the canvas, server, and
  automation hostnames, plus the group allowed to create its ACME EABs.

The caller supplies the Vault backends so this module depends on no specific
mount names.
