---
title: Repository agent system current state
description: Supported agent entry points, evidence boundaries, and limitations
tags:
  - agent
  - architecture
  - repository
---

# Repository agent system current state

This is the maintained guide to the checked-in interfaces. Source links below
identify their owners; runtime availability must be observed in the current
session. Historical audits and acceptance claims remain in the immutable
[goal records](../goals/). Their dates and evidence tiers limit what they prove.

## Start with the affected path

Read applicable `AGENTS.md` instructions and the nearest owner README. Use
the available Cordis repository-context handlers for bounded reads and
searches, or the native tools when Cordis is unavailable. Cordis availability
does not determine whether offline repository work is possible.

The offline context command provides a bounded advisory view:

```sh
bazel_agent bazel run //tools/agents/cmd/agent_system -- \
  --workspace-root "$PWD" --path projects/goal
```

Its [owner documentation](../../../tools/agents/README.md) describes the
supported flags and limitations. Catalog declarations identify source facts;
they do not establish live provider health, action authorization, successful
validation, or a task's association with a goal. Missing observations remain
explicit. A source digest proves identity, not freshness against a live system.

## Use the owning interface

| Need                  | Entry point                                                  | Boundary                                                                                                       |
| --------------------- | ------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------- |
| Policy and layout     | Applicable `AGENTS.md`, owner README, BUILD and MODULE files | The user request owns intended outcome and authorization.                                                      |
| Procedure             | Canonical skill selected from `.agents/skills`               | Load only the procedure needed for the task.                                                                   |
| Build or test         | `bazel_agent bazel <command>`                                | Start with the affected package and avoid redundant compatible invocations.                                    |
| Context               | `//tools/agents/cmd/agent_system`                            | Offline, advisory source projection.                                                                           |
| Stored runtime status | `//tools/agents/cmd/control_status`                          | Stored package observations do not prove present runtime health.                                               |
| Durable work          | `//projects/goal/cmd/goal`                                   | Use workspace records for task coordination and owner-local project records when maintained history is needed. |
| Delivery              | `//tools/repo_delivery`                                      | Prepare, establish applicable validation, then explicitly publish the exact candidate.                         |

The [Bazel runner](../../bazel_agent/) also supports cached control-tool
execution in current builds. If an installed runner lacks that subcommand,
the supported `bazel run` targets remain available. Updating a shared host
installation is a separate operation from changing repository source.

## Keep routine work small

Simple tasks need no durable goal. For work that needs a resume point, the
[goal checkpoint](../../goal/) accepts a short summary, the current candidate,
evidence references, and the next action. The store maintains the record's
identity and digests; `show` includes the active continuation. Detailed attempt
review remains available for repeated failures and maintained project history.
Local progress checkpoints do not trigger publication. Deliver at a meaningful
review boundary or final handoff, with explicit remote backups when needed.

[Delivery](../../../tools/repo_delivery/) can execute a caller-selected Bazel
validation plan and retain its results against the prepared candidate. Its
`continue` command reports the next state; `continue --publish` explicitly
requests publication using those results. The tool carries the candidate and
receipt between steps. It does not decide task ownership, semantic sufficiency,
or whether representative output satisfies the request. Changed candidates and
rebases require fresh validation; uncertain publication outcomes require
inspection before another mutation. Follow the owning README's launcher
instructions when validation itself invokes Bazel.

## Supported guarantees

- [Skill packaging](../../rules_skill/) validates canonical artifacts and
  discovery links. Evaluation payloads are separate from runtime skill bodies.
- [Goal storage](../../goal/) validates resource revisions and artifact
  digests, provides bounded continuation information, and supports recovery
  from interrupted publication. These checks establish record integrity;
  criterion verdicts still require appropriate evidence.
- [Delivery](../../../tools/repo_delivery/) binds the prepared candidate,
  scope, remote expectations, and pull-request identity. Exact leases and
  state checks protect publication. The caller selects and establishes the
  required validation; a preparation receipt is not a test result.
- [Cordis](../../mcp_cordis/) provides bounded repository tools and supervised
  subprocess cleanup. Its package code is trusted. The documented process
  boundary is not a security sandbox.

## Evidence and implementation limits

The [coverage projection](../../../tools/agents/catalogs/skill-coverage.md)
inventories configured cases. It does not measure successful routing or task
completion. Offline configuration checks are useful hygiene; stronger claims
require actual result artifacts identifying the source, fixture, execution,
and judge where applicable.

The [shared control contracts](../../../tools/agents/) include admission and
runtime-control libraries with unit tests. Their existence does not establish
enforcement by production providers. Earlier phase acceptance records describe
broader gateway and runtime guarantees than those unit tests demonstrate.
Treat provider integration, process isolation, and runtime fault containment
as unverified until a concrete consumer and corresponding integration evidence
establish them. The [roadmap](roadmap.md) retains that future work.

Stored control snapshots contain package observations. Their age, missing
writer identity, and lack of a live heartbeat limit their interpretation.
The status reader preserves those limitations and does not create a runtime
or its storage directories merely to inspect them.

## Improve from task outcomes

Measure representative questions, edits, recovery, and delivery tasks before
claiming an instruction or workflow is faster or more reliable. Compare task
success, unauthorized actions, context bytes, commands, duration, and resume
steps against the prior behavior. Keep model and fixture identities with
observations; do not infer measured improvement from configuration checks or
the presence of a new schema.

The [local comparison fixtures](../test/ergonomics/) provide four bounded
scenarios and an observer protocol. They cover question answering, a local
prose edit, recovery after a changed postcondition, and a delivery validation
decision. Their simulated observations do not exercise real publication.
Unobservable routing, command, context, or token metrics remain unavailable;
available instruction bytes are not the same as consumed context. Results
belong with the exact task candidate, not in the configured coverage inventory.

Preserve useful safeguards while reducing repeated consumer instructions.
Retire obsolete workarounds and completed tracking when justified. Clean runs
after a fix are not a reason to remove useful diagnostics or regression tests.
Historical proposals suggesting otherwise do not define current policy.

The [architecture](architecture.md) describes composition responsibilities.
It does not replace these implementation boundaries or authorize new work.
