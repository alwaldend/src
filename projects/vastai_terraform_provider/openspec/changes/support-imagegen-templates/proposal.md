## Why

Reusable image-generation rentals must select existing Vast.ai templates and
inspect their Docker options before accepting a paid offer.

## What Changes

- Add a template lookup data source and creation-only instance template hash.
- Preserve existing image-based instance behavior and imports.
- Exercise template safety and module lifecycle through offline Terraform CLI fixtures.

## Capabilities

### New Capabilities

- `instance-templates`: Explicit template lookup and instance selection.

### Modified Capabilities

None.

## Impact

Provider schemas, acceptance fixtures, and documentation. The reusable module
belongs to `projects/tf_modules/openspec/changes/add-vastai-instance-imagegen`;
protocol support belongs to `internal/vastai_sdk/openspec/changes/read-instance-templates`.
No new external software or live operations.
