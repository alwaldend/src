# Forgejo upgrade acceptance

Run `bazel test //projects/ansible_collection/test/forgejo_upgrade:upgrade_test`.
It executes the packaged role with real pinned Forgejo 15.0.3 and 17.0.1,
disposable SQLite databases and local HTTP/Git operations. The bounded
`systemctl` fixture starts only its private processes, never the host manager.

Seven scenarios cover upgrade and removed-tree restoration, explicit backup
validation and configuration changes, unchanged permission repair, downgrade
rejection, and drain/archive/marker/version/doctor failures. The restore case
includes explicitly declared external repositories, storage, token and Git-home
state, an empty directory and retained root metadata. It also creates a new
organization repository after restoring the old version.

The role does not infer persistence paths or support implicit relocation.
Parser-specific quoting, storage-discovery and path-exchange cases were removed
with that implementation. Caller-declared roots and the configured SQLite file
are authoritative. Scenario failure cases and decisions live in the active
[change](../../openspec/changes/complete-forgejo-recovery-state/design.md).

Each successful scenario emits sanitized JSON in Bazel's undeclared outputs.
Generated credentials, private configuration and archives are removed during
cleanup. Real systemd, cross-account SSH/ACL transfers and actual signing are
outside fixture coverage. Synthetic signing state verifies preservation only.
