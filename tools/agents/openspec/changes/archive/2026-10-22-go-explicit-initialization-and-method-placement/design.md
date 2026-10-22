## Decisions

Avoid init functions. Initialize components explicitly at their owning entry point and return contextual errors. Keep methods in the source file declaring their receiver type.

Canonical instructions remain in tools/agents/skills/repo-go/SKILL.md. Offline Promptfoo validation verifies packaging and configuration, not agent behavior. This PR targets master independently of PR #128.

## Acceptance

Validate the skill aspect, offline evaluation configuration, discovery state, repository quality and semantic lint against the delivery candidate.
