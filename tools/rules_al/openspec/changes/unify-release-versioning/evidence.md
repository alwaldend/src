# Acceptance evidence

## AL consolidation acceptance

The final owner is `tools/rules_al`. Direct AL commands compose the existing
configuration/plugin runtime, release/versioning, guarded delivery, and validated
Bazel runner. Lifecycle procedures live in eight canonical `al-*` skills;
upstream OpenSpec skills remain unchanged. The module has thin startup, local
command adapters and domain packages, and deliberate public APIs under `pkg`.

Before implementation, the built-CLI E2E failed because the lifecycle commands
were absent (`out/unify-release-versioning/al-cli-before.log`). The standalone
workspace then passed all 20 tests, including the CLI E2E, stamped-Bazel exit
status, immutable release fixtures, delivery guards, public bundle rules, and
offline skill configurations (`al-standalone-all-3.log`). Five parent tests passed
for native AL/plugin cleanup, CI, hooks, skill guidance, and exact discovery
(`al-parent-integration.log`). E2E outputs retain version/provider reports,
portable manifests, and literal Bazel argument observations.

The public Go API now resolves through the root module's local first-party
requirement and replacement. A pinned `go list -e -mod=readonly -json` observation
contains neither `Error` nor `DepsErrors` (`al-go-consumer-final.json`). Whole-root
`go mod tidy` remains limited by existing Bazel-generated Anki protobuf imports;
that failed attempt changed no dependency versions. Native Go inspection refreshed
workspace checksums for already declared dependencies. Standalone metadata was
generated successfully through the pinned Go and Bazel tools.

A bounded source comparison found 82 moved Go files unchanged after import and
package normalization. The remaining command, runner, receipt integration, and
fixture differences were reviewed directly (`al-domain-review.diff`). Existing
receipt schema identities and immutable publication guards are preserved.

The complete pre-rebase matrix covered 17 workspaces and 34 build/test phases
(`al-matrix-2.json`, `out/full-repo-check/run.2810578767/report.md`). Every build
and all 16 standalone workspace test phases passed. The root test phase passed
350 of 353 tests; the three failures were stale skill copies and formatting of
two BUILD files. Regeneration and configured formatting closed all three, and
all 28 quality/discovery tests passed (`al-quality-ready.log`). A second targeted
Gazelle run changed none of the 15 consumer BUILD files.

The task commit was rebased onto fetched master
`b8b60ae7ee3ace815ac7d41eb8c267b95e573eb9`. The only conflict was the ComfyUI
README: its new managed systemd launch instructions were retained and their
runner spelling migrated. Newly introduced security BUILD consumers and systemd
and SSH documentation were migrated through the same AL API. Those changes
alter command wiring and documentation, without running maintenance operations.

Preparation also reproduced `AL-CLI-ERROR-OUTPUT`: the unified root printed
Cobra help to stdout after a delivery runtime error. A built-CLI E2E case first
failed with that output (`al-error-output-before.log`). Root usage suppression
keeps diagnostics on stderr and preserves machine-readable stdout; the fixture
retains `al-failure.txt` alongside its successful command artifacts. All 20
standalone tests pass with the correction (`al-standalone-final.log`). The new
systemd workspace also passes its ordinary build and test commands; its manual
runtime probe was not executed.

Final lint caught `AL-ANDROID-REPOSITORY-IMPORT`: global `bazel mod tidy` removed
a platform repository still used explicitly by the Android `sdkmanager_files`
target. Its existing import was restored; no dependency version changed. The
first candidate's five test/build checks passed before lint stopped publication.
Restoring that import passed the complete 213-label lint selection
(`al-restored-import-lint.log`). Final candidate gates are repeated after rebase.

Current matrix and final publication state are recorded in restricted task logs
and guarded delivery receipts under `out/unify-release-versioning/`. No host
installation, live release publication, infrastructure operation, or deployment
was performed. Checked-in provisioning now installs AL instead of the retired
runner.

### Approved synchronization rebase

