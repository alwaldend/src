## Context

PR #128 currently combines Cobra commands, Bazel and shared-site integration, and collection processing. See proposal.md for the split scope.

## Goals / Non-Goals

The command boundary owns flags, cancellation, stdout/stderr routing, and error context. The extracted scaffold does not manage collection data.

## Decisions

Keep the existing Cobra command factories and string-flag definitions. Construct one App instance in the command entry point and pass it into command factories. The internal/app package owns the App type and its placeholder methods in app.go; each method returns an explicit error. PR #128 can implement these methods without changing the command tree.

Use the root Bazel module and existing Cobra dependency. Preserve project and landing registry integration. Document the scaffold's actual status.

Keep the existing CLI workflow test in PR #128: its successful identity-generation case requires business logic. Validate the scaffold with repeatable command invocations recorded under ignored task scratch instead.

## Risks / Trade-offs

- The command names resemble a working tool → help text and documentation identify the placeholder status, and operations fail explicitly.
- Both PRs target master and temporarily overlap → merge the scaffold prerequisite first, then remove shared changes from PR #128 through base synchronization.
