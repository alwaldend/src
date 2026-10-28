---
title: ssh
description: Ssh keys
---

## Links

- CA generation: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates
- OpenPGP ssh key: https://github.com/drduh/YubiKey-Guide/tree/master

## ca_clients

CA for SSH clients

Source: [ca_clients.pub](ca_clients.pub), a checked-in public-key snapshot.

Refresh: the independent Terraform root in [tf/](tf/README.md) reads Vault's
`ssh/clients/config/ca` public-key endpoint using the existing Vault provider.

## ca_servers

CA for SSH servers

Source: [ca_servers.pub](ca_servers.pub), a checked-in public-key snapshot.

Refresh: [tf/](tf/README.md) reads Vault's `ssh/servers/config/ca`
endpoint. Project builds use these local files without contacting Vault or
depending on infrastructure targets. Infrastructure Terraform manages the
issuers but does not update these snapshots.

## simeonwarren

Source [./simeonwarren@alwaldend.com-2026-03-22-openpgp.pub](./simeonwarren@alwaldend.com-2026-03-22-openpgp.pub)

Update: manual export from the yubikey
