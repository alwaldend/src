---
title: Emscripten build toolchain
description: Hermetic JavaScript generation using the pinned official SDK
---

`pkg/defs.bzl` and `:toolchain_type` are public build APIs consumed by production
targets. They invoke the official SDK with declared Python and Node toolchains,
without registering or replacing the repository's C++ or Rust toolchains.
Release inputs belong to
[`third_party/com_github_emscripten_core_emsdk`](../../third_party/com_github_emscripten_core_emsdk/README.md).

`emscripten_javascript` compiles C translation units to its `js` output. With
`wasm = True` it additionally declares the sibling `.wasm` runtime artifact;
`wasm = False` selects upstream's documented JavaScript-only mode. Caller-owned
`arguments` select the module and exported runtime API. Flags generating other
sidecars (source maps, preload files, workers) are outside this narrow rule's
contract and must not be used. All runtime outputs are returned as DefaultInfo
files and runfiles.

Compiler actions are sandboxed without network access. Configuration paths are
resolved inside the execution root; SDK system libraries are frozen inputs,
and temporary files remain inside the action output directory. No host compiler,
host Python/Node, SDK activation, cache mutation or package installation is used.
The initial supported execution platform is Linux x86_64.
