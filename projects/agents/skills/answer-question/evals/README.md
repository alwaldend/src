---
title: Answer-question evaluations
---

# Answer-question evaluations

This Promptfoo suite checks conformance with the reasoning behaviors promised
by the `answer-question` skill. It covers premise testing, calibrated positive
and negative causal conclusions, evidence independence and applicability,
repository versus deployment state, conditional-action answer behavior,
quantitative checks, criteria-based recommendations and critique, and
proportional treatment of a simple stable question. It is not a general
factual-accuracy benchmark.

The Bazel rule creates a temporary workspace containing only the declared
skill under `.agents/skills/`. The subject Codex agent receives one request at
a time; the case's grading criteria are not included in its prompt or staged
workspace.
A separately invoked Codex judge, running in an empty workspace without the
skill, evaluates one consolidated binary checklist per non-deterministic case.
A threshold of one prevents a weak score from passing through averaging, while
consolidation keeps the smoke run to roughly nineteen model calls instead of
more than one hundred. The selected model name does not make the judge
immutable or independent of the subject's model family, so human calibration
still matters.

The separate routing config requires Promptfoo's `skill-used` signal.
Promptfoo currently infers Codex skill use from a successful direct read of
`SKILL.md` rather than a first-class invocation event. Routing is kept out of
the quality configs so the skill and no-skill pass rates use identical
assertions. The conditional-action case evaluates only the final answer; the
read-only sandbox means this suite does not prove that no write was attempted.

The skill and no-skill smoke targets run each quality case once, while the
routing target runs its single routing case once. Paired skill and no-skill
stability targets override the quality suite to three repetitions. All targets
disable caching and persisted threads and use concurrency one, per-row and
total time bounds, isolated subject and judge workspaces, and read-only Codex
sandboxes. They reuse an explicitly supplied local Codex login without minting
or exporting a token. The runner links only the persistent `auth.json` into
distinct subject and judge Codex homes under a mode-0700 temporary directory;
it does not carry over global skills, plugins, configuration, memories, or
MCPs. It holds an advisory lock for the test so Codex can safely persist token
refreshes. Do not use the same login from another Codex process concurrently.
These controls reduce accidental state sharing; they are not a confidentiality
boundary around the host.

Compare the result artifacts from the skill-present and no-skill targets before
claiming improvement. For a release decision, also compare the previous skill
revision and have humans label a representative output sample. Inspect judge
disagreements, false positives, and per-case stability rather than treating one
aggregate score as ground truth. The suite should remain manual until those
baselines establish an acceptable false-reject rate.

Live targets are manual and local because they consume the caller's Codex plan,
use network services, depend on an existing Codex login, and are stochastic.
Pass the absolute Codex home used by a successful `codex login`:

```sh
bazel_agent test //projects/agents/skills/answer-question:eval \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
bazel_agent test //projects/agents/skills/answer-question:eval_no_skill \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
bazel_agent test \
  //projects/agents/skills/answer-question:eval_no_skill_stability \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
bazel_agent test //projects/agents/skills/answer-question:eval_routing \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
bazel_agent test //projects/agents/skills/answer-question:eval_stability \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
```

The config validation targets are offline and safe for ordinary repository
checks:

```sh
bazel_agent test //projects/agents/skills/answer-question:eval_config_test
bazel_agent test \
  //projects/agents/skills/answer-question:eval_no_skill_config_test
bazel_agent test \
  //projects/agents/skills/answer-question:eval_routing_config_test
```
