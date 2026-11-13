# Explicit-path upgrade verification

The user requested simplification of the unmerged upgrade role. The custom
381-line planner is removed. Native Ansible now validates explicitly declared
archive roots, core/startup coverage, symlink targets and backup separation.
The existing collection INI filter is used only for SQLite/HTTP scope and
unchanged database-key bindings; it does not discover filesystem paths.

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
