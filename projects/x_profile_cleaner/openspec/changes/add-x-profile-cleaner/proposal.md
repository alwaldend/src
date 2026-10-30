## Why

The account owner needs repeatable cleanup of posts, replies, and reposts
without expressing policy through a large configuration object. Lua predicates
and preservation functions make account-specific selection practical.

## What Changes

- Add a Go CLI with plan and apply commands, inline planning, confirmation,
  and `-y`/`--yes`.
- Add a loopback-only serve command with a shared-site-theme browser interface
  for editing/loading policy, reviewing and excluding planned actions, applying
  the selected subset, cancelling work, and observing live progress.
- Bind documented Lua predicate, preservation, and graph-query functions.
- Separate the language-independent planner from the Lua adapter, add strict
  arguments, named explanations, relative time helpers, rich metadata, and
  preservation context; share callback-free traversal within each declaration.
- Execute independent deletion endpoints concurrently with separate quotas,
  same-target dependencies, cancellation, and exact action-list authority.
- Allow optional Lua discovery periods for chunked runs, documenting incomplete
  thread context when seeds or connecting posts lie outside the selected period.
- Preserve reply-connected quote branches while ignoring non-author content.
- Resolve forced deletion above preservation above ordinary deletion.
- Supply an example 30-day post/reply/repost policy and editor shims.
- Document bounded timeline discovery; exclude historical discovery work,
  Articles, and Highlights handling.

## Capabilities

### New Capabilities

- `profile-cleanup`: Author-scoped inventory, Lua selection, graph preservation,
  and explicit plan/apply execution.

### Modified Capabilities

None.

## Impact

New project `projects/x_profile_cleaner`, using the existing Cobra and
GopherLua dependencies and the Go standard library. Runtime X access uses
the Article uploader's existing OAuth 1.0a injection configuration and the
shared SDK signer. The uploader's credential-consumer visibility admits the
cleaner runner. Validation uses an offline HTTP
fixture and never performs live X mutations.

Exclude this unreleased CLI from `projects/BUILD.bazel`'s automatic release
deployment aggregation, following the neighboring article uploader's boundary.
