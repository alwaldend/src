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
storage directory and referenced secret file, including symlink targets, an
external Git home, LFS/package storage or token files. The role checks declared
inputs and managed state/startup coverage before an activation. It does not
scan trees for symlinks or discover paths from INI settings. Keep the existing configuration's persistence paths stable;
coordinate relocation and external-database/storage backups separately.
Use absolute file-backed secret references across releases.

Paths must be literal normalized absolute paths. Recovery storage defaults to
`/var/backups/forgejo` and must be outside the archive inputs. Declared paths
pass directly to tar, preserving case-distinct and overlapping roots. The host
must provide GNU `realpath`, GNU `tar` and `gzip`. Before activation, separation
is checked after resolving symlinks, including a not-yet-created recovery directory.
Existing inputs and the configured SQLite file must exist before an upgrade.
Partial installations and downgrades fail before stopping.

The supplied template owns the SQLite/HTTP layout. The role uses the existing
INI filter only to reject changes to `SECRET_KEY` or `SECRET_KEY_URI`.
Key bindings must be literal values; interpolated `%(...)s` expressions are
rejected before downtime rather than resolved by the upgrade role.

Candidate execution is checked once on the managed host. Changed inputs are
staged before this check; unchanged runs inspect the installed executable.

The role manages the account and `acl`, `git`, `git-lfs` packages by default.
Callers disabling `forgejo_manage_user` or `forgejo_manage_packages` supply those
prerequisites. ACL support is needed for unprivileged Ansible module execution.
Serialize upgrades per instance. Check mode cannot verify this transaction.

## Optional custom header

Set `forgejo_custom_header_local_path` to a controller-side Forgejo template
file. The role installs it under `forgejo_work_dir/custom/templates/custom/header.tmpl`
through the guarded activation workflow. The file contents are copied literally.
An empty source path leaves the installed header unmanaged.

## Optional SSH commit-signing pair

Set `forgejo_signing_ssh_key` to a dictionary with `private_key` and `public_key`
strings to install a dedicated Ed25519 pair before service activation. It defaults
to an empty dictionary. The caller owns secure generation and storage and must
supply Git 2.34+ and OpenSSH `ssh-keygen` 8.2p1+ on the managed host. Configure
the service template with `FORMAT = ssh` and
`SIGNING_KEY = {{ forgejo_dir }}/signing/instance.pub`, plus its signing identity.

The role validates the pair, installs private/public files with modes 0600/0644
under a service-owned 0700 directory, and rejects different existing private
content. Matching reruns repair metadata and public content without restarting.
The files remain under `forgejo_dir` and its existing backup coverage. The role
does not generate, rotate, or publish keys. Secret-bearing tasks suppress output;
their temporary candidate is removed even on validation failure.

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
its candidate data trees. Extract the complete archive into `/` (the manifest's
`archive_root`), preserving numeric ownership and permissions. The archive
contains the directory roots and their metadata; no separate reconstruction
step is needed.
Restore the old executable, configuration and unit together with their data,
reload systemd and start that version. Verify service and application behavior
before removing the pending marker. Restoration can discard later writes and
requires explicit operator authorization.

## Verification

`bazel test //projects/ansible_collection/test/forgejo_upgrade:upgrade_test`
runs actual Forgejo 15.0.3 and 17.0.1 with disposable SQLite state and a bounded
service-manager fixture. It checks upgrade/restore, unchanged runs, explicit
backup validation and failure recovery. Sanitized JSON results are retained;
fixture credentials and archives are discarded. The signing case verifies a real
API-generated commit with the installed SSH pair and checks unchanged runs and
private-key mismatch refusal. Real systemd, cross-account SSH and production
data are outside this rehearsal.
