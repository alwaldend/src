## 1. Extract scaffold

- [x] 1.1 Copy the command tree and replace collection calls with explicit placeholders; verify help, argument failures, and valid-operation failures without filesystem writes.
- [x] 1.2 Extract Bazel targets, README, landing content, and project registration; verify source boundaries and landing registry wiring.
- [x] 1.3 Exclude the business-dependent CLI workflow test; verify no test source or collection dependency is present in the scaffold target.

## 2. Deliver

- [x] 2.1 Select exact-candidate OpenSpec, repository quality, buildifier, semantic lint, CLI, and site checks through delivery receipts.
- [x] 2.2 Select guarded publication against master through repo-delivery and record the prerequisite relationship with PR #128 in the commit projection.

Delivery receipts own subsequent validation and publication state; representative command output is recorded in ignored out/anki-cli-scaffold/cli-verification.json.

## 3. Address application boundary review

- [x] 3.1 Move placeholder operations onto an App struct in internal/app and pass one instance to command factories; verify the existing command artifact still reports placeholder errors without filesystem writes.
- [x] 3.2 Update package wiring and layout documentation; select exact-candidate checks and guarded publication/review resolution through delivery receipts.
