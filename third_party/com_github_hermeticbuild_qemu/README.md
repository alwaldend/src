---
title: QEMU static distribution
description: Pinned Linux x86-64 QEMU tools for disposable integration guests
---

The [hermeticbuild distribution](https://github.com/hermeticbuild/qemu-prebuilt)
builds static QEMU system and image tools with KVM, TCG, and libslirp networking.
These are third-party builds, not binaries published by the QEMU project.
The build recipe and release attestations are published with the distribution.

Version 11.0.0.1 packages QEMU 11.0.0. The binary lock and firmware archive pin
use SHA-256 values verified against both the downloaded bytes and the
publisher's release checksum files. Only Linux x86-64 is selected. Firmware
comes from the same release as the executables. Upstream QEMU licensing and
firmware licenses remain applicable; these artifacts are not first-party
products. Source: [QEMU](https://www.qemu.org/).
