---
title: Goal
description: Deprecated goal tool retained for legacy record compatibility
statuses:
  - maintenance
languages:
  - go
  - markdown
  - bzl
tags:
  - agent
  - workflow
---

The goal tool is deprecated. New specifications and resumable work use the
affected owner's `openspec/` workspace through the pinned
[OpenSpec command](../../tools/openspec/README.md). Select its owner directory
with `OPENSPEC_PROJECT=<owner>` when running the command from the root Bazel
workspace. `infra/src/openspec/` is for repository evolution. The goal skill is
disabled and removed from repository skill discovery.

This project retains the portable resource API, deterministic local store,
and existing CLI commands for legacy record compatibility and recovery. The
CLI emits its deprecation notice on stderr, preserving JSON and version
output on stdout. Existing readers, validation, recovery, and legacy writes
remain available; creating new goal records is no longer the supported
repository workflow.

```text
projects/goal/
  api/v1alpha1/      Portable resource types and domain validation
  cmd/goal/          Cobra command wiring and command tests
  docs/              Current architecture and generated diagrams
  internal/fsstore/  Local persistence, per-goal locking, and store tests
  skills/goal/       Disabled skill notice and historical format references
```

Go dependencies flow from `cmd/goal` to `internal/fsstore` to
`api/v1alpha1`. The API package has no filesystem or command dependency.
The Go toolchain is Bazel's `@rules_go//go`; run the package targets below
rather than a host-installed `go`.

The command stores Kubernetes-inspired `goals.alwaldend.com/v1alpha1` YAML
resources and canonical, digest-bound attempt Markdown in ordinary repository
files. Each `attempt.yaml` binds `plan.md`, `result.md`, and `evidence/*.md` by
SHA-256. `Goal.status.plans` records durable plans with one active plan at a
time and states `active`, `accepted`, `rejected`, or `superseded`; attempts can
bind to a plan with `spec.planID`. Each goal record's `README.md` is a generated, bounded,
replaceable projection. The portable API treats resource versions as opaque
and excludes local ownership paths from desired state. The filesystem backend
supplies one cooperative lock per goal, local numeric resource versions,
atomic file replacement by sibling temporary-file rename, session bindings,
promotion, non-destructive unversioned-record import, and recoverable
multi-file publication. Before a multi-file mutation, the backend stages exact
after-images and installs a `.goal-publication.yaml` intent in the goal
record; `goal doctor` classifies the publication state and `goal recover`
replays or discards a pending intent. It is a local file editor with atomic
per-file replacement; it does not claim cross-file transaction semantics
without the publication intent.

These files are not installed CRDs. A future Kubernetes adapter can convert
the API types after adding Kubernetes metadata types, structural schemas,
status-subresource behavior, generated runtime methods, and reconciliation of
API-server-owned metadata.

Run the command through Bazel:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- --help
```

Diagnose and finish an interrupted publication:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- doctor --goal-dir $GOAL_DIR
bazel_agent bazel run //projects/goal/cmd/goal -- recover --goal-dir $GOAL_DIR
```

The historical skill source remains at `skills/goal` for compatibility
documentation. Its package is not registered for discovery; the replacement
skill is [`openspec`](../../tools/openspec/skills/openspec/SKILL.md).
