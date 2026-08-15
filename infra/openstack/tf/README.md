---
title: OpenStack MikroTik Terraform
description: Static LAN addresses for the OpenStack hosts
---

This Terraform package reads the host addresses and MAC addresses from
`../ansible/inventory.yaml`.

It creates permanent `bridge1` DHCP leases for:

- `master1.openstack.alwaldend.com` at `192.168.1.2`.
- `compute1.openstack.alwaldend.com` at `192.168.1.3`.

Set each host's `openstack_mac_address` before applying:

```sh
bazel run //infra/openstack/tf:tf.plan
bazel run //infra/openstack/tf:tf.apply
```

DNS is managed by the repository DNSControl deployment, not by this Terraform
state.
