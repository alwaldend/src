## Context

The existing role has no upgrade transaction. The deployed definition uses a
local SQLite database and local storage beneath `forgejo_dir`. See the proposal
and the upstream [upgrade](https://forgejo.org/docs/v17.0/admin/upgrade/) and
[backup](https://forgejo.org/docs/v17.0/admin/setup/backup/) guidance.

## Goals / Non-Goals

Provide a repeatable local SQLite upgrade and recovery archive. Do not deploy,
automatically downgrade, or implement external-database backup orchestration.

## Decisions

Decision review verdict: **proceed** with an offline archive and guarded
activation. Binary-only restoration is rejected because startup can migrate
the database. Automatic full restoration is also excluded: it could discard
writes accepted between startup and a later failed check. Preserve the archive
and stop a rejected candidate for explicit operator recovery.
Persist a private pending-activation marker before replacing startup files and
clear it only after acceptance. A failed or interrupted activation must not
become an apparently unchanged run that bypasses doctor checks on the next
invocation. Preserve the installed database-encryption-key configuration.

Stage and validate inputs before downtime; compare their checksums to the
installed files. Reject downgrades and partial installations. On an unchanged
run, ensure the service is started without restarting or creating a backup.
Fresh installations have no prior database to archive.

For an existing changed installation, inspect the installed and candidate
configuration without logging secrets. Support SQLite, local file storage,
and local queues. Reject external databases/storage/queues and persistence
paths outside the declared backup roots before stopping the service. Back up
the complete configured Forgejo tree, executable, unit, and explicit additional
paths, into a unique administrator-only directory outside those roots. This
host-local recovery copy does not replace an independent disaster-recovery
backup. Archives remain until an operator removes them.

Drain an active instance with its installed binary and a bounded timeout.
Stop it, create and verify the archive, then install the candidate. Restore
the old service's running state if backup fails before any candidate input is
installed. Once installation begins, stop on failure and retain the archive.
Require HTTP health, the expected running version, and non-repairing doctor
checks before reporting acceptance. Do not use SIGHUP to upgrade the process.

The service manager remains systemd. Test fixtures substitute only its command
interface in PATH, launch actual pinned Forgejo processes, and use temporary
local SQLite trees. They never contact the host service manager, production
Vault, or deployment inventories. Package and account management can be disabled
for caller-managed hosts, with ownership remaining explicit. This follows the
existing Traefik ownership-variable pattern. Collection test sources stay
outside role packaging globs.

## Failure cases before implementation

- Invalid candidate executable/version/configuration: reject before stopping.
- Older candidate version or partial existing installation: reject without
  overwriting startup inputs.
- External database/storage/queue, uncovered data paths, overlapping backup
  destination, or missing required backup input: reject before downtime.
- Queue drain fails or times out: retain the installed executable and service.
- Service stop or archive fails: never install a candidate; recover a previously
  active instance where its inputs are intact.
- Archive would disclose secrets: protect directory and files, suppress secret
  task output, and never fetch archives to the controller.
- Candidate copy, unit reload, startup, health, version, or doctor check fails:
  stop the candidate and retain recovery files; never run an older binary on
  potentially migrated data.
- Fresh install: no invented previous backup; startup failure remains visible.
- Unchanged run: no queue drain, backup, reload, or restart.
- Configuration-only change: use the same backup and acceptance transaction.
- Recovery archive omits SQLite, repositories, keys, or prior executable:
  fail behavioral acceptance; rehearse a complete restoration locally.

## Risks / Trade-offs

Offline archives add downtime and require free space. Retain backups for manual
inspection and restoration; no automatic retention deletion is introduced.
Unsupported storage layouts fail early rather than producing partial backups.
This includes individual `[queue.*]` overrides and queue `CONN_STR` overrides,
which can replace the configured local queue path. The latter require separate
coordination. Relative queue `DATADIR` paths resolve against `APP_DATA_PATH`,
as specified in the upstream
[queue configuration](https://forgejo.org/docs/v17.0/admin/config-cheat-sheet/#queue-queue-and-queue).

## Verification

Write integration scenarios before production code. Run actual Ansible tasks
and pinned Forgejo 15.0.3 and 17.0.1 in isolated local fixtures. Retain a
sanitized result artifact proving version change, repository persistence,
archive permissions/content, no-op lifecycle, restore, downgrade rejection,
and injected failure behavior. The service-manager substitute does not verify
real systemd behavior, OIDC, external storage, or production database contents.
Use repository packaging, semantic lint, and quality gates separately.

## Acceptance refinement from real-binary evidence

The first local service became HTTP-ready, but `doctor check --all` rejected
its missing fresh repository/log directories and disabled LFS garbage
collection. The pinned `services/doctor/lfs.go` confirms `gc-lfs` is optional
and rejects instances with LFS disabled. Decision review verdict: **revise**
the gate to explicit `paths`, `check-db-version`, and `check-db-consistency`
checks without `--fix`, and create the configured runtime directories before
startup. This verifies required activation behavior without running unrelated
maintenance checks or disabling doctor acceptance altogether.
