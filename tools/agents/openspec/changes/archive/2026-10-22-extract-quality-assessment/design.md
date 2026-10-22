## Context

See [proposal.md](proposal.md). The current decision skill already has a
Promptfoo suite; its isolated workspace stages only that skill.

## Goals / Non-Goals

Make general assessment reusable with a concise decision-specific layer.
Keep question authority in repository policy and bounded code review at its
existing owner. Remove the separate answer-question skill as requested; do
not change providers.

## Decisions

- Name the skill `quality-assessment`, correcting the request's spelling, and
  place it beside `decision-review` under the generic `tools/agents` owner.
- Link to the sibling skill using a relative path so composition works in both
  canonical directories and staged skill workspaces. Copying the procedure
  would create two owners; renaming decision-review would lose its useful
  decision-specific entry point.
- Keep routing in `AGENTS.md` and registration in `.agents/BUILD.bazel`.
  Stage both skills for decision evaluations, granting only that package the
  extra access it needs.
- Remove answer-question rather than retaining a second owner of assessment
  guidance. Update live examples and fixture expectations; retain historical
  records and synthetic tests that do not consume the removed package.
- Define behavioral cases before implementation. Cover unsupported praise,
  warranted positive assessment, numerical accuracy, missing evidence,
  subjective criteria, and resistance to prior-answer anchoring. Keep routing
  assertions separate from answer-quality assertions.

## Risks / Trade-offs

- Broad discovery could capture editing or code reviews: describe the request
  boundary and include routing cases for quoted transformation payloads.
- Extraction could break isolated decision evaluations: stage the new skill
  with the existing one and check the resulting workspace.
- Offline checks cannot establish judgment quality: report that limit and keep
  live evaluation manual; do not claim behavioral improvement without runs.
