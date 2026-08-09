---
title: Hermes
description: hermes.alwaldend.com
tags:
  - ansible
  - pve
  - hermes
---

Hermes Agent runs as a systemd service on a single Proxmox VM. Traefik on the
same VM terminates TLS and forwards gateway HTTP and WebSocket traffic to the
loopback-only Hermes API server.

## Deployment

```sh
bazel run //users/simeonwarren/hermes/tf_setup
bazel run //users/simeonwarren/hermes/ansible
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
  bazel run users/simeonwarren:vault.kv_put -- -format json alwaldend.com/vault1/approles/user_simeonwarren/wireguard "@${PWD}/data.json"
  rm data.json *.privatekey.txt *.publickey.txt
  ```
