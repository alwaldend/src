## 1. Upgrade tooling

- [x] 1.1 Upgrade the approved Gazelle pin and reconcile generated module state; inspect the resolution diff.
- [x] 1.2 Build the repository Gazelle binary and a representative Go consumer, inspecting generated dependency outputs.

## 2. Verify reuse

- [x] 2.1 Run equivalent Go repository fetches from fresh task-owned worktrees and preserve observable evidence of completed-repository cache hits or misses.
- [x] 2.2 Verify the selected host cache policy with the same fetched module and record miss behavior.
- [x] 2.3 Run strict OpenSpec validation, repository quality, and affected semantic lint before delivering the regular PR.
