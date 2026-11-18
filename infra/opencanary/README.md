---
title: OpenCanary
---

One Fedora VM on XCP-ng runs Thinkst OpenCanary natively under systemd.
Terraform provisions the guest through Xen Orchestra. The
[OpenCanary Ansible role](../../projects/ansible_collection/roles/opencanary/README.md)
installs the service and its optional read-only Samba decoy.

The profile is defined in [group variables](ansible/group_vars/opencanary.yaml).
It enables all built-in protocol modules, including SMB auditing. The Fedora
profile leaves port-scan detection disabled because OpenCanary requires
legacy iptables for that module. Administrative SSH remains on TCP 22 and
the SSH decoy listens on TCP 2222. The role disables LLMNR in the guest
resolver so the LLMNR decoy can listen on UDP 5355.

Events appear in the guest journal under `opencanary.service`. No external
alert destination or dashboard is configured.

## Deployment inputs

- Apply the [standalone Vault AppRole](../vault/approles/src_infra_opencanary/README.md),
  reconcile its shared Vault memberships, then synchronize its identity through
  the [XCP-ng workflow](../xcp_ng/tf/README.md).
- Assign the Fedora 44 template, Local storage, and an isolated honeypot network
  to the `src_infra_opencanary` resource set. Supply that entry in the XCP-ng
  `resource_set_inventory` input alongside the existing entries.
- Supply the DHCP-enabled network as `TF_VAR_network_name`. DHCP supplies
  the guest's IPv4 address, gateway, and DNS. Use the Terraform `address`
  output for Ansible and `mac_address` if you need a DHCP reservation.
- Export `OPENCANARY_HOST`, `OPENCANARY_MANAGEMENT_SOURCE`, and
  `OPENCANARY_DECOY_SOURCE` before running Ansible. The source values are the
  management and monitored IPv4 CIDRs. Verify the guest's initial SSH host
  key out of band and trust it under the alias `opencanary.alwaldend.com`.

```sh
bazel run //infra/opencanary/tf_setup:tf_setup.plan
bazel run //infra/opencanary/tf_setup:tf_setup.apply
bazel run //infra/opencanary/ansible
```
