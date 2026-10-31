# Native-client validation

Pre-adaptation failure cases are retained at ignored
`out/host-security-updates-20261029/clamav/role/native-failure-cases.txt`.
Superseded Python evidence remains under `clamav/client` and
`clamav/role/previous-python`; it does not validate the native implementation.

## Prospective limit and transport failure cases

Recorded on 2026-10-30 before this documentation revision, using the
primary-approved policy and source findings. These are prospective failure
cases; the parent owns implementation and native acceptance.

- A server stream quota of 104857600 rejects the client's extra byte as a
  protocol error and aborts directory traversal before later files. The parent
  reports an actual HostBot scan became incomplete after about nine seconds.
- Equal client/server stream limits of 104857601 without engine
  `MaxFileSize 104857600`, `HeuristicAlerts yes`, and `AlertExceedsMax yes`
  could accept a truncated prefix. All four controller settings must agree
  with the unchanged client `max_file_bytes + 1` expression.
- An oversized fixture alone cannot prove continued traversal. Acceptance must
  observe `Heuristics.Limits.Exceeded.MaxFileSize FOUND` and a later EICAR
  detection in the same native traversal, with no prefix `OK` verdict.
- Treating the native infected count as a malware-only count misclassifies
  limits and encrypted-content heuristics; successful execution still cannot
  establish host cleanliness or complete coverage.
- Any endpoint command before the controller-only dedicated ControlPath and
  persistent master are established can use a different connection. This
  includes boot identity, tempfile creation, and Ansible module bootstrap.
- Recapturing caller arguments after overriding them, resetting connections,
  or directing mux operations at an implicit master can affect caller-owned
  connections or prevent restoration. Collection on another master can need
  fresh authentication after a long scan.
- Dropping failed or unreachable boot-before/boot-after command results loses
  operational evidence. Retain both command results privately; failed boot
  identity checks cannot establish unchanged boot identity or completion.

## Acceptance cases

These are required boundaries, not a claim that every case has been exercised.
The actual local fixture results and remaining remote checks are recorded below.

| Boundary                                  | Required result                                                           |
| ----------------------------------------- | ------------------------------------------------------------------------- |
| Loader or architecture mismatch           | Explicit unsupported-runtime failure; no fallback or package installation |
| Missing master or rejected forward        | Incomplete receipt; no mux-command authentication fallback                |
| Clean and EICAR fixtures                  | Native logs and summary distinguish detections                            |
| Oversize followed by later EICAR          | Engine limit FOUND, then EICAR FOUND; no prefix OK or traversal abort     |
| Detection plus stream/read error          | Exit 1 does not mask summary errors                                       |
| Deadline or missing wrapper status        | Incomplete workflow; available logs retained                              |
| Certificate expires during intact session | Collection through the existing master before exit                        |
| Lost master or failed transfer            | Remote evidence and recovery path preserved                               |
| First endpoint command                    | Controller-only ControlPath and persistent master already established     |
| Unrelated SSH master exists               | No reset, stop, or exit directed at it; caller arguments restored         |
| Failed or unreachable boot command        | Complete command result retained privately; completion remains false      |
| Symlinks, pseudo mounts, exclusions       | Native exclusions and no following; /dev/shm remains eligible             |

## Source evidence

