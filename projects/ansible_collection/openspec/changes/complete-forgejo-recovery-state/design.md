## Decision

Proceed with the user's requested explicit-path upgrade contract. Remove the
custom configuration interpreter. Existing Ansible modules and the already
pinned `community.general.from_ini` filter are sufficient. Retain that native
filter only for literal database-key binding checks, never to
infer filesystem locations. The caller declares external storage and secret
files in `forgejo_backup_paths` and keeps path migrations separately coordinated.

Keep private staging, version/downgrade checks, stopped-instance archives,
root-directory metadata, startup-file permission repair, the pending marker,
HTTP/version/database acceptance and safe pre-installation failure recovery.
Drop the doctor path check and inferred-directory creation; Forgejo creates its
runtime directories. Expose the known SQLite file as `forgejo_database_path`.

Native Ansible path normalization and escaped prefix matching validate absolute
inputs, core/startup coverage and disjoint backup storage before activation.
The caller owns symlink targets. Pass declared inputs directly to tar without
filesystem discovery or nested-root reduction.
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

## Direct archive simplification

Proceed with a standard GNU tar archive after the service stops. The archive
includes declared directory entries and their numeric ownership and modes,
removing the archive module's separate root-metadata collection and manual
restoration step. Keep the private directory, restricted archive file, checksum
and small version/path manifest. A nonzero tar exit prevents installation and
recovers an intact previously active service. The Linux host supplies GNU tar
and gzip alongside the existing coreutils prerequisite; no new repository
dependency or host installation is introduced by this source change.

Retaining the archive module would preserve the extra bookkeeping. Removing
metadata preservation would break removed-tree restoration. Direct tar meets
the same contract with fewer tasks. Reconsider if the real-binary restore
rehearsal cannot recover directory ownership/modes without manifest metadata.

Before changing implementation, adapt that existing rehearsal to extract the
archive directly and retain its root-metadata assertions. Preserve missing-root,
archive-failure and private-permission checks; do not weaken failure recovery.

## Host-side recovery separation

Resolve the recovery destination and declared archive roots with the managed
host's GNU `realpath --canonicalize-missing` before checking separation. This
follows existing ancestor symlinks while allowing a not-yet-created recovery
directory. Keep declared archive paths and logical coverage unchanged, including
normal system paths such as `/lib` on merged-usr hosts. Reject resolved filesystem
roots and containment in either direction before staging or downtime.

Acceptance covers a recovery ancestor pointing into managed data, an archive
ancestor pointing inside recovery storage, and a disjoint recovery alias that
still upgrades successfully. The original negative reproduction used an unchanged
installation; current negative cases request activation and reject before
stopping. This fixes the existing separation guarantee; no parser or
additional configuration model is introduced.

## Narrow preflight review corrections

Before implementation, extend the existing real-binary scenarios to restore
two explicitly declared storage roots differing only by case, and to reject
unchanged interpolated `SECRET_KEY` and `SECRET_KEY_URI` expressions whose
referenced values change. Rejection must occur without service actions or
archives. These cases cover gaps in the existing recovery/key-preservation
contract, without creating a new test suite or configuration interpreter.

Pass declared roots directly to tar, preserving case-distinct and overlapping
paths without a deduplication algorithm.
Reject interpolation markers in either installed or candidate key bindings;
the upgrade role does not resolve recursive INI expressions. Both corrections
fit the literal-key assertion and direct path list without new preflight tasks.

## Preflight simplification

The user requested the role simplification after reviewing the retained
preflight. Keep content-based change detection, interrupted-upgrade and downgrade
protection, literal key-binding preservation, database/root presence and physical
backup separation. Validate backup inputs only when activation is required.
Run one candidate version probe on the managed host, using the staged file for
changed installations and the installed file for unchanged installations.

Remove recursive symlink discovery, nested-root reduction and repeated
controller/target execution checks. The caller owns external persistence,
including symlink targets. The supplied template owns the supported SQLite/HTTP
layout; the role does not revalidate its schema. Retain the INI filter solely
for comparing the two literal database-key bindings. Declare the fixed startup
file metadata in ordinary role variables, keeping preflight focused on runtime
checks. Existing fixture assertions continue to protect unchanged permission
repair, complete restoration and activation failure boundaries.

Before implementation, adapt negative backup-separation cases to request an
actual configuration change, remove the schema-revalidation expectation, and
exercise overlapping declared roots in the existing removed-tree rehearsal.
An isolated standard-tar rehearsal already restored overlapping entries through
both GNU tar and Python's tar reader. These are source changes only.
