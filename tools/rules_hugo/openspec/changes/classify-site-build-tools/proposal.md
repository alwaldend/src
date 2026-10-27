## Why

Site tools must use the execution configuration so non-tool dependency queries
and cross-platform builds distinguish tooling from site input archives.

## What Changes

- Classify `al_hugo_site.tools` as execution dependencies.
- Preserve the target configuration for site source inputs.

## Capabilities

### New Capabilities

- `site-tool-classification`: Separate execution tools from site inputs.

### Modified Capabilities

None.

## Impact

One attribute transition and its documentation. Validate the existing rule
suite and parent site graph. No dependency or deployment change.
