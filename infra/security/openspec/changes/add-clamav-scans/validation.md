## Controller resource overrides, 2026-10-30

Operational failure cases were recorded in design.md before implementation.
Source now validates thread and memory integer ranges and the shared background
boolean before controller scratch creation. Normalized settings in the
controller object drive the native template and systemd properties and remain
in incomplete, ready, and final lifecycle receipts. Defaults remain two threads,
queue eight, 6 GiB, Nice 19, and idle I/O. The requested eight-thread, 16 GiB,
non-background invocation selects queue 32, Nice 0, and best-effort I/O.

Source diff inspection and `git diff --check` passed. This implementation task
performed no Bazel, authentication, live engine/service, scan, or commit
operations and left captured runtimes unchanged. Parent-owned package/syntax,
OpenSpec CLI, native E2E, controller headroom, and ingress acceptance remain
pending under task 7.2; no native validation of these overrides is claimed.

## Current simplification and validation scope

The current source uses one shared three-play lifecycle with thin fleet and
appliance import callers. Each caller retains its host defaults and connection
variables; inventory and authentication packaging remain separate. The final
play attempts controller cleanup before reporting accumulated host failures,
including when cleanup fails.

One bounded daemon VERSION polling task checks the verified daily version and
a daemon response shape, rejecting the client's local-version fallback. The
local engine-version loop, `engine.json`, separate PING task, and user-manager
status preflight are removed. `signatures.json` retains signature metadata;
lifecycle evidence retains the actual server version without a duplicate
signature copy. Native signature verification/freshness, cache/flock/certificate
paths, minimum RAM/disk checks, private modes, bounded systemd lifetime, and the
workspace socket's MainPID/cwd alias remain.

The completed one-time SSH maintenance playbook, targets, and helper are retired.
The forwarding default remains in the normal shared SSH role. Owner deployments
use its certificate handling and upstream restart handler; the retired
reload-only and rollback workflow is not a current endpoint contract.

Seven isolated polling cases passed: ready daemon-shaped response, local
fallback, wrong daily version, malformed response, delayed readiness, explicit
client failure, and a stalled client. Timeout cases returned 124 in about two
seconds with the fixture's two-second bound; the explicit client failure kept
exit 2. These tests exercise the polling process boundary, not ClamAV, database
verification, Ansible import precedence, systemd, or SSH transport. The
pre-implementation failure plan, results, exact source hashes, and handoff are
under `out/host-security-updates-20261029/simplification/controller/` in
`acceptance.md`, `polling-results.json`, `source-after.json`, and `validation.md`.

Combined package/syntax validation and representative native acceptance of the
simplified candidate remain parent-owned and pending, as do production scan
completion and final delivery. This documentation pass performed no Bazel,
authentication, engine, service, or scan operations and changed no captured
runtime or historical evidence. The sections below preserve earlier recorded
results for their historical candidates; they do not establish acceptance of
the current simplification.

## Historical candidate and scope

The 2026-10-29 native candidate replaced the uncommitted Python controller and
synthetic engine fixtures with shared Ansible tasks, two native configuration
templates, and a bounded transient systemd user service. The obsolete package
is preserved under ignored
`out/host-security-updates-20261029/clamav/superseded-controller/`; its old Bazel
test artifacts were no longer available at their cache links. The earlier
fixture-test pass applies only to the superseded implementation.

Source ownership is `infra/security`, with the diagram's root BUILD changes and
existing Lynis work preserved. Shared roles and dependency pins are owned by
other tasks. No commits, pushes, authentication, live engine execution, or host
scans were performed by this implementation task.

## Historical offline validation

- Both `//infra/security/ansible:clamav_bin` and
  `//infra/security/appliances:clamav_bin` built successfully.
- Both `:clamav_syntax` targets passed without AL authentication or native
  engine execution. These parse YAML/playbook structure, not runtime expressions.
- The final package build completed at 23:30:05 controller local time. Immediate
  runfiles-manifest inspection confirmed both owning playbooks, the same pinned
  payload, shared controller tasks, and the two native templates.
- Private build and mapping evidence is retained in
  `out/host-security-updates-20261029/clamav/validation/` as `packages.log` and
  `package-mappings.txt`.
- The configured Prettier and buildifier tools formatted task-owned files.
  Buildifier's changes to existing maintenance declarations were formatting only.
  `git diff --check` passed after those changes.
- `OPENSPEC_PROJECT=infra/security bazel_agent bazel run //tools/openspec --
validate add-clamav-scans --strict` passed. Its private command log is retained
  alongside the package evidence as `openspec.log`.

