## Why

Delivery repeatedly mistakes sandbox-inaccessible GitHub credentials for a login or Git transport problem. The same host login succeeds outside the sandbox, so retrying with `--git-transport gh` or asking for a fresh login does not address the observed failure.

## What Changes

- Explain API authentication recovery in AL's HTTP 401 preflight diagnostic while preserving the original failure and refusal.
- Teach the delivery skill to compare sanitized authentication status in the failing and approved execution contexts, retain its selected transport, and avoid unnecessary login changes.
- Cover sandbox recovery, genuine credential failure, and approval rejection in delivery evaluations.

## Capabilities

### New Capabilities

- `delivery-authentication`: Distinguish GitHub API credential access from Git transport and recover through an approved execution context.

### Modified Capabilities

None.

## Impact

AL's GitHub preflight, delivery skill, documentation, and existing test/evaluation suites. No new dependencies, credential extraction, automatic escalation, host configuration changes, or transport defaults.
