---
name: git-rebase-remote
description: >-
  Keep a feature branch synchronized with its remote base by fetching current
  refs, rebasing task-owned commits, and updating the remote feature ref with
  an explicit force-with-lease. Use before delivery or whenever an advancing
  remote base requires a rebase; agent-owned branches may be rewritten as long
  as neither local nor remote progress is lost, and shared or human-owned
  history is never rewritten.
---

# Synchronize a feature branch with its remote

Preserve the caller's requirements for commit count, signing, validation, and
pull-request state. Inspection alone does not authorize a rebase or push;
perform the full workflow only when the task authorizes those mutations.

## Establish safe refs

1. Determine the full remote feature ref and the actual target base. Use the
   pull request's configured base when one exists; do not guess from the
   forge's default branch. A fork pull request can have different base-fetch
   and feature-push remotes.
2. Require a clean worktree and index with no Git operation in progress. Do not
   discard, auto-stash, commit, or rewrite unrelated work.
3. Inspect branch relationships and every commit that the rebase would
   rewrite. Stop when the branch appears shared, stacked, human-owned, or mixed
   with unrelated commits, or when ownership is materially uncertain. On an
   agent-owned branch (all commits task or agent generated), rewriting is
   authorized as long as both local and remote progress are preserved: the
   rebase replays every task-owned commit, and the final push uses an exact
   force-with-lease so no remote commit is lost without detection.

## Fetch and rebase

1. Refresh the pull-request metadata, then fetch the exact base and feature
   refs without fetching tags. Distinguish a confirmed missing feature ref
   from a network, authentication, or fetch failure.
2. Record the fetched base OID and feature-ref OID as immutable values; do not
   rely on a remote-tracking ref remaining unchanged. If the feature ref is
   confirmed absent, record that state instead.
3. After fetching, inspect the local and fetched feature tips and both unique
   commit ranges. Establish the expected remote feature state only after
   confirming that replacing it would not overwrite unexpected, shared,
   human-owned, or unrelated work, and that every previously remote task-owned
   commit remains reachable from the replacement.
4. Rebase the intended task-owned commits onto the fetched base OID. Preserve
   required signing and any caller-required commit shape. If a previously
   fetched base changes non-fast-forward, stop and reassess instead of blindly
   replaying onto it.
5. Resolve conflicts only when the correct result is supported by the task and
   repository evidence. Stage explicit resolutions and continue. Abort when a
   conflict cannot be resolved without guessing or importing unrelated work.
6. Rerun checks invalidated by the rewrite and inspect the complete
   base-to-head diff. Reassert every caller-required pre-push invariant,
   including commit count and absence of merge commits. If no feature change
   remains relative to the base, stop without inventing an empty commit.

## Fetch again and push with a lease

Treat fetch, comparison, rebase, and validation as a cycle. After any rebase or
check run, refresh the pull-request metadata and fetch both refs again. If the
base advanced, rebase onto its new fetched OID and repeat invalidated checks.
If the feature ref changed or its existence state changed, stop and inspect
the remote change: when the change is task-owned agent progress, rebase it in
so neither local nor remote progress is lost; otherwise stop instead of
overwriting or recreating it. Push only when a fresh cycle requires no further
rewrite and all pre-push invariants still hold.

For an existing rewritten feature ref, push `HEAD` to the full feature ref
with the exact OID observed by that final fetch:

```sh
git push <feature-remote> HEAD:<feature-ref> \
  --force-with-lease=<feature-ref>:<observed-feature-oid>
```

For a feature ref confirmed absent throughout the workflow, use an explicit
lease that requires it still to be absent:

```sh
git push <feature-remote> HEAD:<feature-ref> \
  --force-with-lease=<feature-ref>:
```

Never use `--force`, a bare `--force-with-lease`, or an expected value read
later from a mutable remote-tracking ref. On a lease failure, stop, fetch, and
inspect the remote change. Do not weaken or refresh the lease merely to
overwrite that change.

After pushing, fetch both refs once more. Verify that the fetched base is an
ancestor of the feature head and that local `HEAD` equals the fetched remote
feature-ref OID. Also run any caller-required commit-count or pull-request
checks.
