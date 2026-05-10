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
bazel run //infra/dc1/vault:deploy # host setup
bazel run //infra/dc1/vault/tf:deploy # vault config
```

## Backup

```sh
bazel run //infra/dc1/vault:backup
```

## Unseal

```sh
bazel run //infra/dc1/vault:unseal
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

## Generate and import a user cilent certificate

```sh
username="username"
bazel run //infra/dc1/vault:gen_client_cert -- --user "${username}" --output_dir "${PWD}"
bazel run //tools/ykman -- piv certificates import 9A "${PWD}/${username}.pfx"
bazel run //tools/ykman -- piv keys import 9A "${PWD}/${username}.pfx"
```

## Generate a host cilent certificate

```sh
bazel run //infra/dc1/vault:gen_client_cert -- --host some-host --output_dir "${HOME}/.al/client_cert"
```

## Generate a user device cilent certificate

```sh
bazel run //infra/dc1/vault:gen_client_cert -- --hostsome-host --user username --output_dir "${HOME}/.al/client_cert"
```

## Unseal

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/dc1/vault:unseal
  ```

## Root token

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/dc1/vault:gen_root_token -- --pgp_key path_to_public_gpg_key_in_base64
  ```

## Revoke all tokens

```sh
bazel run //infra/dc1/vault -- token revoke -mode=path auth
```

## Vault certificates

Vault certificates (`tls_cert_file`, `tls_key_file`) should be updated manually

