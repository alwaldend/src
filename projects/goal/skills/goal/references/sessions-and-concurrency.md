# Sessions and concurrency

Use this reference when resuming, switching goals, delegating, or reconciling
a stale update.

## Many records, one focus

A goal is identified by its stable goal ID and scope, not by a chat thread. A
`GoalSessionBinding` is a replaceable pointer to the current goal. Its
attachment-time generation, revision, and digest fields only detect whether
that pointer's view has become stale. List and inspect stored goals before
attaching when identity is uncertain. Starting a new session must not create a
duplicate merely because the old thread is unavailable.

Attaching does not mutate the goal. Detaching or switching focus does not
pause, close, or delete it. Verify the goal resource version and current
attempt after every fresh attach or context recovery.

## Canonical writer protocol

Use one coordinator for canonical state. For an attempt or lifecycle
`checkpoint` mutation within one existing goal record, the goal tool:

1. acquires the selected goal's cooperative cross-process file lock;
2. rereads canonical state under the lock;
3. checks the caller's expected resource version; when publishing an existing
   attempt, it also verifies its goal reference, lifecycle generation, and
   digest-bound criteria snapshot;
4. validates the complete prospective state and fully stages a new attempt in
   a hidden sibling directory when needed;
5. advances `goal.yaml` as the optimistic-concurrency commit point;
6. publishes the staged directory with one rename, or replaces existing
   attempt files through sibling temporary-file renames;
7. finalizes an immediately closed new attempt in `goal.yaml` at the same
   resource version; and
8. writes the derived projection last and releases the lock.

The intermediate Goal for an immediately closed new attempt keeps an active
pointer until that attempt exists. An interruption before attempt publication
or Goal finalization therefore fails validation rather than silently dropping
the attempt. After the commit point, an error returns the committed Goal
reference and names its advanced resource version. Validate before resuming;
never retry with the caller's stale token.

A `checkpoint --criteria-file` update uses the same goal lock and expected
resource version. It installs the immutable criteria snapshot, replaces the
current criteria, advances `goal.yaml`, and writes the projection last.

Atomic rename prevents torn files; it does not prevent lost updates by itself.
The lock and expected resource version are both required for this
compare-and-swap path. A process that ignores the lock is outside this
cooperative trust boundary.

The lock is keyed by the canonical goal path under
`$XDG_RUNTIME_DIR/alwaldend/goal/locks/`, so it is never tracked and work on
one goal does not block an unrelated sibling goal. If a process exits between
two file renames, validate the record before resuming instead of inferring
state from temporary files or generated Markdown.

The lock file contains the current holder PID for diagnostics. Clean release
truncates the PID while still holding `flock`; a crash releases the kernel
lock automatically and leaves the PID for diagnosis. The next holder
overwrites a stale PID. `flock`, not PID existence or process liveness, is the
authority. Never unlink and recreate a lock path based on PID metadata while
other processes may be waiting on its existing inode.

The key uses the full canonical absolute goal path rather than only its ID.
Same-named goals in parallel workspaces are independent; symlink aliases of
one path share a lock.

Other commands use locks according to the paths they coordinate. Promotion and
migration acquire distinct source and destination locks in canonical-path
order. Attachment reads the goal under its goal lock, then writes the session
binding under a separate path lock.

This is cooperative same-user coordination, not a hostile-process sandbox.
The runtime lock inode is checked without following a final symlink, but the
workspace and goal pathnames must remain stable for the duration of a command.
A process with the same user identity that races workspace renames or symlink
replacement is outside the tool's trust boundary.

On a stale update, retain the isolated attempt output, reread the new canonical
state, and decide whether to rebase the evidence into a new attempt, revise the
plan, or discard publication. Never retry blindly with a newly read version.

## Workers

Workers receive immutable inputs: goal ID, resource version, lifecycle
generation, criteria revision, attempt ID, bounded task, and an isolated
scratch/output location. They do not write into the canonical goal directory.
The coordinator reviews their output and imports selected Markdown evidence
through the canonical checkpoint.

Before publication, recheck that the goal remains open, the execution state
permits publication, and all versions still match. This prevents a late worker
from publishing after pause, completion, abandonment, or supersession.

Parallelize only independent work with disjoint output locations. Bound fanout
and depth; recursive delegation needs its own cost justification. Stop workers
whose output can no longer affect the critical path.
