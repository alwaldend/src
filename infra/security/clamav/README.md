---
title: ClamAV controller lifecycle
---

The fleet and XOA callers import one shared Ansible lifecycle playbook, passing
their host patterns and connection variables while retaining separate inventory
and authentication packaging. Native tasks and configuration templates prepare
an unprivileged controller. The official pinned payload is extracted into a
mode 0700 invocation directory; it is not installed as a host package. No custom
controller helper or database parser is used.

The playbooks require `clamav_scan_report_directory` and
`clamav_database_directory` as existing operator-owned mode 0700 absolute paths
under ignored task scratch. Paths use ASCII letters, numbers, underscores,
periods, slashes, and hyphens. Preparation requires a working systemd user
manager, native tar/flock/timeout tools, at least 4 GiB available RAM, 5 GiB free
disk on both paths, and writable directories.

One native freshclam invocation preserves the existing database cache and its
cooldown state. An exclusive nonblocking flock prevents overlapping update
commands. Use one controller scan workflow at a time against the cache. The
relocated tools receive `LD_LIBRARY_PATH=prefix/usr/local/lib` and
`CVD_CERTS_DIR=prefix/usr/local/etc/certs`; freshclam/clamd configuration and
sigtool's command also specify the certificate directory. The environment is
needed by freshclam's database-test subprocess as well. TLS, code-signature
verification, and `TestDatabases yes` remain enabled.

Native `sigtool --info` verifies daily/main/bytecode containers. Ansible filters
retain their versions, signature counts, build times, and ages. Daily signatures
may be at most 72 hours old by default (`clamav_daily_max_age_seconds=259200`).
Main and bytecode ages are recorded without a daily-age limit because their
release cadence differs; future timestamps beyond five minutes are rejected.
A successful freshclam return code alone is not evidence of freshness.

`systemd-run --user` starts foreground clamd as a uniquely named transient unit
with RuntimeMaxSec (default 172800 seconds), TimeoutStopSec=15,
KillMode=control-group, UMask=0077, and default Nice=19, IOSchedulingClass=idle,
and MemoryMax=6G. `clamav_controller_memory_gib` selects the integer memory
ceiling in GiB (default 6, range 1–64), mapped to systemd's `MemoryMax`.
The shared caller boolean `clamav_scan_background` defaults to `true`;
`false` selects Nice=0 and IOSchedulingClass=best-effort. Supply booleans through
JSON/YAML extra variables, for example
`--extra-vars '{"clamav_scan_background": false}'`. Controller preparation
supplies the default before the endpoint role is loaded. Select overrides only
after checking controller headroom; the existing RAM/disk minimums remain.
The unit owns the process group, its private working directory, and native log
output. One native `clamdscan --version` polling task waits for a daemon response
containing the engine version, verified daily version, and database timestamp;
the client's local-version fallback is rejected. One outer timeout bounds the
whole loop, including a stalled client, with a 15-second forced-stop grace.
The archive pin owns the engine version; the actual server version is retained
in lifecycle evidence. Real clean-file and EICAR validation remain separate
acceptance steps.

The daemon binds a mode 0600 Unix socket inside a mode 0700 native
`RuntimeDirectory`, under the operator's `XDG_RUNTIME_DIR`. Systemd owns its
creation and removal with the transient unit. This gives clients a short
absolute path without process/cwd aliases. Reports, logs, extraction and the
signature cache stay under ignored workspace scratch. The native runtime
root must be operator-owned mode 0700. SSH controls use separate private
runtime subdirectories, retained until authenticated recovery is finished.

`clamav_controller_threads` selects native `MaxThreads` (default 2, integer range
1–32); `MaxQueue` is four times that value (default 8). Thread and memory inputs
accept integers or decimal digit strings; invalid values fail before engine
execution. Lifecycle receipts retain the selected resources and explicit
scheduling settings in their `controller` object.

The fixed engine policy uses no concurrent reload, a 100 MiB file limit,
matching client/server stream limits of 100 MiB plus one byte,
400 MiB aggregate scan size, 10,000 archive members, recursion
17, and 120,000 ms scan time. Heuristic alerts, encrypted-content alerts, and
limit alerts are enabled; heuristic precedence is disabled. ClamAV 1.5.4 has no
`ScanEncrypted` option. Encrypted content and limit alerts are limitations or
heuristics, not confirmed malware. Neither deletion nor quarantine is configured.
The extra stream byte makes oversized inputs produce an explicit file-limit
alert instead of accepting a truncated prefix or aborting directory traversal
at the protocol limit.

The first play prepares the controller; the selected-host play invokes the shared
role and retains failures while continuing later serial batches; the final play
attempts controller cleanup before reporting accumulated failures. Host-failure
reporting still runs when cleanup fails. Every `--limit` must include localhost.
Startup failure stops the named unit before failing the first play. If Ansible
is killed or cleanup is otherwise skipped, RuntimeMaxSec is the
fallback: shutdown is bounded, not immediate on parent-process exit. The private
lifecycle receipt records the exact unit for `systemctl --user stop UNIT` recovery.
No process-name search or unverified PID signal is used. The unit is transient;
no permanent user service or host configuration is installed.

Mode 0600 configuration, signature, freshclam, clamd, and lifecycle artifacts
remain in the invocation directory. Signature metadata is owned by
`signatures.json`; lifecycle evidence retains the actual server version without
duplicating signature metadata. The initial receipt is incomplete before engine
execution. Cleanup confirmation is not an endpoint scan verdict.
Native role reports must be checked for detections, read errors, exclusions, and
coverage limits. For a local Host Bot scan, the caller explicitly excludes the
whole task scratch tree, including the database cache and controller extraction.

Build `//infra/security/ansible:clamav_bin` and
`//infra/security/appliances:clamav_bin` to check packaging. Their separate
`:clamav_syntax` targets parse playbooks locally without AL authentication or
engine/host execution. The parent owns real-engine, clean/EICAR, and fleet
validation. The appliance scan defaults to XOA; XCP-ng is unscanned because of
the native payload's glibc requirement.
