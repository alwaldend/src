---
title: Images
description: Deploy of VM images
---

## Deploy ISO to Proxmox

```sh
bazel run //third_party/images:deploy_proxmox
```

## Deploy ISO to Yandex Cloud

```sh
bazel run //third_party/images:deploy_yc
```
