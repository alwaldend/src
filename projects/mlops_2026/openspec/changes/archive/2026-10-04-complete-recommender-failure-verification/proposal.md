## Why

PR review found that the completed application plan promises failure scenarios
not yet exercised by its E2E harness. Complete that verification before handoff.

## What Changes

- Verify startup rejects either missing CSV and a corrupt packaged model.
- Exercise the training wrapper with malformed source data and confirm it
  fails without publishing an artifact.
- Preserve repeatable JSON evidence for the additional failure scenarios.
- Add the requested client actor and HTTP connection to the README diagram;
  retain build-generated SVG documentation and show a dark 5:2 WebP preview.

## Capabilities

No specification behavior changes; this change completes verification of the
existing course-data and health-api contracts. Specs are explicitly skipped.

## Impact

Project E2E sources, their Bazel runtime dependencies, and diagram source. No new external
dependencies or application behavior changes.
