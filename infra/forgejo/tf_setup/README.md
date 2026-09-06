---
title: Tf setup
description: Terraform setup
tags:
  - terraform
  - xcp-ng
---

This package creates a Xen Orchestra VM in the `src_infra_dc1_forgejo1`
resource set. `infra/forgejo/al.lua` authenticates with Forgejo's own Vault
AppRole and the packaged XO OIDC login plugin. The plugin supplies a temporary
XO token and revokes it on shutdown; no infrastructure administrator token
is loaded. The setup HTTP backend remains owned by the same Forgejo config.

Before running this package, bootstrap the AppRole's XO OIDC user and apply
its resource-set membership and existing VM ACL through `infra/xcp_ng/tf`.
Subjects are matched by immutable Vault entity UUID under the configured
OIDC issuer, not by login name or AppRole group membership. See
[XO authentication](../../xcp_ng/cmd/xo_login/README.md).

Defaults select the imported Fedora 44 template, local storage and wired
network by name from the Forgejo resource set. Override
`TF_VAR_xoa_template_name`, `TF_VAR_xoa_storage_name`, and
`TF_VAR_xoa_network_name` for another assignment. Native provider lookups run
as Forgejo's own identity and reject ambiguous matches. No infrastructure
UUID defaults are required; resolved IDs are passed to the XO API.
The template must contain a single boot disk no larger than 20 GiB, Fedora,
and cloud-init. Bootstrap installs Xen guest tools. The pinned Fedora 44
image uses predictable interface naming; `xoa_guest_interface` defaults to
`enX0` for its first Xen interface and can be overridden for another template.
The selected network must carry `192.168.10.0/24`; gateway and DNS default to
`192.168.10.1` (the wired router) and can be overridden with the Terraform variables.

The PVE-based `infra/cloud_init:xen_linux` target supplies the shared Ansible
user and CA configuration through the repository template rule. Terraform adds
the hostname and static network configuration from `../dnsconfig.json`. After cloud-init completes, verify that `/dev/xvda` is the
boot disk, `/dev/xvdb` is the 40 GiB Forgejo disk, and `/dev/xvdc` is the 5 GiB
Traefik disk before running Ansible, which creates filesystems on the latter
two devices.

The old Proxmox VM was destroyed and its obsolete Terraform state binding
was removed without issuing a destroy. This package now uses only XO.
There is no cross-provider state move or data restoration in this package.

```sh
bazel_agent bazel run //infra/forgejo/tf_setup:tf_setup.plan
bazel_agent bazel run //infra/forgejo/tf_setup:tf_setup.apply
```
