---
title: threexui
description: 3x-ui
tags:
  - xray
  - vpn
---

## Nodes

The data-plane nodes `njalla1` and `work1` are hosted by njal.la; `yc1` is hosted
by Yandex Cloud. Terraform configures each through the shared node module;
Ansible installs the node services.

`work1.nodes.threexui.alwaldend.com` is declared at IPv4 `45.142.142.192` and
IPv6 `2a0a:3840:8078:142::2d8e:8ec0:1337` in `dnsconfig.json`.
The node uses the `njalla1` Mullvad relay selection with `de-fra-wg-001`
replacing the retired `de-dus-wg-003` relay.
Before configuring it, supply its key and address in the existing Vault-backed
`mullvad_keys` map. The `node_xui_tokens` entry is needed for later master
registration, after the node is installed. The node token comes from the `xui_tokens`
field of `alwaldend.com/vault1/approles/src_infra_threexui/nodes`, and its Mullvad
key and address come from the `keys` field of
`alwaldend.com/vault1/approles/src_infra_threexui/mullvad_keys` (mount `secrets`).

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
