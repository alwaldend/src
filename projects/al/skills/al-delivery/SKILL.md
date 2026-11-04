---
name: al-delivery
description: >-
  Finalize task-owned repository changes by verifying, committing, pushing the
  feature branch, maintaining its pull request, and resolving review
  comments. Use when an implementation is ready for delivery; do not use for
  read-only reviews.
---

# Deliver repository changes

Deliver authorized implementation work without asking again whether to commit
or push. Preserve unrelated work and use the repository's dedicated feature
worktree. Load `$al-bazel` and `$repo-bazel` for invocation and validation.
Apply the root `AGENTS.md` standing delivery authorization to its configured
repository destination, including source disclosure and PR creation. An external
forge alone does not require another user approval. Verify the destination and
scope; an automatic approval rejection still requires diagnosis and must not be
bypassed. Describe such a rejection as an approval-system refusal rather than
claiming that the user's standing authorization was absent.

## Choose the pull-request base

Use trunk-based development by default: each new feature branch and pull
request targets the repository trunk (`master` here). Create stacked pull
requests only when the user explicitly asks for them. A request to split a
change into separate PRs, or dependencies between those changes, does not
authorize stacking or an intermediate integration base.

For a split, start each branch from freshly fetched trunk, keep each diff
scoped to its part, and link any merge prerequisites in the PR descriptions.
Publish each PR only after its required validation gates pass. If an
unmerged prerequisite prevents a dependent candidate from passing, publish
the ready prerequisites first and defer the dependent PR. After the
prerequisite merges, synchronize that branch with trunk, validate it, and
then publish. A blocked check grants neither a gate exception nor authority
to merge a prerequisite or create a stack. Inspect an existing PR's actual base before changing
history; do not silently retarget it or rewrite a shared stack.

## Select the entry point once

Run all delivery commands from this feature worktree's Git repository root,
which owns the root `MODULE.bazel` and imports `rules_al`, including when
the change belongs to a nested Bazel module. Use that same worktree's root,
not another checkout. The baseline works with a bare Bazel entry point:

```sh
bazel run @rules_al//:al -- delivery provider
```

For an explicitly installed command, use `al delivery ...` directly. A
delivery refusal or execution failure needs diagnosis, not a change of entry
point. Below, `al delivery <command>` means the selected invocation.

For repeated delivery operations, an installed standalone command or a
task-local launcher is convenient:

```sh
bazel run --script_path=out/task/al \
  @rules_al//:al
out/task/al delivery inspect
```

Refresh the launcher after changing the delivery tool. Preparation performs
Git operations only. Run builds, tests, and lint with ordinary Bazel commands
from the applicable workspace; Bazel owns check configuration and caching.

Use the sanitized `provider` report; never print the raw remote URL. GitHub
requires both `adapter_available` and `delivery_transport_available` to be
true. SSH is the default Git transport. When SSH is unavailable, or the captured
GitHub remote uses HTTPS, select `--git-transport gh` explicitly on `provider`,
`inspect`, `prepare`, `publish`, `verify`, and other commands in that workflow.
This mode uses HTTPS with the selected `--forge-cli` executable (default `gh`)
as the only host-scoped credential helper. Authenticate with `gh auth login`
when needed; do not print tokens or run `gh auth setup-git`. AL does not import
global Git configuration, change remotes, or automatically fall back from SSH.
Receipts bind transport selection: use the same mode when resuming publication
or verification; legacy receipts imply SSH. For example:

```sh
al delivery provider --git-transport gh
al delivery inspect --git-transport gh --base master
al delivery prepare --git-transport gh --base master \
  --path path/to/owned/source --message-file out/task/commit.md \
  --receipt-file out/task/prepare.json
# Run the required checks against the prepared candidate.
al delivery publish --git-transport gh --base master \
  --receipt-file out/task/prepare.json
al delivery verify --git-transport gh --base master \
  --receipt-file out/task/prepare.json
```

Same-repository pull requests are a caller-enforced
precondition: except for `provider`, stop before adapter use for fork or
cross-repository requests, uncertain ownership, or an unknown actual base.

For Forgejo, read [the compatibility workflow](references/forgejo.md); its
unsupported adapter result is expected. Do not create a duplicate pull request
on another forge. For supported GitHub operations, use the adapter instead of
reproducing mutations with Git or `gh` or bypassing a safety refusal.

## Prepare, check, publish

1. Run `inspect`, supplying the intended base when no pull request exists.
   Review the exact refs, OIDs, authorship, status, pull request, refusals, and
   task diff. Decide ownership; never infer it from an author name alone.
   Before consolidation, base synchronization, remote replacement, conflict
   resolution, or recovery, read
   [rewrites and recovery](references/rewrites.md). Never rewrite shared,
   stacked, human-owned, unrelated, or ambiguous history.
2. Run the configured formatters and inspect their changes before preparing
   the candidate. Put the aggregate commit message and receipt under ignored
   `out/<task>/`. Stage fully task-owned paths with repeated `prepare --path`,
   or stage only owned hunks and use `--use-index`; never blanket-stage. The
   scope must include the entire existing feature diff. For one owned feature
   commit, pass `--rewrite <literal-inspect.local_head_oid>`; omit it for a
   new range. Read [commit conventions](references/commits.md) before
   writing the aggregate message; the tool adds its disclaimer.
