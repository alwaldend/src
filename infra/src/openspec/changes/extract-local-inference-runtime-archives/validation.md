## Source acceptance

Extracted the approved runtime declarations from infrastructure candidate
05a5095b0d4626aa7a9579e3a5425c599b25b1c2 onto master
cd000ba4171baed5a566a5f7715df568f0da9386. All twelve copied dependency files
match that source byte-for-byte. Existing Codex review-agent source and
visibility remain unchanged; only its deployment archive alias and declaration
are added relative to master.

The owning module update generated exactly the two expected root includes.
Module dependency/lock reconciliation succeeded without a lockfile change.
The four archive aliases built successfully. Independent SHA-256 checks of the
fetched files match their declared SRI digests, and tar metadata confirms the
llama server/CUDA backend, CUDA runtime libraries, Codex executable, and complete
T3 binary/client assets. The repeatable non-secret result is retained under
out/extraction/archive-verification.json. No archive was installed or executed.

Thirty repository quality, buildifier, repository OpenSpec, and exact-state
skill-discovery tests passed. Semantic lint passed for the three dependency
packages and repository OpenSpec owner. Guarded delivery reruns mandatory gates
against the exact prepared candidate; receipts own publication state.

## Scope and continuation

The independent master-based PR contains only dependency definitions, module
wiring, and this record. The infrastructure PR #233 can be rebased after this
prerequisite merges, dropping the now-merged declarations without changing its
consumer labels. No infrastructure or site was deployed.

## Session ergonomics

ARCHIVE-OUTPUT-RESOLUTION (artifact-verified): cquery reports external paths
relative to Bazel's output base; pure archive aliases do not require execroot
symlinks to be materialized. Inspection used the existing output base and
retained the digest/layout report instead of creating another cache or
unpacking binaries. No subagents or new dependency versions were introduced.
