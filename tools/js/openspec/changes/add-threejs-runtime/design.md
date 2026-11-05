## Context

The merged TypeScript setup provides strict compilation and `@types/three`.
PR #211 currently introduces the runtime alongside the complete renderer.
The tools pnpm workspace owns npm pins and Bazel package links.

## Goals / Non-Goals

Provide an independently validated runtime prerequisite on master. Keep product
rendering, adapters, viewer code, and snapshots in their existing project.

## Decisions

- Add the approved exact runtime pin through `//tools/pnpm`, with lockfile-only
  resolution and lifecycle scripts disabled. Reuse the current package store
  and integrity-checked lock rather than a new download mechanism.
- Extend the existing TypeScript E2E instead of adding another test harness.
  First require the expected length of a Three.js vector in its Node result;
  then import the real runtime in the fixture. Preserve independent compiler,
  browser, package metadata, and declaration checks. Record the runtime revision
  in the existing JSON artifact without making its spelling an acceptance gate.
- Start this PR from current master. The renderer consumes the pin after this
  prerequisite merges; separate PRs do not require a stacked base.

## Risks / Trade-offs

- Runtime and declaration compatibility: compile the actual import under strict
  NodeNext settings and execute it using the pinned Node runtime.
- Lockfile drift: regenerate with the owning pnpm command and inspect the full
  diff for unrelated version changes.
- Supply-chain maintenance: retain the exact approved pin and registry integrity;
  review future upgrades through the existing dependency workflow.
