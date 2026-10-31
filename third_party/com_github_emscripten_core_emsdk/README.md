---
title: Emscripten SDK
description: Official pinned Emscripten compiler, libraries and Binaryen tools
---

This package pins the official Linux x86_64 SDK for Emscripten **6.0.13**.
The [upstream release map](https://github.com/emscripten-core/emsdk/blob/35659964bc2762fffec0d257997fb61475f765b7/emscripten-releases-tags.json)
identifies SDK build `eb132757325583b787fe9414b09afea9442b72ad`.
It was the latest stable release observed on 2026-10-24.

The archive URL, extraction prefix and SHA-256 are published in the release's
[repository rule](https://github.com/emscripten-core/emsdk/blob/35659964bc2762fffec0d257997fb61475f765b7/bazel/emscripten_deps.bzl)
and [revision table](https://github.com/emscripten-core/emsdk/blob/35659964bc2762fffec0d257997fb61475f765b7/bazel/revisions.bzl).
`include.MODULE.bazel` owns the archive integrity and immutable build URL.
The 311,791,988-byte download was independently hashed against the published
checksum before use. Bazel verifies it again when fetching the repository.

The archive provides matching Emscripten, LLVM, Binaryen and prebuilt system
libraries. Repository Python and Node toolchains execute the compiler; no SDK
activation, global installation, npm packages or execution-time downloads are
needed. The dedicated toolchain does not register a C++ or Rust toolchain.
Only the consumed Linux x86_64 execution platform is declared.

Emscripten is MIT-licensed; the SDK bundles components with their own licenses,
including LLVM's Apache-2.0 with LLVM exceptions and Binaryen's Apache-2.0.
Retain upstream license notices in generated output and identify the borrowed
implementation as Emscripten. See [upstream licenses](https://github.com/emscripten-core/emscripten/blob/6.0.13/LICENSE).
