---
title: TruffleHog wrapper command acceptance
---

These failure cases are defined before implementation:

- Only an exact detector/path/commit/line tuple succeeds; each changed coordinate
  and a second unmatched finding must fail without hiding the original records.
- A new commit at the same path and line must not inherit an older suppression.
- Unknown policy fields (including value hashes), malformed or abbreviated commit
  IDs, invalid paths, nonpositive lines, missing reasons, and duplicate entries fail.
- Truncated JSON, missing finding metadata, scanner launch failures, timeout,
  scanner errors after approved output, and contradictory scanner exit codes fail.
- The scanner must receive `--no-verification`, `--no-update`, and error-failure
  flags. HEAD mode skips additional refs; all-ref mode scans all local refs.
- Zero findings succeed; fully approved findings retain original values and counts.
- The E2E test invokes the built wrapper and a controlled scanner subprocess,
  exercising the pinned JSON protocol without contacting credential services.
  It writes `e2e-report.json` to Bazel's undeclared test outputs. The repository
  consumer separately exercises the real pinned scanner over Git history.
