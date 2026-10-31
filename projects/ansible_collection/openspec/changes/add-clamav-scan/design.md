## Context

The user selected official native `clamdscan --stream` for simplicity and reuse.
The controller holds the engine and fresh signatures. Older XCP-ng compatibility
is deferred rather than adding another scanner implementation.

## Goals / Non-Goals

Provide private, bounded native scans and recoverable evidence. Do not install
host packages, expose TCP, reset unrelated SSH connections, reboot, quarantine,
or delete scanned files. Do not implement a custom protocol client, file walker,
or per-file coverage framework.

## Decisions

### Follow-up failure cases

The default native traversal read `/dev/fuse` on both observed LXC endpoints
and aborted before inspecting the host. Discover block and character devices
through Ansible metadata, reject incomplete discovery, exclude only those
nodes, and retain their paths. Regular files, including `/dev/shm`, remain
eligible; do not add another file walker or scanner client.

ClamAV execution already uses Ansible async polling. Read a scanner exit only
from a finished native result with a valid return code; missing, unreachable,
or unfinished results remain incomplete. Exercise detection, operational error,
and timeout outcomes without weakening the existing multiscan completion guard.
Synchronous runtime/connectivity preflight wrappers retain their password-transport boundary. Lynis now uses native async too.

### Maintained behavior

- The caller supplies the official DEB `data.tar.gz`, private clamd socket, and
  controller-owned 0700 report root. The role extracts root-owned scratch and
  supplies bundled library/certificate paths explicitly. Native `--help`
  preflight records execution failures without guessing their cause. Failed or
  unreachable extraction stops before that probe. Report roots, scan roots, and
  exclusions reject `#` to preserve native configuration parsing.
- Native configuration disables symlink following, sets recursion to 1000, and
  escapes exclusions plus `/proc`, `/sys`, and role scratch. `/dev/shm` regular
  files remain eligible.
- Controller `StreamMaxLength 104857601` equals the client stream limit from
  the unchanged `max_file_bytes + 1` expression; `max_file_bytes` remains 104857600. The controller sets `MaxFileSize 104857600`, `HeuristicAlerts yes`,
  and `AlertExceedsMax yes`. Official 1.5.4 clamps the client payload, but the
  extra byte now reaches the engine and produces
  `Heuristics.Limits.Exceeded.MaxFileSize FOUND`. Native traversal continues
  without a successful truncated-prefix verdict. The former 104857600 server
  stream quota caused a protocol error and aborted directory traversal; its
  E2E evidence is superseded. Oversize followed by later EICAR acceptance passed
  locally for the new policy; see validation.md. That execution predates the
  early-master change and does not validate SSH transport.
- Allocate the dedicated ControlPath entirely on the controller and establish
  its persistent master before any endpoint command, including boot identity,
  tempfile creation, and Ansible module bootstrap. Capture caller arguments
  once before overriding them. Preserve injected authentication and aliases;
  explicit mux operations check, forward, cancel, and exit only this master.
  No `meta: reset_connection` is used. Native ping confirms connect/send, not
  receipt of a daemon PONG. ControlPath is limited to 90 bytes, reserving
  17 bytes for OpenSSH's temporary suffix.
  Transport-directory allocation and removal explicitly set the
  task variable `ansible_become: false`, overriding inventory escalation so
  the SSH connection user owns its 0700 listener directory. Scanner tasks
  retain root escalation.
- GNU timeout sends TERM then KILL after 15 seconds; Ansible async adds slack.
  A finished native Ansible async result carries scanner status as data;
  no second status protocol is maintained for the asynchronous scan.
- Recovery cancels the forward, collects logs and cleans retained scratch over
  the existing authenticated master, then exits it and restores caller options.
  Only successful forward setup (`rc=0`) receives explicit cancellation; a
  local, skipped, or never-attempted forward needs no release. Failed or
  uncertain remote forwarding retains conservative cleanup. Long scans do not normally
  require renewed certificates for collection.
