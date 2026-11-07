## 1. Establish acceptance before implementation

- [x] 1.1 Rebase onto merged trunk and record the upgrade decisions and failure
      cases; verify the worktree contains no unrelated changes.
- [x] 1.2 Write the isolated role integration scenarios before implementation;
      verify they use temporary local services, existing dependencies, and a
      retained sanitized result artifact.

## 2. Implement the role transaction

- [x] 2.1 Implement preflight, protected offline archives, controlled activation,
      and fail-closed recovery; verify success, no-op, downgrade, and failure
      behavior through the role integration scenarios.
- [x] 2.2 Document supported persistence, configuration variables, backup
      retention, and explicit full-state recovery; verify the retained archive
      can restore the prior executable and its matching database locally.

## 3. Verify and deliver

- [x] 3.1 Build the role and deployment consumer, run integration acceptance,
      OpenSpec validation, semantic lint, and repository quality checks against
      the implementation; retain results and coverage limitations.
- [x] 3.2 Review the complete source change and prepare delivery inputs without
      deploying. Git and al-delivery receipts own final candidate, publication,
      and remote review state.
