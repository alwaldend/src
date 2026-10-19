---
title: WASI sysroot
description: Pinned WASI C headers and libraries for the repository LLVM toolchain
---

This package acquires only the
[WASI SDK 34.0 sysroot](https://github.com/WebAssembly/wasi-sdk/releases/tag/wasi-sdk-34),
verified by the upstream release SHA-256 digest. The existing repository-pinned
LLVM compiler consumes its headers and libraries for `wasm32-wasip1`; no host SDK
installation is performed. The archive retains upstream license notices.

`tools/cc/include.MODULE.bazel` owns compiler/toolchain integration. The initial
consumer is the SQLite WASI acceptance fixture in
`projects/anki_grpc/test/wasm`.
