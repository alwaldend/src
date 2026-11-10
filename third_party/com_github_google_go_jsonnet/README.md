---
title: Go Jsonnet
description: Pinned Jsonnet compiler and formatter
---

The [binary lock](binary_toolchain.json) pins Google's
[go-jsonnet v0.22.0 release](https://github.com/google/go-jsonnet/releases/tag/v0.22.0)
and its publisher-provided SHA-256. The approved Linux x86_64 archive supplies
jsonnet and jsonnetfmt; only those tools are exposed. Upstream retains its
Apache-2.0 license and ownership. [Repository integration](../../tools/jsonnet/README.md)
uses the existing binary toolchain extension; no additional rules package is added.