The initial source frozen for parent materialization includes Nice=19,
IOSchedulingClass=idle, MemoryMax=6G, and acceptance of a collected unit's
LoadState=not-found as cleanup evidence. Controller start SHA256:
`cd03e91a23f21f6f50b2a501b778bc015bccd7331c108c61d06712e7c4c163dd`.
Controller stop SHA256:
`aee7b5e47aa5f7691e91dda235eb56ac7261fea692d7f9e19a426c0ada5ec262`.

## Historical native E2E acceptance

The parent captured both runtimes (8713 fleet and 8703 appliance entries),
verified their source hashes, and completed the real native fixture run.
`out/host-security-updates-20261029/clamav/e2e/result.json` records `passed=true`
and identifies the three role receipts and controller evidence directory.
Reading those retained receipts confirms:

- Clean fixture: native exit 0, infected files 0, reported errors 0, and
  `complete=true`.
- EICAR fixture: native exit 1, infected files 1, reported errors 0, and
  `complete=true`; a detection is distinct from an execution error.
- Oversized fixture, 104857602 bytes: native exit 2, infected files 0, reported
  errors 1, and `complete=false`. It did not receive a successful prefix verdict.
- All three have `execution_finished=true`, `evidence_retained=true`, retained
  `clamav.log` and `console.log`, unchanged boot identities, and removed verified
  scan scratch. All retain `coverage_complete=false` and `clean=null`.

The successful controller's `lifecycle.json` records `state=stopped`,
`cleanup_complete=true`, no startup failure, and a collected inactive unit.
Its server identity is ClamAV 1.5.4/daily 28168. `signatures.json` records daily
28168, main 63, and bytecode 339. These are observed native results, not results
from the superseded synthetic fixtures.

The initial native run's PID-file failure remains retained at
`out/host-security-updates-20261029/clamav/e2e/reports/clamav-are9aljf`.
Clamd created that redundant file with mode 0644 despite UMask=0077; its socket
was mode 0600 and rescue cleanup succeeded. The parent removed the PidFile
configuration and checks. Systemd MainPID, the checked cwd, and the private
socket now establish identity. The hashes above identify the earlier candidate,
not the corrected native runtime.

## Historical review corrections and remaining scope

The parent applied three focused corrections: native controller VERSION is
wrapped by timeout; report roots, scan roots, and exclusions reject `#`; and
ControlPath accepts at most 90 bytes, reserving OpenSSH's 17-byte temporary
suffix. Source inspection confirms these changes. The parent reports that
focused review returned no code findings and verified the E2E receipts and
cleanup. This does not claim that each negative input or timeout boundary was
exercised by the fixture run.

The fixtures used a local connection. They do not establish SSH forwarding or
recovery behavior, full-host coverage, or production scan completion. Native
limits, exclusions, live-file changes, and parser behavior remain coverage
constraints; successful execution does not establish host cleanliness.

The remote canary was running under parent ownership at this checkpoint;
subsequent XOA and secure-runner results appear below. Fleet scan completion
and final delivery validation remain pending. XOA is the appliance scan default;
XCP-ng remains unscanned because of the official payload's glibc requirement.
The caller must exclude the whole private task tree from a local Host Bot scan
and include localhost in any limit. Host failures remain visible while later
serial batches continue, with aggregate failure reported after controller cleanup.

This follow-up changed documentation only. No Bazel, authentication, engine,
service, or scan commands were run; updated artifact validation belongs to the
parent's delivery checks.

## Historical session ergonomics

- `SEC-CLAMAV-LIFECYCLE-SCOPE` (configured): a custom helper and two synthetic
  fixtures were replaced after the parent confirmed an existing native service
  manager. Future lifecycle design should inspect that capability before adding
  a supervisor. No shared policy change is proposed.
- `SEC-BAZEL-OUTPUT-RETENTION` (observed): built outputs disappeared between
  invocations, independently confirmed by the parent. Building and immediately
  extracting relevant mappings in one invocation produced retained evidence.
  Capture runtime artifacts in the same command rather than repeatedly probing
  stale convenience links.
- `SEC-TOOL-WORKDIR` (observed): one formatter invocation used relative paths and
  found no files. Absolute workspace paths succeeded, as the existing Bazel skill
  already requires. No additional instruction or test is needed.

## Historical native oversized-input continuation

The initial production attempt exposed a limitation missed by the isolated
protocol-rejection fixture: rejecting an oversized stream aborts the current
native directory traversal. That attempt remains incomplete and retained.
The controller now accepts the same 104857601 bytes as the client, with
`MaxFileSize 100M`, `HeuristicAlerts yes`, and `AlertExceedsMax yes`.

