---
title: llama.cpp deployment archives
---

Bazel fetches the existing approved Linux x86-64 deployment archives from
[upstream](https://github.com/ggml-org/llama.cpp/releases). Exact versions, URLs, and SRI digests
are owned by `include.MODULE.bazel`. The deployment packages the original
archives for controller-side Ansible extraction, including supporting libraries
and assets. These are upstream artifacts, not first-party binaries.
