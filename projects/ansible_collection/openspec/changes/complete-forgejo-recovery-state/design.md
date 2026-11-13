## Decision

Proceed with the user's requested explicit-path upgrade contract. Remove the
custom configuration interpreter. Existing Ansible modules and the already
pinned `community.general.from_ini` filter are sufficient. Retain that native
filter only for the SQLite/HTTP scope and database-key binding checks, never to
infer filesystem locations. The caller declares external storage and secret
files in `forgejo_backup_paths` and keeps path migrations separately coordinated.

Keep private staging, version/downgrade checks, stopped-instance archives,
root-directory metadata, startup-file permission repair, the pending marker,
HTTP/version/database acceptance and safe pre-installation failure recovery.
Drop the doctor path check and inferred-directory creation; Forgejo creates its
runtime directories. Expose the known SQLite file as `forgejo_database_path`.

Native Ansible path normalization and escaped prefix matching validate literal
absolute roots, core/startup coverage, disjoint backup storage and covered
symlink targets. Collapse nested roots before archiving to avoid duplicates.
No custom plugin, new dependency or live operation is required.

## Acceptance before implementation

- Upgrade and restore real SQLite/repository data and explicitly declared
  external objects, including directory metadata and file-backed credentials.
- Reject missing declared roots, omitted startup/core coverage, recursive
  backup destinations and a missing database before downtime.
- Preserve no-op permission repair, downgrade rejection and configuration
  changes guarded against an encryption-key binding change.
- Preserve drain, archive, marker-write, running-version and doctor failures.
- Remove parser-specific relocation/quoting/storage-discovery expectations.
  These belonged to the deleted interpreter, not the explicit archive contract.

The existing seven real-binary scenarios remain the acceptance boundary.

## Host-side recovery separation

Resolve the recovery destination and declared archive roots with the managed
host's GNU `realpath --canonicalize-missing` before checking separation. This
follows existing ancestor symlinks while allowing a not-yet-created recovery
directory. Keep declared archive paths and logical coverage unchanged, including
normal system paths such as `/lib` on merged-usr hosts. Reject resolved filesystem
roots and containment in either direction before staging or downtime.

Acceptance covers a recovery ancestor pointing into managed data, an archive
ancestor pointing inside recovery storage, and a disjoint recovery alias that
still upgrades successfully. The negative cases use an unchanged installation
so the prior faulty implementation can be demonstrated without creating a
recursive archive. This fixes the existing separation guarantee; no parser or
additional configuration model is introduced.
