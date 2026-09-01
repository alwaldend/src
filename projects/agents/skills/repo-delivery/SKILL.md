---
name: repo-delivery
description: >-
  Finalize task-owned repository changes by verifying, committing, pushing the
  feature branch, maintaining its pull request, and resolving review
  comments. Use when an implementation is ready for delivery; do not use for
  read-only reviews.
---

# Deliver repository changes

Follow the repository instructions for scope and validation. Once requested
implementation work is ready, deliver it without asking whether to commit or
push.

Before delivery, run
`bazel_agent run //tools/repo_delivery -- provider`. Use its sanitized
repository identity, provider hint, `adapter_available`, sanitized
`git_transport`, and `delivery_transport_available` results; do not print or
inspect the raw configured remote URL because it may contain credentials.
Forge support and delivery transport support are separate: proceed only when
both availability fields are true when using the `repo_delivery` adapter.
Route GitHub to the supported adapter below. Route Forgejo directly to its
compatibility workflow; its expected `adapter_available=false` result is not a
blocker for that fallback, whose own transport preconditions still apply. The
current tool has no Forgejo adapter. Stop rather than guess when the provider
or the pull request's actual base cannot be identified safely, and never
create a duplicate pull request on another forge.

The GitHub adapter is limited to same-repository pull requests. It cannot
reliably discover or map an upstream pull request whose head lives in a fork.
Treat same-repository topology as a caller-enforced precondition. Except for
the sanitized `provider` command, do not invoke `repo_delivery` for a
fork-based or otherwise cross-repository pull request, and stop when remote or
pull-request ownership is uncertain. An adapter refusal for an observed
mismatch does not prove that an undiscovered fork topology is safe.

## Correctness revalidation

Any behavior-changing code edit after the latest correctness or review verdict
invalidates that verdict, independently of exact-commit test invalidation.
Before preparing or republishing the changed candidate, perform fresh,
diff-focused correctness scrutiny against the requested behavior and the
contracts touched by the edit. Actively seek disconfirming cases relevant to
the change, such as alternate and fallback paths, boundary and encoding
semantics, malformed or partial state, concurrency and lifecycle transitions,
and platform or implementation parity. Keep the review proportional to the
changed behavior rather than reopening unrelated accepted code.

Passing tests do not substitute for this scrutiny: tests establish only their
encoded cases. Turn valid findings into focused regression coverage, implement
the corrections, and scrutinize the resulting diff again before proceeding.
Documentation-only delivery records do not invalidate a behavioral verdict
unless they can affect execution or the published interface.

## GitHub adapter

Use `bazel_agent run //tools/repo_delivery -- ...` for `inspect`, `prepare`,
`publish`, `verify`, and the `review` subcommands. The tool owns deterministic
mechanics: exact ref and pull-request discovery, aggregate commit creation,
signing preservation, rebasing, exact-lease pushes, provider-CLI calls,
commit-to-pull-request projection, review mutations, disclaimers, and final
invariants.

For GitHub, do not reproduce supported mutations with direct Git or `gh`
commands, select another provider to work around a failure, or bypass a safety
refusal. Apply the same rule to any future provider supported by a
`repo_delivery` adapter. Task ownership, validation selection, conflict
resolution, and review judgment remain agent responsibilities.

Version 1 Git delivery requires both captured remote endpoints to use SSH,
either SCP-style or `ssh://`. The `gh` CLI remains the GitHub forge adapter;
its mutable global `gh auth setup-git` credential-helper configuration is
deliberately not imported into Git transport. If `provider` identifies a
GitHub repository whose Git endpoint is HTTPS, stop before `inspect` and ask
the user to configure an SSH remote rather than changing their repository
configuration or weakening transport isolation implicitly.

### Inspect and prepare

1. Run `inspect`, supplying the intended base when no pull request exists.
   Review its exact refs, OIDs, pull request, feature-commit authorship,
   status, and refusal reasons alongside the task-scoped diff.
2. Decide which paths or hunks are task-owned. Preserve unrelated changes.
   Pass fully task-owned paths to `prepare`, or pre-stage only safe hunks and
   use `--use-index`; never blanket-stage the worktree.
