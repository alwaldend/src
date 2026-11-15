## 1. Document the upstream limitation

- [x] 1.1 Add the upstream issue, checked versions, integer-only argument workaround, and command/lock recovery guidance to the host README; compare each claim with the linked upstream code and observed parser behavior.
- [x] 1.2 Validate the change record, format the documentation, and pass repository quality and affected-document lint before delivery.

No design artifact is needed: this change documents existing upstream behavior
without changing architecture, dependencies, configuration, or deployment.

## Validation evidence

Verified on 2026-11-15:

- The installed CLI reported `0.167.0`. Upstream source for that release and
  `0.169.0` retains the strict terminal session parser described in the README.
- An isolated check using Codex 0.167.0's pinned Serde versions reproduced
  `invalid type: floating point` for `96230.0`; bare integer parsing succeeded,
  and string, fractional, out-of-range, missing and duplicate inputs were rejected.
  This check demonstrates the existing defect, not a patched executable.
- Strict OpenSpec validation and all 28 selected repository quality/owner
  validation tests passed. Prettier reported the four task-owned files unchanged.
- The affected documentation and OpenSpec source targets built with the
  repository lint configuration. No host configuration or executable was changed.
