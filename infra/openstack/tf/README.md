---
title: OpenStack MikroTik Terraform
description: Static LAN addresses and DNS for the OpenStack hosts
---

This Terraform package reads the host addresses from
`../ansible/inventory.yaml` and the API VIP from `../kolla/globals.yml`.

It creates:

- A static DHCP lease for `master1.openstack.alwaldend.com`.
- A static DHCP lease for `compute1.openstack.alwaldend.com`.
- MikroTik A records for both hosts and `openstack.alwaldend.com`.

Set each host's `openstack_mac_address` in the inventory before applying.

```sh
bazel run //infra/openstack/tf:tf.plan
bazel run //infra/openstack/tf:tf.apply
```
