# Verification protocol

These checklists compare the source-derived description with a running AL.
They do not authorize a live publication. Use synthetic inputs and isolated
repositories or the existing behavior fixtures.

AL's source now lives in `projects/al`, with Bzlmod name `rules_al` unchanged.
The baseline and C2/R1/O1 records below predate that move. C2's retained scratch
directory moved with the project; its current location is listed below, while
raw recorded arguments and source paths retain their original provenance. A
location update is not a new verification pass or a broader portability result.

## Checklists

| File                          | Covers                                                                            |
| ----------------------------- | --------------------------------------------------------------------------------- |
| [Foundations](foundations.md) | Invocation and state boundaries                                                   |
| [Releases](releases.md)       | Bundles, source refs, artifact publication, website operations and their previews |
| [Delivery](delivery.md)       | Candidate, publication, and review                                                |
| [Commands](commands.md)       | Configuration, child execution, affected tests, diagnostics                       |

## Run a pass

1. Build AL from the source snapshot in the [scope](../README.md#scope-decisions),
   or establish that its executable inputs match that snapshot. Record commit,
   executable SHA-256, time, operating system, and whether a TTY is present.
2. Install only into ignored task scratch, or invoke the built executable with
   its declared runfiles when that is the surface being checked. Keep the two
   cases distinct. Use absolute fixture paths.
3. Read each linked claim. P1 covers shared facts and suspected defects, P2
   ordinary behavior, and P3 numeric details. Run P1 before P2, then P3.
4. Use the named isolated setup. Record `pass`, `fail`, or `blocked`, with the
   means of observation. Leave `—` for unrun rows. Preserve stdout, stderr, exit
   status, and relevant resulting files outside tracked documentation.
5. A mismatch can be a document error. Correct that before filing a product
   defect. Deduplicate real suspicions in [triage](../bug-triage.md).
6. Keep coverage `drafted` after an automated pass. Mark a document `verified`
   only after its P1/P2 rows are manually observed or failures are filed and its
   remaining limitations are explicit. Never renumber used checklist IDs.

## Setups

- **Local:** Linux, Bash, pipes unless TTY is specified; source-built binary,
  no user configuration, minimal synthetic environment, task-local output.
- **Git:** fresh scratch repository with a clean commit, trunk `master`, local
  identity and signing disabled for fixtures; optional local bare remote.
  Never run mutation rows in the feature worktree or against GitHub.
- **SSH fixture:** the existing `internal/release/lifecycle_test.go` controlled
  transport fixture or root `infra/download/test/sync` loopback SSH fixture.
  Use fixture content roots and synthetic identities only. Failure injection
  belongs inside that setup; do not disconnect or alter a real host.
- **Delivery fixture:** existing delivery integration fixtures with isolated Git
  repositories and a controlled forge. Actual candidate publication is separately
  evidenced by delivery receipts, not by these test rows.
- **Interrupt fixture:** a controllable slow child/transport or filesystem with
  a visible checkpoint. Send signals only to its recorded PID. Observe state
  before retrying. SIGHUP, repeated signals, and TTY prompts require dedicated
  terminal observation; piped output is not a substitute.

## Historical baseline observations

On 2026-11-02 a scripted pass exercised 30 invocations against the executable
whose SHA-256 is `019074063ca33f92755b608da59da4a1922a385744624ed7fdcac13738f41bf5`.
The executable's AL source matched the original post-rebase snapshot `4a10748`.
These observations predate the fixes and remain historical evidence. Fifteen
observations cover configuration/run/Bazel/doctor; fifteen cover help, version
preparation/publication to a local bare remote, bundle output, and local Git history generation. The generated
bundle's payload hash was independently checked against its manifest.
The history-export command has since been removed; that original observation is
retained only as historical evidence and has no current feature/checklist row.
The original Bazel-forwarding observations likewise describe a removed command;
ordinary builds now invoke Bazel directly.
The original weekly preparation observations included an automatic checkout.
They do not verify the current `al release create` behavior, which leaves the
checkout unchanged and requires ordinary Git to select the release branch.
The delivery observations recorded on 2026-11-02 also predate the check-plan and
results-ledger removal. Their Git/review findings retain their original scope;
they do not establish current check selection or stored validation readiness.

Evidence is task-local under
`out/unify-release-versioning/rebase-description/observations-agent/` and
`observations-release/`. The command checklist cites the stable observation IDs;
release checks cite their IDs likewise. The evidence uses only synthetic inputs.
B-01 and B-02 were reproduced; B-03 was source-derived in this original pass.

This pass did not inspect interactive rendering, terminal prompts, repeated
signals, actual credential providers, or a live release destination. Existing
E2E fixture results from delivery validation are separate evidence. No document
is marked manually verified, and unrun checklist rows remain `—`.

## Later fixture passes

The following automated passes on 2026-11-02 exercised the corrections during
candidate preparation. They supplement the 30 original observations; they do
not replace their provenance or establish a manual pass.

| Record | Scope                                                                                                                                                                                                                                                                                             | Retained evidence                                                                                                                                                                                                                                                                                                                                                                                            |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| C2     | `TestRunStandaloneCLI` (8 cases), `TestRunStreamFiles` (2), and `TestRunPluginEnvironment` (3) passed against a copied binary-only installation and synthetic child/plugin helpers.                                                                                                               | `projects/al/out/unify-release-versioning/run-fixes/after-artifacts/observations/` contains 13 JSON artifacts. Failing-before artifacts remain in the sibling `before-artifacts/` directory.                                                                                                                                                                                                                 |
| R1     | `TestRebasePreservesRemoteUntilPublication` (existing/absent remote), `TestRebaseAlreadyBasedPreservesRemote`, and `TestRebaseRequiresCleanWorktree` passed in isolated repositories and a controlled forge.                                                                                      | `out/unify-release-versioning/bug-fixes/rebase-before.log` and `rebase-after.log`; the passing target emits `rebase-local-false.json` and `rebase-local-true.json` as undeclared outputs.                                                                                                                                                                                                                    |
| O1     | The final complete release package passed, including `TestReleaseObservationCLI`, `TestSelectedDestinationCLI`, and `TestReleaseInspectionCancellationCLI`. The built CLI checked offline planning, destination observations, unchanged destination trees, and cancellation during local hashing. | Under `out/unify-release-versioning/release-observation/`, `final-release-test.log` records the pass and `final-evidence/` contains 25 numbered observation reports (2 plans, 23 inspections), `selected-plan.json`, and `inspection-cancellation.json`. Earlier reports remain in `evidence/`; the missing-command, current-link-newline, and local-hash-cancellation failure logs are retained separately. |

The checklists cite C2, R1, and O1 by exact test/subcase or artifact names.
File-mode and marker observations are limited to those retained fixture outputs.
Source-derived `next_action` wording and later-startup truncation behavior retain
unrun rows. O1 covers missing/corrupt local payloads; absent, matching, conflicting,
corrupt, and unavailable remote publications; retained, active, invalid, unsafe,
and dangling website selections; and failed, malformed, timed-out, or excessive
remote responses. The cancellation fixture records a prompt nonzero exit during
local hashing with no SSH startup; its elapsed time is an observation, not a
product timing guarantee. Cancellation before report encoding need not emit JSON. Timeout and excessive-output cases occupy reports 07/08 in an
unspecified order. Direct receipt reconciliation, third-party inspection, and
unreadable local-payload rows remain unrun.

These passes do not exercise live Vault, GitHub, release destinations, interactive
prompts, repeated signals, or every failure window. Coverage remains `drafted`
throughout.

## Configured workspace and base-selection fixtures

After the move to `projects/al`, the CLI, Bazel-command, and delivery test targets
passed against the changed source. These are bounded fixture checks of the new
doctor inputs and selected-base handling, not a verification of every AL workflow
in another repository. Earlier C1/C2/R1/O1 records retain their original meaning.

| Record | Scope                                                                                                                                                                                      | Retained evidence                                                                                                                                                                                  |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| C3     | `TestDoctorConfiguration` runs seven installed-binary cases: defaults, configured paths, absent/different source binary, missing/invalid pins, and an empty rc path.                       | Under `out/unify-release-versioning/portability/after-artifacts/cli/test.outputs/`, seven `doctor-*.json` files retain arguments, streams, and exit status.                                        |
| R2     | `TestDeliveryInspectionUsesSelectedBase` covers an explicit release base and the base from an existing PR, first with unrelated trunk advancement and then with selected-base advancement. | Under `out/unify-release-versioning/portability/after-artifacts/delivery/test.outputs/`, `selected-base-{false,true}.json` and `selected-base-advanced-{false,true}.json` retain the four reports. |

Failing-before evidence remains under the sibling `before-artifacts/` directory.
No live forge, credential provider, or release destination was used for these
checks. All document coverage remains `drafted`.

## Separate website-operation fixtures

The later **W1** release and CLI package pass exercised the nested
`al release website extract` and `activate` commands. `TestReleaseObservationCLI`
checks each offline dry run with unusable transport executables and a nonexistent
identity path, verifies unchanged destination trees, and rejects a requested
receipt without creating it. `TestLifecycleCLI` checks the retired step selector
and archive flags on commands that do not own them.

Evidence is retained under
`out/unify-release-versioning/api-naming/website-evidence/`:
`observation-01-plan.json` records upload-only planning;
`observation-02-extract.json` and `observation-03-activate.json` record the two
website previews; `al-retired-invocations.json` records unsupported inputs.
`manifest.json` identifies the commands, passing validation scope, and evidence
hashes. The separate real loopback OpenSSH/rsync sync fixture also passed;
`sync-receipt.json` retains its filesystem observations and named checks.
The pass is scripted and isolated. It does not establish a live-destination or
manual-terminal pass. O1 and the original history-export observation retain their
historical scope; they are not rewritten as observations of these new commands.

## Ordinary Bazel release-build fixtures

The **B1** ordinary-Bazel pass retains evidence under
`out/unify-release-versioning/api-naming/`. `plain-bazel-root-help.txt` shows the
six command families present before affected-test selection was added.
`plain-release-consumer/result.json` records a
real tagged Git consumer built with `bazel build --config=release :bundle`, its
version/commit identity, and independently verified payload hash. The associated
dirty build fails, and restoring the source reproduces the manifest. A separate
`missing-al-build.log` records a real Bazel workspace-status failure with
installation guidance when AL is absent. These are local scripted checks;
they do not publish artifacts or exercise credential services.

## Explicit release creation and checkout fixtures

The **C1** pass uses the installed CLI in a real local Git fixture.
`out/unify-release-versioning/api-naming/release-create-evidence/release-creation.json`
records both refs at the original trunk commit and unchanged HEAD after
`al release create`, dry-run ref preservation, and dirty/occupied-ref refusals.
The same fixture checks out the release with ordinary Git, verifies stamped
release identity, and creates its next patch tag. The fixture date is
2026-08-30 (release 2026.35), rather than the example date in the checklist.

`release-create-evidence/manifest.json` records current command help and
rejection of the retired command, switching option, and nonexistent AL checkout.
`release-create-consumer/result.json` records the separate plain-Bazel build
following creation and Git checkout, with verified payload bytes and clean
recovery after the dirty-build refusal. These are isolated scripted observations;
no remote release or website was published.

## Git-only source-publication fixtures

The **D1** source-snapshot pass exercises preparation, publication, and
verification without a check-plan or result ledger. Under
`out/unify-release-versioning/api-naming/`, `source-publication.json` records
equal prepared, published, and remote heads with verified publication.
`publication-review-recovery.json` records an initial
`publication_verification_pending` report followed by successful verification
of the same published head and PR after the read failure clears.

These are isolated Git/forge fixtures from
`TestDeliveryPreparePublishVerify` and
`TestPublicationRetainsEvidenceAfterReviewFailure`. They do not prove that an
operator ran sufficient ordinary Bazel checks or exercise a live forge.
These observations remain separate from the earlier API's dated records;
the shared source pin identifies the checked implementation.

## Affected-test selection with real Bazel

The **B2** scripted pass uses the task-local CLI and an isolated Git repository
with a real Bazel 9.3.1 graph. Under
`out/unify-release-versioning/api-naming/`, `affected-real-results.json` records
nine cases: transitive selection, manual wildcard suites, explicitly named
suites and tests, selected execution, configuration changes, deletion, unmapped
files, and unchanged source. The transitive case selects one dependent test and
excludes its unrelated neighbor. `affected-real-test-target.txt` identifies the
selected executable's retained test output. The `not-manual` tag remains eligible.
`affected-bazel-run-selection.json` separately records a successful selector
invocation through the ordinary `bazel run` entry point.

`affected-bazel-tests.json` records 23 CLI scenarios using real Git and an
instrumented Bazel boundary, including partial-query refusal, test failure,
cancellation, local change collection, and manual-suite selection.

These observations establish selection and execution for the fixture, not a
performance measurement or coverage of every repository workspace. The initial
XML-output compatibility failure is retained in `affected-real-selection.log`;
the final command uses Bazel-native filtering and label output.

## Not checkable by hand

An observation cannot prove the absence of every race or crash window. Source
inspection and adversarial fixtures are required for atomicity, exact lease
checks, unsafe path rejection, and cleanup guarantees. A passing isolated fixture
does not establish behavior under every filesystem, service, or concurrent edit.
