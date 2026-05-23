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
bazel run //infra/dc1/pve1/tf/tf.apply # tf setup
```

## Recreate the test VM

```sh
bazel run //infra/dc1/pve1/tf:tf.apply -- --replace module.vm_cloudinit_test.proxmox_vm_qemu.vm
```

## Cloud-init snippet update

```sh
bazel run //infra/dc1/pve1/ansible -- --tags pve_snippets
```

## Update ACME account

- Create an EAB:
  ```sh
  bazel run //infra/dc1/vault -- write -f pki/ica_servers/roles/src_infra_dc1_pve1_pki_server/acme/new-eab
  ```
- Login as `root@pam`
- Go to Datacenter -> ACME
- Create a new account with the EAB and a directory:
  ```
  https://vault.dc1.alwaldend.com:8200/v1/pki/ica_servers/roles/src_infra_dc1_pve1_pki_server/acme/directory
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
