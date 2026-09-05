# Rewrites and recovery

Read before consolidation, remote replacement, manual conflict resolution, or
recovery. Load `$git-rebase-remote` when synchronizing an advanced base. Commands
use the invocation selected in the skill entry point.

For a prepared, task-owned, single-commit GitHub branch and an advanced base,
prefer `repo_delivery rebase` over manual Git. Establish validation against
the resulting exact candidate before publication. Consolidation uses
`prepare --consolidate`; do not combine it with `--rewrite`.

Ownership determines rewrites: an agent-owned branch (every commit in the
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
Consolidation additionally requires the oldest feature commit to carry the
ownership disclaimer trailer, so add it to that commit before the first
`prepare` when a branch predates the convention; the adapter's error names
the exact commit when it is missing. Divergent remote tips require a typed
`rewrite-authorize` receipt (its `--old-remote-oid` is the freshly
inspected `remote_head_oid`), passed to `prepare --rewrite-authorization`;
do not combine it with `--replace-remote`, which is the untyped, direct
authorization form.
Force-pushing an agent-owned feature branch is allowed only as an
exact-lease update that preserves every previously remote commit's
reachable progress: first fetch the current remote feature tip and record
its OID, verify the replacement preserves every prior remote commit's progress
(amends and rebases of task-owned commits may change ancestry), then push
`HEAD` to the feature ref with
`--force-with-lease=<ref>:<observed-remote-oid>`. Never use a bare
`--force`, and never force-push shared, stacked, human-owned, or
ambiguous history.
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

If any Git history operation loses a commit, tree state, or a resolved
conflict, recover with `git reflog` instead of reconstructing the state
by hand: identify the last known-good `HEAD@{n}` entry, verify it with
`git diff <known-good-oid> <candidate-oid>` before resetting, and preserve any
current unrelated or uncommitted work before restoring the verified state.
Record the OID chosen
for recovery in the task notes, then restart the failed preparation step
from that exact state.

The current adapter aborts and removes an isolated rebase when it encounters a
conflict; it does not expose an interactive resolver. After inspecting the
exact three-way hunks, a manual rebase is allowed only when repository policy
classifies every conflict as minor and task plus repository evidence makes the
combined result unambiguous. Record the fetched base and feature OIDs, preserve
both sides, stage only the explicit resolutions, and rerun invalidated checks.
Stop for direction when any conflict is not minor. Also stop on a changed
lease, ambiguous ownership, human-edited pull-request metadata, an unsafe base,
or any other fail-closed result.

The adjacent lock serializes receipt transitions and reads performed by
cooperating `repo_delivery` processes; stable rechecks detect already-visible
changes but are not a portable atomic pathname-content compare-and-swap
against a same-user writer that ignores the lock. A failed `prepare` may
leave both the receipt and its adjacent `.lock` file. Diagnose the failure and
confirm no preparation is active before recovery. Prefer a new task-owned
receipt path; when reusing the old path is necessary, remove only those exact
task-owned files after preserving needed evidence. Never remove an active lock
or treat deleting a receipt as authorization to bypass a refusal.
