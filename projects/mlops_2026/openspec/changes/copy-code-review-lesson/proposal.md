## Why

Make the imported taxi-fare review exercise runnable and verified in this repository.
The user now requests applying the review suggestions and replacing the CSV sample
with the existing pinned dependency.

## What Changes

- Apply the reviewed typing, feature selection, deterministic training, and artifact changes.
- Remove the checked-in CSV and legacy requirements; use the existing archive and Python lock.
- Add Bazel data, training, model/metrics build, documentation, and end-to-end test targets.
- Remove the vendored/formatter exclusions now that the lesson is maintained here.

## Capabilities

### New Capabilities

- `code-review-training`: Reproducible taxi-fare training with dependency-backed data and saved artifacts.

### Modified Capabilities

None.

## Impact

Changes are scoped to the lesson, its project documentation and change record,
its existing third-party data exposure, and obsolete Git attributes. No new
external dependencies or deployment are required. Preserve the upstream MIT notice.
