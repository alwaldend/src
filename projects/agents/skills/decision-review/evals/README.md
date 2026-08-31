---
title: Decision-review evaluations
---

# Decision-review evaluations

This suite describes the behavioral contract for adversarially reviewing a
material choice before commitment. Its offline Bazel target validates the
Promptfoo configuration, referenced cases, and staged skill without making a
model call. The cases cover separating an outcome from its proposed method,
seeking disconfirming evidence, comparing credible alternatives, issuing an
explicit verdict, preserving authorization boundaries, reopening a repeatedly
failing strategy, and avoiding ceremony for routine reversible choices.

The manual `//projects/agents/skills/decision-review:eval` target exercises the
self-contained retention-window case with the staged skill and a model judge.
The subject runs in the staged skill workspace, while the judge runs in a
separate empty workspace. The runner links only `auth.json` from an explicitly
supplied Codex home into separate otherwise-empty subject and judge homes, so
global skills, plugins, configuration, memories, and MCPs are not inherited.
The target is excluded from ordinary test runs; run it only when live,
credentialed evaluation is intended:

```sh
bazel_agent test //projects/agents/skills/decision-review:eval \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
```

Offline validation is not behavioral proof. Even the live cases do not prove
that a real review with external evidence will find the strongest objection,
assess every alternative accurately, or reach the right verdict.