3. Run `prepare --message-file <path> --receipt-file <path>` with that scope.
   Inspect the resulting base-to-head diff, immutable path scope, receipt, and
   returned `head_oid`. Keep the receipt local, trusted, and unedited; it
   binds candidate, repository, endpoints, refs, remote lease, and expected
   pull-request identity. Preparation neither runs checks nor publishes, and
   its receipt does not prove passing checks or grant authority.
4. Check the prepared aggregate with `git diff --check`, then run the
   applicable ordinary Bazel commands. Both root-workspace
   `bazel test //:repo_quality_test` and semantic lint of affected targets
   are mandatory publish gates. Select targets from the actual diff,
   component owners, BUILD/MODULE structure, and affected consumers. Use
   `bazel build --config=lint <targets>` for semantic lint. Root file changes
   do not imply `//:all`; nested workspaces need their owner's checks.
   Quality alone does not validate BUILD, MODULE, or configuration semantics.
   Return to the same feature worktree's root for delivery commands.
   Verify representative output as well as source. Use a clean checkout of
   the prepared candidate when unrelated dirt could affect checks, or
   establish that the checks cannot consume it. Preserve useful diagnostics
   under ignored task scratch and diagnose every failed required check.
5. Before every PR creation or update, explicitly select `$review-agent` and
   complete a local code review of the entire merge-base-to-prepared-candidate
   diff, including documentation and configuration. Record the comparison base,
   candidate head, verdict, and finding dispositions in the task evidence.
   Evaluate every finding against the request and touched contracts, fix
   in-scope defects, and review the corrected candidate before publication.
   A changed candidate diff invalidates the prior local review. For changed
   behavior, seek relevant disconfirming cases; passing tests cover only their
   encoded cases. Keep scrutiny proportional and do not reopen unrelated
   accepted code for inert documentation records. Tests and remote review do
   not substitute for this local review.
6. After the required checks and local review pass, run
   `publish --receipt-file <path>`.
   The tool compares current HEAD and tree with the prepared candidate,
   rechecks clean source and remote expectations, then owns the exact-lease
   push, commit-to-PR projection, and verification. `--no-pull-request` is
   available only when no PR already exists. The caller owns the decision
   that checks passed; AL does not store or enforce a separate test-results
   ledger. A failed required check stops this workflow even when Git state
   is clean.

Any candidate edit requires appropriate preparation and the required Bazel
checks before publication. Use `prepare --message-only --rewrite <exact-oid>`
when only the aggregate message changes. Invoke the ordinary check commands
and let Bazel reuse results for unchanged declared inputs; a new Git OID alone
is not a reason to disable its cache. Preserve normal stamp-sensitive inputs
and inspect versioned outputs when the commit identity matters. Do not replace
Bazel's cache decisions with a parallel AL check plan or results ledger.

If publication returns `checks_required`, it rebased into a new candidate and
updated the preparation receipt without pushing. Run the required ordinary
Bazel checks on that candidate, and repeat local review if the reviewed diff
changed, then call `publish` with the updated receipt.
Prepare again only if content or message must change. A changed candidate,
receipt, dirty worktree, remote lease, or PR state still requires diagnosis;
never force past a refusal.

After an interrupted or failed publication, use `verify --receipt-file <path>`
and read-only inspection to establish the actual remote state before recovery.
Never blindly repeat the push. When a report contains
`publication_verification_pending`, retain and link its known PR URL and
candidate; `verified: false` means final gates remain incomplete. Use `verify`
for another verification attempt. A concurrent review-read failure alone does
not require another push, PR creation, or repeated builds when their inputs
are unchanged.

## Complete review and handoff

After every PR creation or update, use `review inspect` to inspect remote
feedback; local review does not replace this step. If an enabled
remote review started for the exact final head, wait for an authoritative
terminal result, then reinspect its findings. Poll at modest intervals only
when the selected monitor exposes execution state for that head. Silence or an
unchanged inventory does not prove completion. If no available mechanism can
observe the outcome, make a bounded attempt and report it as unverifiable; do
not wait indefinitely or retrigger the review. Report failed or cancelled runs
explicitly. Publishing a new head invalidates an older review's completion.

When findings need a reply, resolution, or review request, read
[the review workflow](references/reviews.md) before mutation. Prefer `review address`
for thread feedback on a checked, published candidate; it derives guards and
retains reply state for safe resolution recovery. Fix valid in-scope feedback,
prepare and revalidate, complete local review of the corrected candidate, and
publish, then reply with the reasoned change. Inspect remote feedback again for
the new head. Explain inapplicable
feedback before resolving it. Never silently resolve, bypass receipt guards,
or post routine progress/completion comments. Do not mutate actively coedited
PR metadata or threads: provider checks are not atomic compare-and-swap.

Finish with `verify --receipt-file <path>`.
Confirm no requested change remains uncommitted. For substantial work or
recurring friction, use `$agent-ergonomics-review` before final preparation;
keep resulting improvements bounded by the task.

External comments must end exactly with
`LLM-disclaimer: This comment was generated by an LLM.` Issues and other
external artifacts use the corresponding type-specific final disclaimer;
documentation, code, binaries, and conversational responses do not need one.
Report the changed files, verification commands and actual results, and PR or
publication state, including local review and remote feedback outcomes
separately. A linked task report may hold detailed commands. Distinguish
failed checks, unavailable provider capabilities, and unverified review outcomes
from completed work.