3. Ownership determines rewrites: an agent-owned branch (every commit in the
   range is task or agent generated; no shared, stacked, human-owned, or
   unrelated work) may be consolidated or rewritten freely, including with
   `--consolidate <literal-inspect.local_head_oid>`, as long as both local
   and remote progress are preserved — nothing may be lost from either side.
   Review every listed commit and the task-scoped diff before a rewrite. Pass
   the exact reported local OID to `--rewrite` only after that judgment. With
   a pending authorized remote replacement, the pull request must still match
   an exact projectable local or fetched remote commit projection; unrelated
   text remains a refusal. The adapter additionally requires a merge-free
   linear chain, identical author and committer identities, the ownership
   disclaimer on the oldest commit, and a pull-request projection exactly
   matching the requested aggregate message.
   Never use consolidation or other rewrites for shared, stacked, human-owned,
   unrelated, or ambiguous history. Stop and ask the user when ownership is
   uncertain.
   During a rebase, an expected aggregate path may disappear only when the
   prior candidate and fetched base contain the exact same Git tree entry;
   added paths, non-identical loss, and an empty aggregate remain refusals.
   Carry the reduced exact path set in the derived receipt.
   A divergent remote feature tip is refused by default. Pass
   `--replace-remote <literal-inspect.remote_head_oid>` only after
   `$git-rebase-remote` has preserved that exact old remote tip and the new
   tip retains every previous remote commit's reachable progress. Never infer
   the OID or use this authorization for shared, stacked, human-owned,
   unrelated, or ambiguous history. The final fresh snapshot and receipt must
   retain that same OID as the later force-with-lease expectation.
4. Put the aggregate commit-message file and preparation receipt under the
   same ignored `out/<task>/` area. Pass the receipt path with
   `prepare --receipt-file`. Describe the complete change and verification
   actually performed, not merely the latest amendment; the tool adds the
   commit disclaimer. Checks run before `prepare` may guide preparation, but
   they are preliminary and do not establish publish readiness.
5. Run `prepare` and inspect the prepared base-to-head diff, receipt, and
   immutable aggregate path scope. Treat the top-level returned `head_oid` as
   the candidate to validate. Run every check required for delivery after
   `prepare`, and confirm HEAD still equals that exact OID after the checks.
   Record that literal OID; never recompute it from mutable `HEAD` when
   publishing. Run checks from a clean checkout of that exact commit. If
   unrelated dirty files could affect a check, use an ignored isolated linked
   worktree at the literal OID or prove that the check cannot consume them.
   Only this post-prepare validation can establish publish readiness.
6. If the diff, a failed check, a rebase, or conflict resolution requires
   another preparation, repeat the post-prepare validation against the new
   top-level `head_oid`. Use
   `prepare --message-only --rewrite <exact-oid>` when only the aggregate
   message needs refreshing. Every consolidation or message-only amendment
   changes HEAD, so it invalidates the prior exact-OID gate and requires the
   checks to run again against the newly
   returned top-level `head_oid`.

The current adapter aborts and removes an isolated rebase when it encounters a
conflict; it does not expose an interactive resolver. Stop for direction rather
than bypassing the tool with direct Git. Also stop on a changed lease,
ambiguous ownership, human-edited pull-request metadata, an unsafe base, or any
other fail-closed result.

### Publish and review

After the post-prepare checks pass, pass that exact literal
top-level `head_oid` and the emitted receipt to
`publish --validated-head <commit-oid> --receipt-file <path>`. The receipt,
not a separate command-line lease, binds the repository, remote endpoint
digests, forge, refs, base, prepared tree and head, expected remote head,
expected pull-request identity or absence, and aggregate scope. Do not edit or
transfer it. Treat its ignored local path as trusted and do not let another
process write it while `prepare`, `publish`, or receipt-bound `verify` runs.
The adjacent lock serializes receipt transitions and reads performed by
cooperating `repo_delivery` processes; stable rechecks detect already-visible
changes but are not a portable atomic pathname-content compare-and-swap
against a same-user writer that ignores the lock.

