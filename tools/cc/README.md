---
title: Cc
description: Cc
---

`//:refresh_compile_commands` uses the pinned Hedron extractor. Its archive is
declared through `use_repo_rule` and fetched when the extractor's package is
loaded, rather than during module resolution. LLVM toolchain registration
remains automatic.

## Links

- MODULE.bazel config: https://github.com/bazel-contrib/toolchains_llvm/blob/2f4cff8c7d166cdea0fb9facd1caa06aaf894b06/tests/MODULE.bazel#L222
- Custom sysroot: https://steven.casagrande.io/posts/2024/sysroot-generation-toolchains-llvm/
