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

Both skills are staged because decision-review reads the sibling
quality-assessment skill. The separate `eval_routing_config_test` validates a
composition case that requires both skill-use signals; `eval_routing` runs it.
These routing assertions are separate from the behavioral rubrics.

The manual `//tools/agents/skills/decision-review:eval` target exercises the
behavioral cases with the staged skills and a model judge.
The subject runs in the staged skill workspace, while the judge runs in a
separate empty workspace. The runner links only `auth.json` from an explicitly
supplied Codex home into separate otherwise-empty subject and judge homes, so
global skills, plugins, configuration, memories, and MCPs are not inherited.
The target is excluded from ordinary test runs; run it only when live,
credentialed evaluation is intended:

```sh
al bazel test //tools/agents/skills/decision-review:eval \
  //tools/agents/skills/decision-review:eval_routing \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
```

Offline validation is not behavioral proof. Even the live cases do not prove
that a real review with external evidence will find the strongest objection,
assess every alternative accurately, or reach the right verdict.
