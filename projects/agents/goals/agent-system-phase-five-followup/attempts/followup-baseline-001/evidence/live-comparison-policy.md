# Live comparison policy

Live stochastic skill comparisons remain manual or scheduled, outside
ordinary wildcard tests. This policy binds every exact identity required by
Phase 5A.

## Bound identities

| Identity | Binding |
| --- | --- |
| Subject model | `gpt-5.6-sol` with `model_reasoning_effort: high` (set in each `promptfooconfig.yaml`) |
| Subject skill | The exact staged `SKILL.md` and skill resources under `.agents/skills/<id>/` in the isolated `PROMPTFOO_SKILL_WORKSPACE` |
| Catalog | `tools/agents/catalogs/capability.json`, digest-bound in the coverage matrix |
| Fixture cases | `evals/cases.yaml` per skill, referenced by `tests: file://cases.yaml` |
| Judge | A separately invoked Codex judge in an empty workspace (`PROMPTFOO_JUDGE_WORKSPACE`), without the skill, using one consolidated binary checklist per case |
| Subject auth | An isolated `PROMPTFOO_SUBJECT_CODEX_HOME`; the caller passes `CODEX_HOME` explicitly |

## Isolation and execution controls

- Live targets are `bazel_agent test` targets run manually with explicit
  `--test_env` overrides for `CODEX_HOME` and `CODEX_PATH_OVERRIDE`.
- Each live run is bounded by `maxEvalTimeMs`, `maxConcurrency: 1`,
  `--no-cache`, `--no-write`, and `--no-share`.
- Skill-present, no-skill, routing, and stability targets are separate
  configurations (`eval_no_skill`, `eval_routing`, `eval_stability`,
  `eval_no_skill_stability`), so paired pass rates compare like with like.
- Release decisions additionally compare the previous skill revision and
  require human labeling of a representative output sample.
- The offline `eval_config_test` targets validate configuration, cases, and
  skill staging without model calls; these are the only targets allowed in
  ordinary wildcard test runs.

## Prohibition

No ordinary `bazel_agent test //...` or wildcard test target may invoke a
live model. This policy and the offline configuration-validation targets are
the only checked-in live-comparison artifacts.
