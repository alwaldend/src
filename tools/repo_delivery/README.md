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

Run delivery commands from this Git worktree's repository root, which owns the
root `MODULE.bazel` and `tools/repo_delivery`. This also applies when changed
files belong to a nested Bazel module: its workspace does not own the delivery
target. Use the same feature worktree, not another checkout. The baseline works
with older installed runners. Detect the provider without printing a
credential-bearing remote URL:

```sh
bazel_agent bazel run //tools/repo_delivery -- provider
```

The optional cached form is `bazel_agent tool run repo_delivery -- ...`.
Follow the `bazel-agent` skill's one-time capability check before using it:
reuse the known result for this runner binary, or inspect
`bazel_agent tool --help` once. Unsupported `tool` means use the baseline
without updating the host; an execution failure or delivery refusal still
requires diagnosis. All examples below use the baseline.

The sanitized report distinguishes forge support from Git transport support.
`adapter_available` reports whether the forge has an adapter;
`git_transport` reports `ssh`, `https`, or `mixed_or_unsupported`; and
`delivery_transport_available` is true only when both captured endpoints are
canonical SSH endpoints. It never reports either endpoint.

For validation that invokes Bazel, use the cached entry point or generate a
task-local launcher once from the repository root. Ordinary `bazel run` can
retain the Bazel lock while its target runs; `--script_path` releases it before
the generated launcher is executed. Refresh the launcher after changing the
delivery tool. This requires no installed runner update:

```sh
bazel_agent bazel run --script_path=out/task/repo_delivery \
  //tools/repo_delivery
```

For the supported GitHub adapter, the normal workflow is:

```sh
bazel_agent bazel run //tools/repo_delivery -- inspect --base master
bazel_agent bazel run //tools/repo_delivery -- prepare \
  --base master \
  --message-file out/task/commit.md \
  --receipt-file out/task/prepare.json \
  --path path/to/task/file \
  --rewrite <inspect.local_head_oid> # omit when the range has no commit
# For an explicitly reviewed task-owned multi-commit range, use
# --consolidate <inspect.local_head_oid> instead of --rewrite.
out/task/repo_delivery validate \
  --receipt-file out/task/prepare.json --plan-file out/task/checks.json
out/task/repo_delivery continue \
  --receipt-file out/task/prepare.json --publish
out/task/repo_delivery continue \
  --receipt-file out/task/prepare.json
```

`checks.json` is an explicitly selected validation plan. Keep it mode 0600 in
the same ignored `out/<task>/` directory as the preparation receipt. For
example, a change confined to the delivery Go package can use:

```json
{
  "schema": "repo_delivery/validation_plan/v1",
  "checks": [
    {
      "workspace": ".",
      "kind": "test",
      "targets": [
        "//tools/repo_delivery/cmd/repo_delivery:go_test",
        "//:repo_quality_test"
      ],
      "timeout_seconds": 3600
    },
    {
      "workspace": ".",
      "kind": "lint",
      "targets": ["//tools/repo_delivery/cmd/repo_delivery:all"],
      "timeout_seconds": 3600
    }
  ],
  "gap_decisions": []
}
```

Allowed kinds are `test`, `build`, and `lint` (build with `--config=lint`).
The plan accepts 1–32 sequential checks, 1–128 explicit local labels per check,
and a 1–3600 second deadline per check. Package `:all` is supported outside the
root package; recursive target patterns, arbitrary `run` targets, shell
commands, and free-form flags are not. Workspace paths are relative to the Git
worktree root and must name real Bazel workspace roots. Nested checks execute
there; the required `//:repo_quality_test` check executes in `.`. The plan must
lint every suggested non-root affected package and resolve each validation gap
with exactly one `{ "path": "...", "reason": "..." }` decision. The caller
still selects sufficient consumer checks and verifies representative output.

`validate` runs aggregate `git diff --check` first, then the selected checks.
It requires a fully clean worktree and ordinary index flags before and after
checks. It records exact head/tree, preparation revision, plan bytes, inherited
environment digest, check outcomes, and log digests beside the receipt, using
mode 0600 files. Each output stream is capped at 96 KiB; truncation fails the
check. The same receipt lock protects validation and continuation. Keep all
these files trusted and unedited; they are consistency evidence, not security
tokens. Tool binaries, ignored configuration, external services, and other
inputs outside the recorded Git tree/environment still require the caller's
input-stability judgment before publication. No passing checks are reused
automatically by a new validation run.