The real engine then processed one directory containing an oversized input
before an EICAR fixture. Native output recorded
`Heuristics.Limits.Exceeded.MaxFileSize` followed by `Eicar-Test-Signature`,
exit 1, two native detections, zero reported errors, retained evidence, and
confirmed transient-unit cleanup. The limit result denotes omitted content,
not malware. The receipt still has `clean=null` and `coverage_complete=false`.

Evidence: `out/host-security-updates-20261029/clamav/e2e/continuation-result.json`.
An initial harness assertion incorrectly expected the clean filename in normal
native output. That assertion was removed; the unchanged real execution and
its retained output were verified directly. The corrected fixture remains
repeatable under `clamav/runtime/_main/native-limits-e2e.yaml` in task scratch.

The appliance ClamAV caller now disables inherited generic multiplexing;
the shared role owns its private connection lifecycle. This prevents initial
connections from allocating control sockets outside the explicit report root.
At this checkpoint server SSH policy had not changed, and permission for an
account-scoped Unix-forwarding exception was pending separately from scan
authorization. The subsequently approved maintenance results appear below.

## Historical forwarded XOA canary

The final shared role completed a native `/usr/bin` scan on XOA with exit 0,
zero native detections and reported errors, independently retained reports,
unchanged boot identity, released forwarding, restored connection settings,
and absent root, transport, and controller connection scratch. Evidence:
`out/host-security-updates-20261029/clamav/canary-reports/xoa.xcp-ng.alwaldend.com-clamav-gxhx59ky/receipt.json`.
This is bounded canary coverage, not a full appliance scan or proof of absence.
It reused the already validated, running private engine; the canary did not
stop that other invocation's unit.

Earlier failures remain retained: non-multiplexed XOA module bootstrap failed,
then an inventory `ansible_become=true` overrode transport task keywords and
created a root-owned socket directory. The role now establishes its private
master before bootstrap and explicitly disables escalation only for connection-
user transport-directory creation, permissions, and removal. The final canary
required no server SSH policy change. The XOA-only packaged scan target now
injects only XOA authentication, omitting unrelated XCP-ng host credentials.

## Historical approved SSH maintenance implementation

The retired 2026-10-30 maintenance play consumed the inventory package and AL
maintenance selector and invoked the former shared SSH helper. It contained no
copied policy, owner-inventory change, custom SSH transport, or connection reset.
Its failure cases were recorded in the design before implementation. The
current design describes the normal owner-deployment path after retirement.

`ssh_hardening_bin` and `ssh_hardening_syntax` built successfully, and the
packaged syntax target passed without authentication. Prettier and buildifier
accepted the source. Logs are under
`out/host-security-updates-20261029/ssh-hardening/workflow-evidence/`.
The eight task-owned logs were relocated there with identical content hashes;
historical receipts were not rewritten. Syntax success does not validate
runtime expressions or the proposed native configuration. Parent-owned
runtime results are recorded below.

Pinned devsec 10.3.1 source uses distribution, major version, OS family, and
SELinux facts; crypto settings derive from `ssh -V`. The selected service-manager
collector depends on platform and distribution, so the bounded setup includes
system/architecture facts transitively. Source inspection also confirms that
the role reconciles boot enablement, renders revoked keys, and runs SELinux
reconciliation when enabled. At that rollout, the parent reported all nine
services already enabled, active, and reloadable with KillMode=process. The nine-host preview
review below found no ancillary remote changes from these native-role tasks.

`SEC-SSH-DISCOVERY-VOLUME` (observed): the initial full worktree inventory and
combined skill reads produced truncated output. Subsequent commands used
bounded component reads, and the current checkout was verified directly.
Existing workspace instructions already prescribe that approach; no policy
change is needed. The shared Bazel cache required sandbox escalation; builds
used the existing cache after parent coordination, without a duplicate cache.

## Historical SSH preview and secure-runner acceptance

`out/host-security-updates-20261029/ssh-hardening/preview-review.json`, observed
at 2026-10-30 06:25:07 UTC, records nine hosts, nine diff hunks containing only
the approved Match block, no ancillary remote changes, and zero failed Ansible
hosts. Every receipt under `ssh-hardening/preview-reports/` records check mode,
unchanged config, CA trust and boot identity, no host failure, and skipped
reload. They correctly retain `proposed_config_validated=false`.

The preview wrapper exited 1 because Vault token-revocation cleanup failed
after the successful Ansible recap. This is a wrapper cleanup failure, not an
Ansible host failure; successful previews do not establish proposed-config
validity or successful credential cleanup.

