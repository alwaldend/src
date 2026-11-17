# Adoption and validation

On 2026-11-17, the existing TrueNAS 25.10.6 appliance's nine resources were
imported through `//infra/nas/truenas:tf.import`: the `default` RAIDZ2 pool,
its seven filesystem datasets, and Nextcloud lifecycle metadata (stable catalog
version 2.3.60, running). No NFS/SMB shares or periodic snapshot tasks existed.
No NAS configuration apply or encryption-key read was performed.

The pool's five members use verified disk serials. Locally set dataset
properties and encryption inheritance match the appliance. Application
configuration and encryption secrets remain on the NAS. The API key and
Terraform state use the owner's Vault namespace; DNS keeps its existing state.

## Simplified source

The settings are now inline in `storage.tf` and `apps.tf`. The generic variable
schema, separate JSON inventory, unused share/snapshot resources, mock tests,
and unused SSH target were removed at the user's request. The same nine indexed
resource addresses are preserved, so there are no state moves or imports to repeat.
Provider-canonical enum casing avoids the cosmetic updates observed during the
initial adoption.

App dataset properties are defined once on the resource. The inline table
keeps only the pool and encryption-root exceptions; explicit nulls preserve
inherited properties. The deduplicated source also produced a no-change plan.

The approved provider is pinned to 1.9.3. The Linux amd64 archive checksum was
verified against the publisher's SHA256SUMS and GitHub asset digest:
`4a413b46a7524e71cad14ca78237eeb244c92b5d3bc4348ee008612f1db1155d`.

## Verification

The normal-text live plan is the behavioral acceptance check:

```sh
TMPDIR="$PWD/out/nas-iac/t" bazel run //infra/nas/truenas:tf.plan -- -no-color -input=false
```

The simplified source returned exit 0 and reported no changes on 2026-11-17.
The raw plan was not saved. Offline validation uses:

```sh
bazel test //:repo_quality_test //:buildifier_test //infra/nas/truenas:tf_tests.fmt_test //infra/nas/truenas:tf_tests.validate_test //infra/nas/tf:tf_tests.fmt_test
bazel build --config=lint //infra/nas:al //infra/nas/truenas:tf.plan //infra/nas/truenas:tf.import //infra/nas/truenas:tf_tests.validate_test
```

The simplified source passed all 31 repository/Terraform checks and semantic
lint. Final check results, local review, candidate identity, and publication state are
recorded in the task's Git/delivery evidence. Review covers the complete PR diff,
including unchanged resource identities, deletion protection, and secret
exclusion. No external review was requested.