- Retain native console/log files and a private controller receipt. Parse only
  the final summary, including errors that exit 1 can mask. Native infected
  counts include limits and encrypted-content heuristics and are not
  malware-only counts. Retain both boot-before/boot-after command results
  privately, including failed and unreachable errors; failed identity checks
  cannot establish unchanged boot identity or workflow completion. Preserve
  the master-start result privately and save the initial boot result before
  asserting success. Skip boot-after when neither local transport nor a
  successful dedicated master is available. Execution completion, evidence
  handling, detections, and coverage remain separate.
  No host-clean or complete-coverage assertion is emitted. Receipts keep
  `coverage_complete` false and `clean` null even when native execution succeeds.

The shared SSH wrapper defaults `ssh_hardening_match_user` to `Match User
ansible` with `AllowTcpForwarding remote`, `PermitListen none`, and
`AllowStreamLocalForwarding remote`. Its devsec include explicitly maps this
distinct wrapper variable to `ssh_server_match_user`; the upstream default
would otherwise override a same-named wrapper default. Callers can replace
the list or supply `[]`. No deployment inventory needs a copy of the policy.

The existing SSH setup passes its discovered signed certificates, trusted CA,
and forwarding default directly to devsec. Normal owner deployments select
the `ssh` tag and a host limit; devsec retains its native handler. The completed
nine-host rollout's separate helper and maintenance command are retired.
Their private evidence remains available. The scanner never changes SSH policy.

Native `tempfile` owns temporary-directory permissions. Native `unarchive`
transfers and extracts the pinned archives, and native `fetch` validates copied
log bytes. Role-local receipt templates define each evidence field once and
retain the existing allocation, execution, collection, and cleanup checkpoints.

## Risks / Trade-offs

The authorized single-host priority run uses native `--multiscan` to pipeline
streams to the controller worker pool and normal CPU/I/O scheduling. Defaults
remain background scheduling and sequential submission. Both choices are
explicit boolean inputs, retained in receipts; neither changes roots, parser
limits, evidence handling, or SSH policy. Before deployment, verify clean and
EICAR inputs, oversized input with later-file continuation, summary/error
accounting, and cleanup under session mode. Reject non-boolean controls.
ClamAV 1.5.4 can report exit 1 and omit summary errors after a detection followed
by premature session closure. Check native error diagnostics across the retained
console, separately from the summary, and retain that check before accepting
completion. Native limits and absence-of-compromise qualifications remain.
An EOF during traversal can be silent, so any multiscan exit 1 retains native
findings but marks completion uncertain and leaves `complete` false. This avoids
reimplementing the native protocol or pretending that a successful-looking
summary proves all pending streams were processed.
The separate controller retains bounded memory and lifetime; ongoing scans
keep their original captured runtime and controller.

Streaming consumes network bandwidth and sends file contents only to the
private authorized controller. Files remain live; exclusions, recursion and
engine limits constrain coverage. Native logs do not provide per-file coverage
accounting. The caller must keep all four controller stream/engine settings
consistent with the client limit; equal stream limits alone do not prevent
prefix success. Lost masters may require fresh authorized credentials for
recovery.

## Migration Plan

The superseded Python helper, role, and change artifacts are preserved under
ignored `out/host-security-updates-20261029/clamav/role/previous-python`; original
fixtures remain under ignored `clamav/client`. The helper is removed from
packaged source. Parent-owned clean, EICAR, and oversized-file fixtures passed
under the former server quota with retained native evidence and controller
cleanup. That historical E2E is superseded for current-policy acceptance.
The revised local fixture confirmed an oversized-file engine-limit finding
followed by a later EICAR detection in the same traversal, with evidence
retention and cleanup. It predates the early-master change and validates
engine behavior only. Later parent-owned XOA and secure-runner canaries
established bounded forwarded execution and cleanup. The approved shared SSH
default was accepted on all nine rollout hosts; validation.md identifies the
receipts and limits. Full-fleet production scans and delivery remain pending.
These bounded canaries do not establish complete host coverage or
expired-session recovery.

## Native simplification acceptance

Use a short private controller runtime directory for SSH IPC, remove CWD/relative-path tricks, and compute the outcome once from registered results.
Failure cases: unavailable/private runtime storage, daemon startup failure, lost or
unfinished native status, scanner warning/error/timeout, failed artifact transfer,
and failed SSH forward or close. Preserve caller authentication, same-session
recovery, native signature verification, private evidence and no-reboot defaults.
Production scanning was stopped at the user's request; validate only bounded fixtures.
