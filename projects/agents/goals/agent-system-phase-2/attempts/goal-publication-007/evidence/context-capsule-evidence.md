# goal-publication-007 validation evidence

- `bazel_agent test //tools/agents/... //projects/goal/...` — 27/27 pass
  (context-capsule API/CLI tests, control kernel/status tests, catalog drift
  gates, goal store tests).
- `bazel_agent test //:buildifier_test` — pass (new BUILD files for
  `agent_system`, `control`, `control_status`).
- `phase1_check` — criteriaRevision 4, registry digest
  `sha256:b6dbab50076708d39717f0bb57210d9d22e4035f34ddf598da9ce438ae6bdba8`
  bound to baseline; counts: authorities 7, directBinaries 6, generated
  artifacts 5, goals 3, operations 24, projects 28, runtimeTools 10,
  skills 22, workspaces 9; unclassified none;
  requiresMigration: terraform.direct + terraform.state (pre-existing).
- `agent_system --workspace-root $PWD --path projects/agents --json` —
  bounded capsule, `byteSize: 2613`, `Completeness: complete`, 26
  capabilities, zero discovery calls.
- `git diff --check` — clean.
