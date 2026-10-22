## Why

`decision-review` currently owns general assessment guidance that is also useful
for questions such as "is this good?" and "is this accurate?" Extracting that
guidance gives quality assessment one reusable owner without requiring every
critique to become a decision process.

## What Changes

- Add `quality-assessment` for rigorous, evidence-based assessment of claims,
  answers, artifacts, and work against their intended purpose.
- Make `decision-review` use it while retaining decision-specific alternatives,
  trade-offs, verdicts, and reconsideration triggers.
- Remove `answer-question`, its evaluations, discovery entry, and routing;
  update maintained examples and the assessment fixture to use the new skill.
- Register discovery and routing, and cover both standalone assessment and
  decision composition in the evaluation configurations.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-agents`: Provide reusable quality assessment and compose it into
  consequential decision review.

## Impact

Changes belong to `tools/agents/skills` and this OpenSpec workspace, with the
required discovery entry in `.agents/BUILD.bazel` and routing in `AGENTS.md`.
Documentation examples in `rules_skills` and `rules_promptfoo` follow the
remaining skill, and the agent ergonomics fixture expects quality assessment.
No new external dependencies or host configuration changes are needed.
