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

Source: [../../../infra/dc1/vault/tf/output](../../../infra/dc1/vault/tf/output/README.md)

Update: automatically

## Server certificate CA

Source: [../../../infra/dc1/vault/tf/output](../../../infra/dc1/vault/tf/output/README.md)

Update: automatically
