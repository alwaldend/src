---
title: forgejo
description: forgejo.alwaldend.com
tags:
  - ansible
  - xcp-ng
  - forgejo
---

## Architecture

Forgejo runs on an XCP-ng VM behind Traefik. Local Git links redirect to the
canonical Forgejo hostname; public Git links redirect through the download
server to GitHub. The two local Traefik nodes represent routes in the same proxy.
Public Forgejo access passes through the ingress server, which requires a client
certificate and opens a separate HTTPS connection to the service proxy.
See [Deployment](#deployment) for DNS, provisioning, and service configuration.
The [Mermaid source](docs/architecture.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![Forgejo architecture](docs/assets/architecture.svg)](docs/assets/architecture.svg)

## Links

- Docs: https://forgejo.org/docs/latest/admin/installation/binary/
- Config reference: https://forgejo.org/docs/latest/admin/config-cheat-sheet/

## Repository configuration

Service Terraform consumes the shared [repository catalog](../repos/README.md)
for organization-owned repositories and named administrator/developer roles.
Vault continues to own OIDC login identities and service-specific access.
See [service Terraform](tf/README.md) for identity validation, state adoption,
and the preserved automation grants.

## Deployment

Forgejo is recreated on XCP-ng through Xen Orchestra, using the
`src_infra_dc1_forgejo1` resource set provisioned by `infra/xcp_ng`.
VM provisioning authenticates to XO through Vault OIDC as Forgejo's own
AppRole. Its exact synchronized user receives the resource-set membership
and an explicit administration ACL on the existing Forgejo VM; it does not
use the shared infrastructure administrator token. These bindings must be
applied by `infra/xcp_ng/tf` before Forgejo's setup Terraform runs.

The canonical service URL is `https://forgejo.alwaldend.com`. Local
`https://git.alwaldend.com` redirects there, preserving the path and query.
Public DNS points the Git name at the Yandex download server, whose Traefik redirects
ordinary web requests to `https://github.com` with the path and query preserved.
Go discovery requests with `go-get=1` receive metadata from the download
server for `git.alwaldend.com/alwaldend/src` and its descendants, pointing to
`https://github.com/alwaldend/src.git`. Unknown discovery prefixes return 404. The
dc1 A record continues to select this Forgejo VM. The public CNAME remains
owned by this component's `dnsconfig.json` and DNS-only Terraform targets.
Vault's Forgejo URL and OIDC callback, shared AL login configuration, and runner
instance URL use the canonical hostname. Vault's certificate role retains both
names so HTTPS also works on the local Git redirect. These are desired-state
definitions; source changes alone do not apply Vault or Ansible configuration.
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
bazel run //infra/forgejo/tf_setup:tf_setup.apply
bazel run //infra/forgejo/ansible:ansible -- --skip-tags traefik,forgejo,forgejo_oidc
bazel run //infra/forgejo/ansible:ansible -- --tags traefik,forgejo,forgejo_oidc
bazel run //infra/forgejo/tf:tf.apply
```
