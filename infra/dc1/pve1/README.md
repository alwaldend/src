---
title: Pve1
description: Proxmox cluster pve1.dc1.alwaldend.com
tags:
  - ansible
  - pve
---

## Links

- Docs: https://www.proxmox.com/en/products/proxmox-virtual-environment/overview

## Deployment

```sh
bazel run //infra/dc1/pve1/ansible # host setup
```

## Recreate the test VM

```sh
bazel run //infra/dc1/pve1/tf:tf.apply -- --replace proxmox_vm_qemu.cloudinti_test
```

## Cloud-init snippet update

```sh
bazel run //infra/dc1/pve1/ansible -- --tags pve_snippets
```

## Pve token

- Create a token with Privilege Separation
- Grant it required roles
- Create json:
  ```json
  {
    "token_id": "",
    "token_secret": ""
  }
  ```
- Write the data:
  ```sh
  bazel run //infra/dc1/pve1:vault.kv_put alwaldend.com/vault1/approles/src_infra_dc1_pve1/pve_token @"${PWD}/data.json"
  ```
