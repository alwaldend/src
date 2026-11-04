---
title: Repository delivery
description: Guarded Git and pull-request delivery
languages:
  - go
tags:
  - git
  - github
---

`al delivery` implements the deterministic part of repository delivery. It
invokes native Git and a selected forge CLI; it does not infer task ownership,
choose validation commands, resolve conflicts, or judge review feedback.

Status and immutable changed-path inventories, and preparation receipts,
support complete output up to 4 MiB so large source relocations can use the
guarded workflow. Stable reads and writes enforce that
limit; other typed local records retain their 256 KiB limit. Receipt schema,
scope, integrity, and publication checks apply at either size.

## Build and install

`@rules_al//:al` builds the current local source. Install the standalone
command explicitly:

```sh
bazel run //tools/al:install
al delivery provider
```

Installation atomically replaces `~/.local/bin/al` with mode 0755.
Ensure that directory is on `PATH`. It copies the executable rather than
linking to Bazel outputs, so it remains available after those outputs are
removed. Rerun the target to update an installed command after source changes.
No AL executable cache or binary download is involved.

For an isolated installation, pass an absolute destination directory:

```sh
bazel run //tools/al:install -- --directory "$PWD/out/install/bin"
```

## Provider and delivery workflow

Run delivery commands from this Git worktree's repository root, which owns the
root `MODULE.bazel` and imports `rules_al`. This also applies when changed
files belong to a nested Bazel module: its workspace does not own the delivery
target. Use the same feature worktree, not another checkout. The Bazel entry
point builds AL from the current source. Detect the provider without printing a
credential-bearing remote URL:

```sh
bazel run @rules_al//:al -- delivery provider
```

The installed `al delivery` command is an optional direct entry point. An
execution failure or delivery refusal requires diagnosis. All examples below
use the Bazel entry point.

The sanitized report distinguishes forge support from Git transport support.
`adapter_available` reports whether the forge has an adapter;
`git_transport` reports `ssh`, `https`, or `mixed_or_unsupported` by default,
and `gh` when that mode is explicitly selected.
`delivery_transport_available` is true when both captured endpoints are valid
for the selected mode: canonical SSH endpoints by default, or canonical GitHub
SSH/HTTPS endpoints in gh mode. It never reports either endpoint.

For the supported GitHub adapter, preparation, checks, and publication are
separate commands:

```sh
bazel run @rules_al//:al -- delivery inspect --base master
bazel run @rules_al//:al -- delivery prepare \
  --base master \
  --message-file out/task/commit.md \
  --receipt-file out/task/prepare.json \
  --path path/to/task/file
bazel test //:repo_quality_test @rules_al//internal/delivery:go_test
bazel run @rules_al//:al -- delivery publish \
  --receipt-file out/task/prepare.json
bazel run @rules_al//:al -- delivery verify \
  --receipt-file out/task/prepare.json
```

Use `--rewrite <inspect.local_head_oid>` when replacing an existing owned
single-commit candidate, or `--consolidate <inspect.local_head_oid>` for an
explicitly reviewed task-owned multi-commit range. The inspected value must be
carried literally.

`prepare` performs Git preparation and writes the typed receipt. It does not
invoke Bazel, select affected targets, or produce a validation plan. Run the
normal build, test, and lint commands appropriate to the changed components
after preparation. The test command above is an example for delivery changes,
not a complete check selection for every task. Use each target's owning
workspace, ordinary Bazel options, and Bazel's normal cache behavior.

Inspect the check results before publishing. Failed checks leave the prepared
candidate available for diagnosis; AL does not record a check ledger or enforce
their outcome. After fixing source, prepare the changed candidate again. If
the candidate is unchanged, rerun the relevant ordinary Bazel checks without
creating another commit merely to retry them.

`publish --receipt-file PATH` verifies that current HEAD and tree match the
prepared candidate, requires a clean index and unchanged prepared scope, and
checks the remote lease and PR expectations before publication. The receipt binds Git state; it
does not assert that checks passed.

If push and PR synchronization complete but final verification fails,
`publication_verification_pending` retains the known candidate and PR identity
with `verified: false`. Preserve that report, link the known PR, and use
`verify --receipt-file PATH` to inspect the postcondition before deciding
whether any mutation needs recovery. Do not infer that a failed final read
undid a successful push.

To synchronize a prepared, task-owned, single-commit feature branch with an
advanced base before `prepare`, use the guarded adapter workflow:

```sh
bazel run @rules_al//:al -- delivery rebase --base master
```

