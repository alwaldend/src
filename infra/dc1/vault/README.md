---
title: Vault
description: Setup for vault.dc1.alwaldend.com
tags:
  - ansible
  - terraform
---

## Deployment

```sh
bazel run //infra/dc1/vault/ansible:deploy
bazel run //infra/dc1/vault/tf:deploy
```

## Tf

Plan:
```sh
bazel run //infra/dc1/vault/tf:plan
```

Apply:
```sh
bazel run //infra/dc1/vault/tf:deploy
```

Run terraform directly:
```sh
bazel run //infra/dc1/vault/tf -- -chdir="${PWD}" plan
```

## Unseal

- Get the encrypted unseal key
- Decrypt it:
  ```sh
  cat unseal.txt  | base64 --decode | gpg --decrypt; rm -f unseal.txt
  ```
- Run `vault operator unseal`

## Vault certificates

Vault certificates (`tls_cert_file`, `tls_key_file`) should be updated manually
