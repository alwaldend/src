---
title: ingress
description: ingress
tags:
  - ingress
  - public
  - traefik
---

## Run ansible

```sh
bazel run //infra/ingress/ansible
```

## Apply terraform

```sh
bazel run //infra/ingress/tf
```

## Update signed image url

- Create a signed url: https://yandex.cloud/ru/docs/storage/operations/objects/link-for-download
- Update the secret: `alwaldend.com/vault1/approles/src_infra_ingress/image`

## Recreate VMs

```sh
bazel run infra/ingress/tf -- -replace 'yandex_compute_disk.vpc["*"]'
```

## Regenerate wireguard keys

- Regenerate private and public keys:
  ```sh
  wg genkey | tee host1.privatekey.txt | wg pubkey >host1.publickey.txt
  wg genkey | tee host2.privatekey.txt | wg pubkey >host2.publickey.txt
  wg genkey | tee router.privatekey.txt | wg pubkey >router.publickey.txt
  cat - >data.json <<EOF
  {
    "wg_public_keys": {
      "host1": "$(cat host1.publickey.txt)",
      "host2": "$(cat host2.publickey.txt)",
      "router": "$(cat router.publickey.txt)"
    },
    "wg_private_keys": {
      "host1": "$(cat host1.privatekey.txt)",
      "host2": "$(cat host2.privatekey.txt)",
      "router": "$(cat router.privatekey.txt)"
    },
    "wg_preshared_keys": {
      "host1": "$(openssl rand 32 | base64)",
      "host2": "$(openssl rand 32 | base64)"
    }
  }
  EOF
  bazel run infra/ingress:vault.kv_put -- -format json alwaldend.com/vault1/approles/src_infra_ingress/wireguard "@${PWD}/data.json"
  rm data.json *.privatekey.txt *.publickey.txt
  ```
