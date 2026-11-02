---
title: llama.cpp
---

Installs caller-supplied, Bazel-verified llama.cpp and CUDA archives on Linux
x86-64, verifies and stages a GGUF, and manages a loopback-only systemd server.
The caller owns archive pins and packaging, driver provisioning, and GPU/image
requirements. No runtime archive or model is downloaded by this role.

Supply `llama_archive_src`, `llama_cuda_archive_src`, an absolute controller-side
`llama_model_src`, and `llama_model_sha256`. Both runtime archives must have one
top-level directory; extraction strips that level and places the libraries
beside llama-server. The role verifies shared-library resolution and CUDA device
discovery before starting inference. It creates a separate non-login account,
keeps the runtime administrator-owned, and makes the model group-readable.

Changed inputs notify a restart; readiness follows pending restarts. Repeated
unarchive comparison repairs interrupted extraction and preserves unchanged
installations. The role contains no Codex provider or Responses bridge logic.
