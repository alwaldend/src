---
title: Vault
description: Hashicorp Vault
tags:
  - vault
---

## Manual certificate creation

https://arminreiter.com/2022/01/create-your-own-certificate-authority-ca-using-openssl/

## Unseal process

- Get the encrypted unseal key
- Unencrypt it: `cat unseal.txt | base64 --decode | gpg -dq`
- Run `vault operator unseal`