Official `clamav-1.5.4` accepts unsuffixed integer byte limits in
[optparser.c](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/common/optparser.c#L1120-L1186).
The [client loop](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/common/clamdcom.c#L225-L288)
caps payload at its configured limit; the
[server quota check](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/clamd/server-th.c#L786-L821)
rejects excess chunks before the zero terminator dispatches scanning. The
former client 104857601/server 104857600 policy prevented prefix success by
raising a protocol error, but that error aborted directory traversal at the
first oversized file. The historical fixture below demonstrated that rejection;
it does not satisfy the revised continued-traversal acceptance gate.

The primary's official-source research establishes the revised native policy:
controller `StreamMaxLength 104857601` equals the client stream limit from the
unchanged `max_file_bytes + 1` expression, with `max_file_bytes` still 104857600.
Controller `MaxFileSize 104857600`, `HeuristicAlerts yes`, and
`AlertExceedsMax yes` let the extra byte reach the engine and report
`Heuristics.Limits.Exceeded.MaxFileSize FOUND`. Native traversal continues
without a successful truncated-prefix verdict. The revised local native
acceptance below supports this engine behavior; SSH acceptance remains pending.

[Native ping](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/clamdscan/client.c#L176-L222)
prints PONG after connect/send without reading the response. The
[native summary](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/clamdscan/clamdscan.c#L164-L202)
reports infected files and conditionally total errors, not total scanned files.
The native infected count includes limits and encrypted-content heuristics;
it does not mean every finding is malware.

## Historical recorded checks

- Final role and all included task files passed offline syntax checking;
  retained log: `clamav/role/native-final-syntax.log`.
- Final package/default documentation build and strict OpenSpec validation passed;
  logs: `clamav/role/native-build.log` and
  `clamav/role/native-openspec-validation.log`. Exact role hashes are retained in
  `clamav/role/native-source.sha256`.
- Earlier SSH-option and controller-cwd fixtures remain in
  `clamav/role/research/results.json`; they do not validate the corrected
  early-master lifecycle.

## Superseded native E2E

These results used the former server stream quota of 104857600. They remain
historical evidence, superseded for acceptance of the revised native policy.

The retained index is
`out/host-security-updates-20261029/clamav/e2e/result.json`. It reports
`passed=true` and points to the actual role receipts and controller artifacts.
The receipts and oversized fixture size were read directly for this follow-up:

- Clean: native exit 0, infected files 0, reported errors 0, `complete=true`.
- EICAR: native exit 1, infected files 1, reported errors 0, `complete=true`.
- Oversize, 104857602 bytes: native exit 2, infected files 0, reported errors 1,
  `complete=false`. No successful truncated-prefix verdict was returned.

All three record `execution_finished=true`, `evidence_retained=true`, native
logs retained, unchanged boot identities, and verified scratch removal.
Every receipt keeps `coverage_complete=false` and `clean=null`. The oversized
case therefore records a terminal native process without claiming successful
workflow completion or full coverage.

The parent reports the same old quota aborted the actual HostBot directory
scan at its first large file, leaving the run incomplete after about nine
seconds. The single oversized fixture above did not test later-file traversal.

The controller's native lifecycle receipt records `state=stopped` and
`cleanup_complete=true`. Its engine is ClamAV 1.5.4 with daily 28168, main 63,
and bytecode 339. Controller evidence and the preserved initial PID0644 failure
are described by the owning
[controller validation record](../../../../../infra/security/openspec/changes/add-clamav-scans/validation.md).
The redundant PID file was removed; systemd owns daemon identity and cleanup.

## Revised native acceptance: passed locally

The parent retained the revised execution in
`out/host-security-updates-20261029/clamav/e2e/continuation-result.json`, which
references the actual role receipt and controller lifecycle. Those receipts and
the console hash/order were read directly for this update. The console records
the oversized `000` fixture as `Heuristics.Limits.Exceeded.MaxFileSize FOUND`,
then the later `999` EICAR fixture as `Eicar-Test-Signature FOUND` in the same
traversal. This verifies continuation under the revised engine-limit policy.

The receipt records exit 1, infected count 2, reported errors 0,
`execution_finished=true`, `complete=true`, and `evidence_retained=true`.
The scan scratch state is `absent`; the controller is `stopped` with
`cleanup_complete=true`. The retained console matches its indexed SHA256.
`coverage_complete=false` and `clean=null` remain unchanged: the count includes
one engine-limit heuristic and one EICAR detection, not two malware findings.

An initial harness assertion incorrectly required a clean filename to appear
in native output. The parent removed that assumption and revalidated this same
retained execution, preserving the failed assertion log and corrected reusable
fixture source. No repeat scan is implied by the corrected validation result.

This execution predates the early-master movement and used local transport.
It validates native engine behavior only, not the corrected SSH bootstrap,
authentication, forwarding, or remote recovery lifecycle.

## Focused corrections and pending work

The parent applied a timeout around controller VERSION, rejected `#` in report
roots/scan roots/exclusions, and restricted ControlPath to 90 bytes to reserve
OpenSSH's 17-byte temporary suffix. Current source contains these corrections.
The parent reports focused review found no code issues in those changes and
verified the actual E2E receipts and cleanup. The role README now states the
same input and socket-path constraints.

The parent reports the transport correction now selects controller-only
ControlPath allocation, length validation, persistent arguments, and successful
master startup as the first protected-block task, before the initial boot read
or any endpoint tempfile/module bootstrap. Caller arguments are captured once;
explicit mux operations target only that master, without connection resets.
Collection and cleanup use the same master before it exits. Only successful
forward setup (`rc=0`) receives explicit cancellation; a forward that was never
attempted needs no release. The initial boot result is saved before its success
assertion. Full `boot_before_result`, `boot_after_result`, and
`master_start_result` are retained privately, including failed/unreachable
results. Boot-after is skipped when neither local transport nor a successful
dedicated master is available.

The first XOA early-master canary established its dedicated master, read boot
identity, retained evidence, and completed cleanup with unchanged boot identity
and restored caller settings. Its reverse Unix forward failed, so it does not
validate a forwarded scan. Its private receipt is
`out/host-security-updates-20261029/clamav/canary-reports/xoa.xcp-ng.alwaldend.com-clamav-wg5ucyjt/receipt.json`.
Parent read-only inspection then confirmed the transport directory was owned
by root with mode 0700 (`clamav/xoa-transport-owner.log`). Inventory
`ansible_become: true` overrides the task keyword `become: false` in pinned
Ansible. Exactly the transport tempfile, permissions, and removal tasks now
set task variable `ansible_become: false`; root scanner tasks and server SSH
policy are unchanged. Final parent remote-canary acceptance remains pending.

The local skipped-forward guard and connection-user ownership correction
passed the pinned Ansible lifecycle fixture under
`out/host-security-updates-20261029/clamav/role/lifecycle-fixture/`.
`source-hashes.json` identifies the candidate tasks; `execution.log`,
`lifecycle-result.json`, and `ownership-result.json` retain the results.
The fixture exercised the actual task definitions: a skipped local forward
preserves completion without attempting forwarding; an early setup failure
releases its unbound transport scratch without becoming complete; a mux
request against an absent private control socket fails without authentication
fallback and preserves uncertain transport scratch. With inventory escalation
enabled and `/bin/false` as the escalation executable, the actual transport
tasks created a connection-user-owned 0700 directory and removed it. These
offline checks do not establish remote SSH behavior.

Local fixtures bypass SSH setup. They do not validate remote authentication,
forwarding, expired-session recovery, incompatible endpoints, or production
roots. Native exclusions, recursion/stream/parser limits, and live-file changes
remain explicit coverage limits. Remote canary acceptance, production host
scans, and final delivery validation remain pending under parent ownership.

The final early-master, skipped-forward, and ownership corrections passed
pinned Ansible syntax for all task files
(`clamav/role/early-master-final-syntax.log`), role/docs packaging
(`clamav/role/early-master-final-build.log`), and scoped formatting
(`clamav/role/early-master-final-format.log`). Formatting left the lifecycle
fixture's three candidate task files byte-for-byte unchanged. These
checks do not establish remote behavior. The source task performed no live,
authentication, engine, or scan operation; the parent owns remote canaries and
the operational acceptance verdict.

## Final forwarded canary

The final candidate role completed an XOA `/usr/bin` scan: native exit 0,
zero detections and reported errors, reports retained, boot identity unchanged,
forward released, connection restored, and all three owned scratch directories
removed. The retained receipt is
`out/host-security-updates-20261029/clamav/canary-reports/xoa.xcp-ng.alwaldend.com-clamav-gxhx59ky/receipt.json`.
The canary reused the existing private engine and did not manage its lifetime.
No server SSH policy was changed. Full production scans remain a separate
operational outcome; this canary does not establish complete host coverage.

## Shared SSH default: offline acceptance

Failure cases recorded before implementation: the included role shadows a
same-named wrapper default; same-name forwarding recurses; an empty/custom
override is lost; existing certificate or CA references disappear; normal
discovery derives different signed paths; missing certificate input silently
becomes an empty list; or the helper adds wrapper discovery/signing/lifecycle
actions. The offline fixture must load actual roles with pinned Ansible and
render the existing upstream template. It must not execute upstream host tasks.

Offline acceptance passed on 2026-10-30 with pinned Ansible 2.20.7 and devsec
10.3.1. Evidence is retained under
`out/host-security-updates-20261029/ssh-hardening/defaults-evidence/`:

- `result.json` retains the original role-loading and upstream rendering
  results and exact candidate source hashes before the certificate assertion.
  Default, empty/custom overrides, normal signed-path derivation, required
  certificate input, CA preservation, and native restart-handler checks passed.
- `default.diff` shows only the three-rule `Match User ansible` addition and
  its generated heading; normal and maintenance fixtures render identically.
- `fixture.log` records those checks and offline syntax acceptance for both
  entrypoints. `SSH-FIXTURE-INVENTORY`: the initial inventory-setting warning
  is preserved in `fixture-inventory-warning.log`; explicit `-i localhost,`
  corrected the
  harness without source changes or endpoint contact.
- `build.log` records successful `bazel_agent bazel build` of
  `//projects/ansible_collection/roles/ssh:role` and `:docs`. The generated
  defaults document was inspected against source. `SSH-BUILD-CACHE`: the initial
  sandbox denial is retained in `build-sandbox-denied.log`; the authorized
  retry used the existing cache and succeeded.

These checks did not execute upstream host tasks or validate live sshd
configuration. Live SSH acceptance and the approved nine-host rollout remain
parent-owned. No scanner source, owner inventory, or upstream handler changed.

Follow-up failure cases recorded before the guard: a missing or string-valued
certificate input must fail before any upstream tasks, while an explicitly
empty list remains valid for callers intentionally omitting host certificates.
The fixture directory was moved into the original task root without changing
historical receipts or overwriting other maintenance evidence.

The extended `run_fixture.sh` and `fixture.py` passed against the assertion
candidate. `assertion-result.json` records current source hashes and all eight
rendering/input cases; `assertion-fixture.log` also records syntax acceptance
for both entrypoints. The fixture checks that the assertion precedes the native
include, rejects strings and mappings, and accepts an explicit empty list.
`missing-certificates.yaml` executes the actual helper locally with no facts
or certificate input. Its retained `missing-certificates.log` shows failure at
`Validate SSH host certificates`, exit 2, zero changes, and no upstream tasks.

The earlier `result.json`, `fixture.log`, and build receipts remain unchanged.
The build receipt predates the assertion; this follow-up used no Bazel command
and performed no endpoint, authentication, or service operation.

## Shared SSH default: live acceptance on all nine hosts

The parent accepted the approved shared default through the guarded native
maintenance workflow on all nine hosts. The directly read rollout index is
`out/host-security-updates-20261029/ssh-hardening/rollout-result.json`, observed
at `2026-10-30T06:30:05.262307+00:00`. It records `verified_count: 9`, the
approved `ansible` policy, nine per-host receipts, and nine configuration
hashes. Every host records successful fresh SSH, unchanged CA and boot identity,
preserved prior master, rejected reverse and direct TCP forwarding, and closed
owned master. The index records zero reboots and zero remaining owned master
sockets. The parent also confirms effective-policy checks on all nine hosts.

The directly read secure-runner receipt is
`out/host-security-updates-20261029/clamav/canary-reports/secure.forgejo-runner.alwaldend.com-clamav-1vfrqfiq/receipt.json`.
This reverse-Unix canary scanned `/usr/bin/true`, `/usr/bin/bash`, and
`/usr/bin/ssh`: native exit 0, zero detections and reported errors, valid
summary, completed execution, and retained evidence. Boot identity is unchanged,
the forward is released, the connection is restored, and remote, transport,
and controller scratch all have state `absent` without cleanup failures.
There are no collection errors. `coverage_complete: false` and `clean: null`
remain unchanged; this three-binary canary does not establish full-host or
full-fleet scan completion.

The rollout index retains two operational qualifications: reverse-Unix
data-path acceptance is demonstrated here on the secure runner, and the preview
AL wrapper reported a token-revocation cleanup error after successful Ansible
recaps. It also records the recovered ingress pre-play stdin failure and
cleanup of the replaced owned socket. These records are preserved rather than
treated as additional scan coverage. Full-fleet production scans and delivery
remain pending in task 2.3; the scanner itself never changes server policy.

## Simplified role acceptance

The candidate reuses native temporary-directory creation, direct archive
transfer, and checksum-validated fetch, with a single receipt template per role.
The isolated ClamAV fixtures in
`out/host-security-updates-20261029/simplification/scanners/clamav-fixture/summary.json`
passed three scenarios and 21 receipt checkpoints: unreachable extraction stops
before probing, failed transfer preserves recoverable evidence while collecting
the other log, and repeated invocations reset state. Source hashes were unchanged.

The real native shared-workflow fixture in `simplification/native-e2e/result.json`
passed clean input (exit 0), oversized input followed by EICAR (exit 1, two
findings, zero errors), missing input (exit 2), and EICAR plus missing input
(exit 1 with one error). Evidence remained private, scratch cleanup succeeded,
later hosts continued after failure, and controller cleanup preceded the final
aggregate failure. These local inputs do not validate remote SSH or production
coverage; all receipts retain `coverage_complete=false` and `clean=null`.

Normal SSH setup rendering passed default, empty, and custom Match User inputs,
certificate/CA mapping, and tag selection. Native `sshd -t/-T` using disposable
fixture identities confirmed only the three intended ansible-policy changes and
unchanged root policy (`simplification/ssh/native-policy-result.json`). No live
SSH service was changed. Both focused simplification reviews returned no findings.

## Selected-host native performance validation (2026-10-30)

The authorized ingress-only run uses a separate Host Bot controller with eight
workers, a 16 GiB memory ceiling, normal CPU/I/O scheduling, and native
`--stream --multiscan`. Existing active runtimes and units are unchanged.
The captured release runtime matches the maintained source; local syntax and
packaging passed. `priority-scan/validation/native-fixture-summary.json` retains
real-engine clean, oversized-plus-EICAR, missing-file, and mixed-error results
with successful controller cleanup. That full fixture preceded the conservative
multiscan exit-1 completion guard and is not evidence for that final predicate.

The actual native interrupted-session fixture then exercised late EOF, silent
early EOF, and a successful detection through the current collection tasks:
`priority-scan/session-fixture/guard-recheck/summary.json`. All three preserve
terminal status and findings but correctly leave completion uncertain. No
custom scanner is deployed. This intentional native-version limitation prevents
false completion while retaining the faster scan's useful evidence.

`priority-scan/validation/live-controller-status.json` records the separate
production unit's effective resources and the unchanged original unit.
Production endpoint results remain pending in private `priority-scan/reports`;
starting a scan is not completion or a host-clean verdict. Paths are relative
to `out/host-security-updates-20261029`.

## Priority ingress outcome and dependency rebase

The dedicated ingress scan reached native terminal status on 2026-10-30 after
386 seconds (09:42:07–09:48:33 UTC). Retained console output contains three
`Heuristics.Limits.Exceeded.MaxFileSize` alerts and 232 unsupported-file-type
errors, with no other signature matches. Exit 1, reported errors, and native
session uncertainty keep `complete=false`; this is not a host-clean verdict.
Evidence was retained, boot identity remained unchanged, endpoint/transport
scratch was removed, and the dedicated controller stopped with verified cleanup.
The private receipt is under `priority-scan/reports/host1.ingress.alwaldend.com-clamav-3w0vt8mo`.
Original fleet scans were left running.

Dependency PR #215 merged into master. Rebase retains its dependency definitions
and owner documentation; the security PR remains based on master.

## Full Ansible review and native result simplification

On 2026-10-30, the reviewed candidate excludes actual block/character device
nodes discovered through native Ansible metadata. The receipt records the
exclusions and discovery errors. Regular files beneath `/dev/shm` remain
eligible; incomplete discovery fails before scanning. The asynchronous scan
uses the finished native result instead of a second stdout status protocol.
Synchronous preflight and Lynis password-transport wrappers are unchanged.

Four real native-role cases passed: clean input (0), `/dev` with EICAR beneath
`/dev/shm` (1, one detection, no errors), oversized input followed by EICAR
(1, two findings, no errors, conservatively incomplete multiscan), and a
detection plus missing input (1, one detection and one error, incomplete).
All retained evidence and removed only their own terminated scratch. The
separate background-priority fixture engine used verified cached signatures
and was stopped; the running Host Bot scan and its runtime were unchanged.

Pinned native async execution preserved statuses 0, 1, 2 and 124; five missing,
unfinished, unreachable, negative or noninteger result cases supplied no status.
The initial harness used mismatched async directories, and its namespace could
not follow the outside daemon's `/proc` alias. Those failed attempts remain
retained; explicit fixture async paths and a short socket bind corrected the
harness without product changes. Native-source hashes match the accepted run.

Evidence: `out/host-security-updates-20261029/ansible-review/validation-summary.json`,
`native-role-02/`, and `async-results-04/`. These local fixtures do not establish
a new successful password-authenticated remote scan or full host coverage.

## Native result and IPC simplification (2026-10-30)

Derived fields now project registered native results; one final Ansible assertion
sets completion and its always block saves the outcome. Earlier checkpoints
remain incomplete. Native runtime directories replace relative control paths
and controller-CWD discovery. Authentication, same-master recovery, conservative
multiscan status, independent retention and no-reset behavior are preserved.

Four current-source native cases passed: clean (0, complete), device exclusions
with EICAR beneath /dev/shm (1, complete), oversize plus later EICAR (1, two
findings, conservatively incomplete multiscan), and detection plus missing input
(1, one error, incomplete). Native async cases preserve 0/1/2/124 and reject five
missing or invalid results. Three controlled failure-handling invocations pass
21 receipt checkpoints: failed fetch retains remaining evidence and recovery
scratch; unreachable extraction stops before probing; a later successful run
resets state and removes owned scratch.

Native IPC fixtures accept private runtime storage, reject mode 0755, preserve
caller options and show that mux attachment to an absent allocated socket fails
without authentication fallback. The connection fixture is local, not a new
remote SSH canary. Packaged fleet/appliance syntax and collection builds passed.

Evidence under `out/host-security-updates-20261029/native-simplification/`:
`e2e-result.json`, `control-fixture/summary.json`, `socket-result.json`, and
`package-results.json`. The first failure-fixture harness incorrectly expected
a private role projection outside include_role; its failure is retained, and
the corrected harness checks the persisted public receipt. Production scanning
is stopped; no host-clean or complete-coverage conclusion is made.
