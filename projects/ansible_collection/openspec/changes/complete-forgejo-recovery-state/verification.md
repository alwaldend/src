# Explicit-path upgrade verification

The user requested simplification of the unmerged upgrade role. The custom
381-line planner is removed. Native Ansible now validates explicitly declared
archive inputs, core/startup coverage and backup separation before activation.
The existing collection INI filter is used only for unchanged literal database-key
bindings; it does not discover filesystem paths or revalidate the supplied schema.

Seven real-binary scenarios retain stopped-instance upgrade and complete
removed-tree recovery, no-op permission repair, downgrade/configuration checks,
and drain/archive/marker/version/doctor failures. External objects and secrets
are explicitly declared by the fixture. Parser-specific inference tests are
removed with that implementation.

All seven integration scenarios and 125 selected Bazel tests pass. Semantic
lint, role/deployment packaging, explicit Ruff and strict OpenSpec validation
also pass. The first rehearsal caught a fixture that expected archive checks
on an unchanged run; the corrected case triggers a configuration upgrade.
Bounded command logs,
sanitized scenario JSON and exact Git/delivery identities remain under ignored
`out/forgejo-upgrade/`. No deployment or live inventory is used.

The fixture does not cover real systemd, cross-account SSH/ACL transfer,
production data or cryptographic signing. Caller declarations own completeness
of external persistence and separately coordinated path migrations.

## Managed-host path separation

The former lexical comparison accepted recovery storage with a symlinked
ancestor inside managed data. The new scenario reproduced that acceptance
on an unchanged run, without stopping the service or taking an archive.
The corrected preflight resolves paths on the managed host before checking
containment in both directions. It also rejects resolved filesystem roots.

The real-role scenario rejects recovery-inside-data and data-inside-recovery
through ancestor aliases, then completes an actual configuration upgrade with
a disjoint recovery alias and a not-yet-created destination. All seven native
scenarios and 125 selected Bazel tests pass. Archives keep their declared path
identity; ordinary system-directory aliases do not become unsupported roots.

## Direct archive recovery

The backup now uses one GNU tar invocation and seven Ansible tasks instead of
eleven. Directory entries carry their own numeric ownership and permission
bits, so the separate root-metadata manifest list and Python restore helper
are removed. The small manifest retains versions, paths and archive checksum.

All seven real-binary scenarios pass with direct extraction into the recorded
archive root. The removed-tree rehearsal checks directory ownership and modes,
the SQLite database, committed Git contents, external objects, file-backed
credentials and an empty storage root before starting the prior executable.
Archive permissions and failure recovery remain covered. This validation uses
temporary local processes and standard host tar/gzip; it does not deploy or
establish production service health.

## Reduced preflight and review corrections

The role no longer recursively audits symlinks, reduces nested archive roots,
executes candidates on the controller, or checks the supplied SQLite/HTTP schema.
It probes candidate execution once on the managed host, validates backup inputs
only for changed startup contents and keeps startup-file metadata in role variables.
Downgrade, unfinished activation, database presence, physical backup separation
and literal key-binding protection remain in place.

All seven real-binary scenarios pass. The removed-tree rehearsal restores both
`Storage` and `storage`, plus an explicitly repeated nested signing directory,
without dropping case-distinct inputs or reconstructing directory metadata.
The configuration scenario rejects unchanged interpolated `SECRET_KEY` and
`SECRET_KEY_URI` expressions with changed referents without stopping the process
or taking an archive. Both review findings are corrected without a new parser.

An earlier narrow-fix test run was cancelled when the user expanded the task to
simplify the role. The completed rehearsal above covers the final flow; the
cancelled run is not passing evidence. Sanitized receipts and bounded logs
remain under ignored `out/forgejo-upgrade/`.
