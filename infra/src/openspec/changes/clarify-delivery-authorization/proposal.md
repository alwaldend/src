## Why

An automatic approval review rejected the task's GitHub publication despite the standing delivery instruction. Clarify that an authorized repository change includes disclosure of reviewed task-owned source to the configured repository destination and PR creation, rather than merely offering to create one.

## What Changes

- Make the root delivery authorization explicit about push, PR creation, destination, and source disclosure.
- Route the delivery skill to that policy and add an evaluation case for avoiding redundant destination approval.
- Preserve validation, secret handling, task ownership, and automatic approval boundaries.

## Capabilities

### New Capabilities

None. Clarification of existing authorization; `skip_specs: true`.

### Modified Capabilities

None.

## Impact

Root AGENTS.md and the delivery skill in tools/repo_delivery. No host, infrastructure, dependency, or provider configuration changes. The user's current request explicitly authorizes publication of this change.
