# Implementation evidence

## Password transport correction

The first shared-role XCP-ng canary retained a finished 245-test report, but
its command was classified unreachable. Pinned Ansible's command module calls
`fail_json` for a nonzero scanner exit; the resulting module exit 1 is interpreted
as an sshpass error by its SSH connection plugin. Before changing this path,
require the wrapper to return the scanner status as a validated stdout record
while exiting successfully at the module boundary. Preserve warning, failure,
timeout, missing-status, and genuinely unreachable outcomes distinctly. Verify
the wrapper with a real password-authenticated appliance audit; local transport
fixtures do not establish SSH behavior.

The corrected candidate separates the module wrapper status from a strictly
validated `LYNIS_EXIT_CODE` record. Role task SHA-256
`329bf551e944dfb57c09826662a653f7164990a5558d2bfbd689bf9a06b5b062`
passed 12 fixture cases in `suite-20261029T193613Z-066oeizr`; the other six
initially hit a stale fixture assertion equating wrapper and scanner statuses.
After correcting only that assertion and stdout instrumentation, all six
passed in `suite-20261029T194119Z-fb7w9446`, preserving scanner statuses 78,
124, and 2. The candidate role hash was identical across both runs.

The authorized password-authenticated live canary then passed on XCP-ng and
XOA. Both module wrappers returned 0 and preserved actual scanner status 78.
Retained reports contained `finish=true` and respectively 245 and 257 tests;
both receipts marked completion and evidence retention true. All result files
were mode 0600. Before/after boot IDs matched, and both play recaps had zero
failed or unreachable tasks. Private artifacts are under
`out/host-security-updates-20261029/lynis-review/live-v2/` with wrapper log
`lynis-review/live-audit-v2.log`. This verifies transport and evidence handling;
it does not establish absence of compromise.

## Parent-reviewed correction failure cases (before implementation)

Recorded on 2026-10-29 against `0ce1dbb4` before changing the role:

- `LYNIS-EVIDENCE-LOSS-001`: a failed or unreachable report stat/fetch must
  never delete the remote directory containing the only evidence. Attempt
  collection of the other files, retain the private recovery path, and fail.
- `LYNIS-UNREACHABLE-RECEIPT-001`: an unreachable scanner command must retain
  its useful transport diagnosis privately, keep a controller receipt written
  before scanning, and still attempt collection. A missing command `rc` must
  not prevent receipt creation or cause a misleading completion result.
- A setup failure or transport loss before scanning must leave an incomplete
  receipt; once allocation succeeds, record the remote path before proceeding.
- Exit 0 or 78 with an absent, truncated, unfinished, or zero-test report must
  fail. Completion requires the fully retained report's `finish=true` and a
  positive integer `lynis_tests_done`, as well as an accepted scanner exit.
- A scanner error or timeout must retain available partial console, log, and
  report files before failing. Successful collection permits cleanup even
  when the scan is incomplete; failed collection does not.
- Missing remote files must not stop collection of existing files. Failed
  local receipt persistence must not permit remote deletion. A failed remote
  cleanup must leave a private recovery path and an unsuccessful role result.
- Completed exit 0 and warning exit 78 with a finished positive-test report
  must retain private evidence and remove only role-owned remote scratch.
- The shared profile must skip the DNF and YUM queries that can mutate package
  caches while retaining `audit system` with all groups and `--cronjob`.

Validation must use isolated fixtures and local transport only. No live scanner,
inventory, authentication, or appliance changes are authorized in this phase.

## Original implementation validation

Observed on 2026-10-29 from the task worktree based on `efda7f748`.

- Official stable release: Lynis 3.1.7. The publisher's detached signature
  verified against its published software key. The exact archive fetched by
  Bazel matches SHA256
  `b5314a07fd85fa3ffc7da57b508f0108ec3280d84e4af823f805d95cbbc2428c`
  and contains the expected `lynis/` source directory.
- The owning root-module generator added the dependency include. `bazel mod
deps --lockfile_mode=update` completed without a lockfile change. It reported
  the existing unrelated `com_android_commandlinetools_linux_x86_64` import
  classification warning; this task did not alter that dependency.
- Narrow Bazel build passed for the archive alias and the role's `role` and
  `role.docs` targets. The generated defaults page matches source defaults.
- Role YAML parsing and the pinned Ansible `--syntax-check` passed with a
  localhost-only fixture. No playbook tasks, scanner, authentication, or live
  host operations were executed.
- The owning OpenSpec change passed strict validation. `git diff --check`
  passed for the candidate.
- Raw receipts are under ignored
  `out/host-security-updates-20261029/lynis-build/`. The parent owns canary
  execution and `infra/security` integration.

## Correction validation and parent handoff

The correction candidate is an uncommitted diff on
`ops/host-security-updates-20261029`, based on
`0ce1dbb4fd989015a5d5e668095d9f3923d581d7`. Only the Lynis role and this change
record are owned by this phase; updater and appliance changes remain parent-owned.
The source archive and existing role inputs/output facts are preserved.

