## Context

See proposal.md. The existing collection store owns Anki restore, export and
lifecycle; archive-build policy belongs outside that lifecycle owner.

## Goals / Non-Goals

The Builder owns only filesystem staging and publication. Native services,
configs, planning and state are supplied through the apply callback.

## Decisions

- Keep `New(apply)` and `Build(ctx, base, output)` as the complete package API.
  A callback avoids coupling this prerequisite to #128's executor or App wiring.
- Stage beside the destination so rename remains atomic. Resolve source/output
  aliases and reject equal files before copying. Do not interpret archive bytes.
- Use existing App state operations in the external-package native E2E test to
  change real archives without requiring unmerged reconciliation code.
- Name the package `rebuild` for its copy-and-apply behavior. It needs no exception
  to the repository's generic build-output ignore rules.
- Preserve both `ctx.Err()` and `context.Cause(ctx)` at startup and before
  publication so callers can match cancellation and its original cause.

## Risks / Trade-offs

- Concurrent external writers → callers own exclusive access; path checks are
  not a security boundary against races.
- Cleanup can fail after publication → preserve the cleanup error; publication
  is atomic visibility, not a power-loss durability guarantee.
- Main #128 temporarily includes prerequisite files → link this independent
  master-based PR, then remove the shared diff by rebasing after merge.
