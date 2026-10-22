## Context

Source: `infra/download/ansible/playbook_deploy.yaml` at
`adce867e6200dd9a617b9cc6b66d3ee2913e93b9` (PR #102).
The independent branch starts at trunk
`b84d302a7a5b15b48a6a79d0795271a394e60d51`.

## Decisions

Keep `download_*` variable names and move the two systemd templates to the role's
`templates/` directory. Create only publication roots; publishers create projects
and sites. Keep the mount root root-owned and avoid privileged ownership changes
inside publisher-writable directories. The role uses explicit `host`, `nginx`,
and `traefik` imports so attachment checks precede host setup. Filesystem tasks
live in `tasks/filesystem.yaml`. Shared roles own the service units; Traefik
uses its system-disk defaults. The second review removes custom service
overrides and undeployed migration cleanup. The caller retains `force_handlers: true`, inventory,
privilege escalation, and domain-specific templates.

## Failure cases and validation plan

Before implementation, identify extraction failures: lost handlers or ordering;
missing role files in the collection; consumer templates resolving relative to
the wrong role; accidental filesystem mutation during checks; changed systemd
unit bytes; missing-device safety checks bypassed. Verify packaged role discovery,
consumer syntax and ordered task listing, both rendered unit templates, and a
localhost missing-device failure that stops before host mutation. Use only
synthetic input and task-local output. Compare moved tasks and handlers with the
source revision. Full VM provisioning is outside this refactor's validation.

## Delivery

Publish the role against `master`. Keep the dependent consumer patch ready for
PR #102 after this role merges; no stacked PR or deployment is authorized.

## Verification evidence

The role is packaged in the existing `//infra/xcp_ng/ansible:ansible_bin`
consumer, which includes the full collection. A task-local copy of download's
inventory and routing configuration uses the thin role-based play. Both
`download_local` and `download_yandex` pass `--syntax-check` and `--list-tasks`.
A localhost execution with a nonexistent task-local device fails at the block
device assertion, before host mutation. Ansible renders both packaged systemd
templates successfully; their output matches the original templates. Task and
handler comparison preserves order and behavior, except the documented template
lookup and site-variable substitutions in the initial candidate. No full VM run or deployment is claimed.

The initial fixture pointed Ansible at the source configuration directory, so
its relative collection path could not resolve the role. Using the packaged
configuration and collection path fixed the fixture; the role was present in
the package. Logs and JSON receipts are under ignored `out/download-host/`.
The final prepared candidate's quality, lint, and OpenSpec checks are recorded
by the delivery receipt rather than assumed from these prepublication probes.

`out/download-host/consumer-conversion.patch` applies to the source revision
above. It reduces the play to the role call, removes moved templates and duplicate
defaults, and updates consumer/fixture documentation. Keep it deferred until the
role merges so PR #102 continues to validate independently against trunk.

## Review adjustment

[Review on PR #109](https://github.com/alwaldend/src/pull/109#discussion_r4110655561)
requests removal of the exact 100 GiB size assertion. Remove both `blockdev`
capacity discovery and that assertion; Terraform owns provisioning size.
Retain the attachment/block-device check before base-host configuration.
The initial preservation comparison predates this intentional change.
Rerun syntax/task expansion, the safe missing-device probe, and delivery gates
on the adjusted candidate. Update the prepared consumer documentation to stop
claiming that Ansible enforces disk capacity.

## Directory and proxy review

The publisher must be able to create any number of sites. Root ownership
operations inside publisher-writable site directories allow symlink traversal
on a subsequent run. Verdict: revise. Create only the top-level `projects/`,
`sites/`, and private `staging/` directories, with a root-owned mount root and
`follow: false`. Do not manage per-site paths. Checking a nested symlink and
then changing its owner would still race with the publisher, so avoid that
privileged traversal entirely. Keep root-owned service state outside the
publisher-owned roots.

The first directory/proxy revision removed Traefik's mount dependency and kept
Nginx's override. The second review removes that Nginx override and the obsolete
Traefik cleanup. Traefik data stays on the shared role's default `/opt/traefik`
system disk location. The prepared consumer
patch must remove its prior content-disk override. Existing ACME data migration
is not automated; this service has not been deployed in this task.

Failure cases to check before this revision: nested publisher symlinks must
never receive privileged ownership changes; a top-level symlink must not have
its target modified; the mount root must not be publisher-writable; and no
Traefik configuration/state path may require the content mount. Validate the
actual directory tasks against a task-local sentinel symlink, repeat execution,
and inspect packaged task expansion and consumer configuration. The SELinux
label declaration remains: the custom content path needs `httpd_sys_content_t`
for read access by confined Nginx, and restorecon applies that persistent policy.

The actual directory tasks passed a local Ansible fixture in an unprivileged
user namespace. A nested `sites/alwaldend.com` symlink left the sentinel target's
mode and contents unchanged; a top-level `sites` symlink was rejected without
changing its target. The mount root remained mode 0755, staging 0700, and state 0711. Repeating the normal tasks reported `changed=0`. Evidence and the exact
extracted tasks are in `out/download-host/directory-fixture/`. This checks path
handling, not a deployed VM's SELinux enforcement or multi-user isolation.

The size-check revision `531fc165` was pushed successfully, but delivery's final
verification reported newly arrived unresolved review threads. The remote head
was inspected before continuing; no failed publication was blindly retried.

## Second review pass

Remove the reusable role's single-host assertion; inventory and wrapper limits
own deployment selection. Keep the nonempty publisher-key precondition.
Move filesystem creation/mounting, content directories, SELinux labels, and
maintenance tasks into `tasks/filesystem.yaml`. Keep attachment preconditions
before base-host configuration. Use the pinned community.general 11.3.0
`btrfs_info` module to discover UUIDs, matching its device list against the
canonical device path so Yandex's by-id symlink works. The pinned implementation
uses `btrfs filesystem show -d`, including unmounted filesystems.

Remove the Nginx service override and obsolete Traefik cleanup: the user confirms
this service has not been deployed. Shared roles own their service units and
handlers. Retain the maintenance unit's own mount requirement.

Before implementing, cover these failures: selection of another Btrfs filesystem;
failure to match a by-id symlink; no matching filesystem must fail before mount;
per-host variables must work with multiple selected hosts; missing attachment
must still stop before host mutation; and moving tasks must preserve packaged
imports and the directory safety behavior already verified. Use the pinned
module with synthetic read-only Btrfs/findmnt responses for discovery/selection,
not a real host filesystem operation. Live mkfs/mount is outside this check.

Second-pass probes pass for direct-device and by-id symlink selection, with an
unrelated filesystem listed first. Missing matches fail without producing a
mount source. A two-host localhost fixture passes the publisher-key check on
both hosts and fails at the attachment assertion before shared host tasks.
Both consumer syntax checks and expanded task lists include the packaged
filesystem import. The directory tasks are structurally identical to the
previous sentinel fixture, so its evidence remains applicable. The refreshed
consumer patch applies to the source revision. Evidence is under
`out/download-host/review-pass2/`; no real disk was formatted or mounted.
