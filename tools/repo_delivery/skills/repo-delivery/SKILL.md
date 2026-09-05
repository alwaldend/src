---
name: repo-delivery
description: >-
  Finalize task-owned repository changes by verifying, committing, pushing the
  feature branch, maintaining its pull request, and resolving review
  comments. Use when an implementation is ready for delivery; do not use for
  read-only reviews.
---

# Deliver repository changes

Deliver authorized implementation work without asking again whether to commit
or push. Preserve unrelated work and use the repository's dedicated feature
worktree. Load `$bazel-agent` and `$repo-bazel` for invocation and validation.

## Select the entry point once

Run all delivery commands from this feature worktree's Git repository root,
which owns the root `MODULE.bazel` and `tools/repo_delivery`, including when
the change belongs to a nested Bazel module. Use that same worktree's root,
not another checkout. The baseline works with older installed runners:

```sh
bazel_agent bazel run //tools/repo_delivery -- provider
```

For the cached path, follow `$bazel-agent` capability detection once per runner
binary: if support is unknown, inspect `bazel_agent tool --help`; use
`bazel_agent tool run repo_delivery -- ...` only when supported. Retain the
result. An unsupported `tool` command uses the baseline above without a host
update. A delivery refusal or execution failure needs diagnosis, not a change
of entry point. Below, `repo_delivery <command>` means the selected invocation.

For `validate`, use the cached entry point or generate one task-local launcher
from the repository root, then invoke that launcher. This releases Bazel's
lock before validation starts another Bazel command:

```sh
bazel_agent bazel run --script_path=out/task/repo_delivery \
  //tools/repo_delivery
out/task/repo_delivery validate --receipt-file out/task/prepare.json \
  --plan-file out/task/checks.json
```

Refresh the launcher after changing the delivery tool. The owning README
defines the strict plan format: explicit workspace, `test`/`build`/`lint`,
targets, deadline, and gap decisions. No automatic check selection or arbitrary
shell/run command is authorized by the plan mechanism.

Use the sanitized `provider` report; never print the raw remote URL. GitHub
requires both `adapter_available` and `delivery_transport_available` to be
true. Both Git endpoints must use SSH. An HTTPS endpoint requires explicit SSH
configuration before adapter use; do not import mutable credential helpers or
silently change remotes. Same-repository pull requests are a caller-enforced
precondition: except for `provider`, stop before adapter use for fork or
cross-repository requests, uncertain ownership, or an unknown actual base.

For Forgejo, read [the compatibility workflow](references/forgejo.md); its
unsupported adapter result is expected. Do not create a duplicate pull request
on another forge. For supported GitHub operations, use the adapter instead of
reproducing mutations with Git or `gh` or bypassing a safety refusal.

## Prepare, validate, publish

1. Run `inspect`, supplying the intended base when no pull request exists.
   Review the exact refs, OIDs, authorship, status, pull request, refusals, and
   task diff. Decide ownership; never infer it from an author name alone.
   Before consolidation, base synchronization, remote replacement, conflict
   resolution, or recovery,
   read [rewrites and recovery](references/rewrites.md). Never rewrite shared,
   stacked, human-owned, unrelated, or ambiguous history.
2. Put the complete aggregate commit message and receipt under ignored
   `out/<task>/`. Stage fully task-owned paths with repeated `prepare --path`,
   or stage only owned hunks and use `--use-index`; never blanket-stage. The
   scope must include the entire existing feature diff. For one owned feature
   commit, pass `--rewrite <literal-inspect.local_head_oid>`; omit it for a
   new range. Keep the aggregate message accurate; the tool adds its disclaimer.
3. Run `prepare --message-file <path> --receipt-file <path>` with that scope.
   Inspect the resulting base-to-head diff, immutable path scope, receipt, and
   literal returned `head_oid`. Keep the receipt local, trusted, and unedited;
   it binds candidate, repository, endpoints, refs, remote lease, and expected
   pull-request identity. It does not prove validation or grant authority.
