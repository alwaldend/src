## Context

The user authorized adding Qwen3.8-27B after discussing a single 80 GB GPU.
The isolated checkout is on `t3/gpu-vram-ranking`. Existing model packages
use `http_file` with SRI, manual targets, and upstream attribution.

## Goals / Non-Goals

Provide reproducible BF16 model inputs and an optional flat snapshot archive.
Runtime integration, inference benchmarking, and deployment are separate work.

## Decisions

- Follow the existing Hugging Face package pattern and reuse `rules_pkg` for
  an uncompressed snapshot archive; add no serving or packaging dependency.
- Pin the public, ungated upstream model to revision
  `1d4bf0f2ff6012fd82039f2fa52739d0dd7c60c0`, observed on 2026-11-01.
- Keep the model pin and SRI declarations in `include.MODULE.bazel`.
  The archive preserves upstream filenames for offline directory loading.
- Include 18 shards and 12 supporting files, including license and model card.
  The index references all 18 shards and 1,199 tensors.
- Verify downloaded support files against upstream Git blob or LFS identifiers
  before computing SHA-256. Weight checksums are publisher LFS identifiers;
  full weight contents are not downloaded during lightweight validation.
- Keep all download and packaging targets manual. Validate metadata without
  requiring 55.6 GB of downloads or an available GPU.

## Risks / Trade-offs

BF16 weights require about 51.8 GiB before runtime and context-cache overhead.
Fitting in VRAM does not establish agent quality, supported context, or speed.
The snapshot archive requires substantial extra disk space. Document runtime
compatibility and offline use without adding unrequested dependencies.

## Evidence and continuation

The publisher file identities and shard-index coverage have been verified;
raw evidence is under ignored `out/qwen38`. Full weight download and GPU
inference remain unperformed. The user approved existing Bazel cache access
on 2026-11-01, and the pinned OpenSpec wrapper now runs successfully.
The package is installed in source. Root module regeneration adds one include;
module reconciliation leaves the lock unchanged. Bazel builds metadata and
upstream documentation successfully. All 12 fetched support files match the
verified SHA-256 values, the shard index covers all 18 shards and 1,199 tensors,
and repository formatting leaves no changes. The bounded receipt is
`out/qwen38/bazel-output-verification.json`.
The prepared implementation passes all 120 selected repository quality and
OpenSpec tests, metadata build checks, and semantic lint. Guarded publication
created PR 230 and verified the remote head. Final archival changes receive
new exact-candidate validation before the final push. Git and local delivery
receipts own commit and publication state; no GPU inference is claimed.

Session learning: `BAZEL-CACHE-ACCESS-001` identifies the premature escalated
retry. Diagnose sandbox and ownership evidence before seeking access; do not
route around a rejection. This record proposes no shared policy changes.