The user approved substantive resolution against fetched trunk
`5cfe3ae3d8299448ace8e7190d7b4b86bc5f847a`, preserving task commit
`bea2b2ab05306985a58502a3a52e50f8c74669db` and remote feature head
`f78e7ec70c21e0a51a07364e5f4340693c8cc08a`. Release publishing remains immutable
in both namespaces, identical retries remain valid, and activation retains sites
until explicit deletion. The incoming Ansible mirror remains separate and
authoritative. Incoming release records moved to the AL owner with provenance
notes identifying the superseded publication policies.

Before integrating the destination checks, the built CLI reproduced activation
through symlinked ancestry (`al-rebase-safety-before.log`). The same E2E now passes
with both namespace collisions, retained-site rollback, ancestry refusals, and
mount refusal at operands, ancestors, and descendants (`al-rebase-safety-after.log`).
The native fixture controls `mountpoint`; it does not claim actual mount coverage.
The parent loopback SSH fixture now consumes AL's native binary and generated
bundles, preserves identical concurrent retries, and tests mounted explicit
deletion instead of automatic site cleanup. All 20 standalone AL tests passed
after integration (`al-standalone-rebased.log`). The actual SSH/Ansible fixture
passed its publication, ancestry and bind-mount checks plus the authoritative
mirror cases (`al-rebase-sync-2.log`, `sync-receipt-approved.json`). Its first run
identified `AL-SYNC-SSH-WRAPPER`: applying a private-home shell wrapper to rsync's
argument protocol malformed its command. Limiting that fixture wrapper to
selection and deletion restored the intended coverage; production transport
was unchanged.

The incremental correctness review compared the preserved task commit with the
rebased release code and reviewed the parent fixture, documentation and incoming
record migration. No remaining actionable findings were identified. Destination
checks preserve literal shell arguments, inspect both deletion trees before
removal, allow the configured content-root mount, and retain per-project locking.
The complete pre-rebase implementation review and matrix remain bounded evidence;
final quality, affected lint, consumer, and selected E2E checks run on the exact
prepared candidate.

### Observed workflow corrections

- `AL-MOVE-OUTPUT-SYMLINK`: a former workspace convenience symlink caused recursive
  discovery of external source. The task-owned stale symlink was removed; the
  nested-module procedure now records that rename step.
- `AL-MOVE-WORKSPACE-BOUNDARY`: target discovery caught the old `.bazelignore`
  entry; the owner was corrected to `tools/rules_al` before the full matrix.
- `AL-ROOT-GO-RESOLUTION`: a direct module-resolution observation reproduced the
  missing first-party reference and a handwritten protobuf wrapper package
  mismatch. Both were corrected without adding external software.
- `AL-MATRIX-IDLE-SERVERS`: completed workspaces retained idle Bazel servers.
  Available host memory was observed at 2.6 GiB before cleanup and 12.4 GiB after
  shutting down 15 completed task-workspace servers, retaining their disk caches.
  Three bounded root commands timed out during pressure; serialized retries
  succeeded after cleanup. A future operating procedure can conditionally stop
  completed task-owned servers when memory pressure is observed, without touching
  active or unrelated workspaces.
- Two formatter attempts were premature: unstaged retired files remained in the
  tracked-file inventory. The final workflow stages reviewed moves before
  formatting and regenerating discovery. Raw failure logs are retained locally.

## Earlier release implementation

The task branch was rebased from `a9ac58eda` onto fetched trunk
`4e10a19e32f51c93f9c56a3f872951725a8a3f57` before implementation.

The built release CLI is exercised by `//cmd/release:go_test` against
isolated destination trees. Coverage includes failed-transfer isolation,
concurrent identical publication, conflicting well-formed releases, local and
remote corruption, independent projects, archive extraction and selection,
active deletion refusal, idempotent deletion, different-byte republishing after
deletion, unsafe namespaces, source identity failures, and mutable OCI refusal.
The SSH authentication fixture separately uses actual OpenSSH with an in-memory
agent. Lifecycle file transfer is simulated; it does not claim live rsync or
host deployment verification. Test outputs include publication manifests.