4. Select a validation plan and run `validate` with `--receipt-file <path>`
   and `--plan-file <path>` using the launcher or cached tool. It records candidate,
   plan, environment, checks, and logs beside the receipt; no OID copying is
   needed for continuation. It runs aggregate `git diff --check` before checks.
   Both
   the root-workspace `bazel_agent bazel test //:repo_quality_test` and semantic
   lint of affected targets are mandatory publish gates. Start lint with
   `bazel_agent bazel build --config=lint <affected_bazel_labels>` when labels
   are returned. These derive from the prepared aggregate changed paths and
   select nearest non-root packages, not dependency impact analysis: shared
   inputs may require consumer checks. Resolve every reported
   `bazel_validation_gaps` item with explicit target selection or a recorded
   reason no target check applies. Root files never imply `//:all`; nested or
   ignored workspaces need their owner's checks. Return to the repository root
   for further delivery commands. Quality alone cannot validate
   BUILD, MODULE, or configuration semantics. Verify representative output as
   well as source. Use a clean exact-candidate checkout when unrelated dirt
   could affect checks, or establish that the checks cannot consume it.
5. Review changed behavior against the requested result and touched contracts.
   A behavior-changing edit invalidates the prior correctness verdict; seek
   relevant disconfirming cases, fix valid findings, and scrutinize that diff
   again. Passing tests cover only their encoded cases. Keep this proportional
   and do not reopen unrelated accepted code for inert documentation records.
6. Run `continue --receipt-file <path> --publish` for explicit publication.
   The tool verifies recorded inputs and passing results under the receipt
   lock, then owns the exact-lease push, commit-to-PR projection, and
   verification. Without `--publish`, continuation only reports readiness or
   verifies a previous attempt. Manual `publish` remains available with
   `--receipt-file <path>` and `--validated-head <literal-head_oid>` for other
   validation workflows and diagnosed recovery. Never populate its flag from
   mutable `HEAD` or infer validation from a preparation receipt.
   `--no-pull-request` is available only when no PR already exists.

Any candidate-tree edit requires new preparation and invalidates its affected
checks and both mandatory gates. By default, validate after preparation.
For a tree-preserving amendment, reuse a passing check only with recorded prior
and new OIDs, matching tree OID, command, tools, configuration, environment, and
evidence that every relevant input is unchanged. Rerun commit-, history-, or
stamp-sensitive checks and any check with uncertain inputs. Matching trees
alone do not establish reuse. Use `prepare --message-only --rewrite <exact-oid>`
when only the aggregate message changes, then bind validation to its new OID.

If publication returns `revalidation_required`, it rebased into a new candidate
and updated the receipt before pushing. Rerun `validate`, then explicitly
`continue --publish`; prepare again only if content or message must change.
The structured path refuses stale candidate, receipt, plan, environment, logs,
and dirty worktrees. The caller still judges changes to external tools,
ignored configuration, and other unrecorded inputs before publication.
The optional message-amend `deliver` wrapper also stops at
`revalidation_required` before publication. Do not rerun it to continue.
An interrupted or failed publication is verification-only on continuation;
diagnose it before manual recovery. Never force past a refusal or blindly
retry. A new validation run cannot reset that state for the same candidate.

## Complete review and handoff

After publication, use `review inspect` to inspect feedback. If an enabled
remote review started for the exact final head, wait for an authoritative
terminal result, then reinspect its findings. Poll at modest intervals only
when the selected monitor exposes execution state for that head. Silence or an
unchanged inventory does not prove completion. If no available mechanism can
observe the outcome, make a bounded attempt and report it as unverifiable; do
not wait indefinitely or retrigger the review. Report failed or cancelled runs
explicitly. Publishing a new head invalidates an older review's completion.

When findings need a reply, resolution, or review request, read
[the review workflow](references/reviews.md) before mutation. Fix valid feedback,
revalidate and publish, then reply with the reasoned change. Explain inapplicable
feedback before resolving it. Never silently resolve, bypass receipt guards,
or post routine progress/completion comments. Do not mutate actively coedited
PR metadata or threads: provider checks are not atomic compare-and-swap.

Finish with `verify --receipt-file <path>`; it takes no `--validated-head` flag.
Confirm no requested change remains uncommitted. For substantial work or
recurring friction, use `$agent-ergonomics-review` before final preparation;
keep resulting improvements bounded by the task.

External comments must end exactly with
`LLM-disclaimer: This comment was generated by an LLM.` Issues and other
external artifacts use the corresponding type-specific final disclaimer;
documentation, code, binaries, and conversational responses do not need one.
Report the changed files, verification commands and actual results, and PR or
publication state. A linked task report may hold detailed commands. Distinguish
failed checks, unavailable provider capabilities, and unverified review outcomes
from completed work.
