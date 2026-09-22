# Final validation evidence

## Exact subject

All source and rendering checks below apply to the 20-file candidate manifest
`d6e29c5644ca9f834f3a70d2856f3069673d88a19f59ee38dee599a06c40c8ca`.
The only later source mutation in this checkpoint is the goal tool's canonical
attempt closure, which receives a separate post-closure validation.

## Source, formatting, and links

- `sha256sum --check out/repo_agent_system/final-candidate-manifest.txt`:
  **PASS**, 20 of 20 files reported `OK`.
- `git diff --check`: **PASS**.
- Repository-pinned Prettier `--check` over all 15 changed maintained Markdown
  files: **PASS** after a focused four-file formatter pass.
- Local Markdown link checker: **PASS**, 40 local links across 15 maintained
  files resolve without leaving the repository.
- Maintained-policy scan, excluding the immutable closed audit attempt:
  **PASS**, no positive “sensitive operational detail” or
  `operational_sensitive` class remains.
- Rendered-site checker: **PASS**, five agent-system pages and 2,751 internal
  `/docs/` links resolve; no relative `.md` content link remains.

The link checker parses inline and reference Markdown destinations, ignores
external URLs and same-page anchors, resolves decoded paths from each source,
and rejects repository escapes or missing paths. The rendered checker parses
the five generated HTML pages, rejects non-URL `.md` destinations, maps every
`/docs/` destination to the site output, and requires the target to exist.

## Bazel tests

```text
bazel_agent test //:buildifier_test \
  //projects/agents/skills/answer-question:eval_config_test \
  //projects/agents/skills/answer-question:eval_no_skill_config_test \
  //projects/agents/skills/answer-question:eval_routing_config_test \
  //projects/agents/skills/repo-secrets:eval_config_test
```

**PASS:** five of five tests. The later formatter pass touched only the agent
landing page and the three design documents; none is an input to these five
tests.

## Bazel builds and site render

```text
bazel_agent build //projects/agents:docs \
  //projects/agents/docs:docs \
  //projects/agents/goals:docs \
  //projects/agents/skills/answer-question:skill \
  //projects/agents/skills/repo-secrets:skill \
  //projects/alwaldend.com:site
```

**PASS** against the final formatted source. Hugo 0.152.2 rendered 7,737
pages. The build emitted one pre-existing duplicate `_index.md` warning from
the external `rules_template` and `rules_binary_toolchain` documentation;
neither source is part of this candidate, and the agent-system output passed
the rendered-link check.

## Referenced Bazel labels

```text
bazel_agent query 'set(//infra/vault/tf:tf.apply \
  //infra/yandex_cloud/org1/tf:tf.apply \
  //infra/pve/tf:tf.apply \
  //third_party:publish_helm.io_goharbor_helm_harbor)'
```

**PASS:** all four exact labels resolve.

## Goal integrity

```text
bazel_agent run //projects/goal/cmd/goal -- validate \
  --goal-dir projects/agents/goals/repo-agent-system
```

**PASS before closure:** `valid: true`, count 1. The same validation is rerun
after the goal tool publishes this close review and result.

## Scope judgment

The changed behavior is documentation, policy, docs packaging, and skill
routing text. Focused format, skill, Buildifier, package, full-site, link,
label, and goal-integrity checks are more discriminating than a repository-wide
product test. No state-changing infrastructure, live authentication, full
repository test matrix, or unrelated subproject check was run.
