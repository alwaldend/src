---
title: org_kernel_headers
description: Linux kernel headers
---

## Links

- Site: https://www.kernel.org/
- Header install: https://www.kernel.org/doc/Documentation/kbuild/headers_install.txt
- Reference: https://github.com/DanielEliad/rules_linux_kernel_headers/tree/main

## Generate the config

```sh
zcat /proc/config.gz >third_party/org_kernel_headers/config.txt
```