`//internal/versioning:versioning_test` exercises real Git against
an isolated bare remote: patch-zero publication, a patch diverged from trunk,
shared-ref reuse, exclusive operation locking, and a verified patch receipt.
The existing version and command tests preserve calendar and stamping behavior.

The shared bootstrap built `//projects/sri/releases:bundle` under release build
configuration. An independent read of its output verified filename, byte count,
and SHA-256. This task's feature checkout correctly produced a development
bundle; live release publication was not performed. Projects with nonempty
payload declarations provide independently selectable bundle targets.

Final candidate validation, publication, and PR identity are owned by the
repository delivery receipts under ignored `out/unify-al-versioning/`.

## Standalone module and review corrections

The branch was subsequently rebased onto trunk
`c52acd876632eaf017a2f2ccea046816df37aae7`. At that stage the release owner was the
standalone `rules_release` module. Its public rules and providers are unprefixed,
and its canonical skill is `al-versioning`. The used Git schema, parser,
and toolchain moved into this module; the separate Git project and unused
changelog wrappers were removed. Parent tools and website data use this owner.

Five focused standalone test targets passed. New E2E coverage verifies selected
destination publication through activation with a regular-file symlink payload,
release-history generation through the built CLI, and actual portable bundle
outputs from the public Bazel API. Test artifacts include
`selected-publication.json`, `release-history.json`, and
`standalone-release.json`.

The full-workspace build exposed an indirectly generated development-environment
label and missing development metadata in unstamped bundle builds. The label
now points to `@rules_al//:al`. A new public-rule E2E case reproduced
the metadata failure before implementation; after reusing the existing pinned
`bazel_lib` stamping helper, it passed and retained `development-release.json`.
Ordinary bundles use constant development status; stamped bundles still require
valid supplied source metadata. The CLI's missing-status refusal is preserved.
The broad test build also identified five empty documentation payload sets.
Four documentation-only projects now omit invalid bundle declarations, while
the website's bundle label aliases its existing deployment bundle. Empty bundles
remain rejected.

Parent checks passed for merged Git metadata, the raw Git executable, a generated
release page, the two Git-invoking tools, skill discovery, and offline skill
evaluation configuration. The migrated bootstrap built the SRI bundle; an
independent payload read matched its manifest SHA-256 and byte count. Final
candidate gates and the complete workspace matrix are recorded in task-local
delivery and full-check reports.

The full matrix passed normal build and test commands in every nested workspace,
including all six rules_release test targets. After the integration corrections,
the complete root rebuild passed all 6,664 normal targets. The exploratory root
test pass ran 370 tests: 368 passed and the two generated-state quality checks
failed. Their owning generators were run, and exact-candidate quality, lint,
consumer, and standalone checks are mandatory guarded-delivery gates. The
task-local validation report retains the full command matrix and diagnosis;
manual/incompatible targets and optional configurations remain outside that
normal matrix.

Final representative-output inspection found that the inherited Git reader
returned empty metadata for a linked worktree. A native Git fixture reproduced
this through the built `git-info` CLI before implementation. Enabling the pinned
Go Git reader's common-directory option fixed object and remote lookup. The
fixture then passed and retained `linked-git-info.json`; the real parent output
was independently checked for its current source commit and configured remotes.
Final candidate delivery gates are repeated after this correction.

## Session observations

- `managed-metadata-access` (live): the sandbox initially made Git metadata and
  the configured Bazel output base read-only. Approved operations used those
  managed locations; no duplicate cache or host reconfiguration was introduced.
- `formatter-validation-overlap` (live): one exploratory test run overlapped
  a whole-repository formatter and reported formatting drift/concurrent source
  changes. It is not final-candidate evidence. Final delivery runs formatting
  before preparation and repeats required checks on the clean committed tree.
