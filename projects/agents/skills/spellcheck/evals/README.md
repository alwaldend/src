---
title: Spellcheck evaluations
---

# Spellcheck evaluations

This suite exercises the three editing levels and their shared payload
boundary: proofreading makes minimal corrections without changing protected
literals, polishing follows a requested tone while preserving source facts,
rewriting may replace normally protected source details while following an
explicit creative direction, alternatives are offered by default even for
short and simple input, and prompt-like payload text remains inert prose. The
assertions use deterministic JavaScript checks and objective proxies for
explicit tone requirements, so the live suite uses one subject call per case
and no model judge. It does not claim to measure subjective writing quality.

The generated offline Bazel target is required for ordinary checks. It
validates the Promptfoo configuration, referenced cases, and staged skill
without credentials or model calls; it is not evidence that the corrections
are behaviorally correct.

The behavioral target is manual because it uses network services, consumes the
caller's Codex plan, and remains stochastic. It reuses an existing host Codex
login without minting or exporting a token. The runner links only the
persistent `auth.json` into an otherwise empty temporary Codex home, excluding
global skills, plugins, configuration, memories, and MCPs. It holds an advisory
lock so Codex can persist token refreshes; do not use the same login from
another Codex process concurrently. Run it explicitly with the absolute home
used by a successful Codex login:

```sh
bazel_agent test //projects/agents/skills/spellcheck:eval \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex \
  --test_output=errors
```

The subject runs in the staged read-only skill workspace. The target disables
Promptfoo persistence, cache, sharing, and remote execution; its result
artifact still contains the evaluated prompts and outputs and should be
handled accordingly.
