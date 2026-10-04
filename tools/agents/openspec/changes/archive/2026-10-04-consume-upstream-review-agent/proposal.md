## Why

The user refined the request from strengthening review-change to consuming the
authoritative Codex review-agent directly and removing the repository substitute.
OpenAI publishes the requested skill under Apache 2.0 in its Codex source tree.

## What Changes

- Pin the upstream source archive by immutable commit and SHA-256 integrity.
- Package the upstream review-agent instructions and metadata without edits and
  generate their repository discovery projection using the existing workflow.
- Remove review-change and require review-agent in repository review routing.
- Retain upstream licensing and provenance; move offline evaluation fixtures to
  the repository consumer rather than modifying the upstream skill.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- project-agents: Use pinned upstream review-agent for all bounded change reviews.

## Impact

Adds the user-requested OpenAI Codex skill dependency, not the Codex executable.
The 17 MB source archive is fetched through Bazel; only the two skill files are
installed. Updating the pin and regenerating discovery is the maintenance cost.
The earlier local strengthening implementation was superseded before delivery.
