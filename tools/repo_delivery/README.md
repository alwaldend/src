---
title: Repository delivery
description: Guarded Git and pull-request delivery
languages:
  - go
tags:
  - git
  - github
---

`repo_delivery` implements the deterministic part of repository delivery. It
invokes native Git and a selected forge CLI; it does not infer task ownership,
choose validation commands, resolve conflicts, or judge review feedback.

## Provider and delivery workflow

Detect the provider without printing a credential-bearing remote URL:

```sh
bazel_agent run //tools/repo_delivery -- provider
```

The sanitized report distinguishes forge support from Git transport support.
`adapter_available` reports whether the forge has an adapter;
`git_transport` reports `ssh`, `https`, or `mixed_or_unsupported`; and
`delivery_transport_available` is true only when both captured endpoints are
canonical SSH endpoints. It never reports either endpoint.

For the supported GitHub adapter, the normal workflow is:

```sh
bazel_agent run //tools/repo_delivery -- inspect --base master
bazel_agent run //tools/repo_delivery -- prepare \
  --base master \
  --message-file out/task/commit.md \
  --receipt-file out/task/prepare.json \
  --path path/to/task/file \
  --rewrite <inspect.local_head_oid> # omit when the range has no commit
# For an explicitly reviewed task-owned multi-commit range, use
# --consolidate <inspect.local_head_oid> instead of --rewrite.
# Run every required validation against the top-level literal head_oid.
bazel_agent run //tools/repo_delivery -- publish \
  --base master \
  --receipt-file out/task/prepare.json \
  --validated-head <literal-head_oid>
bazel_agent run //tools/repo_delivery -- verify \
  --receipt-file out/task/prepare.json
```

`prepare` refuses a divergent remote feature tip by default. Use
`--replace-remote <literal-inspect.remote_head_oid>` only after the
`git-rebase-remote` workflow has preserved that exact old tip and established
that it is task-owned history the user authorized rewriting. Never infer the
OID or use this escape hatch for shared, stacked, human-owned, unrelated, or
ambiguous history. A mismatched, malformed, absent, or unnecessary
authorization is refused. The final fresh snapshot and preparation receipt
bind the same old remote OID for the later exact force-with-lease push.

`prepare --consolidate <literal-inspect.local_head_oid>` is the explicit
ownership authorization for replacing a multi-commit feature range with one
aggregate commit. It is not inferred from author names. The adapter still
requires a merge-free linear chain, one author and committer identity across
the range, the ownership disclaimer on its oldest commit, pull-request
metadata matching the requested aggregate message, and the exact inspected
head.
It preserves any signature requirement found in the range and keeps the
fetched remote feature tip as the publication lease. `--consolidate` and
`--rewrite` are mutually exclusive.

When a clean replay onto an advanced base makes an expected aggregate path
disappear from the resulting diff, the adapter accepts that shrink only when
the prior candidate and the new base have byte-identical Git tree entries at
that path. A new path, a non-identical disappearance, or an entirely empty
aggregate remains a refusal. The derived receipt records the reduced exact
aggregate path set.

Never populate `--validated-head` by resolving the current `HEAD` during
publication. Carry the literal OID returned by `prepare`; otherwise a checkout
change after validation could authorize an unvalidated commit.

Version 1 delivery accepts only SCP-style or `ssh://` Git fetch and push
endpoints. The `gh` CLI is still used for the GitHub forge API, but the tool
does not import mutable global credential-helper configuration installed by
`gh auth setup-git`. An HTTPS Git remote is therefore rejected before
`inspect`, preparation, or publication performs network access. Configure an
SSH remote explicitly before using delivery; the tool will not rewrite the
repository's remote configuration for you.

`prepare` writes a strict versioned receipt under an ignored
`out/<task>/` directory. The receipt binds the prepared head and tree to the
exact worktree and Git directories, fetch and push endpoint digests, sanitized
repository identity, forge adapter, base and head refs, fetched base, expected
remote feature ref, expected pull-request identity or absence, and immutable
aggregate path scope. The file contains no
remote URL, must remain untracked, and is installed with mode 0600 by atomic
rename. It is a consistency record, not an unforgeable authorization token.

The adjacent receipt lock serializes receipt transitions and reads performed
by cooperating `repo_delivery` processes. Keep the ignored receipt path
trusted and do not edit, replace, copy over, or otherwise write it while
`prepare`, `publish`, or receipt-bound `verify` is running. Stable byte and
revision checks detect changes visible before their checks, but operating
systems do not provide a portable atomic
pathname-content compare-and-swap. A same-user process that ignores the lock
can race the final comparison and rename; the tool does not claim to exclude
that writer. A persistent exclusive revision-claim sidecar would have the same
trust boundary, while a crash before installation could strand the receipt
without a safe automatic cleanup decision, so it is not used.

Use repeated `--path` flags for fully task-owned paths, or pre-stage only
task-owned hunks and use `--use-index`. Both modes bind the complete existing
feature diff, not merely the paths newly staged by that invocation.
`--message-only --rewrite <exact-oid>` preserves the tree but changes the
commit OID. Consolidation also changes the commit OID and parent structure.
Every prepare, consolidation, or message-only amendment therefore requires
fresh validation against its returned exact head.
When an authorized remote replacement is pending, rewrite evidence accepts an
existing pull request only if its metadata matches the exact projectable local
or fetched-remote commit projection. A legacy remote tail that lacks an
aggregate disclaimer cannot block a matching local aggregate, and unrelated
pull-request text remains a refusal.

Publish and receipt-bound verify require a clean index and reject staged,
unstaged, or untracked changes in the prepared task scope. They preserve
unrelated unstaged and untracked files when no rebase is needed. A required
rebase demands a fully clean worktree and index. Run validation from a clean
checkout at the literal `head_oid` whenever unrelated dirty files could affect
the check; the tool cannot infer a check's input set.

