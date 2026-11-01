---
title: Qwen3.8-27B model inputs
description: Integrity-checked BF16 model snapshot for self-hosted inference
---

[Qwen3.8-27B](https://huggingface.co/Qwen/Qwen3.8-27B) is an upstream
Qwen vision-language model with coding, reasoning, and tool-use support.
The weights and accompanying files retain their upstream Apache-2.0 license
and ownership. This package provides model inputs; a compatible inference
runtime and an agent client are supplied separately.

`include.MODULE.bazel` owns the immutable upstream revision and SHA-256
integrity values. All 18 BF16 weight shards, the shard index, tokenizer,
chat template, model and processor configuration, upstream README, and
license come from that same revision. Weight checksums use the publisher's
Hugging Face LFS SHA-256 object identifiers; small files are hashed after
verification against their upstream Git blob identifiers.

## Targets

- `:metadata` downloads the configuration, tokenizer, index, and upstream
  documentation without the weight shards.
- `:weights` provides the 18 weight shards.
- `:model` provides every file required for the upstream model snapshot.
- `:snapshot` creates a flat, uncompressed model archive with original upstream
  filenames, suitable for extraction into a local model directory.

Model targets are marked `manual` to exclude these large downloads from
wildcard builds. Build an explicit target to fetch the model:

```sh
bazel_agent bazel build \
  //third_party/co_huggingface_qwen_qwen3_8_27b:snapshot
mkdir -p out/qwen3_8_27b/model
tar -xf bazel-bin/third_party/co_huggingface_qwen_qwen3_8_27b/qwen3.8-27b.tar \
  -C out/qwen3_8_27b/model
```

The complete snapshot is about 55.6 GB. Allow additional disk space for the
Bazel download cache, archive, and extracted directory. The weight data alone
occupies about 55.6 GB (51.8 GiB); inference also needs GPU memory for runtime
buffers and the context cache. An 80 GB GPU is a reasonable starting point
for one agent, but context length and actual performance need validation on
the chosen runtime and GPU. Start with a bounded context such as 32K tokens.

Point an already provisioned, compatible runtime at the extracted directory
and enable its offline mode to prevent fallback downloads. Follow the
[upstream serving guidance](https://huggingface.co/Qwen/Qwen3.8-27B#serving-qwen38)
for architecture and tool-call support. This package does not provision a GPU,
install a serving runtime, start a service, or change host configuration.
