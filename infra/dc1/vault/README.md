---
title: Vault
description: Setup for vault.dc1.alwaldend.com
tags:
  - ansible
  - terraform
---

## Links

- Intermediate CA: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca

## Deployment

```sh
bazel run //infra/dc1/vault/ansible:deploy
bazel run //infra/dc1/vault/tf:deploy
```

## Tf

Plan:
```sh
bazel run //infra/dc1/vault/tf:tf.plan
```

Apply:
```sh
bazel run //infra/dc1/vault/tf:tf.apply
```

Run terraform directly:
```sh
bazel run //infra/dc1/vault/tf:tf.direct -- -chdir="${PWD}" plan
```

## Backup

```sh
bazel run //infra/dc1/vault/backup
```

## Generate and import a client certificate

```sh
username="username"
bazel run //infra/dc1/vault/gen_client_cert -- --username "${username}" --output_dir "${PWD}"
bazel run //tools/ykman -- piv certificates import 9A "${PWD}/${username}.pfx"
bazel run //tools/ykman -- piv keys import 9A "${PWD}/${username}.pfx"
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

