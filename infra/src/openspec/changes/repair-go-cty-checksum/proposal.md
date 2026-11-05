## Why

The root module selects github.com/zclconf/go-cty v1.13.1 but records only
its go.mod checksum. Gazelle cannot fetch its source archive, blocking the
semantic lint gate for PR 211 after rebasing onto master.

## What Changes

- Use the pinned Go module command to record the missing source checksum
  for the already selected version.
- Verify the original lint command and the existing Mermaid consumers.

## Capabilities

No runtime or specification behavior changes. This repairs dependency metadata.

## Impact

Root go.sum and this maintenance record. Dependency versions remain unchanged.
Acceptance requires successful broad semantic lint and normal delivery gates.
