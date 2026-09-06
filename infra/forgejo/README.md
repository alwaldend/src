---
title: forgejo
description: git.alwaldend.com
tags:
  - ansible
  - xcp-ng
  - forgejo
---

## Links

- Docs: https://forgejo.org/docs/latest/admin/installation/binary/
- Config reference: https://forgejo.org/docs/latest/admin/config-cheat-sheet/

## Deployment

Forgejo is recreated on XCP-ng through Xen Orchestra, using the
`src_infra_dc1_forgejo1` resource set provisioned by `infra/xcp_ng`.
VM provisioning authenticates to XO through Vault OIDC as Forgejo's own
AppRole. Its exact synchronized user receives the resource-set membership
and an explicit administration ACL on the existing Forgejo VM; it does not
use the shared infrastructure administrator token. These bindings must be
applied by `infra/xcp_ng/tf` before Forgejo's setup Terraform runs.

The canonical service URL is `https://git.alwaldend.com`. The previous
`forgejo.alwaldend.com` name remains available for existing clients.
The VM retains `192.168.10.40` and `host1.forgejo.alwaldend.com`. Its boot, Forgejo,
and Traefik disks are 20, 40, and 5 GiB respectively.

See [Terraform setup](tf_setup/README.md) for the template, network, storage,
and previous-state prerequisites. These commands create and configure a
fresh service; they do not restore the destroyed Proxmox VM's data.

Fresh host hardening can outlast the ten-minute Vault token. Run host setup
and service deployment separately to obtain fresh authentication for each phase. Ansible
creates the missing Vault OIDC login source after starting Forgejo; existing
active Vault sources are preserved.

Before configuring Forgejo with Terraform, set `TF_VAR_vault_oauth_source_id`
to the verified ID of the active `vault` OIDC source, as described in
[service Terraform](tf/README.md). Its narrowly scoped Vault entity-read policy
must be provisioned first so the required external users can be created.

The guest disables IPv4 redirect acceptance so its traffic uses the configured
gateway. LAN clients must also retain symmetric routes: same-interface router
redirects can otherwise make connections fail. Router redirect policy and stale
client route caches must be resolved before relying on LAN access.

```sh
bazel_agent bazel run //infra/forgejo/tf_setup:tf_setup.apply
bazel_agent bazel run //infra/forgejo/ansible:ansible -- --skip-tags traefik,forgejo,forgejo_oidc
bazel_agent bazel run //infra/forgejo/ansible:ansible -- --tags traefik,forgejo,forgejo_oidc
bazel_agent bazel run //infra/forgejo/tf:tf.apply
```
