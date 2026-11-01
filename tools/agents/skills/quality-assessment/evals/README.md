---
title: Quality-assessment evaluations
---

# Quality-assessment evaluations

`eval_config_test` validates the behavioral configuration, fixtures, and staged
skill offline. The cases cover unsupported praise, warranted positive judgment,
numerical and causal accuracy, correction of a prior answer, source independence
and relevance, missing verification, and subjective fitness for purpose.

`eval_routing_config_test` validates the separate routing configuration. Its
cases use Promptfoo's skill-use signal for "is this good?" and "is this accurate?"
and reject activation for a quoted spelling-only payload. Answer quality and
skill routing are separate observations.

The manual `eval` and `eval_routing` targets run those suites with the isolated
Codex provider described in [rules_promptfoo](../../../../rules_promptfoo/README.md).
Behavioral judgments run in an empty workspace; the subject receives this skill.
Use the credential and scratch isolation requirements from that owning README:

```sh
al bazel test //tools/agents/skills/quality-assessment:eval \
  //tools/agents/skills/quality-assessment:eval_routing \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
```

The cases use self-contained public fixtures. They do not exercise live research
or domain-specific tool use and cannot establish the quality of real-world source
selection. Offline validation is not evidence of successful answers or routing.
Live judge results are regression evidence, not proof; no improvement over an
unskilled or previous-version control is claimed by this suite.
