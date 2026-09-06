---
title: Cloud-init
description: Shared PVE-based infrastructure VM bootstrap configuration
---

`assets/cloud_init.yaml` is the canonical configuration moved from PVE. PVE
and Yandex Cloud consume it directly. It owns the Ansible sudo user, SSH keys,
local TLS root CA, and common bootstrap policy. The Ansible SSH CA is selected by its certificate-authority principal options;
user and key ordering do not affect the Xen configuration. Missing or duplicate
Ansible users or CA keys fail template rendering.
`assets/cloud_init_min.yaml` is the existing minimal PVE configuration with
QEMU guest tools. Both original PVE files retain their contents.

The `xen_linux` target uses the repository's `rules_template` Go template rule
to derive `xen_linux.json` from that same PVE configuration. Its small template
selects the shared Ansible CA key, uses Fedora's `users,wheel` groups, and adds
Xen guest tools. It omits PVE's explicit package-reboot and public metadata-key
flags to preserve the existing Xen bootstrap behavior. All other base fields,
including the TLS CA, come from the PVE configuration. There is no separate
copy of the Xen user or CA configuration and no Terraform template file.

Forgejo merges its hostname into the generated object and uses Terraform's
`yamlencode` to serialize it. Its static network structure stays in Terraform,
where interface, address, MAC, gateway, and DNS are runtime values. The rendered
Forgejo user-data and network configuration remain byte-identical to the
previous configuration.

PVE Ansible packages the moved snippets under the existing
`files/cloud_init.yaml` and `files/cloud_init_min.yaml` names. Those packaging
aliases preserve `local:snippets/cloud_init.yaml` references on PVE hosts.
Subsequent persistent VM configuration belongs in Ansible. Targets are
repository-internal runtime inputs, not published artifacts.
