# Attempt history

[Back to durable goal](../)

- [Attempt 11](011.md): consolidate the owned range, reconcile the advanced
  base, validate and publish the exact rebased candidate, then incorporate and
  reconcile the final hosted review. Complete.
- [Attempt 10](010.md): keep the thin wrapper and patch pinned Cordis HMR to
  serialize reloads and drain writes arriving during an in-flight reload.
  Complete locally and carried into Attempt 11.
- [Attempt 9](009.md): simplify source updates to validated atomic persistence
  plus Cordis HMR, deleting the private acknowledgement transaction. Refined
  after independent review reproduced a lost overlapping update.
- [Attempt 8](008.md): correct the final hosted review's fallback filtering,
  byte-offset, and UTF-8 body-preview findings. Complete locally.
- [Attempt 7](007.md): replace the custom manifest/version store and worker
  generations with official Cordis Loader, Include, HMR, standard
  `cordis.yaml`, and normal modules. Delivered, then refined by Attempt 8.
- [Attempt 6](006.md): make bounded execution backward-compatible, close
  process/lifecycle admission races, publish exact package completeness, and
  replace worker-side spawning after independent review. Rejected because its
  package persistence model was custom rather than standard Cordis.
- [Attempt 5](005.md): make every output-loss signal and bounded Git result
  exact, then cover the imported skill behaviors. Refine after independent
  review found lifecycle, compatibility, and policy defects.
- [Attempt 4](004.md): import PR 24's scoped agent guidance, then correct the
  three valid PR 32 review findings. Refine after independent review.
- [Attempt 3](003.md): replace a lifecycle test's elapsed-time inference with
  a deterministic started/release handshake. Published, then superseded by
  review findings.
- [Attempt 2](002.md): retained the proven runtime and added a bounded search
  fallback for hermetic portability. Starter packages pass; refine because
  the integrated lifecycle evidence was timing-dependent.
- [Attempt 1](001.md): direct Cordis runtime with MCP v2, worker-isolated
  generations, two-tier immutable persistence, and three starter packages.
  Rejected because one starter package required an unavailable executable.
