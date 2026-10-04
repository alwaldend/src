## Why

The user requested the skill name review-change. The existing change-review
identity must be renamed consistently across source, discovery, and routing.

## What Changes

- Rename the canonical skill directory and frontmatter to review-change.
- Update Bazel discovery, generated links, and AGENTS.md routing.
- Update the maintained agent specification while preserving review behavior
  and historical archived records.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- project-agents: Identify the bounded change review skill as review-change.

## Impact

The old skill invocation name is replaced by review-change. The review procedure,
authority boundaries, and offline evaluation cases remain unchanged.
