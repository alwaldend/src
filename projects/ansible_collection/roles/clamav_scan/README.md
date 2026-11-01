---
title: ClamAV scan
description: Native endpoint scans through a private controller ClamAV socket
---

`alwaldend.main.clamav_scan` unpacks the official pinned ClamAV payload into
private endpoint scratch and runs native `clamdscan --stream`. The caller owns
the controller daemon, fresh signatures, parser settings, and its lifecycle.
The role installs no host package, service, schedule, or persistent SSH
configuration. It does not reboot, quarantine, delete scanned files, or expose
TCP. Check mode is unsupported.

Supply `clamav_scan_archive` as the controller path to
`//third_party/com_github_cisco_talos_clamav:payload_tar`, the compressed DEB data
payload containing `usr/local`. Endpoints need Linux x86_64, compatible glibc,
tar, GNU timeout, nice, ionice, root escalation, and Ansible's existing Python
runtime. Failed or unreachable extraction stops preparation before the native
execution preflight; both results are retained privately. A failed preflight
does not by itself identify an architecture or library mismatch. Older XCP-ng
compatibility is deferred; there is no custom scanner or compatibility fallback.

Supply `clamav_scan_socket` as an existing absolute controller Unix socket and
`clamav_scan_report_directory` as an existing absolute controller-owned 0700
directory under ignored task scratch. Socket paths use ASCII letters, digits,
underscores, dots, slashes, and hyphens. IPC sockets use the controller's native
`XDG_RUNTIME_DIR`, required to be operator-owned mode 0700. Each SSH master
gets its own temporary subdirectory there, independent of the daemon's runtime
directory so daemon failure cannot remove a connection needed for recovery.
The absolute control path is limited to 90 bytes, reserving OpenSSH's temporary
suffix. Report and scan paths must not contain `#`.

Remote runs allocate the dedicated ControlPath entirely on the controller and
establish a persistent SSH master before any endpoint command, including boot
identity, tempfile creation, and Ansible module bootstrap. Caller arguments are
captured once before overrides; the master preserves injected authentication,
host aliases, and caller options. Endpoint operations then create a
connection-user-owned 0700 Unix transport directory and separate root-owned
0700 scan scratch. Transport-directory creation and removal set
the task variable `ansible_become: false` so inventory escalation cannot change
their owner; scanner tasks retain root escalation. Explicit mux operations target only this master to attach
one reverse Unix forward, cancel it, and exit the master. No
`meta: reset_connection` is used. Local runs use the direct controller socket.
Native `--ping=1:1` is a connectivity preflight: ClamAV 1.5.4 prints PONG after
connect/send without reading a daemon reply. Actual scan fixtures establish
that the daemon processes requests.

The native configuration disables file and directory symlink following, limits
recursion to 1000, and excludes `/proc`, `/sys`, role scratch, and
`clamav_scan_exclude_paths` with escaped, anchored subtree expressions. It does
not exclude all of `/dev`. Ansible inspects metadata beneath `/dev` without
following symlinks or opening file contents, and excludes the actual block and
character device nodes. Incomplete discovery prevents scanning. Their paths and
discovery failures are retained in the receipt; regular files under `/dev/shm`
remain eligible. Roots default to `[/]`.

The caller's daemon must use `StreamMaxLength 104857601`,
`MaxFileSize 104857600`, `HeuristicAlerts yes`, and `AlertExceedsMax yes`.
`clamav_scan_max_file_bytes` stays fixed at 104857600; the unchanged client
`max_file_bytes + 1` expression produces a stream limit of 104857601, equal to
the controller stream limit. This lets the extra byte reach the engine, where
the lower file limit reports `Heuristics.Limits.Exceeded.MaxFileSize FOUND`.
Native traversal continues to later files without a successful truncated-prefix
verdict. All four controller settings are required; the role cannot attest
caller-owned daemon configuration.

Scanning uses `timeout --signal=TERM --kill-after=15s`, `nice -n 19`, and
`ionice -c 3` by default. Set `clamav_scan_background: false` for Nice 0 and
best-effort I/O priority 4. Set `clamav_scan_multiscan: true` to pipeline native
streams for concurrent processing by the controller workers. It remains one
uploader; increasing daemon workers alone does not parallelize sequential
submission. Both controls are booleans recorded in the private receipt, and
neither changes coverage or parser limits. The deadline defaults to 86400 seconds, bounded to seven days;
Ansible async allows 60 seconds of additional supervision and polls every
15 seconds. Bundled library/certificate paths are explicit. The role takes the
native return code from a finished Ansible async result; missing, unreachable,
or unfinished results cannot supply a scanner status. Async polling returns
the result as data independently of the polling process exit. Synchronous
runtime/connectivity preflights retain their password-transport status wrappers.
Timeout remains incomplete.

Native exit 0 reports no detection, 1 reports detections, and 2 reports an
operational error. Exit 1 can mask errors, so the role also reads the final
summary's `Infected files` and `Total errors`. Those counters do not account
for every session failure: the role also checks the entire
retained console for native `ERROR:` diagnostics and file-error records before
accepting completion. The check's result is retained; failure to read the log
or an error diagnostic keeps the receipt incomplete. The native infected count
includes limits and encrypted-content heuristics; it is not a malware-only
count. The role does not parse every filename or invent a scanned-file count.
In ClamAV 1.5.4, a multiscan detection can also mask a silent session failure.
Consequently, multiscan exit 1 retains findings and terminal status but sets
`native_completion_uncertain: true` and `complete: false`, even if no errors
are printed. This conservative outcome includes heuristic-only detections;
it does not classify them as malware. Sequential scans retain their existing
completion rules, and multiscan exit 0 still requires the other checks.
Raw 0600 `console.log`, `clamav.log`, and
controller `receipt.json` remain in the unique 0700
`clamav_scan_result_directory`.

The initial receipt is incomplete and precedes endpoint operations. It records
boot identity before/after, recovery paths, native runtime diagnostics, exit,
summary counts, retention, and cleanup. The master-start result and both
boot-command results, including failed and unreachable errors, are retained
privately. The initial boot result is saved before its success assertion;
boot-after is skipped when neither local transport nor a successful dedicated
master is available. A failed identity check cannot establish an unchanged
boot identity. `execution_finished` records native terminal status, including
exit 2. `complete` additionally requires exit 0 or 1,
a valid summary without reported errors, unchanged boot identity, verified logs,
and successful cleanup. Neither means the host is clean: `coverage_complete`
remains false and `clean` remains null. Live files can change, exclusions and
recursion limits omit content, and parser settings constrain inspection.
Scanning a host disk image does not establish guest-filesystem coverage.

Native registered results are the source of status and summary fields. One
final Ansible assertion validates execution, evidence and cleanup, and its
result supplies completion. Receipts before that assertion remain incomplete;
there is no incrementally updated completion flag. `native_extraction` retains the full unarchive result.
Artifact transfer relies on Ansible fetch's checksum validation; collection
failure preserves endpoint evidence and its recovery path.

Recovery explicitly cancels only a role-owned forward whose setup succeeded
with exit code 0. Local, skipped, and never-attempted forwards need no release;
failed or uncertain remote forwarding remains conservative. It then
collects logs and cleans verified scratch through the existing authenticated
master. Long scans need no fresh certificate on this normal recovery path.
A nested `always` exits only that
master and restores caller arguments. Available logs are collected independently
on failure or unreachable results. Scan scratch is deleted only after both logs
are verified retained, a receipt is saved, and terminal process status is known.
Lost connections or failed retention preserve recovery paths; fresh authorized
authentication may be needed to recover them.
