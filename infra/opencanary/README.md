---
title: OpenCanary
date: 2026-11-18
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

## Network architecture

[![OpenCanary network diagram](docs/assets/architecture.svg)](docs/assets/architecture.svg)

The [Mermaid source](docs/architecture.mmd) shows public ingress and the local
access path labeled with the decoys in the [guest profile](ansible/group_vars/opencanary.yaml).
Traefik requires a trusted client certificate and uses HTTPS over WireGuard
to reach the web decoy. Traefik verifies the Vault-issued backend certificate
for `opencanary.alwaldend.com`.

## Deployment inputs

- Apply the [standalone Vault AppRole](../vault/approles/src_infra_opencanary/README.md),
  reconcile its shared Vault memberships, then synchronize its identity through
  the [XCP-ng workflow](../xcp_ng/tf/README.md).
- The Xen Orchestra resource set uses the regular VM network. The
  [DNS declaration](dns/dnsconfig.json) owns the static IPv4 address and the
  [DNS root](dns/README.md) publishes the name in the site-local view.
- Management SSH accepts connections from any routed source. Decoy traffic
  defaults to the guest subnet. The ingress WireGuard peer may also reach the
  HTTPS decoy. Verify the guest SSH host key out of band before connecting.

```sh
bazel run //infra/opencanary/tf_setup:tf_setup.plan
bazel run //infra/opencanary/tf_setup:tf_setup.apply
bazel run //infra/opencanary/dns:tf.apply
bazel run //infra/opencanary/ansible
```

For an existing HTTP upstream, apply the AppRole changes, then use these
Ansible stages to switch the route without retiring its HTTP access early:

```sh
bazel run //infra/opencanary/ansible -- --tags opencanary_tls
bazel run //infra/opencanary/ansible -- --tags firewall --extra-vars '{"firewall_delete_rules":[]}'
bazel run //infra/ingress/ansible -- --tags traefik
bazel run //infra/opencanary/ansible -- --tags firewall
```

Verify the HTTPS listener from ingress before switching Traefik. The guest
keeps its private key. Run the `opencanary_tls` stage weekly to replace the
30-day certificate when fewer than seven days remain.
