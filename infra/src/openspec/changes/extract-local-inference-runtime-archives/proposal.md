## Why

The approved runtime archive declarations in PR #233 can be reviewed and merged
independently of the local inference host infrastructure.

## What Changes

- Extract the llama.cpp b11960 and CUDA 12.8 archives, Codex 0.166.0 archive,
  and T3 Code 0.1.6 archive into a separate master-based change.
- Retain the existing immutable upstream URLs, SRI digests, aliases, and
  repository-internal visibility without version or dependency changes.
- Generate the root module includes through the owning update workflow.

## Capabilities

No spec-level behavior changes: this separates existing dependency declarations
from their infrastructure consumer. skip_specs is true.

## Impact

Only the three third_party packages, generated root module wiring, and this
repository-owned extraction record. The existing Codex review-agent source pin
remains unchanged. PR #233 consumes these definitions after this prerequisite
merges. No deployment, installation, or first-party runtime changes are included.
