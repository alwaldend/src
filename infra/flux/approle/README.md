---
title: Flux Vault AppRole
---

This Terraform root manages the Flux AppRoles, their DNS and SSH signing
policies, and the Flux PKI, OIDC, and SOPS Transit configuration. Shared Vault
manages the mounts and cross-component memberships, including the Forgejo and
Harbor grants. The [DNS root](../dns/README.md) authenticates with the
`src_infra_flux` AppRole and uses a separate state backend.

Review changes with `//infra/flux/approle:tf.plan`. The
[state transfer record](../openspec/changes/retain-flux-dns-and-approle/design.md#state-transfer)
describes the completed migration from shared Vault state.
