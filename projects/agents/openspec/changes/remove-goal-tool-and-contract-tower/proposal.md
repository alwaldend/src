## Why

`tools/agents` carried repository-internal control contracts, seven catalog
compilers, an offline context capsule, and typed admission/control/plan/evidence
libraries. The goal tool was its only real consumer: the admission, control,
plan, and evidence packages never gained a production consumer, no CI or
`AGENTS.md` invocation used the catalogs, and the goal tool is being removed.
Keeping the tower costs maintenance for projections nothing reads.

## What Changes

- **BREAKING** Remove `tools/agents` entirely: catalog compilers, context
  capsule, control-status and evidence tools, registry declarations, generated
  catalogs, and the admission/control/plan/evidence libraries.
- Relocate the repository skill-discovery declaration to `.agents/BUILD.bazel`,
  next to the directory it generates, and repoint skill `BUILD.bazel`
  visibility from `//tools/agents:skill_discovery` to `//.agents:skill_discovery`.
  Skill content, link targets, and the 27 existing symlinks are unchanged.
- Delete the architecture, current-state, and roadmap documents. This project is
  a collection of skills; its documentation should describe skills, not an
  architecture program for tooling that no longer exists.
- Leave the removed bytes in git history; the change neither copies nor
  restates them.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-agents`: Discovery packaging moves to `//.agents`, the deprecated
  goal skill is removed rather than excluded, and the guidance that derived
  catalogs supply orientation is withdrawn.

## Impact

`tools/agents` is deleted; `.agents/BUILD.bazel` is added and 29 skill
`BUILD.bazel` files change one visibility label. `projects/agents/docs` is
deleted, and the `openspec` and `bazel-rules-skill` skills drop their
`tools/agents` references. `.gitattributes` no longer ignores removed catalog
JSON.