The command refuses dirty trees, divergent remote feature tips, multi-commit
or merge-containing ranges, and incompatible pull-request metadata. It fetches
fresh refs, replays the commit in an isolated worktree, preserves signature
requirements, and reports the literal local head and captured remote lease.
It does not push or update the PR. Prepare the resulting candidate, run its checks,
then publish through the receipt-bound workflow. When the remote feature exists,
authorize its exact captured replacement only after confirming that all prior
task-owned progress is preserved. The later publisher refreshes both refs and
uses its exact lease; a rebase report is not validation evidence.

Remote rewrites use `rewrite-authorize` to write a typed, non-authorizing
authorization receipt (old remote OID, new head OID, owner root, task paths,
provider ownership), then `prepare --rewrite-authorization <file>` instead of
raw `--replace-remote` handoff. Review replies accept `--goal-ref`,
`--delivery-ref`, and `--defect-id` durable join references.

`prepare` refuses a divergent remote feature tip by default. Use
`--replace-remote <literal-inspect.remote_head_oid>` only after the
`al-rebase` workflow has preserved that exact old tip and established
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

For a message-only amendment, use
`prepare --message-only --rewrite <exact-oid>` with the new message and a
preparation receipt. Run the checks appropriate to that changed commit before
publishing through the same receipt-bound workflow.

SSH is the default Git transport, accepting canonical SCP-style or `ssh://`
fetch and push endpoints. Select `--git-transport gh` explicitly to use GitHub
CLI credentials over HTTPS. This mode accepts canonical GitHub SSH or HTTPS
remotes, derives the HTTPS endpoint for the captured repository, and uses the
selected `--forge-cli` executable (default `gh`) as a host-scoped credential
helper. It never reads tokens into AL, imports global Git configuration,
runs `gh auth setup-git`, rewrites remotes, or falls back automatically.
Authenticate using `gh auth login` when needed, then pass the same transport
selector throughout inspection, preparation, publication and verification:

```sh
al delivery provider --git-transport gh
al delivery inspect --git-transport gh --base master
al delivery publish --git-transport gh --base master \
  --receipt-file out/task/prepare.json
```

Preparation receipts also bind transport mode; legacy receipts imply SSH.
Changing modes requires preparing a new receipt. HTTPS transport verifies TLS,
disables redirects and proxies, resets ambient credential helpers, disables
prompts and credential diagnostics, and rejects credentials in endpoint URLs.
`gh` authentication uses its configured credential store or supported token
environment variables; secrets do not belong in command arguments or output.

GitHub API preflight uses the selected forge CLI in both Git modes. An HTTP
401 there is independent of `--git-transport`. A sandbox can prevent the CLI
from retrieving a valid stored credential: the same login can succeed outside
it. AL's HTTP 401 diagnostic points to an approved-context comparison; it does
not retry, extract tokens, change login state, or escalate itself. Follow the
delivery skill's [authentication recovery](../../skills/al-delivery/references/authentication.md)
procedure before requesting a new login or changing transport.

`prepare` writes a strict versioned receipt under an ignored
`out/<task>/` directory. The receipt binds the prepared head and tree to the
exact worktree and Git directories, fetch and push endpoint digests, sanitized
repository identity, forge adapter, base and head refs, fetched base, expected
remote feature ref, expected pull-request identity or absence, and immutable
aggregate path scope. The file contains no
remote URL, must remain untracked, and is installed with mode 0600 by atomic
rename. It is a consistency record, not an unforgeable authorization token.

The adjacent receipt lock serializes receipt transitions and reads performed
by cooperating `al delivery` processes. Keep the ignored receipt path
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
Every prepare, consolidation, or message-only amendment therefore requires a
validation decision bound to its returned exact head. Run required checks
after preparation by default. For a tree-preserving amendment, prior passing
evidence may be reused only after recording its exact prior candidate, the
matching tree OID, the new OID, and unchanged inputs relevant to each check.
Those inputs include commands, tools, configuration, and environment. Rerun
checks affected by commit identity, history, or stamping, and any check whose
input stability is unknown. The caller owns and records this applicability
judgment; the receipt neither proves validation nor authorizes reuse.
When an authorized remote replacement is pending, rewrite evidence accepts an
existing pull request only if its metadata matches the exact projectable local
or fetched-remote commit projection. A legacy remote tail that lacks an
aggregate disclaimer cannot block a matching local aggregate, and unrelated
pull-request text remains a refusal.