The secure-runner apply receipt at
`ssh-hardening/apply-reports/secure.forgejo-runner.alwaldend.com-ssh-2ifjsjzn/receipt.json`
records a validated candidate, successful reload, no failure or rollback,
unchanged CA content and boot identity, and identical root effective policy.
The ansible effective policy differs only in the approved values:
`allowtcpforwarding remote`, `permitlisten none`, and
`allowstreamlocalforwarding remote`. Retained files under
`ssh-hardening/host-operations/secure-runner/` record wrapper exit 0, the
original master still running after apply, successful fresh login, and native
reverse-TCP and direct-TCP rejections (both exit 255 with the expected denial).

The native Unix-forwarding canary receipt at
`clamav/canary-reports/secure.forgejo-runner.alwaldend.com-clamav-1vfrqfiq/receipt.json`
records a completed scan of `/usr/bin/true`, `/usr/bin/bash`, and `/usr/bin/ssh`,
exit 0, zero native detections and reported errors, retained evidence, released
forwarding, restored connection settings, and removed remote, transport, and
controller connection scratch. It retains `clean=null` and
`coverage_complete=false`; this establishes bounded canary behavior only.

Paths above are relative to `out/host-security-updates-20261029/` unless given
in full.

## Historical completed nine-host SSH policy rollout

`out/host-security-updates-20261029/ssh-hardening/rollout-result.json`, observed
at 2026-10-30 06:30:05 UTC, verifies all nine policy applications. Inspection
of every referenced apply receipt confirms native candidate validation,
identical root effective policy, only the approved ansible directive values,
unchanged CA content and boot identity, and final configuration checksums
matching the validated candidates and rollout record. No apply receipt records
a failure or rollback.

The rollout record and its referenced operation logs retain successful fresh
signed SSH logins, rejected reverse-TCP and direct-TCP forwarding, and the
original task master usable after each successful reload and subsequently
closed. All nine boot identities were preserved, with zero reboots and zero
remaining owned master sockets. The native Unix data-path canary remains the
secure-runner result above; full scans on other hosts provide additional
transport evidence as they finish.

Ingress attempt 1 failed locally before Ansible play execution because stdin
was nonblocking. Reopening stdin allowed the successful retry. The original
attempt's master PID was absent and its socket was replaced and closed; the
failure and retry evidence remain retained. This recovered controller failure
and the earlier preview token-revocation cleanup failure are not erased by
the successful rollout.

That rollout's fleet SSH policy acceptance completed. It does not validate
the current normal owner-deployment path. Full scan completion and final
delivery validation remain pending; the historical rollout receipt does not
establish complete scan coverage or host cleanliness.

## Simplified workflow acceptance

Both native caller packages built and passed syntax checks. Their inventory
listings retain the fleet selection and XOA-only appliance default. The initial
appliance listing used the wrong inventory filename; the corrected check used
the packaged `inventory.yaml` and selected XOA successfully.

The isolated real-engine run at
`out/host-security-updates-20261029/simplification/native-e2e/result.json` used the
shared three-play workflow, a separate copied signature cache, and an
invocation-owned controller. Clean, size-limit-then-EICAR, missing-file, and
detection-plus-error cases all produced the expected receipts. Failed hosts did
not prevent later hosts running. The controller stopped before the expected
aggregate failure; lifecycle records `cleanup_complete=true`. The first result
checker used a nonexistent cleanup key; correcting the checker reused the
original execution and artifacts. Production runtime and scan state were unchanged.

The role fixtures and native SSH policy comparison are recorded in the
collection's `add-clamav-scan/validation.md`. Appliance boot-evidence checks
passed both unchanged and changed-identity cases under `simplification/boot-check`;
the initial local fixture lacked a writable Ansible remote-temp path, corrected
to task-owned scratch. Seven isolated VERSION polling cases and both focused
source reviews passed. No remote execution of the simplified candidate is claimed.

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

## Native socket lifecycle simplification (2026-10-30)

Systemd now creates and removes the mode 0700 daemon RuntimeDirectory under
the operator's private XDG runtime root. Reports, extracted binaries and the
reusable signature cache remain under ignored task scratch. No PID/CWD socket
alias is used. SSH masters use separate private runtime directories, so stopping
the engine cannot remove a connection needed for evidence recovery.

The current source controller started a real ClamAV 1.5.4 engine against verified
cached signatures; its native VERSION response matched the daily database.
Four bounded local native cases passed: clean, device exclusions with EICAR in
shared memory, oversized input followed by EICAR, and detection plus missing
input. The source stop tasks confirmed termination and RuntimeDirectory removal.
Fleet and appliance packages and all three packaged syntax invocations passed.

Evidence: `out/host-security-updates-20261029/native-simplification/e2e-result.json`,
`native-role-02/engine/cleanup.json`, and `package-results.json`. These fixtures
do not establish a new remote production scan. The user-requested production
stop is recorded under `hostbot-background-20261030-122450/operator-stop/stop.json`;
partial logs were retained, and production scans were not restarted.