- `bounded-read-size` (live): several early combined file reads were truncated.
  Subsequent inspection used bounded component files and diagnostic extracts.
  No shared-policy changes are proposed for these task-local execution issues.
- `indirect-consumer-discovery` (live): bounded literal-label searches missed one
  label generated from the development environment's tool-name list. The full
  build identified it; the owning BUILD now declares the external tool directly.
- `formatter-launch-context` (live): directly invoking the built formatter lacked
  the Bazel runtime workspace context. The final formatter uses its owning
  `al bazel run //:format` entry point.

## PR review follow-up

The built-CLI cases first failed against the prior implementation
(`out/unify-release-versioning/review-fixes-before.log`): version operations
ignored SIGINT/SIGTERM and identical retries still required a payload transfer.
The corrected focused version, release, and CLI suites pass. Native Cobra
commands accept shared flags before or after subcommands, preserve `--repo`,
and expose `--repository`. Cancellation reaches Git process groups and the
stamped Bazel client with a bounded interrupt wait.

Publication now verifies existing manifest identity and payload hashes under
its destination lock before requiring local payloads or staging. The final
post-transfer locked check still protects concurrent publication. Deletion
preflight checks every descendant type, including bind-mounted files.

The maintained `tests/consumer` templates were built in
`out/unify-release-versioning/external-consumer` with the fixture's independent
Bazel profile and only explicit first-party module overrides. Both
`//:bundle` and `@rules_al//:al` build without parent configuration imports.
The output manifest records `2026.44.0`, the supplied source commit, a 28-byte
payload, and verified SHA-256
`365275bcd7c278147f582e849598ec7e5ab5a33f09d02d2268eeeca0b24474c5`.
The manifest is retained in
`out/unify-release-versioning/external-consumer-release.json` and its repeatable
procedure is in `tests/consumer/README.md`.

All 20 standalone AL tests pass after the complete changes. The real
`//infra/download/test/sync:sync_test` fixture also passes, including deletion
refusal for bind-mounted `index.html` and archive payload files without changing
either retained tree. The SSH fixture completed in 179 seconds; raw evidence is
in `out/unify-release-versioning/review-fixes-ssh.log`.

### Terminal and artifact evidence follow-up

`AL-GIT-TERMINAL-GROUP` reproduced in the built CLI before implementation: a
controlling-terminal fixture stopped at authentication and reached its eight-second
deadline (`terminal-receipt-before.log`). The same fixture now completes input,
SIGINT, and SIGTERM cases and verifies that the original foreground group is
restored (`terminal-receipt-complete.log`). Noninteractive cancellation remains
covered by the existing CLI fixture. The terminal code reuses the existing pinned
x/sys dependency; no new external software was introduced.

The release CLI initially rejected --receipt in the same failing E2E run. Its
completed receipt now matches the destination manifest digest and resolved identity.
Partial multi-bundle progress and failed/cancelled archive retrieval preserve
completed upload, unknown extraction and pending activation. The fixture inspects
the running checkpoint before cancellation, private permissions, and refusal before
remote effects when receipt persistence fails. Undeclared test outputs retain
artifact-completed.json, artifact-running.json and artifact-incomplete.json.

The automated `tests/consumer:check` command built an independent workspace and
verified version 2026.44.0, the supplied clean source identity, and the 28-byte
payload with SHA-256
365275bcd7c278147f582e849598ec7e5ab5a33f09d02d2268eeeca0b24474c5
(`consumer-automated/result.json`). Its source-built AL runfile removes executable
location discovery from the ordinary invocation. It runs outside build actions.

The session ergonomics review identified `AL-CONSUMER-RUNNER-LOCATION`: a Bazel
run launcher changes directory, so using one as a native runner failed a read-only
OpenSpec invocation. The new checker selects its native AL runfile by default and
documents the override contract. A separate relative-path invocation error was
corrected using the existing native runner. These were routing errors, not source
defects. Final acceptance and publication are bound by the task-local delivery
receipts; no host installation or live release publication was performed.