If the base advances after preparation, `publish` rebases the exact captured
commit in an isolated ignored worktree. It then exits nonzero before pushing
and emits a `revalidation_required` JSON report containing the new exact head,
tree, and derived receipt. Validate that returned head directly, then retry
`publish` with the same receipt file and the new literal OID. Do not run
`prepare` again unless content or the commit message must change. A receipt
also supports a guarded retry when the remote already equals its prepared head
but pull-request creation or metadata synchronization stopped partway through.
Replacement pull-request identities and any state other than the exact prior
or desired projection are refused. Multi-commit consolidation likewise
requires an existing pull request to equal the requested aggregate message's
projection, so consolidation cannot silently replace independently edited
pull-request text.

The tool creates commits with Git plumbing and pushes one exact ref with hooks
disabled. It rejects shallow, promisor, or grafted history and ignores
replacement objects. Each network operation runs in a fresh private,
config-free bare Git directory bound to the captured SSH endpoint and the
repository's canonical object database. It does not load mutable main-repo
local/worktree URL rewriting, proxy, TLS, SSH-command, or remote
configuration. Normal DNS resolution and SSH host-key verification still
apply; the tool does not claim to pin an IP address or server key. Ordinary
local Git operations still intentionally use repository identity and signing
configuration. Fetches use isolated temporary refs without pruning, then
install exact private refs locally; pushes use an exact force-with-lease.
Repository-required hook checks remain the caller's responsibility and must
complete before recording the validated OID.

The feature-ref lease and base-ref checks are separate because Git forges do
not offer a cross-ref compare-and-swap. If the base advances in the narrow push
window, the feature ref may be visible briefly before the tool detects the
race and attempts an exact guarded rollback. Coordinate consumers that react
immediately to branch-update webhooks.

## Review workflow

Start with a bounded structured inventory:

```sh
bazel_agent run //tools/repo_delivery -- review inspect
```

Carry values from the latest inspection literally. All mutations require the
pull-request node ID, expected head, and pull-request expectation digest. A
top-level reply also requires the
reported last-comment sentinel and top-level inventory digest:

```sh
bazel_agent run //tools/repo_delivery -- review comment \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --expected-last-top-level-comment \
    <expected_last_top_level_comment> \
  --expected-top-level-comments-digest \
    <expected_top_level_comments_digest> \
  --body-file out/task/comment.md
```

Thread replies and resolutions additionally bind the thread, its last comment,
and its complete expectation digest:

```sh
bazel_agent run //tools/repo_delivery -- review reply \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --body-file out/task/reply.md \
  --reply-receipt-file out/task/reply.json

bazel_agent run //tools/repo_delivery -- review resolve \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --reply-receipt-file out/task/reply.json

bazel_agent run //tools/repo_delivery -- review request \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --reviewer <login>
```

Comment and reply bodies must be ignored files under `out/<task>/`; the tool
adds the exact comment disclaimer. Resolution requires the strict reply
receipt written by `review reply`, and the receipt must still match the exact
pull request, head, thread, reply, body, and complete bounded review inventory
projection. The parent pull-request `UpdatedAt` is a nondecreasing floor,
rather than an exact value, because reply side effects can advance it
asynchronously; every review, top-level comment, thread and thread comment, and
review request remains bound.
Reinspect after every mutation and use the newly reported IDs and digests for
the next one.
The one-use reply authority expires no more than five minutes after issuance.
Resolution verifies the authority window and atomically consumes the receipt
before contacting the provider. On expiration, failure, or an outcome-unknown
result, including any full-inventory mismatch, the authority is consumed;
never recreate or reuse the receipt. Reinspect the remote state and, if
resolution remains appropriate, leave a fresh reasoned reply to obtain new
one-use authority.

GitHub provides neither compare-and-swap resolution nor a documented monotonic
review-thread epoch. A human resolve followed by unresolve can therefore be
unobservable if it restores identical thread and pull-request state inside the
five-minute authority window. The exact-state checks, short expiration, and
one-use consumption narrow that risk; they do not prove that such an ABA or a
change to provider metadata omitted from the bounded projection did not occur.
Do not resolve a thread that may be concurrently moderated.

GitHub does not expose atomic compare-and-swap mutations for pull-request
metadata, comments, replies, resolutions, or review requests. The adapter
checks exact expectations before each mutation and verifies the complete
result afterward, but a concurrent human change can still occur inside that
narrow provider-side window. Do not mutate an actively coedited pull request.
Post-mutation inspection failures are reported as outcome unknown; reinspect
instead of retrying blindly.

## Forge adapters

GitHub is supported through the `gh` CLI. Forge behavior is selected behind a
Go interface and uses argument-vector subprocess calls, so another structured
CLI adapter can be added without changing delivery policy. `--forge auto`
selects GitHub for `github.com`; use
`--forge github --forge-cli <gh-path>` to select an explicit executable.
Unknown or unsupported forges fail before any commit rewrite or push.

Version 1 is limited to same-repository pull requests. It cannot reliably
discover or map an upstream pull request whose head lives in a fork. Treat
same-repository topology as a caller-enforced precondition: do not use the
delivery or review commands for a fork-based or otherwise cross-repository
pull request, and stop when remote or pull-request ownership is uncertain.
The adapter fails closed when it observes an unsupported or inconsistent
topology, but that refusal is not a substitute for this caller check.

The read-only `provider` command recognizes known Forgejo hosts, but there is
currently no Forgejo delivery adapter. Use the repository's documented
Forgejo compatibility workflow instead of pretending GitHub operations are
portable.
