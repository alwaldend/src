## Why

The user requested Qwen3.8-27B for self-hosted agent inference. The repository
already pins Hugging Face model inputs, but does not provide this model.

## What Changes

- Add the explicitly requested Qwen3.8-27B model under
  `third_party/co_huggingface_qwen_qwen3_8_27b` using existing `http_file` rules.
- Pin every BF16 shard and support file to one upstream revision with SHA-256
  integrity, preserving upstream filenames and license.
- Expose manual metadata, weights, complete-model, and snapshot-archive targets.
- Regenerate the root module include and document explicit downloads and
  offline use with a separately provisioned compatible inference runtime.

## Capabilities

### New Capabilities

None. This is a model dependency under the existing reproducibility contract.

### Modified Capabilities

None. Use `skip_specs: true` for the dependency integration.

## Impact

The explicit complete-model target downloads about 55.6 GB. Wildcard builds
must not fetch weight shards. No serving dependency, host installation,
service, GPU provisioning, or deployment is included. Acceptance requires
verified upstream file identities, complete shard-index coverage, metadata
build output inspection, formatting and build validation, and delivery.
