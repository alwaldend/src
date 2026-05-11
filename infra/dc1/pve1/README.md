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
bazel run //infra/dc1/pve1/ansible:deploy # host setup
```
