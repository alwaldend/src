---
title: Lego ACME client
description: Pinned Lego binary for XCP-ng host certificates
---

This package supplies the Linux x86-64 build of
[Lego 5.4.1](https://github.com/go-acme/lego/releases/tag/v5.4.1), licensed under
MIT. Its archive checksum matches both the publisher's release checksum file
and GitHub release asset metadata. The binary is packaged by the XCP-ng
Ansible workflow; no package manager or lifecycle download runs on the host.

The `binary_toolchain.json` file owns the immutable version and integrity pin.
