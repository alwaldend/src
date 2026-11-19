---
title: threexui
description: 3x-ui
tags:
  - xray
  - vpn
---

## Architecture

The [architecture source](docs/architecture.mmd) is one diagram with three separate
sections: overview, control plane, and data plane. It describes maintained
architecture, not live inventory or service health. Bazel renders light and dark SVGs.

The overview shows Master in the Control Plane hierarchy and full node hostnames
in provider groups, with Vault, Let's Encrypt, and Mullvad as external services.
Overview links attach to the Control Plane and Data plane boundaries.
The control-plane section expands the 3x-ui VM's Traefik and panel with one simple
data-plane node. The data-plane section expands one node's Traefik, 3x-ui, and Xray
with one simple control-plane endpoint. Both njalla1 and yc1 share this implementation.

[![3x-ui architecture](docs/assets/architecture.svg)](docs/assets/architecture.svg)

## Nodes

The data-plane node `njalla1` is hosted by njal.la. Yandex Cloud hosts `yc1`.
Terraform configures both through the shared node module, and Ansible
installs their services.

## Links

- Site: https://docs.sanaei.dev/
- Releases: https://github.com/MHSanaei/3x-ui/releases

## Fetch and fix subs for a particular subscription id

```sh
bazel run //infra/threexui:fix_subs -- --hosts njalla1.nodes.threexui.alwaldend.com,yc1.nodes.threexui.alwaldend.com --sub_id subid
```

## Fix subs from a local file

```sh
bazel run //infra/threexui:fix_subs -- --sub_file path_to_file
```
