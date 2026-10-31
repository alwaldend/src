## 1. Guidance

- [x] 1.1 Add evaluation cases for preserving agreed design semantics and distinguishing generated build outputs from maintained inputs; inspect the expected outcomes against the observed failures.
- [x] 1.2 Update the two canonical skills and validate their packaged evaluation configurations and discovery.
- [x] 1.3 Prepare the guidance for delivery with the runtime fixes after formatting, quality, affected-package lint, offline configuration validation, discovery, and strict OpenSpec validation pass. The delivery receipt records exact-commit validation and publication.

## Evidence and limits

The selected suite passed all 32 tests, including the affected offline skill
configurations and discovery. Semantic lint and strict OpenSpec validation
passed. These are configured coverage, not evidence of improved agent behavior.
Measure future design correction turns and generated-output reversals before
claiming an improvement.

ERG-003 is already being corrected in [PR #223](https://github.com/alwaldend/src/pull/223);
this change does not duplicate that active work. ERG-004's original-output
retention contract is already present in the canonical `repo-secret-scan` skill.
