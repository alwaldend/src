---
title: Tf setup
description: Terraform setup
tags:
  - terraform
  - xcp-ng
---

This package creates Xen Orchestra VMs in the `src_infra_openhands` resource
set for the three OpenHands components that run off host-bot: Agent Canvas,
the agent server, and the automation server. `infra/openhands/al.lua`
authenticates with OpenHands' own Vault AppRole and the packaged XO OIDC login
plugin. The plugin supplies a temporary XO token and revokes it on shutdown;
no infrastructure administrator token is loaded. The setup HTTP backend
remains owned by the same OpenHands config.

Before running VM setup, bootstrap the AppRole's XO OIDC user and apply
its resource-set membership through `infra/xcp_ng/tf`. Subjects are matched by
immutable Vault entity UUID under the configured OIDC issuer, not by login
name or AppRole group membership. See
[XO authentication](../../xcp_ng/cmd/xo_login/README.md).

`vms` holds the component definitions, and each entry names the `dnsconfig`
record prefix that owns its hostname and address. Adding a component means
adding a record to `../dnsconfig.json` and an entry here; the guest hostname,
static address, and derived MAC all come from that one record. The agent
server is sized larger because it stores every backend conversation,
workspace, and bash event.

Defaults select the imported Fedora 44 template, local storage and wired
network by name from the OpenHands resource set. Override
`TF_VAR_xoa_template_name`, `TF_VAR_xoa_storage_name`, and
`TF_VAR_xoa_network_name` for another assignment. Native provider lookups run
as OpenHands' own identity and reject ambiguous matches. The selected network
must carry `192.168.10.0/24`; gateway and DNS default to `192.168.10.1` (the
wired router) and can be overridden.

The PVE-based `infra/cloud_init:xen_linux` target supplies the shared Ansible
user and CA configuration through the repository template rule. Terraform adds
the hostname and static network configuration from `../dnsconfig.json`. After
cloud-init completes, verify that `/dev/xvda` is the boot disk before running
Ansible.

```sh
bazel_agent bazel run //infra/openhands/tf_setup:tf_setup.plan
bazel_agent bazel run //infra/openhands/tf_setup:tf_setup.apply
```

Use this package's `dns.plan`, `dns.show`, and `dns.apply` targets for the
[scoped DNS adoption workflow](../../dns/README.md).
They select `dns=1` and target `module.dns` in the same root, AppRole, and
backend without starting XO login. Apply requires the reviewed saved plan;
the [shared procedure](../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md#prepare-the-owner)
owns preparation and reconciliation.

`dns_enabled` defaults to `true` after
[verified initial provisioning](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md).
Keep it enabled to retain the managed records. Unchanged declarations and
provider records produce a no-change DNS plan. Successful DNS checks do not
establish VM, OpenHands service, or TLS readiness; the separate
[service deployment change](../openspec/changes/add-openhands-deployment/design.md)
retains that acceptance work.
