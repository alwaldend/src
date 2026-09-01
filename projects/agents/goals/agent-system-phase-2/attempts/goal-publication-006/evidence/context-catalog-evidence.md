# goal-publication-006 validation evidence

- `bazel_agent test //tools/agents/... //projects/goal/...` — 23/23 pass
  (compiler unit tests, schema validation, checked drift gates).
- `bazel_agent test //:buildifier_test` — pass (registry/BUILD edits).
- `git diff --check` — clean.
- `goal validate --goals-root projects/agents/goals` — valid (3 goals).
