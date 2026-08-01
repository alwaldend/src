---
title: threexui
description: 3x-ui
tags:
  - xray
  - vpn
---

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
