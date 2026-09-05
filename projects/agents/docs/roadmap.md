---
title: Repository agent-system roadmap
description: Remaining work ordered by observed task friction and evidence
tags:
  - agent
  - architecture
  - roadmap
---

# Repository agent-system roadmap

This document records intentional future work. The
[current-state guide](current-state.md) owns the supported interface summary;
component documentation and executable tests establish implementation details.
The [architecture](architecture.md) describes composition responsibilities.
Historical phase plans and acceptance records remain in the immutable
[goals](../goals/). A completed historical phase does not establish a current
runtime guarantee.

## First: measure and simplify ordinary tasks

Use representative question, edit, recovery, and delivery scenarios to compare
the current workflow with a candidate. Keep source, instructions, fixture,
model, observation method, and execution identity with each result. Measure
task success, unauthorized actions, commands, context actually consumed,
duration, and resume steps where the observer can establish them. Record
unavailable measurements explicitly. Configuration validation and shorter
documents do not establish behavioral improvement.

Local simulated delivery scenarios can exercise candidate selection and
validation decisions without publishing anything. They do not prove forge
integration, remote concurrency, or successful real publication. Use the
owning delivery integration tests for their encoded cases, and distinguish
those tests from model behavior.

Reduce repeated common-path instructions and load exceptional procedures only
when their triggering condition occurs. Keep exact candidate, task ownership,
authority, leases, recovery, and review completion requirements. Improve
entrypoint capability detection and workspace-relative path handling before
introducing another execution interface.

Acceptance requires measured task outcomes with no new critical boundary
failures. A small exploratory comparison is useful evidence, not a universal
performance claim. Retain useful diagnostics and regressions after clean runs;
remove obsolete workarounds and completed tracking when their purpose expires.

## Next: reduce manual discovery and validation work

Extend bounded offline context only with facts that have an identifiable
owner and observation method. Prefer observed Git state and direct routes to
applicable policy, owner documentation, build declarations, and explicit goal
records. Never infer authorization or task-to-goal association from proximity,
ownership, a catalog entry, or an available tool.

Check suggestions must explain their selection basis and uncovered scope.
Nearest-package selection is not dependency analysis. Root configuration,
build-graph changes, ignored paths, and nested workspaces can require explicit
checks beyond the returned labels. Avoid accidental broad builds while keeping
the mandatory repository quality and affected semantic checks.

Delivery preparation binds the candidate and publication expectations. The
caller selects applicable validation. The delivery tool can execute that plan,
record candidate-bound results, and carry them into explicitly requested
publication. The manual publication path also remains supported. Neither a
preparation receipt nor an asserted commit ID proves that tests ran. Broader
evidence imports must specify their trust boundary and preserve caller
responsibility for semantic sufficiency.

Reuse observations only while their relevant inputs remain unchanged. Reuse
test evidence across a message-only amendment only when tree, configuration,
toolchain, environment, and every other relevant input are stable. Rerun
commit-sensitive or uncertain checks. A receipt alone never establishes reuse.

Acceptance requires representative edits to select sufficient checks, explain
omissions, and invalidate evidence when a relevant input changes. Preserve
cheap read-only discovery when optional providers are unavailable.

## Expand behavioral evidence when it answers a concrete question

The configured coverage inventory is not a successful-run ledger. Extend
routing scenarios across positive, adjacent-negative, inert-payload, exclusion,
conflict, and composition cases when the corresponding skill is changed or a
recurring failure demonstrates a gap.

Stronger aggregated evidence needs actual execution artifacts bound to the
source, fixture, model or provider, and judge where applicable. Keep live
stochastic evaluations manual or scheduled and separate from normal wildcard
tests. Use isolated writable fixtures for Git, Bazel, forge, infrastructure,
secret-handling, and runtime trajectories before making claims about them.
Never reuse shared authentication in a mode that requires exclusive access
while concurrent sessions are active.

Promote a repeated lesson only with an owner, minimized public reproducer,
regression, supported contract, fallback, measured cost, delivered revision,
and a reasoned retirement rule. A schema or proposal does not establish that
this learning loop runs automatically.

## Runtime integration remains conditional work

Shared admission and control libraries have unit-test coverage. Production
provider enforcement, process isolation, and runtime fault containment remain
unverified. Build further runtime machinery only for an identified consumer
whose required behavior and failure modes justify it.

For that consumer, test bounded status during a stuck extension, package
deadlines, desired versus observed revision, cancellation and cleanup,
cross-task storage isolation, and the actual filesystem, process, network,
credential, and output boundaries. Trusted in-process code is not a security
sandbox. A source contract cannot substitute for consumer integration evidence.

Keep bounded native and offline fallbacks available. A missing optional
provider must produce an explicit unavailable result rather than preventing
ordinary repository work.

## Constraints on further work

- Facts remain at their owning sources; generated views carry provenance and
  limits rather than becoming new mutation authorities.
- Keep goals, delivery, release, review, and provider state separate. Add joins
  through typed references only when a real workflow needs them.
- Do not add a mandatory daemon, database, vector store, network call, or
  central orchestrator for basic orientation.
- Do not relocate components merely to match conceptual diagrams.
- Keep reviewed evidence bounded and free of credentials and personal data.
- Use workspace goals for temporary coordination and owner-local maintained
  records only when durable project history is justified.
- Choose the least costly check that answers acceptance; broad audits remain
  separately scoped work.
