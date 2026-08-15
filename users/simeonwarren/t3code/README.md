---
title: T3 Code
description: t3code.simeonwarren.users.alwaldend.com
tags:
  - ansible
  - ai
  - t3code
---

T3 Code runs as a systemd service on a single Yandex Cloud VM. Traefik on the
same VM terminates TLS and forwards HTTP and WebSocket traffic to the
loopback-only T3 Code server.

## Deployment

```sh
bazel run //users/simeonwarren/t3code/tf_setup
bazel run //users/simeonwarren/t3code/ansible
```

## Regenerate wireguard keys

- Regenerate private and public keys:
  ```sh
  wg genkey | tee host1.privatekey.txt | wg pubkey >host1.publickey.txt
  wg genkey | tee router.privatekey.txt | wg pubkey >router.publickey.txt
  cat - >data.json <<EOF
  {
    "wg_public_keys": {
      "host1": "$(cat host1.publickey.txt)",
      "router": "$(cat router.publickey.txt)"
    },
    "wg_private_keys": {
      "host1": "$(cat host1.privatekey.txt)",
      "router": "$(cat router.privatekey.txt)"
    },
    "wg_preshared_keys": {
      "host1": "$(openssl rand 32 | base64)"
    }
  }
  EOF
  bazel run users/simeonwarren/t3code:vault.kv_put -- -format json alwaldend.com/vault1/approles/user_simeonwarren/t3code/wireguard "@${PWD}/data.json"
  rm data.json *.privatekey.txt *.publickey.txt
  ```
