## Decisions

Use the existing checkout mirror at its previously used v4.4.0 commit. Its live
tag resolves to that exact pin. Keep checkout in the workflow because the local
action is loaded from the workspace; a remote action would add separate source
selection. Disable credential persistence and check out the event SHA with Git
LFS assets enabled, since normal repository targets consume binary source data.

Use only a push branch filter for master. Do not add release branches, PRs or
dispatch triggers. Preserve the unprivileged host runner boundary. CA installation
belongs to host Ansible; the workflow does not override Node's certificate trust.

Keep workspace discovery with its shared owner. Repair the existing subprocess
boundary to invoke Bazel rather than the removed AL command. Test the actual
compiled command against a private fake Bazel executable and retain a JSON
receipt; this checks process launching as well as phase/failure orchestration.

## Validation and activation

Run the existing action/command tests, executable acceptance fixture, pinned
runner workflow validation, repository quality and semantic lint. A branch PR
must not trigger trunk-only CI. Publication does not authorize merging: the
workflow becomes active when its reviewed change reaches master.
