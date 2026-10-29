## Context

See proposal.md. The shared role owns native endpoint scanning and SSH transport.
The historical native runs established a working Host Bot systemd user manager
and verified all three signature databases. Those observations apply to their
captured invocations; the parent owns live execution and current acceptance.

## Goals / Non-Goals

Use native Ansible, ClamAV, and systemd behavior without a custom supervisor or
database parser. No persistent service, package installation, reboot, quarantine,
or maintenance-role execution. XCP-ng scanning and compatibility work are deferred.

## Decisions

Controller resource overrides map directly to native settings:
`clamav_controller_threads` (default 2, range 1–32) sets MaxThreads and four
times that value sets MaxQueue; `clamav_controller_memory_gib` (default 6,
range 1–64) sets MemoryMax in GiB. `clamav_scan_background` is a shared caller
boolean, with a controller fallback of true before the endpoint role is loaded.
True uses Nice 19 and idle I/O; false uses Nice 0 and best-effort I/O. Validate
inputs before scratch creation and retain selected values in the controller
object carried by every lifecycle receipt. Existing runtime, signature,
private-socket, and engine-limit policy remains unchanged.

One shared three-play workflow prepares the controller, scans selected endpoints,
and stops its unique transient user unit. Thin native imports supply each
inventory/auth boundary's host pattern and connection variables. Startup failures
stop the same unit in a rescue path. The final play attempts cleanup after host
failures, and its `always` block reports accumulated host failures even when
cleanup fails. RuntimeMaxSec, TimeoutStopSec=15,
and KillMode=control-group bound daemon lifetime if playbook cleanup is skipped;
there is no immediate Ansible-parent-death guarantee.

Extract the pinned payload into an invocation-owned mode 0700 directory and
retain the cache across runs, including freshclam cooldowns. Native flock
prevents simultaneous update commands. Explicit library and certificate
paths apply to all tools and freshclam's database-test subprocess. Verify
daily/main/bytecode with sigtool, extract its metadata using Ansible filters,
check daily age, and use one daemon VERSION polling task matching the verified
daily version and a server response shape. One outer timeout bounds the whole
loop, including a stalled client, with a 15-second forced-stop grace. The local
client-version fallback is rejected. Main and bytecode ages are recorded without
a daily age limit because their release cadence differs; future timestamps are
still rejected. The RAM and disk minimum checks remain.

The pinned archive owns engine identity. Remove the separate local engine-version
loop and `engine.json`; retain the actual daemon version in lifecycle evidence.
`signatures.json` owns signature metadata without a duplicate lifecycle copy.
Starting the actual transient unit establishes user-manager availability; there
is no separate `is-system-running` preflight.

The foreground daemon binds `clam.sock` in its native systemd RuntimeDirectory,
mode 0700 beneath the operator's private XDG runtime directory. Systemd removes
the directory when the unit stops. Separate invocation-owned runtime directories
hold SSH control sockets, so engine termination cannot interrupt evidence
recovery. Reports, logs, extraction and cached databases remain under ignored
`out/`. No process/CWD aliases are needed.

## Failure cases

Resource override failure cases recorded before implementation (2026-10-30):

- Reject non-integral thread or memory values, threads outside 1–32, memory
  outside 1–64 GiB, and a non-boolean background setting before controller
  scratch creation or engine execution.
- Preserve omitted-value behavior: two threads, queue eight, 6 GiB memory
  ceiling, Nice 19, and idle I/O. Controller preparation runs before shared
  role defaults are loaded and must supply its own background fallback.
- Explicit foreground scheduling must select Nice 0 and best-effort I/O;
  changing threads must change MaxQueue to four times the thread count.
- Record selected resources and scheduling in lifecycle evidence. Larger
  ceilings do not establish available headroom; the parent checks the selected
  controller and validates the exact candidate with native E2E execution.
- Preserve private-unit ownership, signatures, Unix-only transport, runtime
  bounds, scan limits, and every running scan's captured runtime.