- The Bazel wrapper built `//projects/ansible_collection/roles/lynis_audit:role`,
  its `role.docs`, and `//tools/ansible:ansible_playbook` successfully. The scoped
  `--config=lint` build also passed; its configured aspects do not provide
  additional Ansible behavioral coverage.
- The packaged Ansible localhost-only `--syntax-check` passed without execution.
  Configured Prettier formatting and scoped `git diff --check` passed. Strict
  OpenSpec validation passed after the evidence-handling requirements were added;
  the later malware-coverage sentence is documentation clarification.
- Pinned source inspection confirmed `--quick` only sets quick mode, while
  `--cronjob` also enables it; `WaitForKeyPress` omits pauses in that mode.
  `include/tests_ports_packages` supplies the documented skipped checks, and
  `include/tests_malware` mainly checks tool presence and processes. No runnable
  plugins are bundled beyond the template. No actual scanner was run here.
- The parent reports existing retained real audits with 242–257 tests,
  `finish=true`, and all groups selected. These are previous observations, not
  acceptance thresholds or runs of this correction candidate. The parent also
  reports successful native XCP Python bootstrap/module checks; this phase did
  not perform or independently repeat those operations.
- Logs and isolated fixture artifacts remain private under ignored
  `out/host-security-updates-20261029/lynis-review/`. No live operation,
  authentication, commit, push, or upstream source change was performed.

The existing isolated E2E run passed all 18 cases on the unchanged candidate:
`validation/suite-20261029T191922Z-bbjaad4s/summary.json` beneath that scratch
directory records `mode=run`, `role_execution_performed=true`,
`source_unchanged=true`, and `pass=true`. Its manifest identifies all role and
fixture bytes by SHA256. The task file hash is
`6047cef7cbb2a7fb81a968acacbfd3fc2acfbf752b57d8647be643976e31696d`.

Cases cover valid exits 0/78, command unreachability without an exit code,
stat/fetch failures and unreachable results, empty/unfinished/zero-test reports,
malformed finish/count values, actual GNU timeout exit 124, and scanner exit 2.
The suite checked exact retained evidence hashes, 0700/0600 permissions, useful
private diagnostics, remote recovery paths, scratch preservation on collection
failure, cleanup after successful retention, and a neighboring sentinel.
Representative completed, unreachable, fetch-failed, and timed-out receipts were
inspected. The runnable command and fixture boundaries are in the scratch
`validation/README.md`; no additional checks followed this bounded run.

The fixture uses the existing packaged Ansible runtime, unchanged role snapshots,
local transport, and a synthetic archive. It normalizes root ownership to the
ordinary test user and redirects remote temporary storage into private scratch.
Filesystem containment is enforced; kernel network isolation was unavailable,
and the fixture uses no network. This verifies role control flow and evidence
retention, not live Lynis accuracy, SSH behavior, root ownership, or appliance
compatibility. Role source is stable for the parent's packaged XCP/XOA canary;
that live operation remains entirely parent-owned.

## Session review

`LYNIS-DISCOVERY-001` (routed): the initial combined skill/worktree read exceeded
the output budget, causing two bounded rereads. Future continuation should read
this record and the role README, then only the needed owning skill sections.

`LYNIS-LOCAL-RUNTIME-001` (configured): host Python had YAML but lacked Jinja.
Validation switched to the existing pinned Ansible syntax-check target instead
of installing dependencies. Bazel required approved access to its existing
cache; no alternate cache or credentials were introduced.

`LYNIS-CORRECTION-DISCOVERY-001` (routed): three initial combined reads exceeded
their output budgets, including an unnecessary full worktree inventory. Bounded
path reads resolved this; continuation needs this record and the fixture summary,
not another repository inventory. No shared policy changes are proposed.

`LYNIS-CORRECTION-CACHE-001` (configured): the initial Bazel invocation failed
with `Read-only file system` at its existing shared output root. Approved access
to the same cache resolved it; no alternative cache or authentication was used.

## PR simplification acceptance

The simplified role uses native temporary-directory permissions and direct
archive transfer, with one receipt template at the existing checkpoints.
The retained fixture at
`out/host-security-updates-20261029/simplification/scanners/lynis-fixture/summary.json`
passed six cases: warning exit 78, failed report fetch, unreachable report fetch,
unreachable execution, an unfinished report with exit 0, and timeout. The
recorded source hashes match the simplified role. The same local fixture
limitations described above apply; these results do not establish live scan
coverage or SSH behavior.

## Native async simplification (2026-10-30)

The scanner runs through native async polling with job data in private audit
scratch. Finished native results supply exit status; the stdout exit-marker
protocol is removed. Cleanup additionally requires known terminal status, so
lost or unfinished execution leaves its recovery directory intact.

Six current-source isolated cases passed: completion 0, warning completion 78,
timeout 124, missing finish marker, failed report fetch, and unreachable
execution. Receipts preserve native status, partial evidence and incomplete
results; cleanup occurs only after retention and known termination. The
fixture scanner is controlled and does not run another production audit.
Collection packaging and fleet/appliance syntax checks passed. Evidence:
`out/host-security-updates-20261029/native-simplification/lynis-result.json`
and `lynis-fixtures/`. Production scanning remains stopped.
