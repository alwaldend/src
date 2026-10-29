---
title: Alwaldend.com
description: Alwaldend.com certificates
tags:
  - pki
---

## Links

- Generation: https://arminreiter.com/2022/01/create-your-own-certificate-authority-ca-using-openssl/
- Vault: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca

## Root CA

Source: [./root_ca.crt](./root_ca.crt)

Update: manually

## Intermediate CA1

Source: [./ica1.crt](./ica1.crt)

Update: manually

## Client certificate CA

Source: [pki_ca_clients.crt](pki_ca_clients.crt), a checked-in public certificate.

Refresh: the independent Terraform root in [tf/](tf/README.md) reads Vault's
`pki/ica_clients/cert/ca` endpoint using the existing Vault provider.

## Server certificate CA

Source: [pki_ca_servers.crt](pki_ca_servers.crt), a checked-in public certificate.

Refresh: [tf/](tf/README.md) reads Vault's public `pki/ica_servers/cert/ca`
endpoint. Project builds use these local files without contacting Vault or
depending on infrastructure targets. Infrastructure Terraform manages the
issuers but does not update these snapshots.