Publish and receipt-bound verify require current HEAD and tree to match the
receipt, a clean index, and no unstaged or untracked changes in the prepared
task scope. Unrelated unstaged and untracked files are preserved when no rebase
is needed. A required rebase demands a fully clean worktree and index.
Preparation receipts do not make later source edits part of the candidate.

If the base advances after preparation, `publish` rebases the exact captured
commit in an isolated ignored worktree and updates the preparation receipt.
It exits nonzero before pushing with status `checks_required`, identifying
the new head, tree, and derived receipt. Run the normal Bazel checks against
that candidate, then retry `publish --receipt-file PATH` with the same receipt.
Do not run `prepare` again unless content or the commit message must change.

A receipt also supports a guarded retry when the remote already equals its
prepared head but PR creation or metadata synchronization stopped partway
through. Inspect the current remote state first. Replacement PR identities and
any state other than the exact prior or desired projection are refused.
Multi-commit consolidation likewise requires an existing PR to equal the
requested aggregate message's projection, so consolidation cannot silently
replace independently edited PR text.

The tool creates commits with Git plumbing and pushes one exact ref with hooks
disabled. It rejects shallow, promisor, or grafted history and ignores
replacement objects. Each network operation runs in a fresh private,
config-free bare Git directory bound to the captured endpoint, selected
transport, and repository's canonical object database. It does not load mutable main-repo
local/worktree URL rewriting, proxy, TLS, SSH-command, or remote
configuration. Normal DNS resolution and SSH host-key verification still
apply; the tool does not claim to pin an IP address or server key. Ordinary
local Git operations still intentionally use repository identity and signing
configuration. Fetches use isolated temporary refs without pruning, then
install exact private refs locally; pushes use an exact force-with-lease.
Repository-required hook checks remain the caller's responsibility and must
complete before publication.

The feature-ref lease and base-ref checks are separate because Git forges do
not offer a cross-ref compare-and-swap. If the base advances in the narrow push
window, the feature ref may be visible briefly before the tool detects the
race and attempts an exact guarded rollback. Coordinate consumers that react
immediately to branch-update webhooks.

## Review workflow

Start with a bounded structured inventory:

```sh
bazel run @rules_al//:al -- delivery review inspect
```

For authorized feedback on a clean exact prepared candidate whose publication
is verified, use one command:

```sh
bazel run @rules_al//:al -- delivery review address \
  --receipt-file out/task/prepare.json \
  --thread-id <thread-id> --body-file out/task/reply.md
```

`address` checks the exact prepared candidate and verified remote publication,
derives current provider guards, posts the reasoned reply, saves its receipt,
and resolves the thread. It does not run checks or read validation logs.
The default reply receipt path is derived from the preparation receipt
and thread; `--reply-receipt-file` selects an explicit path in the same task
directory. A durable attempt checkpoint beside that file prevents duplicate
replies after interruption, success, or an unknown reply outcome. Retry with
the same paths after an explicitly reported read-failure restoration to resume
resolution using the unchanged receipt and expiry. Expired, consumed, or
changed authority stops without reposting. Inspect the actual remote state
before explicitly selecting a fresh receipt path for recovery. Keep checkpoints
trusted and do not delete them to bypass this protection.

Thread inventory includes each thread's first comment page in its enclosing
request. Only additional comment pages need separate queries. All pages retain
shared byte, node, request, cursor, and metadata consistency checks.

The low-level commands remain available. Carry values from the latest inspection literally. All mutations require the
pull-request node ID, expected head, and pull-request expectation digest. A
top-level reply also requires the
reported last-comment sentinel and top-level inventory digest:

```sh
bazel run @rules_al//:al -- delivery review comment \
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
bazel run @rules_al//:al -- delivery review reply \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --body-file out/task/reply.md \
  --reply-receipt-file out/task/reply.json

bazel run @rules_al//:al -- delivery review resolve \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --reply-receipt-file out/task/reply.json

bazel run @rules_al//:al -- delivery review request \
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
before contacting the provider. If a provider read fails or the review epoch advances before any resolution
mutation is attempted, the tool can restore the original unexpired receipt
without replacing another file or extending its authority window. Only an
explicit restoration report permits retry with that receipt; no additional
public reply is needed for that read failure.
Expiration, semantic or full-inventory mismatch, mutation failure, and an
unknown mutation outcome still consume authority. Never recreate the receipt
yourself. Reinspect the remote state and, if resolution remains appropriate,
leave a fresh reasoned reply to obtain new one-use authority.

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