The tool refetches before pushing. If it exits nonzero with a structured
`revalidation_required` report, it has rebased the exact candidate in an
isolated worktree, updated the receipt, and stopped before the push. Validate
the report's exact new head directly, then retry `publish` with that literal
OID and the derived receipt. Do not run `prepare` again unless content or the
message must change, because another prepare creates another unvalidated OID.
Do not post routine progress or completion comments.

The feature-ref lease and base check are separate rather than a cross-ref
compare-and-swap. A base advance in the narrow push window can expose the new
feature head briefly before detection and an exact guarded rollback. Do not
describe the whole operation as atomic, and account for webhook consumers.

Every network Git operation runs in a fresh, private, config-free bare
transport context bound to the captured SSH endpoint and the repository's
canonical object database. It does not load the main repository's mutable
local/worktree routing, trust, URL-rewrite, or remote configuration. Normal
DNS resolution and SSH host-key verification still apply; do not describe the
transport as pinning an IP address or server key.

Use `repo_delivery review inspect` to read a bounded structured inventory of
all threads, reviews, review requests, and top-level comments. GitHub does not
offer an atomic multi-connection snapshot, so reinspect before mutation and
after every result. Evaluate each unresolved item against the request,
repository evidence, and current diff. Carry the exact pull-request node ID,
head OID, and pull-request expectation digest from the latest inspection into
every review mutation. For a top-level reply, also carry the last-comment
sentinel and top-level inventory digest into `review comment`. For a thread
reply or resolution, carry its thread ID, last-comment ID, and expectation
digest. Put comment bodies under ignored `out/<task>/` and pass them with
`--body-file`. `review reply` also writes a strict reply receipt in that task
directory; pass that exact receipt to `review resolve`. Use
`review request --reviewer <login>` for explicit user accounts. Never invoke
`gh` directly for these mutations.

When inspection shows that an enabled remote review has started or is still
running for the exact final head, wait for it to reach a terminal state before
declaring delivery complete. Use the product's authoritative wait or monitoring
mechanism when available. Poll `review inspect` at a modest interval only when
its selected adapter exposes the review execution state needed to distinguish
running, completed, failed, and cancelled outcomes for that exact head. Do not
infer completion from the absence of a pending review, a quiet interval, or an
unchanged review inventory. If no available mechanism exposes that state, make
a bounded observation attempt, keep the user informed, and report the remote
review result as unverifiable instead of waiting indefinitely or claiming it
passed. Do not retrigger the review merely because it is still running. After
it finishes, reinspect and evaluate every finding through the workflow below.
If the review fails, is cancelled, or cannot reach a terminal state because of
an external blocker, report that state explicitly rather than claiming that
review passed. A new published head invalidates completion observed for an
older head and requires waiting for any review started for the new exact head.

- For valid feedback, implement the fix, prepare the aggregate commit again,
  rerun all checks that establish publish readiness against the exact new
  top-level `head_oid`, publish, then reply with what changed and why before
  resolving the thread.
- For invalid or inapplicable feedback, reply with the concrete reason no
  change is needed, then resolve it when supported.
- Never resolve a comment silently. Resolution must use the receipt from a
  reasoned tool reply with the exact comment disclaimer. Preserve a top-level
  comment that cannot be resolved after leaving the reasoned reply.

The reply receipt grants one-use resolution authority for no more than five
minutes. `review resolve` checks that window and consumes the receipt before
attempting the provider mutation, preventing reuse after an attempted
resolution. The receipt binds the complete bounded review inventory
projection: every top-level comment, review, thread and thread comment, and
review request. Its parent
pull-request `UpdatedAt` is a nondecreasing floor because reply side effects
can advance that timestamp asynchronously without changing the bound
inventory. If resolution fails, expires, reports an unknown outcome, or finds
an inventory mismatch, authority is consumed; do not recreate or reuse that
receipt. Reinspect, determine the actual remote state, and, only when
resolution is still appropriate, leave a fresh reasoned reply to obtain new
one-use authority.

