---
title: Forgejo
description: Forgejo installation and controlled local SQLite upgrades
---

The role manages Forgejo's executable, configuration and systemd unit. When
startup contents change, it stages the candidate, drains and stops the old
service, archives the declared files, installs and starts the candidate, then
checks HTTP health, the running version and database consistency. Unchanged
runs repair file ownership and modes without restarting or archiving.

## Declare the deployment layout

This workflow supports a local SQLite database, local storage and an HTTP
backend, including HTTP behind a separately managed TLS proxy. Set
`forgejo_database_path` to the SQLite file used by the configuration template;
it defaults to `{{ forgejo_work_dir }}/data/forgejo.db`.

`forgejo_backup_paths` is the complete archive contract. Defaults cover
`forgejo_dir`, the configuration, executable and unit. Add every external local
storage directory and referenced secret file, including an external Git home,
LFS/package storage or token files. The role checks declared roots, managed
state/startup coverage and symlink targets. It does not discover paths from
INI settings. Keep the existing configuration's persistence paths stable;
coordinate relocation and external-database/storage backups separately.
Use absolute file-backed secret references across releases.

Paths must be literal normalized absolute paths. Recovery storage defaults to
`/var/backups/forgejo` and must be outside the archive roots. Nested roots are
collapsed before archiving. The managed host must provide GNU `realpath`;
separation is checked after resolving ancestor symlinks, including when the
recovery directory does not exist yet. Existing roots and the configured SQLite file must
exist before an upgrade. Partial installations and downgrades fail before
stopping. The existing Ansible INI filter checks SQLite/HTTP scope and rejects
changes to `SECRET_KEY` or `SECRET_KEY_URI`; it does not interpret storage paths.

The role manages the account and `acl`, `git`, `git-lfs` packages by default.
Callers disabling `forgejo_manage_user` or `forgejo_manage_packages` supply those
prerequisites. ACL support is needed for unprivileged Ansible module execution.
Serialize upgrades per instance. Check mode cannot verify this transaction.

## Recovery

Every changed existing installation retains a private recovery directory with
`snapshot.tar.gz` and `manifest.json`. Directories use mode `0700`; archives and
manifests use `0600`. They remain on the managed host and contain secrets.
Provision sufficient space and independent disaster-recovery backups; retention
cleanup is an operator responsibility.

If draining fails, the old instance stays running. If preparation or archiving
fails before startup-file replacement, an intact previously active service is
restarted. After replacement begins, failed acceptance leaves the candidate
stopped. There is no automatic downgrade against potentially migrated data.
A `pending.json` marker blocks retries until explicit recovery or a deliberate
corrected retry; it is removed only after successful acceptance.

To restore, stop Forgejo, verify the selected archive/checksum, and move aside
its candidate data trees. Extract the complete archive into the manifest's
`archive_root`, preserving ownership and permissions. Recreate omitted roots
from `directory_roots`, applying numeric UID/GID before their octal mode.
Restore the old executable, configuration and unit together with their data,
reload systemd and start that version. Verify service and application behavior
before removing the pending marker. Restoration can discard later writes and
requires explicit operator authorization.

## Verification

`bazel test //projects/ansible_collection/test/forgejo_upgrade:upgrade_test`
runs actual Forgejo 15.0.3 and 17.0.1 with disposable SQLite state and a bounded
service-manager fixture. It checks upgrade/restore, unchanged runs, explicit
backup validation and failure recovery. Sanitized JSON results are retained;
fixture credentials and archives are discarded. Real systemd, cross-account
SSH, production data and cryptographic signing are outside this rehearsal.
