---
title: AL product description
description: Commands, observable behavior, interruption, and recovery
---

# AL product description

AL gives a repository operator one command for inspecting configuration, running
commands with selected plugins, selecting tests affected by Git changes,
choosing a shared source version,
publishing individual project releases, extracting and activating websites, and
delivering a branch for review.
The operator chooses each project and each publication step. No CI scheduler or
automatic weekly release is part of this surface.

## Purpose

This is a source-linked description of the command-line experience: what an
invocation reads, changes, reports, and leaves behind after interruption. It is
for operators and people changing AL. The executable and its tests remain the
source of behavior; these documents are a projection of the source snapshot,
not an additional implementation contract. The [project README](../../README.md)
owns build API documentation and operator recipes.

## Conventions

Use the [glossary](glossary.md) consistently. Describe observable behavior in
plain language. Put mechanisms in a `Technical note` only when they change what
an operator should expect. Link to the owning foundation instead of copying its
rules. A source footer identifies the code examined; it does not claim a manual
runtime pass. Suspected defects belong in [bug triage](bug-triage.md).

## The work to be done

### Document template

Features use eight sections: Summary; The simple case; The interaction, event by
event; Modifiers; Cancel and interrupt; Interactions with other systems; Edge
cases; Open questions and verification. Foundations may omit inapplicable phases.

The interaction is one invocation. Its phases are **Invoke**, **Exits
immediately**, **Begins running**, **While running**, and **Finishes**. Each feature
includes a Mermaid state diagram for repository viewers that render Markdown.

Every modifier table uses these rows, in order: **Arguments and flags**,
**Environment and configuration**, **Terminal and streams**. Columns distinguish
behavior from output. Arguments do not change after invocation.

Every interrupt table uses these rows, in order: **Ctrl+C once or twice**,
**SIGTERM**, **Terminal closed (SIGHUP)**, **Input or output pipe closes**,
**Network lost**, **Disk full or file inaccessible**, **Required state changes**,
**Concurrent invocation**, **Forced kill**. Columns are **Before work** and
**During work**. Unknown behavior is identified, never inferred as rollback.

Cross-cutting concerns appear in this order: **Configuration**, **Output**,
**Permissions**, **Concurrency**, **Persistence and recovery**, **Network**.

### Method

Read the owning command, state handling, and behavior tests; draft the pilot,
foundations, and release publication first; then draft other features. Check
terminology, links, table coverage, and source footers across the set. Run safe
local invocations and retain their outputs. The checklists distinguish those
observations from fixture coverage and from unperformed terminal or live checks.

### Verification

[Verification](verification/README.md) owns setup, results, and limitations.
Coverage remains `drafted` after a scripted pass; manual observation is required
before marking a document `verified`. No release destination is changed to
verify these documents.

### Scope decisions

- **Surface:** Linux AL CLI in this repository, with Bash invocation, defaults
  unless flags are stated. Both an installed binary and Bazel-run commands are
  distinguished where their behavior differs. No user configuration is assumed.
- **Owner and integration:** AL lives in `projects/al`; its Bzlmod name remains
  `rules_al`. Repository-specific configuration supplies its integration inputs
  and its validation policy. The observations here cover this
  repository and its isolated fixtures; moving the owner does not broaden that
  evidence into a verification of every workflow in another repository.
- **Source:** `alwaldend/src`, commit `7cd4714139984bf4f9b2ae1de9189470b336e12e`,
  following the authorized rebase and implementation work on PR 226. All feature
  footers use one source snapshot, refreshed together after authorized fixes.
  Earlier observations retain their original source and executable identities.
- **Runtime:** source-built AL installed into ignored task scratch; behavior
  tests exercise isolated repositories, plugins, and SSH destinations.
- **Placement:** this owner-local directory, delivered with the existing change;
  no separate Git repository. Drafting reads the implementation; separately
  authorized implementation fixes require a coordinated source and evidence
  refresh rather than silently changing the meaning of earlier observations.
- **Excluded:** Go/Starlark API reference, third-party programs' own behavior,
  actual credential values, live deployments, and infrastructure synchronization.
  Relevant limitations and boundaries are described where the operator meets them.
- **Completeness:** all seven command families are covered, including the nested
  `release website` extraction and activation commands. Detailed delivery
  flags remain in its existing command reference; this set describes their
  user-visible lifecycle and state handoffs.

## Structure

- `README.md`: purpose, conventions, scope, and coverage.
- `goal.md`: drafting and continuation instructions.
- `glossary.md`: shared vocabulary.
- `AGENTS.md`, `CLAUDE.md`: local documentation entry points.
- `BUILD.bazel`: documentation packaging.
- `foundations/invocation.md`: command entry, output, and interruption.
- `foundations/state.md`: local, source, publication, and receipt state.
- `releases/bundle.md`: pilot; generate an inspectable project bundle.
- `releases/publication.md`: artifact planning, inspection, upload and deletion;
  website extraction, activation, dry runs, and recovery.
- `releases/versioning.md`: source versions, refs, and stamped builds.
- `delivery/workflow.md`: candidate preparation through publication and review.
- `commands/configuration.md`: inspect and merge explicit configuration.
- `commands/run.md`: run a child with selected plugins and cleanup.
- `commands/bazel-doctor.md`: select affected tests and inspect local setup.
- `verification/README.md`: protocol, evidence, and pass limits.
- `verification/foundations.md`: foundation checklists.
- `verification/releases.md`: version, bundle, and publication checklists.
- `verification/delivery.md`: delivery checklist.
- `verification/commands.md`: configuration, execution, and diagnostic checklists.
- `bug-triage.md`: deduplicated suspected defects and decisions.

## Coverage

| Document                    | Status  |
| --------------------------- | ------- |
| glossary.md                 | drafted |
| foundations/invocation.md   | drafted |
| foundations/state.md        | drafted |
| releases/bundle.md          | drafted |
| releases/publication.md     | drafted |
| releases/versioning.md      | drafted |
| delivery/workflow.md        | drafted |
| commands/configuration.md   | drafted |
| commands/run.md             | drafted |
| commands/bazel-doctor.md    | drafted |
| verification/README.md      | drafted |
| verification/foundations.md | drafted |
| verification/releases.md    | drafted |
| verification/delivery.md    | drafted |
| verification/commands.md    | drafted |
| bug-triage.md               | drafted |

## Reference

Paths below are relative to `projects/al`:

- `cmd/al/main.go`, `internal/cli`: startup, streams, configuration, child execution.
- `pkg/al`, `pkg/al_plugin`, `pkg/lifecycle`: configuration and plugin lifetime.
- `internal/versioncmd`, `internal/versioning`: calendar and Git operations.
- `internal/release`: bundle creation and SSH release transitions.
- `internal/delivery`: candidate, publication, and review guards.
- `internal/bazelcmd`: affected-test selection and local AL/Bazel diagnostics.
- `cmd/al/lifecycle_e2e_test.go`, `cmd/al/run_e2e_test.go`,
  `internal/release/lifecycle_test.go`, `internal/release/observation_test.go`,
  `internal/versioning/publication_integration_test.go`,
  `internal/delivery/delivery_integration_test.go`: behavior fixtures.
- Root `tools/al/test/command`, `tools/al/test/contract`, and
  `infra/download/test/sync`: plugin and real SSH/Ansible fixtures.
