---
title: consul1
description: Consul setup for consul1.dc1.alwaldend.com
tags:
  - ansible
  - pve
  - consul
---

## Links

- Docs: https://developer.hashicorp.com/consul/tutorials/production-vms/deployment-guide
- Docs: https://developer.hashicorp.com/consul/tutorials/get-started-vms/virtual-machine-gs-deploy

## Deployment

```sh
bazel run //infra/dc1/pve1/tf_setup:tf.apply # Create hosts
bazel run //infra/dc1/pve1/ansible # Configure hosts
bazel run //infra/dc1/pve1/tf:tf.apply # Configure Consul
```

## Generate a gossip key

```sh
bazel run //infra/dc1/consul1:gen_gossip_key
```
