---
title: OpenStack MikroTik Terraform
description: Static LAN addresses for the OpenStack hosts
---

This Terraform package creates permanent MikroTik DHCP leases for the two
OpenStack hosts.

It reads:

- Hostnames from `../ansible/inventory.yaml`.
- Addresses and MAC addresses from
  `../ansible/group_vars/openstack_master.yaml` and
  `../ansible/group_vars/openstack_compute.yaml`.

DNS is managed separately by DNSControl through
`../dnsconfig.json`.

Set each host's `openstack_mac_address` before applying.

```sh
bazel run //infra/openstack/tf:tf.plan
bazel run //infra/openstack/tf:tf.apply
```