GitHub exposes neither compare-and-swap resolution nor a documented monotonic
review-thread epoch. A human resolve followed by unresolve is therefore
unobservable if it restores identical state before the receipt's first use
inside its five-minute window. Exact-state checks, expiration, and one-use
consumption narrow this risk; they do not prove that such an ABA or a change
to omitted provider metadata did not occur.

GitHub review mutations do not provide an atomic compare-and-swap. The tool
checks exact expectations immediately before each mutation and verifies the
complete result afterward, but a concurrent human edit may still land inside
that narrow provider window. Avoid an actively coedited pull request and
reinspect after every mutation or ambiguous failure; never retry blindly with
stale IDs or digests.

After a review-driven rewrite, use `repo_delivery review request` to request
review for the new commit SHA, or use `review inspect` to verify that review
already covers that exact SHA. After the final publish, inspect the review
state, wait for any review started for the exact final head to finish, then
reinspect and handle its findings. Run `verify --receipt-file <path>` and
confirm no requested change remains uncommitted.

## Forgejo compatibility workflow

Except for its required read-only `provider` command, do not invoke
`repo_delivery` for a Forgejo remote until it has a Forgejo adapter. Use
`$git-rebase-remote` for exact ref discovery, synchronization, rebasing, and
an exact-lease push. Use
`bazel_agent run //tools/fj -- ...` for only the basic pull-request operations
that the installed `fj` v0.5.0 exposes: `pr search`, `pr create`, `pr view`,
`pr view ... comments`, `pr edit ... title`, `pr edit ... body`, and
`pr comment`. Check the command's help rather than assuming newer capabilities.

The Forgejo fallback requires a completely clean worktree and index before
`$git-rebase-remote`; it does not preserve dirty files or autostash them.
Preserve unrelated work by stopping until its owner moves or commits it. Stage
only task-owned paths or hunks, and maintain exactly one non-merge task-owned
feature commit. Keep its aggregate message and signing correct. A commit,
amend, rebase, or conflict resolution invalidates
affected earlier checks; rerun them against the resulting tree before the
lease-protected push. Follow `$git-rebase-remote` through its final fetch,
explicit lease, push, and Git-level verification. Stop rather than rewrite
shared, stacked, human-owned, or ambiguous history.

Use `fj pr search` and `pr view` to inspect candidate pull requests before
creating one. Create the exact-head pull request only when its absence is
clear. For an agent-owned pull request, set the title from Git's `%s` and the
body from `%b`, replacing only its final commit-disclaimer line with the pull
request disclaimer. Use `pr edit` for those fields; preserve human-authored
metadata by stopping for direction. Treat `fj`'s human-oriented output as
observational context only: never scrape or parse it as structured proof, and
do not claim exact metadata or review postconditions that the command cannot
verify mechanically.

The exact ownership markers are
`LLM-disclaimer: This commit was generated by an LLM.` for the commit and
`LLM-disclaimer: This pull request was generated by an LLM.` for the pull
request. Do not improvise their spelling or padding.

`fj` v0.5.0 can display and add top-level pull-request comments, but this path
does not provide structured review-thread inspection, thread replies, thread
resolution, or review requests. Evaluate feedback visible through
`pr view ... comments`. After a valid fix, amend, revalidate, synchronize, and
push through `$git-rebase-remote`, then leave a reasoned top-level
`pr comment`. For invalid feedback, leave the concrete reason as a top-level
comment. Do not claim that either comment is an in-thread reply, and leave the
original thread or unresolved state intact. Report these review limitations
explicitly rather than silently skipping or falsely completing them.

## Disclaimers and handoff

End every external comment with
`LLM-disclaimer: This comment was generated by an LLM.` End an issue
description with
`LLM-disclaimer: This issue was generated by an LLM.` Give any other
LLM-generated user-facing artifact its corresponding type-specific final
`LLM-disclaimer:` line. Documentation, code, binaries, and conversational
responses do not need one.

In the final report, name affected files and every verification command with
its actual result. Distinguish failures and provider limitations from completed
checks, and do not claim broader coverage.
