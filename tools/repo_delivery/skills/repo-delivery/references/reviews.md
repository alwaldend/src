# Review findings and guarded mutations

Read when feedback requires a reply, resolution, or review request. Commands
use the invocation selected in the skill entry point.

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
`gh` directly for these mutations. After a reply, resolve with the receipt's
`result_thread_digest` and `reply_comment_id` (the post-reply thread state
and appended comment), not the pre-reply inspection values.

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
inventory. For a failed provider read conclusively before any resolution
mutation, the tool may restore the original receipt within its unchanged
expiration window. Retry without another public reply only when the tool
explicitly reports that restoration; never restore or extend it yourself.
Expiry, semantic or inventory mismatch, mutation failure, and an unknown
mutation outcome consume authority. Reinspect, determine the actual remote
state, and, only when resolution is still appropriate, leave a fresh reasoned
reply to obtain new one-use authority.

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
