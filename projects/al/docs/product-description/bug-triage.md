# AL defect triage

These entries retain the original findings and record their corrections. Current
behavior belongs to the feature documents; historical reproductions below must
not be read as current expectations. No external issues have been filed. Source
paths below are relative to `projects/al`. The [verification protocol](verification/README.md)
separates baseline observations from later fixture passes.

| ID   | Severity | Behavior                                                        | Decision | Status                               |
| ---- | -------- | --------------------------------------------------------------- | -------- | ------------------------------------ |
| B-01 | medium   | Installed `run` cannot start without Bazel runfiles             | fixed    | C2 installed-command fixture passes  |
| B-02 | medium   | Output overrides preserve stale trailing bytes                  | fixed    | C2 stream-file fixture passes        |
| B-03 | medium   | Conflicting plugin environment exports have no defined winner   | fixed    | C2 plugin-environment fixture passes |
| B-04 | medium   | Direct delivery rebase pushes before validation                 | fixed    | R1 local-rebase fixtures pass        |
| B-05 | medium   | Selected release base is refused after unrelated trunk advances | fixed    | R2 selected-base fixtures pass       |

## B-01: Installed run requires runfiles

- **Raised by:** [Run](commands/run.md#open-questions-and-verification).
- **Where:** A user installs the standalone executable and invokes
  `al run -- /usr/bin/true`, without configuration or plugins.
- **Original behavior:** AL exited 1 with a runfiles discovery error before starting the
  child. The same invocation can run with an explicitly supplied empty runfiles
  directory, demonstrating that no actual plugin runtime files were needed here.
- **Expected:** A standalone command with no selected plugin should execute, or
  the installation should supply/document its required runtime layout.
- **Baseline reproduction:** 1. Use the original task-local installed binary. 2. Unset `RUNFILES_DIR`
  and `RUNFILES_MANIFEST_FILE`. 3. Run the command above. 4. Observe exit 1 and
  `runfiles: no runfiles found`. Do not confuse a Bazel-run invocation with this case.
- **Original cause:** `internal/cli/run.go:44` always prepares runfiles environment;
  `pkg/al/command.go:40` requires discovery. The plugin manager also constructs a
  runfiles resolver at `pkg/al_plugin/manager.go:37`. Installation copies only AL.
- **Resolution:** Missing optional runfiles no longer prevent standalone child
  execution. Absolute plugin paths work without runfiles; selected relative
  identifiers still require them. Invalid declared runfiles remain errors.
- **History:** C1 `run-installed-without-runfiles` reproduced the failure on
  2026-11-02 against the baseline identified in the verification protocol.
- **Status:** Fixed; C2 `TestRunStandaloneCLI` passes all eight isolated cases.
  No manual terminal pass or live credential use is claimed.

## B-02: Output redirection retains old bytes

- **Raised by:** [Run](commands/run.md#open-questions-and-verification).
- **Where:** `al run --stdout OUTPUT -- COMMAND`, or the corresponding stderr flag.
- **Original behavior:** Existing output was opened at offset zero without truncation.
  A child writing `new` into `previous-content` leaves `newvious-content`.
  Missing output files are refused, and existing permissions are retained.
- **Expected:** Shell-like replacement would create/truncate output; intentional
  in-place writes need a clear contract because trailing bytes look like child output.
- **Baseline reproduction:** 1. In isolated scratch create `OUTPUT` containing
  `previous-content`. 2. Supply the empty runfiles fixture to bypass B-01. 3. Run `al run --stdout OUTPUT -- /usr/bin/printf new`. 4. Read the whole file.
- **Original cause:** `internal/cli/run.go:108` and `:117` used write-only opening without
  create/truncate flags. This differs from config dump's private truncated output.
- **Resolution:** Both output overrides create missing files or truncate existing
  ones. New files use mode 0600 subject to the process's permission mask; existing
  permissions remain unchanged. Later startup failure does not restore old bytes.
- **History:** C1 `run-existing-output` and `run-missing-output` reproduced the
  failure on 2026-11-02; the existing file retained mode 0644.
- **Status:** Fixed; C2 `TestRunStreamFiles` passes the missing/existing output
  cases for stdout and stderr. Retained artifacts show new files at 0600 and
  existing files at 0644. The input file is unchanged. No real secrets were used.

## B-03: Conflicting plugin environment exports

- **Raised by:** [Run](commands/run.md#open-questions-and-verification).
- **Where:** Two selected plugins export the same environment variable.
- **Original behavior:** The child environment was assembled from map iteration, so
  conflicting exports have no documented stable precedence.
- **Expected:** Reject duplicates or define a deterministic precedence that the
  operator can inspect.
- **Baseline reproduction:** 1. Prepare two synthetic plugins exporting different values for
  the same key. 2. Select both. 3. Invoke a child that prints only that synthetic
  key repeatedly. 4. Record values and order. Repeated agreement alone cannot
  prove ordering is guaranteed.
- **Original cause:** `pkg/al_plugin/manager.go:64` traversed the configured plugin map
  while collecting environment exports.
- **Resolution:** Equal selected exports are accepted once. Different values for
  one key fail before child launch, with a deterministic diagnostic naming the
  key and plugin names but omitting values. Started plugins proceed through
  cleanup. Selected exports can still override inherited environment values.
- **History:** The original description recorded a source-derived suspicion;
  its initial 30-invocation pass did not exercise conflicting plugins. The later
  failing-before fixture artifacts are retained separately from C1.
- **Status:** Fixed; C2 `TestRunPluginEnvironment` passes equal-value and both
  conflict-order cases, including child refusal, omitted values, and plugin
  shutdown markers. This is synthetic fixture evidence, not a live-provider pass.

## B-04: Direct delivery rebase pushes before validation

- **Raised by:** [Delivery](delivery/workflow.md#open-questions-and-verification).
- **Where:** An operator runs `al delivery rebase` on an eligible branch after its base advances.
- **Original behavior:** The command rebased and pushed with an exact ref lease
  before the operator could check the changed candidate. The then-current
  receipt-bound publication workflow instead stopped for revalidation.
- **Expected:** Replaying source locally must leave publication to the separate
  preparation, validation, and explicit publication workflow.
- **Baseline reproduction:** 1. Use the former `TestRebaseReplaysAdvancedBaseAndLeasePushes` in the isolated
  delivery fixture. 2. Advance the base. 3. Run direct rebase. 4. Observe the remote
  feature head and `pushed: true` without a validation invocation.
- **Original cause:** `internal/delivery/rebase.go:140` called the guarded push directly.
- **Resolution:** Direct rebase changes the local candidate only. Its report
  retains the observed remote head and lease, says `pushed: false`, and directs
  preparation, ordinary checks, and explicit publication after a replay.
  Receipt-bound publication now reports `checks_required` after a rebase;
  the operator runs checks and then publishes with the updated preparation receipt.
- **History:** Observed during the task's original authorized rebase;
  retained report `out/unify-release-versioning/rebase-description/rebase.json`.
  Final delivery validation is separate evidence, not a claim that rebase ran it.
- **Status:** Fixed; R1 `TestRebasePreservesRemoteUntilPublication` checks existing
  and absent remote refs. `TestRebaseAlreadyBasedPreservesRemote` checks the no-op
  case, and `TestRebaseRequiresCleanWorktree` checks preserved dirty inputs.
  The `next_action` wording is source-derived and remains an unrun checklist row.

## B-05: Selected base compared with unrelated trunk

- **Raised by:** The isolated selected-base fixture during the project move.
- **Original behavior:** Delivery inspection compared the selected release base
  with local `origin/master`. Advancing unrelated trunk marked that selected base
  stale and refused both an explicitly selected base and one supplied by a PR.
- **Expected:** Ancestry and rebase decisions use the selected base's fresh ref.
- **Original cause:** The additional `origin/master` comparison in
  `internal/delivery/delivery.go` applied independently of the selected base.
- **History:** R2 failing-before `selected-base-false.json` and
  `selected-base-true.json` report `base_stale: true` and a refusal for the
  unchanged release base.
- **Resolution and status:** Fixed. `TestDeliveryInspectionUsesSelectedBase`
  passes both cases: unrelated advancement leaves the selected base usable,
  while selected-base advancement sets `needs_rebase`. Local HEAD is preserved.
  The [R2 protocol record](verification/README.md#configured-workspace-and-base-selection-fixtures)
  identifies the retained before/after evidence; no manual or live-forge pass is claimed.

## Verification boundary

The observations use synthetic values and task-private files. No live release,
credential, or infrastructure configuration was changed. Other open questions
in feature documents are coverage gaps unless evidence establishes a defect.