The original implementation and subsequent simplification recorded their failure
cases before editing source. The simplification plan is retained under
`out/host-security-updates-20261029/simplification/controller/acceptance.md`.

- Reject root execution, unavailable user manager, non-private or symlinked
  report/cache paths, insufficient available RAM/disk, or unwritable paths.
- Refuse a concurrent update and preserve all existing cooldown state.
- Retain an incomplete receipt before engine execution; preserve native logs
  when extraction, configuration, daemon VERSION validation, or freshclam fail.
- Reject absent/corrupt databases, malformed sigtool metadata, future build
  times, or stale daily signatures even when freshclam returns zero.
- Reject daemon startup failure, a local VERSION fallback, a wrong daily version,
  a malformed response, or a stalled VERSION request. Reject unsafe runtime
  directories and private-socket violations.
- Stop only the unique transient unit on startup or host-task failure. If the
  playbook is interrupted, the unit's bounded runtime is the cleanup fallback.
- Report an unconfirmed stop as failure while still reporting accumulated host
  failures; never claim lifecycle cleanup proves endpoint coverage or absence
  of malware.
- Let systemd own process identity and socket-directory lifecycle; no
  application PID file or process/CWD alias is needed.
- Keep ordinary maintenance targets, inventory owners, and auth selectors intact;
  default the appliance ClamAV play to XOA without changing XCP-ng maintenance.

## SSH policy

The shared SSH hardening role owns the approved forwarding default. Normal
owner deployments apply it with existing role variables, tags, and host limits.
The completed nine-host rollout's dedicated playbook and certificate helper
are retired; their private before/after and recovery evidence remains retained.
The scanner does not apply server policy. Normal hardening retains its upstream
restart handler; the former reload-only and rollback workflow is no longer a
current endpoint contract.

## Risks / Trade-offs

Linux, systemd user services, a private XDG runtime directory, and native
timeout/flock/tar are controller prerequisites. A transient unit may remain until its runtime limit if Ansible
is killed. Use one controller scan workflow at a time against the shared cache.
The cache is not deleted to recover from signature download failures. Native
scan completion does not establish full coverage: exclusions, live filesystem
changes, recursion, stream limits, and parser limits remain. Role receipts keep
`coverage_complete=false` and `clean=null`, including successful clean fixtures.

## Migration Plan

The superseded controller helper and synthetic fixtures were preserved under
ignored task scratch and removed from source. The current simplification shares
the lifecycle, removes redundant probes/evidence, and retires the completed SSH
rollout command/helper. Preserve all historical receipts and the running scans'
captured runtime; source changes do not alter those invocations.

Seven isolated polling cases passed for the simplified VERSION task. They do
not establish native engine behavior, import-variable precedence, packaged
execution, or SSH transport. Build and syntax-check both simplified packages,
then complete parent-owned native acceptance against the exact candidate.
Historical clean/EICAR/oversize, canary, and nine-host SSH rollout results remain
in validation.md. Production scan completion and final delivery stay pending.

### Oversized input continuation

The initial 100 MiB server stream quota caused the native client to abort a
whole directory at the first oversized file. Both native stream quotas now
allow 104857601 bytes, while the engine file limit remains 104857600 bytes.
With heuristic and limit alerts enabled, that extra byte produces
`Heuristics.Limits.Exceeded.MaxFileSize` and traversal continues. The limit
alert denotes omitted inspection, not malware. Acceptance requires a real
oversized input followed by the EICAR fixture, zero protocol errors, retained
native results, and confirmed controller cleanup. No manifest or custom
scanner is introduced.

## Native simplification acceptance

Use systemd RuntimeDirectory with mode 0700 for the daemon socket; remove PID/cwd aliases and retain evidence under ignored out.
Failure cases: unavailable/private runtime storage, daemon startup failure, lost or
unfinished native status, scanner warning/error/timeout, failed artifact transfer,
and failed SSH forward or close. Preserve caller authentication, same-session
recovery, native signature verification, private evidence and no-reboot defaults.
Production scanning was stopped at the user's request; validate only bounded fixtures.