`continue` reports readiness without publishing. Only `continue --publish`
uses the captured passing result to call the existing guarded publication and
verification path; it never substitutes the current mutable `HEAD` as evidence.
Changed candidate, receipt, plan, environment, incomplete results, altered
logs, and dirty inputs refuse publication. A base rebase records
`revalidation_required` and stops before pushing: rerun `validate` against the
updated receipt, then explicitly continue. An interrupted or failed publication
remains `publication_attempted`; further continuation only verifies its remote
postcondition and never blindly repeats a push or metadata mutation. Diagnose
an incomplete result before using the existing manual recovery API. A new
validation run cannot reset an uncertain attempt for the same candidate.

`prepare` derives `affected_bazel_labels` from the prepared aggregate changed
paths, selecting each nearest non-root Bazel package. It never infers `//:all`
from a root-owned file. The reported
`bazel_selection_basis` is `nearest_non_root_package_without_dependency_analysis`:
this bounds discovery but does not establish downstream impact. Shared build
inputs may need explicit consumer checks.

`bazel_validation_gaps` lists paths requiring a separate decision, with reasons
`root_package_requires_explicit_targets`, `ignored_by_root_workspace`,
`nested_workspace`, or `no_bazel_package`. For each gap, select appropriate
targets or record why no target check applies. Validate nested or ignored
workspaces through their owner, then return to the repository root for delivery.
The suggested commands always include the root-workspace gate
`bazel_agent bazel test //:repo_quality_test`; that gate does not establish
semantic correctness for BUILD, MODULE, or configuration changes. Run semantic
lint for the selected affected targets in addition to that mandatory gate.

To synchronize a prepared, task-owned, single-commit feature branch with an
advanced base before `prepare`, use the guarded adapter workflow:

```sh
bazel_agent bazel run //tools/repo_delivery -- rebase --base master
```

The command refuses dirty trees, divergent remote feature tips, multi-commit
or merge-containing ranges, and pull-request metadata changes. It fetches
fresh refs, replays the commit in an isolated worktree, preserves signature
requirements, pushes only with the captured remote feature lease, verifies
base advancement, and reports the literal resulting head. Run required
validations against that head before `prepare` and `publish`.

Remote rewrites use `rewrite-authorize` to write a typed, non-authorizing
authorization receipt (old remote OID, new head OID, owner root, task paths,
provider ownership), then `prepare --rewrite-authorization <file>` instead of
raw `--replace-remote` handoff. Review replies accept `--goal-ref`,
`--delivery-ref`, and `--defect-id` durable join references.

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

The manual `publish --receipt-file <path> --validated-head <literal-head_oid>`
API remains available for validations outside the structured plan and diagnosed
recovery. Never populate that flag by resolving current `HEAD`: carry the
literal candidate that the checks covered. A preparation receipt alone does
not establish validation. The structured continuation path supplies that
literal value from its recorded passing results instead.

For message-only amendments, `deliver --message-file <path> --receipt-file
<path> --owner-root <root> --task-path <path>` combines inspection, any needed
rewrite authorization, and preparation. It stops before publication with a
nonzero exit and a `revalidation_required` report containing the exact head,
tree, receipt, affected labels, and suggested checks. Establish validation for
that candidate using `validate`, then use `continue --publish`, or use the
manual `publish` API and `verify`. Do not repeat `deliver` to continue, because
that prepares another candidate.

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
bazel_agent bazel run //tools/repo_delivery -- review inspect
```

Carry values from the latest inspection literally. All mutations require the
pull-request node ID, expected head, and pull-request expectation digest. A
top-level reply also requires the
reported last-comment sentinel and top-level inventory digest:

```sh
bazel_agent bazel run //tools/repo_delivery -- review comment \
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
bazel_agent bazel run //tools/repo_delivery -- review reply \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --body-file out/task/reply.md \
  --reply-receipt-file out/task/reply.json

bazel_agent bazel run //tools/repo_delivery -- review resolve \
  --pull-request-id <pull_request.id> \
  --expected-head <pull_request.head_ref_oid> \
  --expected-pull-request-digest <pull_request_expectation_digest> \
  --thread-id <review_threads[index].id> \
  --expected-last-comment-id <review_threads[index].comments[-1].id> \
  --expected-thread-digest <review_threads[index].expectation_digest> \
  --reply-receipt-file out/task/reply.json

bazel_agent bazel run //tools/repo_delivery -- review request \
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
before contacting the provider. If a provider read fails before any resolution
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
