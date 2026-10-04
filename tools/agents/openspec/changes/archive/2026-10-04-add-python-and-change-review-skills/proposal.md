## Why

The repository needs Python-specific guidance with mandatory type hints and a
repeatable change review procedure that finds bugs and violations of applicable
repository standards. The initial PR omitted OpenSpec; this record captures its
actual scope and subsequent user correction without implying prior planning.

## What Changes

- Package and discover repo-python with mandatory type hints and owning Python
  dependency, configuration, Gazelle, and validation workflows.
- Correct its Python support guidance after splitting dependency metadata from
  shared tool settings.
- Add change-review for evidence-based bug and standards review of exact candidate and baseline revisions.
- Update openspec and repo-workspace to apply the AGENTS.md requirement for
  every change without adding redundant approvals.
- Register both skills in Bazel discovery and repository routing, with offline
  Promptfoo coverage and clearly stated behavioral evaluation limits.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- project-agents: Support Python development guidance and bounded change-set bug review.

## Impact

Canonical skills live in tools/agents/skills and are projected by .agents/BUILD.bazel.
AGENTS.md owns skill routing. The linked repository change is
../../../../../../infra/src/openspec/changes/archive/2026-10-04-organize-python-course-development/.
No review request alone authorizes source fixes, posting comments, or merging.
