# Controlled Forgejo upgrade verification

| Dimension    | Evidence                                                                                                                                                              |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Completeness | Six implementation tasks, three requirements and six specification scenarios are covered by the role, documentation and integration target.                           |
| Correctness  | Packaged Ansible tasks run real Forgejo 15.0.3 and 17.0.1 with temporary SQLite instances; all six integration scenarios pass.                                        |
| Coherence    | Activation follows preflight, drain, stop, offline archive, install and acceptance. Failures retain recovery material and never downgrade the database automatically. |

No unresolved implementation findings. Final candidate identity, publication
and remote review state belong to Git and the local delivery receipts.

## Behavioral acceptance

| Scenario                                | Observed result                                                                                                                                                                                                                                                                                    |
| --------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Upgrade and complete restoration        | An account and committed repository README remain accessible after 15.0.3 to 17.0.1. The protected archive restores the old executable, matching database and Git contents; 15.0.3 starts again.                                                                                                   |
| Preflight and configuration transaction | External database and per-queue Redis settings, queue connection strings, uncovered queue paths, overlapping backup storage and encryption-key changes fail without downtime. A configuration-only activation serves its new application name; its repeated run does not restart or archive again. |
| Rejected running version                | The candidate is stopped and one completed recovery archive remains.                                                                                                                                                                                                                               |
| Queue drain failure                     | The old instance remains available; no candidate installation or archive begins.                                                                                                                                                                                                                   |
| Recovery storage failure after stop     | Installed inputs stay intact and the old instance starts again.                                                                                                                                                                                                                                    |
| Doctor rejection                        | The candidate stops, recovery remains available and an unchanged retry is blocked by the pending-activation marker.                                                                                                                                                                                |

Both release binaries report database migration markers `version=305` and
`forgejo_version=44` in these fixtures. Acceptance verifies version transition,
native database compatibility checks and restoration; it does not assume that
every release transition adds a schema migration.

Each scenario retains a sanitized JSON record in Bazel's undeclared outputs.
Only versions, schema markers, lifecycle actions, archive counts and acceptance
are retained. Generated credentials, raw application logs and recovery archives
are discarded with the private fixture. Raw build logs and copied sanitized
results remain under ignored `out/forgejo-upgrade/`.

## Validation

The following repository commands pass:

```sh
bazel test //projects/ansible_collection/test/forgejo_upgrade:upgrade_test \
  //:repo_quality_test //:buildifier_test \
  //infra/src/openspec/validation:validate_test \
  //infra/src/openspec/validation:archive_test --test_output=errors
bazel build //projects/ansible_collection/roles/forgejo:role \
  //infra/forgejo/ansible:ansible_bin
bazel build --config=lint //projects/ansible_collection/roles/forgejo:role \
  //projects/ansible_collection/plugins:plugins \
  //projects/ansible_collection/test/forgejo_upgrade:upgrade_test \
  //third_party/org_codeberg_forgejo_bin //infra/forgejo/ansible:ansible_bin \
  //tools:transitive_consumers
OPENSPEC_PROJECT=projects/ansible_collection bazel run //tools/openspec -- \
  validate --specs
```

The combined test selection contains 124 Bazel tests. Explicit pinned Ruff
checks cover the filter and both Python fixture sources, because the packaged
filter is loaded by Ansible rather than declared as a Python library. Pinned
Prettier and buildifier format the collection paths excluded by the repository's
global formatter attributes. Module dependency refresh introduces no lockfile
change; the historical binary fixture reuses its previously recorded checksum.

## Scope and limits

No deployment or live inventory was used. The service-manager fixture launches
only disposable processes and implements the systemctl operations Ansible uses.
Actual systemd notification, production data, OIDC, LFS/package workflows and
external storage are outside behavioral acceptance. Unsupported persistence
layouts fail preflight; independent controllers must serialize upgrades.

The initial `doctor --all` gate rejected optional disabled LFS maintenance.
The final gate uses nonrepairing paths, database-version and consistency checks;
the design records that evidence and decision. Full source review also corrected
individual queue override validation and queue path resolution before delivery.
